package updater

import (
	"bufio"
	"context"
	"encoding/json"
	"eq-expansion-switcher/internal/download"
	"eq-expansion-switcher/internal/env"
	"eq-expansion-switcher/internal/unzip"
	"fmt"
	"github.com/google/go-github/v41/github"
	"github.com/mattn/go-isatty"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// Service is a service that checks for updates to the app
type Service struct {
	// this is the package.json embedded in the binary which contains the app version
	packageJson []byte
}

// NewService creates a new updater service
func NewService(packageJson []byte) *Service {
	return &Service{
		packageJson: packageJson,
	}
}

// EnvResponse is the response from the env endpoint
type EnvResponse struct {
	Env     string `json:"env"`
	Version string `json:"version"`
}

// PackageJson is the package.json file
type PackageJson struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
}

// getAppVersion gets the app version from the package.json embedded in the binary
func (s Service) getAppVersion() (error, EnvResponse) {
	var pkg PackageJson
	err := json.Unmarshal(s.packageJson, &pkg)
	if err != nil {
		return err, EnvResponse{}
	}

	return nil, EnvResponse{
		Env:     env.Get("APP_ENV", "local"),
		Version: pkg.Version,
	}
}

const (
	appFileName       = "eq-expansion-switcher"
	appExecutableName = "ProjectEQ Expansion Switcher"
	appName           = "ProjectEQ Expansion Switcher"
	org               = "Akkadius"
	repo              = "eq-expansion-switcher"
)

func (s Service) IsUpdateAvailable() bool {
	// get executable name and path
	executableName := filepath.Base(os.Args[0])
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	executablePath := filepath.Dir(ex)

	// check if a .old version exists, delete it if does
	oldExecutable := filepath.Join(executablePath, fmt.Sprintf("%s.old", executableName))
	if _, err := os.Stat(oldExecutable); err == nil {
		e := os.Remove(oldExecutable)
		if e != nil {
			log.Fatal(e)
		}
	}

	// if being ran from go run main.go
	if executableName == "main.exe" || executableName == "main" {
		fmt.Println("[Update] Running as go run main.go, ignoring...")
		return false
	}

	// internet connection check
	if !isconnected() {
		fmt.Printf("[Update] Not connected to the internet\n")
		return false
	}

	fmt.Printf("[Update] Checking for updates...\n")
	fmt.Printf("[Update] Running as binary [%v]\n", executableName)
	debug(fmt.Sprintf("[Update] Running as executablePath [%v]\n", executablePath))

	// get releases
	client := github.NewClient(&http.Client{Timeout: 5 * time.Second})
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), org, repo)
	if err != nil {
		log.Println(err)
		return false
	}

	// get app version
	err, e := s.getAppVersion()
	if err != nil {
		log.Println(err)
	}

	localVersion := fmt.Sprintf("v%v", e.Version)
	releaseVersion := *release.TagName

	// already up to date
	if releaseVersion == localVersion {
		fmt.Printf("[Update] %s is already up to date @ [%v]\n", appName, localVersion)
		return false
	}

	fmt.Printf("Local version [%s] latest [%v]\n", localVersion, releaseVersion)

	return localVersion != releaseVersion
}

