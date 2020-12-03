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
