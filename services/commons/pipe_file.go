package commons

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
)

func PipeDownloadFile(cryptoSvc intf.CryptoService, driveApi intf.DriveAPI, fileId string, writer io.Writer) (*drive.File, string, string, error) {
	var target *drive.File

	fileName := ``
	md5sum := ``
	if cryptoSvc.IsEnabled() {
		err := driveApi.RetrieveFile(fileId, func(file *drive.File, reader io.Reader) error {
			var err error

			target = file
			fileName, md5sum, err = cryptoSvc.DecryptFileNameWithMd5(file)

			if err != nil {
				return err
			}

			return cryptoSvc.DecodeViaPipe(reader, writer)
		})
		return target, fileName, md5sum, err
	} else {
		err := driveApi.RetrieveFile(fileId, func(file *drive.File, reader io.Reader) error {
			target = file
			fileName = file.Name
			md5sum = file.Md5Checksum
			written, err := io.Copy(writer, reader)
			if err == nil {
				fmt.Fprintf(os.Stderr, "file: %s, %d bytes\n\n", file.Name, written)
			}
			return err
		})
		return target, fileName, md5sum, err
	}
}
