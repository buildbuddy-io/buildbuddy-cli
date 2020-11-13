package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bazelbuild/bazelisk/core"
	"github.com/bazelbuild/bazelisk/repositories"
	"github.com/buildbuddy-io/buildbuddy-cli/sidecar"
)

func die(exitCode int, err error) {
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func main() {
	ctx := context.Background()

	// Make sure we have a home directory to work in.
	bbHome := core.GetEnvOrConfig("BUILDBUDDY_HOME")
	if len(bbHome) == 0 {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			die(-1, err)
		}
		bbHome = filepath.Join(userCacheDir, "buildbuddy")
	}

	updated, err := sidecar.MaybeUpdateSidecar(ctx, bbHome)
	if err != nil {
		log.Printf("Error updating sidecar: %s", err.Error())
	}
	log.Printf("Updated sidecar: %t", updated)

	sidecarSocket, err := sidecar.RestartSidecarIfNecessary(ctx, bbHome)
	if err == nil {
		os.Args = append(os.Args, fmt.Sprintf("--bes_backend=unix://%s", sidecarSocket))
	}

	// Run bazel.
	gcs := &repositories.GCSRepo{}
	gitHub := repositories.CreateGitHubRepo(core.GetEnvOrConfig("BAZELISK_GITHUB_TOKEN"))
	// Fetch releases, release candidates and Bazel-at-commits from GCS, forks from GitHub
	repos := core.CreateRepositories(gcs, gcs, gitHub, gcs, true)

	exitCode, err := core.RunBazelisk(os.Args[1:], repos)
	die(exitCode, err)

}
