package progress_service

import (
	"github.com/linuzilla/gdrive/services/interrupt_service"
	"github.com/linuzilla/gdrive/utils"
	"io"
)

type ProgressServiceFactory interface {
	NewInstance(uploadOrDownload string) ProgressService
}

type ProgressService interface {
	Close()
	ExecAndWait(goroutine func(ProgressService, chan error)) error
	SetName(name string)
	SetSize(dataSize int64)
	WrapperReader(reader io.Reader) io.Reader
}

func NewFactory() ProgressServiceFactory {
	return &progressServiceFactory{}
}

type progressServiceFactory struct {
	InterruptSvc interrupt_service.InterruptService `inject:"*"`
}

func (factory *progressServiceFactory) NewInstance(uploadOrDownload string) ProgressService {
	progressReader := utils.NewProgressReader(nil)
	errChan := make(chan error, 1)

	return &progressServiceImpl{
		reader:         progressReader,
		doneChannel:    errChan,
		uploadDownload: uploadOrDownload,
		signalChan:     factory.InterruptSvc.Channel(),
	}
}
