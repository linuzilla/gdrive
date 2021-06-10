package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	cmdline_service "github.com/linuzilla/go-cmdline"
)

type NewFolderCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*NewFolderCommand)(nil)

func (NewFolderCommand) Command() string {
	return "create-folder"
}

func (cmd *NewFolderCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: create-folder <folder-name>")
	} else {
		cmd.DriveAPI.CreateFolder(utils.StripQuotedString(args[0]))
	}
	return 0
}
