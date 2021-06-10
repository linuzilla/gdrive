package intf

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type CryptoService interface {
	IsEnabled() bool
	GetPassword() string
	SetPassword(password string)
	//EncryptFile(source string, target string)
	//DecryptFile(source string, target string)
	//EncryptFileName(fileName string) string
	//DecryptFileName(fileName string) (string, error)
	//JustDecryptFileName(fileName string) string
	//DecodeViaPipe(reader io.Reader, writer io.Writer) error
	//EncodeViaPipe(reader io.Reader, writer io.Writer) error
	DecodeViaPipe(reader io.Reader, writer io.Writer) error
	DecodeAndDownloadFromReader(fileName string, reader io.Reader) error
	EncodeAndUploadReader(fileName string, updateWriter func(writer io.Writer) io.Writer, callback func(reader io.Reader) error) error
	Close()

	EncryptFileNameWithMd5(fileName string, md5 string) string
	DecryptFileNameWithMd5(file *drive.File) (string, string, error)

	GetEncoder() string
	GetDecoder() string
}
