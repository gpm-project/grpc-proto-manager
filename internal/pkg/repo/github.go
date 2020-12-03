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
