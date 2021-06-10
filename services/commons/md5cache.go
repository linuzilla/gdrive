package commons

import (
	"database/sql"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/utils"
	"os"
	"path"
)

func Md5WithCache(connection intf.DatabaseBackendConnection, baseDir string, relativeDir string, fileInfo os.FileInfo,
	data *models.SyncFileInfo) (string, error) {
	fullPath := path.Join(baseDir, fileInfo.Name())
	relativePath := path.Join(relativeDir, fileInfo.Name())

	if err := connection.FindFirstById(relativePath, data); err != nil {
		if md5sum, err := utils.Md5sum(fullPath); err != nil {
			return "", err
		} else { // create a new entry
			data = &models.SyncFileInfo{
				Id:         relativePath,
				RemoteName: sql.NullString{Valid: false}, //*gds.encryptFileName(fileInfo.Name(), md5sum),
				CheckSum:   md5sum,
				ModTime:    fileInfo.ModTime(),
				FileSize:   fileInfo.Size(),
				Uploaded:   false,
			}
			// return md5sum, gds.xdb.Gorm.Create(data).Error
			return md5sum, connection.Persist(data)
		}
	} else {
		if data.FileSize != fileInfo.Size() || data.ModTime.Unix() != fileInfo.ModTime().Unix() {
			if md5sum, err := utils.Md5sum(fullPath); err != nil {
				return "", err
			} else {
				data.FileSize = fileInfo.Size()
				data.ModTime = fileInfo.ModTime()
				data.CheckSum = md5sum
				data.Uploaded = false
				// gds.xdb.Gorm.Save(data)
				connection.SaveOrUpdate(data)
				return md5sum, nil
			}
		}
		return data.CheckSum, nil
	}
}
