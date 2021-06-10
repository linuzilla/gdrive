package built_in

import (
	"fmt"
	"github.com/linuzilla/go-cmdline"
	"os"
)

type PwdCommand struct {
}

var _ cmdline_service.CommandInterface = (*PwdCommand)(nil)

func (PwdCommand) Command() string {
	return `pwd`
}

func (cmd *PwdCommand) Execute(args ...string) int {
	if currentDir, err := os.Getwd(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(currentDir)
	}
	return 0
}
