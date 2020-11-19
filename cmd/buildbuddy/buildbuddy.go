package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/bazelbuild/bazelisk/core"
	"github.com/bazelbuild/bazelisk/repositories"
	"github.com/buildbuddy-io/buildbuddy-cli/commandline"
	"github.com/buildbuddy-io/buildbuddy-cli/parser"
	"github.com/buildbuddy-io/buildbuddy-cli/sidecar"
)

func die(exitCode int, err error) {
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func runBazelAndDie(args []string) {
	// Now run bazel.
	gcs := &repositories.GCSRepo{}
	gitHub := repositories.CreateGitHubRepo(core.GetEnvOrConfig("BAZELISK_GITHUB_TOKEN"))
	// Fetch releases, release candidates and Bazel-at-commits from GCS, forks from GitHub
	repos := core.CreateRepositories(gcs, gcs, gitHub, gcs, true)

	exitCode, err := core.RunBazelisk(args, repos)
	die(exitCode, err)
}

func main() {
	// Parse any flags (and remove them so bazel isn't confused).
	filteredOSArgs := commandline.ParseFlagsAndRewriteArgs(os.Args[1:])
	ctx := context.Background()

	// Make sure we have a home directory to work in.
	bbHome := os.Getenv("BUILDBUDDY_HOME")
	if len(bbHome) == 0 {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			die(-1, err)
		}
		bbHome = filepath.Join(userCacheDir, "buildbuddy")
	}
	if err := os.MkdirAll(bbHome, 0755); err != nil {
		die(-1, err)
	}

	bazelFlags := commandline.ExtractBazelFlags(filteredOSArgs)
	log.Printf("bazelFlags: %+v", bazelFlags)
	rcFiles := make([]string, 0)
	if !bazelFlags.NoSystemRC {
		rcFiles = append(rcFiles, "/etc/bazel.bazelrc")
		rcFiles = append(rcFiles, "%ProgramData%\bazel.bazelrc")
	}
	if !bazelFlags.NoWorkspaceRC {
		rcFiles = append(rcFiles, ".bazelrc")
	}
	if !bazelFlags.NoHomeRC {
		usr, err := user.Current()
		if err == nil {
			rcFiles = append(rcFiles, filepath.Join(usr.HomeDir, ".bazelrc"))
		}
	}
	if bazelFlags.BazelRC != "" {
		rcFiles = append(rcFiles, bazelFlags.BazelRC)
	}
	opts, err := parser.ParseRCFiles(rcFiles...)
	if err != nil {
		log.Printf("Error parsing .bazelrc file: %s", err.Error())
	}

	// Determine if cache or BES options are set.
	subcommand := commandline.GetSubCommand(filteredOSArgs)
	besBackendFlag := parser.GetRCFlagValue(opts, subcommand, bazelFlags.Config, "--bes_backend")
	remoteCacheFlag := parser.GetRCFlagValue(opts, subcommand, bazelFlags.Config, "--remote_cache")

	if besBackendFlag == "" && remoteCacheFlag == "" {
		runBazelAndDie(filteredOSArgs)
	}

	log.Printf("besBackendFlag was %q", besBackendFlag)
	log.Printf("remoteCacheFlag was %q", remoteCacheFlag)

	// Maybe update the sidecar? If we haven't recently.
	updated, err := sidecar.MaybeUpdateSidecar(ctx, bbHome)
	if err != nil {
		log.Printf("Error updating sidecar: %s", err.Error())
	}
	log.Printf("Updated sidecar: %t", updated)

	// Re(Start) the sidecar if the flags set don't match.
	sidecarArgs := make([]string, 0)
	if besBackendFlag != "" {
		sidecarArgs = append(sidecarArgs, besBackendFlag)
	}
	// TODO(tylerw): enable once cache is supported by sidecar.
	if remoteCacheFlag != "" {
		sidecarArgs = append(sidecarArgs, remoteCacheFlag)
	}

	if len(sidecarArgs) > 0 {
		sidecarSocket, err := sidecar.RestartSidecarIfNecessary(ctx, bbHome, sidecarArgs)
		if err == nil {
			filteredOSArgs = append(filteredOSArgs, fmt.Sprintf("--bes_backend=unix://%s", sidecarSocket))
		}
	}
	runBazelAndDie(filteredOSArgs)
}
