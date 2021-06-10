package dot_gdrive_service

import (
	"github.com/linuzilla/gdrive/intf"
	"google.golang.org/api/drive/v3"
)

type DotGDriveService interface {
	LoadEnvironment() *DotEnvironment
	LoadTopLevelEnvironment(targetDir string) *DotEnvironment
	LoadTopLevelEnvByFolder(folderId string) *DotEnvironment
	PrepareBeforeSync(folderName string, folderId string,
		callback func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId string))
}

type DotEnvironment struct {
	Dir          string
	FolderId     string
	Password     string
	DatabaseFile string
	Encoder      string
	Decoder      string
	DirectArgs   string
}
