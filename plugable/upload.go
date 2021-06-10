package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type UploadCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*UploadCommand)(nil)

func (UploadCommand) Command() string {
	return "upload"
}

func (cmd *UploadCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: upload <local-file-name>")
	} else {
		if err := cmd.DriveAPI.Upload(utils.StripQuotedString(args[0])); err != nil {
			fmt.Println(err)
		}
	}
	return 0
}
