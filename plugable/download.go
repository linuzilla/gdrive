package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"os"
)

type DownloadCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DownloadCommand)(nil)

func (DownloadCommand) Command() string {
	return "download"
}

func (cmd *DownloadCommand) Execute(args ...string) int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	switch fileName := "."; len(args) {
	case 2:
		fileName = utils.StripQuotedString(args[1])
		fallthrough
	case 1:
		if err := cmd.DriveAPI.DownloadFile(utils.StripQuotedString(args[0]), currentDir, fileName); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("usage: download <file-id> [local-file-name]")
	}

	return 0
}
