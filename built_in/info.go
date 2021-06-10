package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
)

type InfoCommand struct {
	DriveAPI      intf.DriveAPI                `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*InfoCommand)(nil)

func (InfoCommand) Command() string {
	return "info"
}

func (cmd *InfoCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: info <fileId>")
	} else {
		fileName, file, err := cmd.DriveAPI.FileInfo(args[0], func(file *drive.File) string {
			return commons.FileNameDecoder(file, cmd.CryptoFactory.GetInstance())
		})

		if err != nil {
			fmt.Println(err)
		} else {
			if file.MimeType == google_drive.FolderMimeType {
				fmt.Printf("Folder: %s\n", fileName)
				fmt.Printf("Folder ID: %s\n", file.Id)
			} else {
				fmt.Printf("File: %s\n", fileName)
				fmt.Printf("File ID: %s\n", file.Id)
				fmt.Printf("File Size: %d\n", file.Size)
				fmt.Printf("Check Sum: %s\n", file.Md5Checksum)
				fmt.Printf("Mime Type: %s\n", file.MimeType)
			}

			parent := ``
			if len(file.Parents) > 0 {
				parent = file.Parents[0]
			}
			fmt.Printf("Parent ID: %s\n", parent)
			fmt.Printf("Team Drive ID: %s\n", file.TeamDriveId)

			if localTime, err := google_drive.ToLocalTime(file.ModifiedTime); err == nil {
				fmt.Printf("Modify Time: %s\n", localTime.String())
			}
		}
	}
	return 0
}
