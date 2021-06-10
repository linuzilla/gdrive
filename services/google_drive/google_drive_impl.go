package google_drive

import (
	"context"
	"fmt"
	"github.com/linuzilla/gdrive/config"
	"github.com/linuzilla/gdrive/constants"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type googleDriveImpl struct {
	Conf              *config.Config `inject:"*"`
	credentialFile    string
	impersonate       string
	domainAdminAccess bool
}

const FolderMimeType = `application/vnd.google-apps.folder`

func (googleDrive *googleDriveImpl) PostSummerConstruct() {
	googleDrive.credentialFile = googleDrive.Conf.GoogleDrive.Credential

	if runtime.GOOS != `windows` {
		if !strings.HasPrefix(googleDrive.Conf.GoogleDrive.Credential, `/`) {
			if currentDir, err := os.Getwd(); err == nil {
				googleDrive.credentialFile = currentDir + `/` + googleDrive.Conf.GoogleDrive.Credential
			}
		}
	}

	if googleDrive.Conf.GSuite.UseDomainAdminAccess {
		googleDrive.domainAdminAccess = true

		if googleDrive.Conf.GSuite.Impersonate != `` {
			googleDrive.impersonate = googleDrive.Conf.GSuite.Impersonate
		}
	}
}

func (googleDrive *googleDriveImpl) RootFolder() string {
	return googleDrive.Conf.GoogleDrive.FolderId
}

func (googleDrive *googleDriveImpl) SetImpersonate(impersonate string) {
	googleDrive.impersonate = impersonate
}

func (googleDrive *googleDriveImpl) GetImpersonate() string {
	return googleDrive.impersonate
}

func (googleDrive *googleDriveImpl) SetDomainAdminAccess(useDomainAdminAccess bool) {
	googleDrive.domainAdminAccess = useDomainAdminAccess
}

func (googleDrive *googleDriveImpl) GetDomainAdminAccess() bool {
	return googleDrive.domainAdminAccess
}

func (googleDrive *googleDriveImpl) ReadTeamDrives(service *drive.Service, callback func(teamDrive *drive.TeamDrive) bool) error {
	pageToken := ""

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		fields := service.Teamdrives.List().
			UseDomainAdminAccess(googleDrive.domainAdminAccess).
			Fields("nextPageToken, teamDrives(id, name)")

		r, err := fields.PageToken(pageToken).Do()

		if err != nil {
			return err
		}

		if len(r.TeamDrives) > 0 {
			for _, teamDrive := range r.TeamDrives {
				if !callback(teamDrive) {
					return nil
				}
			}
		}
		pageToken = r.NextPageToken
	}

	return nil
}

func (googleDrive *googleDriveImpl) PermissionInfo(service *drive.Service, driveId string, callback func(permission *drive.Permission)) error {
	pageToken := ""

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		r, err := service.Permissions.List(driveId).
			UseDomainAdminAccess(googleDrive.domainAdminAccess).
			SupportsAllDrives(true).
			Fields("nextPageToken, permissions(kind,id,type,emailAddress,role,deleted)").
			Do()

		if err != nil {
			return err
		}

		for _, permission := range r.Permissions {
			callback(permission)
		}

		pageToken = r.NextPageToken

		break
	}
	return nil
}

func (googleDrive *googleDriveImpl) MoveFile(service *drive.Service, source *drive.File, newParentId string) (*drive.File, error) {
	parentFolderId := ``

	if len(source.Parents) > 0 {
		parentFolderId = source.Parents[0]
	}

	return service.Files.Update(source.Id, &drive.File{
		Name: source.Name,
	}).
		SupportsAllDrives(true).
		AddParents(newParentId).
		RemoveParents(parentFolderId).
		Fields("id, parents").Do()
}

func (googleDrive *googleDriveImpl) Rename(service *drive.Service, fileId string, newFileName string) (*drive.File, error) {
	file, err := service.Files.Get(fileId).
		SupportsAllDrives(true).
		Fields("id, name").
		Do()

	if err != nil {
		return nil, err
	}

	return service.Files.Update(file.Id, &drive.File{Name: newFileName}).
		SupportsAllDrives(true).
		Fields("id, name, parents").Do()
}

