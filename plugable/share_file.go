package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type ShareFileCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ShareFileCommand)(nil)

func (ShareFileCommand) Command() string {
	return "share-file"
}

func (cmd *ShareFileCommand) Execute(args ...string) int {
	if len(args) != 3 {
		fmt.Println("usage: share-file <file-id> <email> <role>")

		fmt.Println("Roles:")
		for _, r := range models.AllRoles {
			fmt.Printf("  %s\n", string(r))
		}

		return 0
	}

	if err := cmd.DriveAPI.ShareFileWith(
		utils.StripQuotedString(args[0]),
		utils.StripQuotedString(args[1]),
		models.Role(utils.StripQuotedString(args[2]))); err != nil {
		fmt.Println(err)
	}

	return 0
}
