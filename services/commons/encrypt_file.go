package commons

//func EncryptFile(workingDir string,
//	cryptoService crypto_service.CryptoService,
//	fullPath string,
//	md5sum string,
//	callback func(encryptedFile string, encryptedFileName string) error) error {
//
//	fileName := filepath.Base(fullPath)
//
//	if cryptoService.IsEnabled() {
//		encryptFileName := cryptoService.EncryptFileNameWithMd5(fileName, md5sum)
//
//		if reader, err := cryptoService.EncodeAndUploadReader(fullPath); err != nil {
//			return err
//		} else {
//
//		}
//		tempFile, err := ioutil.TempFile(workingDir, `upload`)
//		if err != nil {
//			return err
//		}
//
//		defer os.Remove(tempFile.Name()) // clean up
//
//		cryptoService.EncryptFile(fullPath, tempFile.Name())
//
//		return callback(tempFile.Name(), encryptFileName)
//	} else {
//		return callback(fullPath, fileName)
//	}
//}
