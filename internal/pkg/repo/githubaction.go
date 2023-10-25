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

import "github.com/rs/zerolog/log"

// GitHubActionProvider structure implementation to work with GitHub from within
// a GitHub action environment. Some changes related to authentication makes this
// provider different from the one for classic GitHub.
type GitHubActionProvider struct {
	GHCommon
}

// NewGitHubActionProvider creates a new provider for GitHub when executed from within a GitHub action docker environment.
func NewGitHubActionProvider() (Provider, error) {
	log.Debug().Msg("Using GitHubActionProvider")
	return &GitHubCmdProvider{
		GHCommon: GHCommon{
			UseHTTPS:          true,
			SetPusherUserName: false,
			SetPusherEmail:    false,
		},
	}, nil
}
