package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBuildStatus(t *testing.T) {

	t.Run("Succeeded", func(t *testing.T) {

		// act
		githubAPIClient := newGithubAPIClient()
		err := githubAPIClient.SetBuildStatus("v1.1c85af3bbfbe393055c74a22b638b4a52238d05d", "estafette/estafette-extension-bitbucket-status", "a66df8a71b566e712cd9e4096469ae903d05ad12", "succeeded")

		assert.Nil(t, err)
	})

	// t.Run("Failed", func(t *testing.T) {

	// 	// act
	// 	githubAPIClient := newGithubAPIClient()
	// 	err := githubAPIClient.SetBuildStatus("v1.1c85af3bbfbe393055c74a22b638b4a52238d05d", "estafette/estafette-extension-bitbucket-status", "a66df8a71b566e712cd9e4096469ae903d05ad12", "failed")

	// 	assert.Nil(t, err)
	// })

	// t.Run("Pending", func(t *testing.T) {

	// 	// act
	// 	githubAPIClient := newGithubAPIClient()
	// 	err := githubAPIClient.SetBuildStatus("v1.1c85af3bbfbe393055c74a22b638b4a52238d05d", "estafette/estafette-extension-bitbucket-status", "a66df8a71b566e712cd9e4096469ae903d05ad12", "pending")

	// 	assert.Nil(t, err)
	// })
}
