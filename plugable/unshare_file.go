package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type UnshareFileCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*UnshareFileCommand)(nil)

func (UnshareFileCommand) Command() string {
	return "unshare-file"
}

func (cmd *UnshareFileCommand) Execute(args ...string) int {
	if len(args) != 2 {
		fmt.Println("usage: unshare-file <file-id> <email>")

		fmt.Println("Roles:")
		for _, r := range models.AllRoles {
			fmt.Printf("  %s\n", string(r))
		}

		return 0
	}

	if err := cmd.DriveAPI.UnshareFileWith(
		utils.StripQuotedString(args[0]),
		utils.StripQuotedString(args[1])); err != nil {
		fmt.Println(err)
	}

	return 0
}
