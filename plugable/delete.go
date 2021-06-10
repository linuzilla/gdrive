package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type DeleteCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DeleteCommand)(nil)

func (DeleteCommand) Command() string {
	return "delete"
}

func (cmd *DeleteCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: delete <fileId>")
	} else {
		cmd.DriveAPI.Delete(args[0])
	}
	return 0
}
