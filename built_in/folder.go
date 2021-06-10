package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type FolderCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*FolderCommand)(nil)

func (FolderCommand) Command() string {
	return `cd`
}

func (cmd *FolderCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: cd <folderId>")
	} else {
		cmd.DriveAPI.Pushd()
		cmd.DriveAPI.ChangeFolder(utils.StripQuotedString(args[0]))
	}
	return 0
}
