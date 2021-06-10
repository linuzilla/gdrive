package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
)

type PermissionInfoCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*PermissionInfoCommand)(nil)

func (PermissionInfoCommand) Command() string {
	return `permission-info`
}

func (cmd *PermissionInfoCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("permission-info <drive-id>")
	} else {
		if err := cmd.DriveAPI.PermissionInfo(utils.StripQuotedString(args[0]), func(permission *drive.Permission) {
			fmt.Printf("  %s (%s) [ %s ]\n", permission.EmailAddress, permission.Role, permission.Id)
		}); err != nil {
			fmt.Println(err)
		}
	}
	return 0
}
