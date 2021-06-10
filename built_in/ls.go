package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/interrupt_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
	"strconv"
	"strings"
)

type LsCommand struct {
	DriveAPI     intf.DriveAPI                      `inject:"*"`
	InterruptSvc interrupt_service.InterruptService `inject:"*"`
	terminate    bool
}

var _ cmdline_service.CommandInterface = (*LsCommand)(nil)

func (LsCommand) Command() string {
	return "ls"
}

func (cmd *LsCommand) Execute(args ...string) int {
	allDrive := false
	folderId := cmd.DriveAPI.CurrentFolder()
	numberOfFile := -1

	for len(args) > 0 {
		if args[0] == `-a` {
			allDrive = true
			args = args[1:]
			continue
		}
		if strings.HasPrefix(args[0], `-`) {
			if counter, err := strconv.Atoi(args[0][1:]); err != nil {
				fmt.Println("usage: ls [-a] [-number] <folderId>")
				return 0
			} else {
				numberOfFile = counter
				args = args[1:]
				continue
			}
		}

		folderId = utils.StripQuotedString(args[0])
		break
	}

	counter := 0

	fmt.Printf("Reading folder: %s\n", folderId)

	interrupted, err := cmd.InterruptSvc.Exec(func(errChan chan error) {
		errChan <- cmd.DriveAPI.ReadFolderDetail(folderId, allDrive, func(file *drive.File) bool {
			if file.MimeType == google_drive.FolderMimeType {
				fmt.Printf("  [ %s ] <dir> %s\n", file.Id, file.Name)
			} else {
				fmt.Printf("  [ %s ] ----- %s   [ %d ]\n", file.Id, file.Name, file.Size)
			}
			counter++
			return !cmd.terminate && (numberOfFile < 0 || counter < numberOfFile)
		})
	})

	if interrupted {
		cmd.terminate = true
	} else if err != nil {
		fmt.Println(err)
	}
	return 0
}
