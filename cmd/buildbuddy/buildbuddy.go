package main

import (
	"context"
	"flag"
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

	bblog "github.com/buildbuddy-io/buildbuddy-cli/logging"
)

var (
	disable = flag.Bool("disable_buildbuddy", false, "If true, disable buildbuddy functionality and just run bazel.")
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

func parseBazelRCs(bazelFlags *commandline.BazelFlags) []*parser.BazelOption {
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
		bblog.Printf("Error parsing .bazelrc file: %s", err.Error())
		return nil
	}
	return opts

}

func main() {
	// Parse any flags (and remove them so bazel isn't confused).
	filteredOSArgs := commandline.ParseFlagsAndRewriteArgs(os.Args[1:])

	if *disable {
		bblog.Printf("Buildbuddy was disabled, just running bazel.")
		runBazelAndDie(filteredOSArgs)
	}

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
	opts := parseBazelRCs(bazelFlags)

	// Determine if cache or BES options are set.
	subcommand := commandline.GetSubCommand(filteredOSArgs)
	besBackendFlag := parser.GetRCFlagValue(opts, subcommand, bazelFlags.Config, "--bes_backend")
	remoteCacheFlag := parser.GetRCFlagValue(opts, subcommand, bazelFlags.Config, "--remote_cache")

	if besBackendFlag == "" && remoteCacheFlag == "" {
		runBazelAndDie(filteredOSArgs)
	}

	bblog.Printf("besBackendFlag was %q", besBackendFlag)
	bblog.Printf("remoteCacheFlag was %q", remoteCacheFlag)

	// Maybe update the sidecar? If we haven't recently.
	updated, err := sidecar.MaybeUpdateSidecar(ctx, bbHome)
	if err != nil {
		bblog.Printf("Error updating sidecar: %s", err.Error())
	}
	if updated {
		bblog.Printf("Updated sidecar: %t", updated)
	}

	// Re(Start) the sidecar if the flags set don't match.
	sidecarArgs := make([]string, 0)
	if besBackendFlag != "" {
		sidecarArgs = append(sidecarArgs, besBackendFlag)
	}
	if remoteCacheFlag != "" {
		sidecarArgs = append(sidecarArgs, remoteCacheFlag)
	}

	if len(sidecarArgs) > 0 {
		sidecarSocket, err := sidecar.RestartSidecarIfNecessary(ctx, bbHome, sidecarArgs)
		// TODO(tylerw): test the sidecar connection before passing it to bazel.
		if err == nil {
			if besBackendFlag != "" {
				filteredOSArgs = append(filteredOSArgs, fmt.Sprintf("--bes_backend=unix://%s", sidecarSocket))
			}
			// // TODO(tylerw): enable once cache is supported by sidecar.
			// if besBackendFlag != "" {
			// 	filteredOSArgs = append(filteredOSArgs, fmt.Sprintf("--remote_cache=unix://%s", sidecarSocket))
			// }
		}
	}
	bblog.Printf("Rewrote bazel command line to: %s", filteredOSArgs)
	runBazelAndDie(filteredOSArgs)
}
