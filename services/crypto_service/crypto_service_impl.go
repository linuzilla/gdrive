package crypto_service

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/execution_service"
	"github.com/linuzilla/go-logger"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type cryptoServiceImpl struct {
	envSvc          environment_service.EnvironmentService
	enabled         bool
	opensslPassword string
	passwordFile    string
	config          *models.Codec
	//encoder         execution_service.CommandExecutor
	//decoder         execution_service.CommandExecutor
	fileNameEncoder execution_service.CommandExecutor
	fileNameDecoder execution_service.CommandExecutor
}

var encryptedFileWithChecksumPattern = regexp.MustCompile(`^([a-f0-9]+)-(.*)$`)

func (service *cryptoServiceImpl) setExecutionService(data interface{}) {
	if executionService, ok := data.(execution_service.CommandExecutionService); ok {
		//service.executionService = executionService

		encoder := constants.EncoderDefaultCommand
		decoder := constants.DecoderDefaultCommand

		if service.config.Encoder != `` {
			encoder = service.config.Encoder
		}

		if service.config.Decoder != `` {
			decoder = service.config.Decoder
		}

		//if service.config.DirectArgs != `` {
		//	arg = service.config.DirectArgs
		//}

		//service.encoder = executionService.Executor(encoder + ` ` + arg)
		//service.decoder = executionService.Executor(decoder + ` ` + arg)
		service.fileNameEncoder = executionService.Executor(encoder)
		service.fileNameDecoder = executionService.Executor(decoder)
	} else {
		panic(fmt.Errorf("SetExecutionService(%v(", data))
	}
}

func (service *cryptoServiceImpl) updateEnvironment() {
	service.envSvc.SetEncoding(service.enabled)
}

func (service *cryptoServiceImpl) PostSummerConstruct() {
	service.updateEnvironment()
}

func (service *cryptoServiceImpl) IsEnabled() bool {
	return service.enabled
}

func (service *cryptoServiceImpl) SetPassword(password string) {
	service.createPasswordFile(password)
}

func (service *cryptoServiceImpl) GetPassword() string {
	return service.opensslPassword
}

func (service *cryptoServiceImpl) createPasswordFile(password string) {
	service.Close()

	if password == `` {
		if service.enabled || len(service.opensslPassword) == 0 {
			service.enabled = false
			service.updateEnvironment()
			fmt.Println(`Encryption: off`)
			return
		}
	} else {
		service.opensslPassword = password
	}

	if service.passwordFile == `` {
		tempFile, err := ioutil.TempFile(service.envSvc.GetWorkingDirectory(), `key`)
		if err != nil {
			panic(err.Error())
		}

		service.passwordFile = tempFile.Name()

		_, err = tempFile.WriteString(service.opensslPassword)
		if err != nil {
			panic(err.Error())
		}

		tempFile.Close()

		fmt.Printf("Password temp file: %s (Encryption: on)\n", service.passwordFile)
		service.enabled = true
		service.updateEnvironment()
	}
}

//func (service *cryptoServiceImpl) EncryptFile(sourceFile string, targetFile string) {
//	if service.enabled {
//		service.encoder.Exec(map[string]string{
//			constants.ArgPasswordFile: service.passwordFile,
//			constants.ArgSourceFile:   sourceFile,
//			constants.ArgTargetFile:   targetFile,
//		}, nil)
//	}
//}
//
//func (service *cryptoServiceImpl) DecryptFile(sourceFile string, targetFile string) {
//	if service.enabled {
//		service.decoder.Exec(map[string]string{
//			constants.ArgPasswordFile: service.passwordFile,
//			constants.ArgSourceFile:   sourceFile,
//			constants.ArgTargetFile:   targetFile,
//		}, nil)
//	}
//}

func (service *cryptoServiceImpl) encryptFileName(fileName string) string {
	if service.enabled {
		if buffer, _ := service.fileNameEncoder.Exec(map[string]string{
			constants.ArgPasswordFile: service.passwordFile,
		}, strings.NewReader(fileName)); buffer != nil {
			return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(
				buffer.Bytes())
		}
	}
	return fileName
}