func (googleDrive *googleDriveImpl) Duplicate(service *drive.Service, sourceFileId string, targetFileName string) (*drive.File, error) {
	sourceFile, err := service.Files.Get(sourceFileId).
		SupportsAllDrives(true).
		Fields("id, parents").
		Do()

	if err != nil {
		return nil, err
	}

	parentId := ``

	if sourceFile.Parents != nil && len(sourceFile.Parents) > 0 {
		parentId = sourceFile.Parents[0]
		newFile := &drive.File{
			Name:    filepath.Base(targetFileName),
			Parents: []string{parentId},
		}
		return service.Files.Copy(sourceFile.Id, newFile).
			SupportsAllDrives(true).
			Fields("id, name, parents").Do()
	} else {
		return nil, fmt.Errorf(`file %s has no parent`, sourceFile.Id)
	}
}

func (googleDrive *googleDriveImpl) FindWithQuery(service *drive.Service, query string, callback func(file *drive.File) bool) error {
	pageToken := ""

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		fields := service.Files.List().
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, files(id, name, parents, md5Checksum, teamDriveId, size, mimeType, modifiedTime)")

		if len(query) > 0 {
			//fmt.Println("Query: " + query)
			fields = fields.Q(query)
		}

		r, err := fields.PageToken(pageToken).Do()
		if err != nil {
			return err
		}

		if len(r.Files) > 0 {
			for _, f := range r.Files {
				if !callback(f) {
					return nil
				}
			}
		}
		pageToken = r.NextPageToken
	}
	return nil
}

func (googleDrive *googleDriveImpl) FindDirectoriesInFolder(service *drive.Service, folderId string, filePath string, callback func(file *drive.File, path string) bool) error {
	return googleDrive.FindWithQuery(service, fmt.Sprintf(
		`mimeType = '%s' and '%s' in parents and not trashed`,
		FolderMimeType, folderId,
	), func(file *drive.File) bool {
		return callback(file, filePath)
	})
}

func (googleDrive *googleDriveImpl) RecursiveFindDirectoriesInFolder(service *drive.Service, folderId string, filePath string, callback func(file *drive.File, path string)) error {
	var list []*drive.File

	if err := googleDrive.FindDirectoriesInFolder(service, folderId, filePath, func(file *drive.File, path string) bool {
		list = append(list, file)
		return true
	}); err != nil {
		return err
	}

	for _, f := range list {
		callback(f, filePath)
		if err := googleDrive.RecursiveFindDirectoriesInFolder(service, f.Id, filePath+"/"+f.Name, callback); err != nil {
			return err
		}
	}

	return nil
}

func (googleDrive *googleDriveImpl) FindFilesInFolder(service *drive.Service, folderId string, callback func(file *drive.File) bool) error {
	return googleDrive.FindWithQuery(service, fmt.Sprintf(
		`mimeType != '%s' and '%s' in parents and not trashed`,
		FolderMimeType, folderId,
	), callback)
}

func (googleDrive *googleDriveImpl) FindFilesByNameInFolder(service *drive.Service, folderId string, fileName string, callback func(file *drive.File) bool) error {
	return googleDrive.FindWithQuery(service, fmt.Sprintf(
		`name = '%s' and '%s' in parents and not trashed`,
		fileName, folderId,
	), callback)
}

func (googleDrive *googleDriveImpl) FileInfo(service *drive.Service, fileId string) (*drive.File, error) {
	return service.Files.Get(fileId).
		SupportsAllDrives(true).
		Fields("id, name, parents, md5Checksum, teamDriveId, size, mimeType, modifiedTime").
		Do()
}

func (googleDrive *googleDriveImpl) ShareFileWith(service *drive.Service, fileId string, email string, role string) (*drive.Permission, error) {
	newPerm := &drive.Permission{
		Type:         `user`,
		Role:         role,
		EmailAddress: email,
	}

	return service.Permissions.Create(fileId, newPerm).
		UseDomainAdminAccess(googleDrive.domainAdminAccess).
		SupportsAllDrives(true).
		Fields(`id`).
		Do()
}

