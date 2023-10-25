/**
 * Copyright 2023 GPM Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repo

import (
	"fmt"
	"strings"
)

// RepositoryType defines a type for all supported repositories.
type RepositoryType int

const (
	// GitHub cloud.
	GitHub RepositoryType = iota
	// GitHubAction corresponds to a git provider being executed from within a GitHub Action.
	GitHubAction
)

// RepositoryTypeToString map associating type to its string representation.
var RepositoryTypeToString = map[RepositoryType]string{
	GitHub:       "github",
	GitHubAction: "githubaction",
}

// RepositoryTypeToEnum map associating string representation with type.
var RepositoryTypeToEnum = map[string]RepositoryType{
	"github":       GitHub,
	"githubaction": GitHubAction,
}

// Provider defines the common interface for different repository managers (e.g., GitHub)
type Provider interface {
	// ConfigurePusher prepares the system to use a particular username/email to appear as the pusher of the commits.
	ConfigurePusher(username string, email string, accessToken string) error
	// GetRepoURL builds the URL require for clone and commit operations.
	GetRepoURL(organization string, repoName string) (string, error)
	// Clone a given repository to a path
	Clone(repoURL string, outputPath string) error
	// GetLastVersion obtains the latest version of the repo.
	GetLastVersion(repoPath string) (*Version, error)
	// Publish the changes and create a new version tag.
	Publish(repoPath string, newVersion *Version) error
}

// NewRepoProvider factory method to instantiate a repository provider for a given system.
func NewRepoProvider(repoProviderName string) (Provider, error) {
	provider, exists := RepositoryTypeToEnum[strings.ToLower(repoProviderName)]
	if !exists {
		return nil, fmt.Errorf("Provider not found for %s", repoProviderName)
	}
	switch provider {
	case GitHub:
		return NewGitHubCmdProvider()
	case GitHubAction:
		return NewGitHubActionProvider()
	}
	return nil, fmt.Errorf("No provider implementation found for %s", repoProviderName)
}