func (service *cryptoServiceImpl) decryptFileName(fileName string) (string, error) {
	if service.enabled {
		if data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(fileName); err != nil {
			return fileName, err
		} else {
			buffer, errExec := service.fileNameDecoder.Exec(map[string]string{
				constants.ArgPasswordFile: service.passwordFile,
			}, bytes.NewReader(data))

			if errExec != nil {
				return fileName, errExec
			} else if buffer != nil {
				return buffer.String(), nil
			} else {
				return fileName, fmt.Errorf("empty buffer")
			}
		}
	} else {
		return fileName, nil
	}
}

func (service *cryptoServiceImpl) DecodeViaPipe(reader io.Reader, writer io.Writer) error {
	return service.fileNameDecoder.Pipe(map[string]string{
		constants.ArgPasswordFile: service.passwordFile,
	}, reader, writer)
}

func (service *cryptoServiceImpl) encodeViaPipe(reader io.Reader, writer io.Writer) error {
	return service.fileNameEncoder.Pipe(map[string]string{
		constants.ArgPasswordFile: service.passwordFile,
	}, reader, writer)
}

//func (service *cryptoServiceImpl) JustDecryptFileName(fileName string) string {
//	if name, err := service.DecryptFileName(fileName); err != nil {
//		return fileName
//	} else {
//		return name
//	}
//}

func (service *cryptoServiceImpl) DecodeAndDownloadFromReader(fileName string, reader io.Reader) error {
	f, errOpenFile := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)

	if errOpenFile != nil {
		return errOpenFile
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	err := service.DecodeViaPipe(reader, writer)
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

func (service *cryptoServiceImpl) EncodeAndUploadReader(fileName string, wrapWriter func(writer io.Writer) io.Writer, callback func(reader io.Reader) error) error {
	fileReader, errOpenFile := os.Open(fileName)

	if errOpenFile != nil {
		logger.Error(errOpenFile)
		return errOpenFile
	}

	reader, writer := io.Pipe()

	progressWriter := wrapWriter(writer)

	go func() {
		defer writer.Close()

		err := service.encodeViaPipe(fileReader, progressWriter)
		if err != nil {
			logger.Error(err)
		}
	}()

	return callback(reader)
}

func (service *cryptoServiceImpl) EncryptFileNameWithMd5(fileName string, md5 string) string {
	if service.enabled {
		encryptedFileName := service.encryptFileName(md5 + `-` + fileName)
		return encryptedFileName
	} else {
		return ``
	}
}

func (service *cryptoServiceImpl) DecryptFileNameWithMd5(file *drive.File) (string, string, error) {
	if service.enabled {
		if decrypted, err := service.decryptFileName(file.Name); err != nil {
			return file.Name, ``, err
		} else {
			if matched := encryptedFileWithChecksumPattern.FindStringSubmatch(decrypted); matched != nil {
				return matched[2], matched[1], nil
			}
		}
		return file.Name, file.Md5Checksum, nil
	} else {
		return file.Name, file.Md5Checksum, nil
	}
}

func (service *cryptoServiceImpl) GetEncoder() string {
	return service.config.Encoder
}

func (service *cryptoServiceImpl) GetDecoder() string {
	return service.config.Decoder
}

func (service *cryptoServiceImpl) Close() {
	if service.passwordFile != `` {
		os.Remove(service.passwordFile)
		service.passwordFile = ``
	}
}

func newCryptoService(codec *models.Codec, envSvc environment_service.EnvironmentService, execSvc execution_service.CommandExecutionService) intf.CryptoService {
	service := &cryptoServiceImpl{
		enabled: false,
		config:  codec,
		envSvc:  envSvc,
	}

	service.setExecutionService(execSvc)

	//if service.enabled {
	//	if codec.Password == `` {
	//		service.enabled = false
	//	} else {
	//		decodeString, err := base64.StdEncoding.DecodeString(codec.Password)
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//
	//		service.createPasswordFile(string(decodeString))
	//	}
	//
	//	//codec.Encoder
	//}

	return service
}
