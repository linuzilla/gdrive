package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type RemovePermissionCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*RemovePermissionCommand)(nil)

func (RemovePermissionCommand) Command() string {
	return "remove-permission"
}

func (cmd *RemovePermissionCommand) Execute(args ...string) int {
	if len(args) != 2 {
		fmt.Println("usage: remove-permission <file-id> <permission-id>")
		return 0
	}

	if err := cmd.DriveAPI.RemovePermission(
		utils.StripQuotedString(args[0]),
		utils.StripQuotedString(args[1])); err != nil {
		fmt.Println(err)
	}

	return 0
}