// CheckForUpdates checks for updates to the app
func (s Service) CheckForUpdates() {
	// get executable name and path
	executableName := filepath.Base(os.Args[0])
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	executablePath := filepath.Dir(ex)

	// check if a .old version exists, delete it if does
	oldExecutable := filepath.Join(executablePath, fmt.Sprintf("%s.old", executableName))
	if _, err := os.Stat(oldExecutable); err == nil {
		e := os.Remove(oldExecutable)
		if e != nil {
			log.Fatal(e)
		}
	}

	// if being ran from go run main.go
	if executableName == "main.exe" || executableName == "main" {
		fmt.Println("[Update] Running as go run main.go, ignoring...")
		return
	}

	// internet connection check
	if !isconnected() {
		fmt.Printf("[Update] Not connected to the internet\n")
		return
	}

	fmt.Printf("[Update] Checking for updates...\n")
	fmt.Printf("[Update] Running as binary [%v]\n", executableName)
	debug(fmt.Sprintf("[Update] Running as executablePath [%v]\n", executablePath))

	// get releases
	client := github.NewClient(&http.Client{Timeout: 5 * time.Second})
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), org, repo)
	if err != nil {
		log.Println(err)
		return
	}

	// get app version
	err, e := s.getAppVersion()
	if err != nil {
		log.Println(err)
	}

	localVersion := fmt.Sprintf("v%v", e.Version)
	releaseVersion := *release.TagName

	// already up to date
	if releaseVersion == localVersion {
		fmt.Printf("[Update] %s is already up to date @ [%v]\n", appName, localVersion)
		return
	}

	fmt.Printf("Local version [%s] latest [%v]\n", localVersion, releaseVersion)

	for _, asset := range release.Assets {

		assetName := *asset.Name
		downloadUrl := *asset.BrowserDownloadURL
		targetFileNameZipped := fmt.Sprintf("%s-%s-%s.zip", appFileName, runtime.GOOS, runtime.GOARCH)
		// targetFileName := fmt.Sprintf("%s-%s-%s", appFileName, runtime.GOOS, runtime.GOARCH)

		debug(fmt.Sprintf("[Update] Looping assets assetName [%v] targetFileNameZipped [%v]\n", assetName, targetFileNameZipped))

		if assetName == targetFileNameZipped {
			fmt.Printf("Found matching release [%s]\n", assetName)

			// download
			file := path.Base(downloadUrl)
			downloadPath := filepath.Join(os.TempDir(), file)
			err := download.WithProgress(downloadPath, downloadUrl)
			if err != nil {
				log.Println(err)
			}

			// unzip
			tempFileZipped := filepath.Join(os.TempDir(), targetFileNameZipped)
			uz := unzip.New(tempFileZipped, os.TempDir())
			err = uz.Extract()
			if err != nil {
				log.Println(err)
			}

			// rename running process to .old
			err = os.Rename(
				filepath.Join(executablePath, executableName),
				filepath.Join(executablePath, fmt.Sprintf("%s.old", executableName)),
			)
			if err != nil {
				log.Fatal(err)
			}

			// relink
			tempFile := filepath.Join(os.TempDir(), appExecutableName)

			newAppExecutableName := appExecutableName
			if runtime.GOOS == "windows" {
				newAppExecutableName = appExecutableName + ".exe"
			}
			newExecutable := filepath.Join(executablePath, newAppExecutableName)

			newExecutableTemp, err := exec.LookPath(tempFile)
			if err != nil {
				log.Println(err)
			}

			fmt.Println("tempFile " + tempFile)
			fmt.Println("newExecutableTemp " + newExecutableTemp)
			fmt.Println("newExecutable " + newExecutable)

			err = moveFile(newExecutableTemp, newExecutable)
			if err != nil {
				log.Println(err)
			}

			err = os.Chmod(newExecutable, 0755)
			if err != nil {
				log.Println(err)
			}

			// if terminal, wait for user input
			if isatty.IsTerminal(os.Stdout.Fd()) {
				fmt.Println("")
				fmt.Printf("[Update] [%s] updated to version [%s] you must relaunch [%s] manually\n", appName, releaseVersion, appName)
				fmt.Println("")
				fmt.Printf("Press [Enter] to exit [%s]...", appName)
				fmt.Println("")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			} else {
				fmt.Printf("[Update] %s updated to version [%s] you must relaunch %s manually\n", appName, releaseVersion, appName)
			}

			os.Exit(0)
		}
	}
}

func debug(msg string) {
	if len(os.Getenv("DEBUG")) > 0 {
		fmt.Printf("[Debug] " + msg)
	}
}
