package main

import (
	"encoding/json"
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
	apiTokenJSON         = kingpin.Flag("credentials", "Github api token credentials configured at the CI server, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_GITHUB_API_TOKEN").Required().String()
	gitRepoSource        = kingpin.Flag("git-repo-source", "The source of the git repository, github.com in this case.").Envar("ESTAFETTE_GIT_SOURCE").Required().String()
	gitRepoFullname      = kingpin.Flag("git-repo-fullname", "The owner and repo name of the Github repository.").Envar("ESTAFETTE_GIT_FULLNAME").Required().String()
	gitRevision          = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	estafetteBuildStatus = kingpin.Flag("estafette-build-status", "The current build status of the Estafette pipeline.").Envar("ESTAFETTE_BUILD_STATUS").Required().String()
	statusOverride       = kingpin.Flag("status-override", "Allow status property in manifest to override the actual build status.").Envar("ESTAFETTE_EXTENSION_STATUS").String()
	ciBaseURL            = kingpin.Flag("estafette-ci-server-base-url", "The base url of the ci server.").Envar("ESTAFETTE_CI_SERVER_BASE_URL").Required().String()
	estafetteBuildID     = kingpin.Flag("estafette-build-id", "The build id of this particular build.").Envar("ESTAFETTE_BUILD_ID").Required().String()
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

	// get api token from injected credentials
	var credentials []APITokenCredentials
	err := json.Unmarshal([]byte(*apiTokenJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling injected credentials: ", err)
	}
	if len(credentials) == 0 {
		log.Fatal("No credentials have been injected")
	}

	// set build status
	githubAPIClient := newGithubAPIClient()
	err := githubAPIClient.SetBuildStatus(credentials[0].AdditionalProperties.Token, *gitRepoFullname, *gitRevision, status)
	if err != nil {
		log.Fatalf("Updating Github build status failed: %v", err)
	}

	log.Println("Finished estafette-extension-github-status...")
}
