package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type DropCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DropCommand)(nil)

func (DropCommand) Command() string {
	return `drop`
}

func (cmd *DropCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: drop <fileId>")
	} else {
		cmd.DriveAPI.DropFile(args[0])
	}
	return 0
}
