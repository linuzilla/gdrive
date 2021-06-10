package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type ShareInfoCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ShareInfoCommand)(nil)

func (ShareInfoCommand) Command() string {
	return "share-info"
}

func (cmd *ShareInfoCommand) Execute(args ...string) int {

	if len(args) != 1 {
		fmt.Println("usage: share-info <file-id>")
		return 0
	}

	if err := cmd.DriveAPI.ShareInfo(utils.StripQuotedString(args[0])); err != nil {
		fmt.Println(err)
	}

	return 0
}
