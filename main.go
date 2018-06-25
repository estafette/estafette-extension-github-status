package main

import (
	"log"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"
)

var (
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	githubAPIAccessToken = kingpin.Flag("github-api-token", "The time-limited access token to access the Github api.").Envar("ESTAFETTE_GITHUB_API_TOKEN").Required().String()
	gitRepoFullname      = kingpin.Flag("git-repo-fullname", "The owner and repo name of the Github repository.").Envar("ESTAFETTE_GIT_NAME").Required().String()
	gitRevision          = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	estafetteBuildStatus = kingpin.Flag("estafette-build-status", "The current build status of the Estafette pipeline.").Envar("ESTAFETTE_BUILD_STATUS").Required().String()
	statusOverride       = kingpin.Flag("status-override", "Allow status property in manifest to override the actual build status.").Envar("ESTAFETTE_EXTENSION_STATUS").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting estafette-extension-github-status version %v...", version)

	// check if there's a status override
	status := *estafetteBuildStatus
	if *statusOverride != "" {
		status = *statusOverride
	}

	// set build status
	githubAPIClient := newGithubAPIClient()
	err := githubAPIClient.SetBuildStatus(*githubAPIAccessToken, *gitRepoFullname, *gitRevision, status)
	if err != nil {
		log.Fatalf("Updating Github build status failed: %v", err)
	}

	log.Println("Finished estafette-extension-github-status...")
}