func (googleDrive *googleDriveImpl) RemovePermission(service *drive.Service, fileId string, permissionId string) error {
	return service.Permissions.Delete(fileId, permissionId).
		UseDomainAdminAccess(googleDrive.domainAdminAccess).
		SupportsAllDrives(true).
		Do()
}

func (googleDrive *googleDriveImpl) GetPermissionIdByEmail(service *drive.Service, fileId string, email string) (string, error) {
	file, err := service.Files.Get(fileId).
		SupportsAllDrives(true).
		Fields("id, permissions, permissionIds").
		Do()

	if err != nil {
		return ``, err
	}

	if file.Permissions != nil && len(file.Permissions) > 0 {
		for _, permission := range file.Permissions {
			if permission.EmailAddress == email {
				return permission.Id, nil
			}
		}
	}

	if file.PermissionIds != nil && len(file.PermissionIds) > 0 {
		for _, permissionId := range file.PermissionIds {
			if permission, err := service.Permissions.Get(fileId, permissionId).
				UseDomainAdminAccess(googleDrive.domainAdminAccess).
				SupportsAllDrives(true).
				UseDomainAdminAccess(googleDrive.domainAdminAccess).
				Fields(`id, emailAddress, role, type, teamDrivePermissionDetails`).
				Do(); err != nil {
				fmt.Println(err)
			} else {
				if permission.EmailAddress == email {
					return permission.Id, nil
				}
			}
		}
	}

	return "", fmt.Errorf("permission not found for %s", email)
}

func (googleDrive *googleDriveImpl) ShareInfo(service *drive.Service, fileId string) (*drive.File, error) {
	file, err := service.Files.Get(fileId).
		SupportsAllDrives(true).
		Fields("id, name, parents, md5Checksum, teamDriveId, size, mimeType, modifiedTime, permissions, permissionIds, shared, webViewLink").
		Do()

	if err != nil {
		return nil, err
	}

	fmt.Printf("WebViewLink: %s\n", file.WebViewLink)
	fmt.Printf("Share: %s\n", constants.OnOffMap[file.Shared])

	if file.Permissions != nil && len(file.Permissions) > 0 {
		for _, permission := range file.Permissions {
			fmt.Printf("  %s (%s) *\n", permission.EmailAddress, permission.Role)
		}
		fmt.Println()
	}

	if file.PermissionIds != nil && len(file.PermissionIds) > 0 {
		for _, permissionId := range file.PermissionIds {
			if permission, err := service.Permissions.Get(fileId, permissionId).
				UseDomainAdminAccess(googleDrive.domainAdminAccess).
				SupportsAllDrives(true).
				UseDomainAdminAccess(googleDrive.domainAdminAccess).
				Fields(`emailAddress, role, type, teamDrivePermissionDetails`).
				Do(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("  %s (%s) [ %s ]\n", permission.EmailAddress, permission.Role, permission.Type)
			}
		}
	}
	return file, err
}

func (googleDrive *googleDriveImpl) ReadFolder(service *drive.Service, folderId string, callback func(file *drive.File)) error {
	pageToken := ""

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		fields := service.Files.List().
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, files(id, name, parents, md5Checksum, teamDriveId, size, mimeType)")

		if len(folderId) > 0 {
			fields = fields.Q("'" + folderId + "' in parents and not trashed")
		}

		r, err := fields.PageToken(pageToken).Do()
		if err != nil {
			return err
		}
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				callback(i)
			}
		}
		pageToken = r.NextPageToken
	}
	return nil
}

func (googleDrive *googleDriveImpl) ReadFolderDetail(service *drive.Service, folderId string, allDrive bool, callback func(file *drive.File) bool) error {
	pageToken := ""

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		fields := service.Files.List().
			IncludeItemsFromAllDrives(allDrive).
			SupportsAllDrives(allDrive).
			Fields("nextPageToken, files(id, name, parents, md5Checksum, teamDriveId, size, mimeType)")

		if len(folderId) > 0 {
			fields = fields.Q("'" + folderId + "' in parents and not trashed")
		}

		r, err := fields.PageToken(pageToken).Do()
		if err != nil {
			return err
		}
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				if !callback(i) {
					return nil
				}
			}
		}
		pageToken = r.NextPageToken
	}
	return nil
}

