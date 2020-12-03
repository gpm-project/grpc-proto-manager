package repo

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// CmdUtils structure with helper methods to execute commands.
type CmdUtils struct {
}

// execCmd executes a given command and returns the output if successful.
func (cu *CmdUtils) execCmd(cmd string, args []string, workingDir string) (string, error) {
	toExecute := exec.Command("git", args...)
	toExecute.Dir = workingDir
	stdoutStderr, err := toExecute.CombinedOutput()
	if err != nil {
		return string(stdoutStderr), fmt.Errorf("unable to execute command %s due to %w, %s", cmd, err, string(stdoutStderr))
	}
	log.Debug().Str("output", string(stdoutStderr)).Msg("execution finished")
	return string(stdoutStderr), nil
}
