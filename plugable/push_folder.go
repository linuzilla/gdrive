package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/sync_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"os"
	"strings"
)

type PushFolderCommand struct {
	SyncService sync_service.SyncService               `inject:"*"`
	EnvSvc      environment_service.EnvironmentService `inject:"*"`
	DotService  dot_gdrive_service.DotGDriveService    `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*PushFolderCommand)(nil)

func (PushFolderCommand) Command() string {
	return `push-folder`
}

func (cmd *PushFolderCommand) pushFolder(folderName string, uploadFolderId string, speedUpload bool, deleteFlag bool) {
	cmd.SyncService.PushUpload(folderName, uploadFolderId, speedUpload, deleteFlag)
}

func (cmd *PushFolderCommand) help() {
	on := constants.OnOffMap[cmd.EnvSvc.Encoding()]

	fmt.Printf("usage: push-folder [ -f | -h | -d ] <local-folder-name> <folder-id>    [Encryption: %s]\n", on)
}

func (cmd *PushFolderCommand) Execute(args ...string) int {
	uploadFolderId := ``
	speedUpload := true
	deleteFlag := false

	for len(args) > 0 {
		if args[0] == `-f` {
			args = args[1:]
			speedUpload = false
		} else if args[0] == `-d` {
			args = args[1:]
			deleteFlag = true
		} else if args[0] == `-h` {
			cmd.help()
			return 0
		} else {
			break
		}
	}

	switch len(args) {
	default:
		if dotEnv := cmd.DotService.LoadEnvironment(); dotEnv != nil {
			cmd.pushFolder(dotEnv.Dir, dotEnv.FolderId, speedUpload, deleteFlag)
		} else {
			cmd.help()
		}
	case 2:
		uploadFolderId = args[1]
		fallthrough
	case 1:
		folderName := utils.StripQuotedString(args[0])

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return 0
		}

		if folderName == `.` {
			folderName = currentDir
		} else if !strings.HasPrefix(folderName, `/`) {
			folderName = currentDir + `/` + folderName
		}

		if strings.HasSuffix(folderName, `/`) {
			folderName = folderName[:len(folderName)-1]
		}

		cmd.pushFolder(folderName, uploadFolderId, speedUpload, deleteFlag)

	}
	return 0
}
