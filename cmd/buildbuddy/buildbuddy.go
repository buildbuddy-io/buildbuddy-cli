package main

import (
	"log"
	"os"

	"github.com/bazelbuild/bazelisk/core"
	"github.com/bazelbuild/bazelisk/repositories"
)

func main() {
	// For now, just do what bazelisk does!
	gcs := &repositories.GCSRepo{}
	gitHub := repositories.CreateGitHubRepo(core.GetEnvOrConfig("BAZELISK_GITHUB_TOKEN"))
	// Fetch releases, release candidates and Bazel-at-commits from GCS, forks from GitHub
	repos := core.CreateRepositories(gcs, gcs, gitHub, gcs, true)

	exitCode, err := core.RunBazelisk(os.Args[1:], repos)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}
