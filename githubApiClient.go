package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sethgrid/pester"
)

// GithubAPIClient allows to communicate with the Github api
type GithubAPIClient interface {
	SetBuildStatus(string, string, string, string) error
}

type githubAPIClientImpl struct {
}

func newGithubAPIClient() GithubAPIClient {
	return &githubAPIClientImpl{}
}

type buildStatusRequestBody struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

// SetBuildStatus sets the build status for a specific revision
func (gh *githubAPIClientImpl) SetBuildStatus(accessToken, repoFullname, gitRevision, status string) (err error) {

	// https://developer.github.com/v3/repos/statuses/
	// estafette status: succeeded|failed|pending
	// github stat: success|failure|error|pending

	state := "success"
	switch status {
	case "succeeded":
		state = "success"

	case "failed":
		state = "failure"

	case "pending":
		state = "pending"
	}

	logsURL := fmt.Sprintf(
		"%vpipelines/%v/%v/builds/%v/logs",
		*estafetteBuildID,
		*gitRepoSource,
		repoFullname,
		*estafetteBuildID,
	)

	params := buildStatusRequestBody{
		State:     state,
		TargetURL: logsURL,
	}

	_, err = callGithubAPI("POST", fmt.Sprintf("https://api.github.com/repos/%v/statuses/%v", repoFullname, gitRevision), params, "token", accessToken)

	return
}

func callGithubAPI(method, url string, params interface{}, authorizationType, token string) (body []byte, err error) {

	// convert params to json if they're present
	var requestBody io.Reader
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return body, err
		}
		requestBody = bytes.NewReader(data)
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("Authorization", fmt.Sprintf("%v %v", authorizationType, token))
	request.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Printf("Deserializing response for '%v' Github api call failed. Body: %v. Error: %v", url, string(body), err)
		return
	}

	log.Printf("Received successful response for '%v' Github api call", url)

	return
}
