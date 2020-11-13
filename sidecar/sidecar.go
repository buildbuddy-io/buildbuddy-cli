package sidecar

import (
	"context"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/mod/semver"

	"github.com/buildbuddy-io/buildbuddy-cli/download"
)

const (
	windowsOSName        = "windows"
	windowsFileExtension = ".exe"

	// The name of the directory that will contain
	// tag/version dirs (which each contain a sidecar binary)
	// inside of the buildbuddy dir.
	sidecarsSubdir = "sidecars"

	// The name of the file that contains a timestamp indicating
	// when we last checked for an update.
	lastCheckedForUpdateFileName = "last_checked_for_update"

	// How long to wait between checks for a new sidecar version.
	timeBetweenUpdateChecks = 24 * time.Hour

	sockPrefix = "sidecar-"
)

func getSidecarBinaryName() string {
	extension := ""
	if runtime.GOOS == windowsOSName {
		extension = windowsFileExtension
	}
	sidecarName := fmt.Sprintf("sidecar-%s-%s%s", runtime.GOOS, runtime.GOARCH, extension)
	return sidecarName
}

func getLatestInstalledSidecarVersion(sidecarDir, sidecarName string) string {
	files, err := ioutil.ReadDir(sidecarDir)
	if err != nil {
		return ""
	}
	latestVersion := ""
	for _, f := range files {
		version := f.Name()
		binPath := filepath.Join(sidecarDir, version, sidecarName)
		if _, err := os.Stat(binPath); !os.IsNotExist(err) && f.IsDir() {
			if semver.Compare(version, latestVersion) > 0 {
				latestVersion = version
			}
		}
	}
	return latestVersion
}

func getlastUpdateCheck(bbHomeDir string) (time.Time, error) {
	lastCheckedForUpdateFilePath := filepath.Join(bbHomeDir, lastCheckedForUpdateFileName)
	content, err := ioutil.ReadFile(lastCheckedForUpdateFilePath)
	if err != nil {
		return time.Time{}, err
	}
	i, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

func setLastUpdateCheck(bbHomeDir string) error {
	lastCheckedForUpdateFilePath := filepath.Join(bbHomeDir, lastCheckedForUpdateFileName)
	f, err := os.OpenFile(lastCheckedForUpdateFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10))); err != nil {
		return err
	}
	return f.Close()
}

func MaybeUpdateSidecar(ctx context.Context, bbHomeDir string) (bool, error) {
	sidecarDir := filepath.Join(bbHomeDir, sidecarsSubdir)
	if err := os.MkdirAll(sidecarDir, 0755); err != nil {
		return false, err
	}

	// Figure out appropriate os/arch for this machine.
	sidecarName := getSidecarBinaryName()

	// Check what is the latest sidecar we have installed.
	latestInstalledVersion := getLatestInstalledSidecarVersion(sidecarDir, sidecarName)

	// We're done If:
	//  1) we already have a version
	//  2) we've checked recently and
	//  3) checking is not being forced
	forceUpdateCheck := os.Getenv("BB_ALWAYS_CHECK_FOR_UPDATES") != ""
	lastChecked, _ := getlastUpdateCheck(bbHomeDir)
	if latestInstalledVersion != "" && time.Since(lastChecked) < timeBetweenUpdateChecks && !forceUpdateCheck {
		log.Printf("Not checking for update, last checked at %s", lastChecked)
		return false, nil
	}

	// Check what is the latest sidecar on github.
	bin, err := download.GetLatestSidecarFromGithub(ctx, sidecarName)
	if err != nil {
		return false, err
	}

	setLastUpdateCheck(bbHomeDir) // ignore error; if it doesn't work we'll check again.

	// If there is an update available, download it.
	if semver.Compare(bin.Version(), latestInstalledVersion) > 0 {
		log.Printf("Version %q is available, downloading...", bin.Version())
		sidecarOutputDir := filepath.Join(sidecarDir, bin.Version())
		if err := os.MkdirAll(sidecarOutputDir, 0755); err != nil {
			return false, err
		}
		sidecarOutputPath := filepath.Join(sidecarOutputDir, sidecarName)
		if err := bin.Download(ctx, sidecarOutputPath); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func hashStrings(in []string) string {
	data := []byte{}
	for _, i := range in {
		data = append(data, []byte(i)...)
	}
	u := crc32.ChecksumIEEE(data)
	return fmt.Sprintf("%d", u)
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func startBackgroundProcess(cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	return c.Start()
}

func RestartSidecarIfNecessary(ctx context.Context, bbHomeDir string) (string, error) {
	sidecarName := getSidecarBinaryName()
	sidecarDir := filepath.Join(bbHomeDir, sidecarsSubdir)
	latestInstalledVersion := getLatestInstalledSidecarVersion(sidecarDir, sidecarName)
	binPath := filepath.Join(sidecarDir, latestInstalledVersion, sidecarName)

	cmd := binPath
	args := []string{}

	sockName := sockPrefix + hashStrings(append(args, cmd)) + ".sock"
	sockPath := filepath.Join("/tmp/", sockName)

	// Check if a process is already running with this sock.
	// If one is, we're all done!
	if pathExists(sockPath) {
		return sockPath, nil
	}

	// This is where we'll listen for bazel traffic
	args = append(args, fmt.Sprintf("--listen_addr=unix://%s", sockPath))
	if err := startBackgroundProcess(cmd, args); err != nil {
		return "", err
	}
	return sockPath, nil
}
