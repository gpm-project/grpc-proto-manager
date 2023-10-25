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
	"github.com/rs/zerolog/log"
)

// GitHubCmdProvider structure with the implementation to manage a GitHub system.
// This implementation relies on leveraging the existing git command on the system. In the
// future another provider making use of the API or golang SDK should be added :)
type GitHubCmdProvider struct {
	GHCommon
}

// NewGitHubCmdProvider creates a new provider connecting to GitHub.
func NewGitHubCmdProvider() (Provider, error) {
	log.Debug().Msg("Using GitHubCmdProvider")
	return &GitHubCmdProvider{
		GHCommon: GHCommon{
			UseSSH:            true,
			SetPusherUserName: false,
			SetPusherEmail:    false,
		},
	}, nil
}
