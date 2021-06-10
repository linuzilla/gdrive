package drive_api

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/progress_service"
)

type driveApiImpl struct {
	DriveApi        google_drive.API                        `inject:"*"`
	CryptoFactory   crypto_service.CryptoFactory            `inject:"*"`
	ProgressFactory progress_service.ProgressServiceFactory `inject:"*"`
	folderId        string
	folderStack     stack
	filesCache      map[string]string
}

type stack []string

func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *stack) Push(v string) stack {
	*s = append(*s, v)
	return *s
}

func (s *stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		i := len(*s) - 1
		x := (*s)[i]
		*s = (*s)[:i]

		return x, true
	}
}

var _ intf.DriveAPI = (*driveApiImpl)(nil)

func New() intf.DriveAPI {
	return &driveApiImpl{filesCache: make(map[string]string)}
}

func (api *driveApiImpl) PostSummerConstruct() {
	rootFolder := api.DriveApi.RootFolder()
	if rootFolder != `` {
		api.folderId = rootFolder
	}
}

func (api *driveApiImpl) CurrentFolder() string {
	return api.folderId
}

func (api *driveApiImpl) ChangeFolder(folderNameOrId string) {
	folderId := api.cachedName(folderNameOrId)
	fmt.Println("Change folder to: " + folderId)
	api.folderId = folderId
}

func (api *driveApiImpl) Pushd() {
	fmt.Println("Push current folder to stack: " + api.folderId)
	api.folderStack.Push(api.folderId)
}

func (api *driveApiImpl) Popd() {
	if folderId, found := api.folderStack.Pop(); found {
		fmt.Println("Pop previous folder: " + folderId)
		api.folderId = folderId
	}
}

func (api *driveApiImpl) Dirs() {
	fmt.Println("Folder Stack:")

	for i := len(api.folderStack) - 1; i >= 0; i-- {
		fmt.Println("    " + api.folderStack[i])
	}
}

func (api *driveApiImpl) cachedName(fileNameOrId string) string {
	if cachedId, found := api.filesCache[fileNameOrId]; found {
		return cachedId
	} else {
		return fileNameOrId
	}
}

func (api *driveApiImpl) Version() {
	fmt.Println(constants.VERSION)
}
