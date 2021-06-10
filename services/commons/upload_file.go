package commons

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
	"time"
)

func UploadProgress(fileName string, fileSize int64, done chan error, progressReader utils.ProgressReader) error {
	started := time.Now()
	last := time.Now()
	var lastRead int64

	for i := 0; ; i += 1 {
		select {
		case err := <-done:
			uploadSize := progressReader.N()
			bps := float64(uploadSize) / float64(time.Now().Sub(started)/time.Second) / 1024.0
			fmt.Printf("%s uploaded (%d bytes, %.2f Kbytes/sec)\n",
				fileName, uploadSize, bps)
			return err

		default:
			time.Sleep(time.Millisecond * 100)

			if i%100 == 99 {
				now := time.Now()
				elapsedTotal := now.Sub(started) / time.Second
				elapsedLast := now.Sub(last) / time.Second

				currentRead := progressReader.N()
				periodRead := currentRead - lastRead

				avgRead := float64(currentRead) / float64(elapsedTotal) / 1024.0
				currRead := float64(periodRead) / float64(elapsedLast) / 1024.0

				fmt.Printf("%s: upload %d (current: %.2f Kbytes/sec, avg: %.2f Kbytes/sec, %.2f%%)\n",
					fileName, currentRead, currRead, avgRead,
					float64(currentRead)*100/float64(fileSize))

				lastRead = currentRead
				last = now
			}
		}
	}
}

func UploadFile(api google_drive.API, driveSvc *drive.Service, cryptoSvc intf.CryptoService, fileFullPath string, folderId string, fileName string, md5sum string) error {
	progressReader := utils.NewProgressReader(nil)
	done := make(chan error)

	fileState, err := os.Stat(fileFullPath)

	if err != nil {
		return err
	}
	fmt.Printf("(%d) ... ", fileState.Size())

	if cryptoSvc.IsEnabled() {
		encryptFileName := cryptoSvc.EncryptFileNameWithMd5(fileName, md5sum)
		logger.Debug("Encrypted Name: %s", encryptFileName)

		progressWriter := utils.NewProgressWriter(nil)

		go func() {
			done <- cryptoSvc.EncodeAndUploadReader(fileFullPath, func(writer io.Writer) io.Writer {
				return progressWriter.SetWriter(writer)
			}, func(reader io.Reader) error {
				progressReader.SetReader(reader)

				file, err := api.PushFile(driveSvc, folderId, encryptFileName, progressReader)

				if err == nil {
					logger.Debug("file pushed: %s", file.Id)
				}
				return err
			})
		}()

		return UploadProgress(fileName, fileState.Size(), done, progressReader)

		//return cryptoSvc.EncodeAndUploadReader(fileFullPath, func(reader io.Reader) error {
		//	file, err := api.PushFile(driveSvc, folderId, encryptFileName, reader)
		//
		//	if err == nil {
		//		logger.Debug("file pushed: %s", file.Id)
		//	}
		//	return err
		//})
		//}); err != nil {
		//	logger.Error(err)
		//	//} else {
		//	//	fmt.Printf("%d bytes (encrypted)\n", node.fileInfo.Size())
		//	//	traversal.uploadCount++
		//	//	traversal.uploadSize += node.fileInfo.Size()
		//	//	traversal.handler.MarkUploaded(node.data)
		//}
	} else {
		go func() {
			_, err := api.UploadFile(driveSvc, fileFullPath, folderId, fileName, func(reader io.Reader) io.Reader {
				return progressReader.SetReader(reader)
			})
			done <- err
		}()

		return UploadProgress(fileName, fileState.Size(), done, progressReader)

		//} else {
		//	fmt.Printf("%d bytes\n", node.fileInfo.Size())
		//	traversal.uploadCount++
		//	traversal.uploadSize += node.fileInfo.Size()
		//	traversal.handler.MarkUploaded(node.data)
	}
}
