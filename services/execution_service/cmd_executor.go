package execution_service

import (
	"bytes"
	"github.com/linuzilla/go-logger"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type commandExecutorImpl struct {
	cmd []string
}

var pattern = regexp.MustCompile(`^(.*)#\{([-a-zA-F]+)\}(.*)$`)

func executeCommand(cmd []string, dict map[string]string) *exec.Cmd {
	command := make([]string, len(cmd))

	for i, arg := range cmd {
		command[i] = arg
		if matched := pattern.FindStringSubmatch(arg); matched != nil {
			if val, found := dict[matched[2]]; found {
				command[i] = matched[1] + val + matched[3]
			}
		}
	}

	logger.Debug(command)
	return exec.Command(command[0], command[1:]...)
}

func (exe *commandExecutorImpl) Pipe(dict map[string]string, reader io.Reader, writer io.Writer) error {
	cmd := executeCommand(exe.cmd, dict)
	cmd.Stdin = reader
	cmd.Stdout = writer

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (exe *commandExecutorImpl) Exec(dict map[string]string, reader io.Reader) (*bytes.Buffer, error) {
	cmd := executeCommand(exe.cmd, dict)
	//command := make([]string, len(exe.cmd))
	//for i, arg := range exe.cmd {
	//	command[i] = arg
	//	if matched := pattern.FindStringSubmatch(arg); matched != nil {
	//		if val, found := dict[matched[2]]; found {
	//			command[i] = matched[1] + val + matched[3]
	//		}
	//	}
	//}
	//cmd := exec.Command(command[0], command[1:]...)

	if reader != nil {
		cmd.Stdin = reader
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	//if err := cmd.Wait(); err != nil {
	//	if exiterr, ok := err.(*exec.ExitError); ok {
	//		// The program has exited with an exit code != 0
	//
	//		// This works on both Unix and Windows. Although package
	//		// syscall is generally platform dependent, WaitStatus is
	//		// defined for both Unix and Windows and in both cases has
	//		// an ExitStatus() method with the same signature.
	//		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
	//			log.Printf("Exit Status: %d", status.ExitStatus())
	//		}
	//	} else {
	//		log.Printf("cmd.Wait: %v", err)
	//	}
	//} else {
	//	fmt.Println("success");
	//}
	return &out, nil
}

func newExecutor(cmd string) CommandExecutor {
	return &commandExecutorImpl{
		cmd: strings.Fields(cmd),
	}
}