func (googleDrive *googleDriveImpl) Read(service *drive.Service, folderId string, fileNameDecoder func(file *drive.File) string) error {
	pageToken := ""

	fmt.Println("Files:")

	for haveNextPage := true; haveNextPage; haveNextPage = pageToken != "" {
		fields := service.Files.List().
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, files(id, name, parents, md5Checksum, teamDriveId, size, mimeType)")

		if len(folderId) > 0 {
			fmt.Println("folder: " + folderId)
			fields = fields.Q("'" + folderId + "' in parents and not trashed")
		}

		r, err := fields.PageToken(pageToken).Do()
		if err != nil {
			return err
		}
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				fileName := i.Name

				if fileNameDecoder != nil {
					fileName = fileNameDecoder(i)
				}

				if i.MimeType == FolderMimeType {
					fmt.Printf("  [ %s ] <dir> %s\n", i.Id, fileName)
				} else {
					fmt.Printf("  [ %s ] ----- %s\n", i.Id, fileName)
				}
			}
		}
		pageToken = r.NextPageToken
		//fmt.Println("PageToken: " + pageToken)
	}
	return nil
}

func (googleDrive *googleDriveImpl) UploadFile(service *drive.Service, localFile string, folderId string, fileName string, readerWrapper func(reader io.Reader) io.Reader) (*drive.File, error) {
	newFile := &drive.File{
		Name:    filepath.Base(fileName),
		Parents: []string{folderId},
	}

	if reader, err := os.Open(localFile); err != nil {
		return nil, err
	} else {
		var r io.Reader = reader

		if readerWrapper != nil {
			r = readerWrapper(reader)
		}

		file, err := service.Files.Create(newFile).
			SupportsAllDrives(true).
			Fields("id, parents").
			Media(r).
			Do()

		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func (googleDrive *googleDriveImpl) ReplaceFile(service *drive.Service, localFile string, source *drive.File, fileName string) error {
	if reader, err := os.Open(localFile); err != nil {
		return err
	} else {
		file, err := service.Files.Update(source.Id, &drive.File{
			Name:    fileName,
			Parents: source.Parents,
		}).
			SupportsAllDrives(true).
			Fields("id, parents").
			Media(reader).
			Do()

		if err != nil {
			return err
		} else {
			fmt.Println(localFile + ": file replaced (id=" + file.Id + ")")
		}
	}
	return nil
}

func (googleDrive *googleDriveImpl) RetrieveFile(service *drive.Service, fileId string, fileRetriever func(*drive.File, io.Reader) error) error {
	file, err := service.Files.Get(fileId).
		SupportsAllDrives(true).
		Fields("id, name, parents, md5Checksum, size, teamDriveId").
		Do()

	if err != nil {
		// log.Printf("No such file (fileId: %s): %v", fileId, err)
		return err
	} else {
		response, err := service.Files.Get(fileId).
			SupportsAllDrives(true).
			Download()

		if err != nil {
			return err
		} else {
			return fileRetriever(file, response.Body)
		}
	}
}

func (googleDrive *googleDriveImpl) PushFile(service *drive.Service, folderId string, fileName string, reader io.Reader) (*drive.File, error) {
	newFile := &drive.File{
		Name:    filepath.Base(fileName),
		Parents: []string{folderId},
	}

	return service.Files.Create(newFile).
		SupportsAllDrives(true).
		Fields("id, parents").
		Media(reader).
		Do()
}

func (googleDrive *googleDriveImpl) DownloadFile(service *drive.Service, fileId string, localFile string, overwrite bool) (string, error) {
	file, err := service.Files.Get(fileId).
		SupportsAllDrives(true).
		Do()

	if err != nil {
		// log.Printf("No such file (fileId: %s): %v", fileId, err)
		return ``, err
	} else {
		fileName := file.Name

		if len(localFile) == 0 || localFile == "." {
			localFile = fileName
		}

		response, err := service.Files.Get(fileId).
			SupportsAllDrives(true).
			Download()

		flag := os.O_WRONLY | os.O_CREATE | os.O_EXCL

		if overwrite {
			flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		}

		if err != nil {
			return file.Name, err
		} else {
			file, err := os.OpenFile(localFile, flag, 0600)
			defer file.Close()

			if err != nil {
				log.Printf("Unable to write file: %v", err)
			} else {
				written, err := io.Copy(file, response.Body)
				if err != nil {
					log.Printf("Unable to write file: %v", err)
				} else {
					fmt.Printf("written %d byte(s)", written)
				}
			}
		}

		return file.Name, nil
	}
}

func (googleDrive *googleDriveImpl) CreateFolder(service *drive.Service, parentFolder string, folderName string) (string, error) {
	newFolder := &drive.File{
		Name:     folderName,
		Parents:  []string{parentFolder},
		MimeType: FolderMimeType,
	}

	file, err := service.Files.Create(newFolder).
		SupportsAllDrives(true).
		Do()

	if err != nil {
		// log.Printf("Unable to create folder: %v", err)
		return "", err
	} else {
		return file.Id, nil
	}
}

func (googleDrive *googleDriveImpl) Delete(service *drive.Service, fileId string) error {
	return service.Files.Delete(fileId).SupportsTeamDrives(true).Do()
}

func (googleDrive *googleDriveImpl) FindOrCreateTrashFolder(service *drive.Service, fileId string) (string, error) {
	if file, err := googleDrive.FileInfo(service, fileId); err != nil {
		fmt.Println(err)
		return ``, err
	} else {
		trashParent := file.TeamDriveId

		if file.TeamDriveId == `` {
			if len(file.Parents) == 0 {
				return ``, fmt.Errorf("ERROR!! not parents or team drive")
			} else {
				trashParent = file.Parents[0]
			}
		}

		trashFolderId := ``
		googleDrive.FindFilesByNameInFolder(service, trashParent, constants.TrashFolderName, func(file *drive.File) bool {
			if file.MimeType == FolderMimeType {
				trashFolderId = file.Id
				fmt.Println("Trash folder found, ID: " + trashFolderId)
				return false
			} else {
				return true
			}
		})

		if trashFolderId == `` {
			if folder, err := googleDrive.CreateFolder(service, trashParent, constants.TrashFolderName); err != nil {
				return ``, err
			} else {
				fmt.Println("Create a new trash folder, folder-id=" + folder)
				return folder, nil
			}
		} else {
			return trashFolderId, nil
		}
	}
}

func (googleDrive *googleDriveImpl) GetCredentialEmail() (string, error) {
	if config, err := googleDrive.loadConfig(); err != nil {
		return ``, err
	} else {
		return config.Email, nil
	}
}

func (googleDrive *googleDriveImpl) loadConfig() (*jwt.Config, error) {
	if b, err := ioutil.ReadFile(googleDrive.credentialFile); err != nil {
		log.Fatalf("Unable to read client secret file: %v\n", err)
		return nil, err
	} else {
		// If modifying these scopes, delete your previously saved token.json.
		return google.JWTConfigFromJSON(b,
			drive.DriveMetadataReadonlyScope,
			drive.DriveScope)
	}
}

func (googleDrive *googleDriveImpl) Connect(callback func(service *drive.Service) error) error {
	config, err := googleDrive.loadConfig()
	if err != nil {
		return err
	}

	if googleDrive.impersonate != `` {
		config.Subject = googleDrive.impersonate
	}

	ctx := context.Background()

	service, err := drive.NewService(ctx,
		option.WithHTTPClient(config.Client(ctx)),
		option.WithUserAgent(constants.UserAgent))

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return err
	}

	return callback(service)
}

// var _ GoogleDriveAPI = (*googleDriveImpl)(nil)

func New() API {
	return &googleDriveImpl{}
}
