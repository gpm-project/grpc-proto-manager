package repo

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// NoTagsFoundErrorMsg with the error to be reported if no previous tags are found.
const NoTagsFoundErrorMsg = "No names found, cannot describe anything"

// GHCommon structure with common operation over GitHub. Notice that depending on the environment.
// some options may apply.
type GHCommon struct {
	CmdUtils
	// UseSSH determines if clone operations will use SSH credentials.
	UseSSH bool
	// UseHTTPS determines if clone operations will use HTTPS credentials.
	UseHTTPS bool
	// PersonalAccessToken required to operate with GitHub. This is used in conjunction with UseHTTPS if non empty.
	PersonalAccessToken string
	// SetPusherUserName determines if the user name of the pusher actor needs to be set.
	SetPusherUserName bool
	// PusherUserName with the name to use as commiter.
	PusherUserName string
	// SetPusherEmail determines if the email of the pusher actor needs to be set.
	SetPusherEmail bool
	// PusherEmail with the email of the commiter.
	PusherEmail string
}

// ConfigurePusher prepares the system to use a particular username/email to appear as the pusher of the commits.
func (ghc *GHCommon) ConfigurePusher(username string, email string, accessToken string) error {
	log.Debug().Str("username", username).Str("email", email).Str("accessToken", strings.Repeat("*", len(accessToken))).Msg("setting pusher information")
	// Notice that in GitHub, it is recommended to setup this information per repository, therefore this action
	// will be executed on per-repo basis before the commit & push information.

	// To setup this, the following commands will be issued:
	// git config -f /tmp/grpc-internal-agenda-go/.git/config user.name \"Your Name\"
	// git config -f /tmp/grpc-internal-agenda-go/.git/config user.email "my.name@server.com"

	// Alternatively, each git command may be configured with a given user name. However, given that this is executed from
	// a local temporal copy, there should not be any collateral impact in configuring the local repo. The alternative would be:
	// git -c user.name="Your name" -c user.email="my.name@server.com" commit -m "Commit message" ...
	if username != "" {
		ghc.SetPusherUserName = true
		ghc.PusherUserName = username
	}
	if email != "" {
		ghc.SetPusherEmail = true
		ghc.PusherEmail = email
	}
	ghc.PersonalAccessToken = accessToken
	return nil
}

// GetRepoURL builds the URL require for clone and commit operations.
func (ghc *GHCommon) GetRepoURL(organization string, repoName string) (string, error) {
	if ghc.UseSSH {
		// git@github.com:dhiguero/go-template.git
		return fmt.Sprintf("git@github.com:%s/%s.git", organization, repoName), nil
	} else if ghc.UseHTTPS {
		if ghc.PersonalAccessToken == "" {
			// https://github.com/dhiguero/go-template.git
			return fmt.Sprintf("https://github.com/%s/%s.git", organization, repoName), nil
		}
		// https://{TOKEN}@github.com/dhiguero/go-template.git
		return fmt.Sprintf("https://%s@github.com/%s/%s.git", ghc.PersonalAccessToken, organization, repoName), nil
	}
	return "", fmt.Errorf("cannot obtain target repo URL. Set useSSH or UseHTTPS")
}

// Clone a given repository to a path
func (ghc *GHCommon) Clone(repoURL string, outputPath string) error {
	log.Debug().Str("repoURL", repoURL).Str("outputPath", outputPath).Msg("cloning repository")
	// TODO Check output path exists.
	cmdArgs := []string{"clone", repoURL, outputPath}

	cmd := exec.Command("git", cmdArgs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to clone repo %s due to %w", repoURL, err)
	}

	log.Debug().Str("output", string(stdoutStderr)).Msg("repo successfully cloned")
	return nil
}

// GetLastVersion obtains the latest version of the repo.
// git describe --abbrev=0 --tags
func (ghc *GHCommon) GetLastVersion(repoPath string) (*Version, error) {
	log.Debug().Str("repoPath", repoPath).Msg("obtaining latest tag")
	// TODO Check output path exists.
	cmdArgs := []string{"describe", "--abbrev=0", "--tags"}

	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = repoPath
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(stdoutStderr), NoTagsFoundErrorMsg) {
			return EmptyVersion(), nil
		}
		return nil, fmt.Errorf("unable to obtain latest tag from repo %s due to %w, %s", repoPath, err, string(stdoutStderr))
	}

	log.Debug().Str("output", string(stdoutStderr)).Msg("repo successfully cloned")
	return FromTag(string(stdoutStderr))
}

// SetPusherInfo sets the pusher information of the local repository.
func (ghc *GHCommon) SetPusherInfo(repoPath string) error {

	if ghc.SetPusherUserName {
		// Set the name
		// git config -f /tmp/grpc-internal-agenda-go/.git/config user.name \"Your Name\"
		userNameArgs := []string{"config", "user.name", ghc.PusherUserName}
		_, err := ghc.execCmd("git", userNameArgs, repoPath)
		if err != nil {
			return err
		}
	}

	if ghc.SetPusherEmail {
		// Set the email
		// git config -f /tmp/grpc-internal-agenda-go/.git/config user.email "my.name@server.com"
		userEmailArgs := []string{"config", "user.email", ghc.PusherEmail}
		_, err := ghc.execCmd("git", userEmailArgs, repoPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Publish the changes and create a new version tag.
func (ghc *GHCommon) Publish(repoPath string, newVersion *Version) error {
	log.Debug().Str("repoPath", repoPath).Str("version", newVersion.String()).Msg("publishing version")

	err := ghc.SetPusherInfo(repoPath)
	if err != nil {
		return err
	}

	// TODO improve commit message
	// Add all new files
	addCmdArgs := []string{"add", "-A"}
	_, err = ghc.execCmd("git", addCmdArgs, repoPath)
	if err != nil {
		return err
	}
	// Commit changes
	commitCmdArgs := []string{"commit", "-a", "-m", fmt.Sprintf("gpm automatic publish")}
	_, err = ghc.execCmd("git", commitCmdArgs, repoPath)
	if err != nil {
		return err
	}

	// Push changes
	pushCmdArgs := []string{"push"}
	_, err = ghc.execCmd("git", pushCmdArgs, repoPath)
	if err != nil {
		return err
	}
	// Create new tag
	// tag -a v1.4 -m "my version 1.4"
	tagCmdArgs := []string{"tag", "-a", newVersion.String(), "-m", fmt.Sprintf("new version %s generated by GPM", newVersion.String())}
	_, err = ghc.execCmd("git", tagCmdArgs, repoPath)
	if err != nil {
		return err
	}
	// Push the tags
	// git push origin --tags
	pushTagCmdArgs := []string{"push", "origin", "--tags"}
	_, err = ghc.execCmd("git", pushTagCmdArgs, repoPath)
	if err != nil {
		return err
	}
	return nil
}
