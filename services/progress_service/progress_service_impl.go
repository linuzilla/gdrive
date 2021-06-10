package progress_service

import (
	"fmt"
	"github.com/linuzilla/gdrive/utils"
	"io"
	"os"
	"time"
)

type progressServiceImpl struct {
	name           string
	reader         utils.ProgressReader
	dataSize       int64
	doneChannel    chan error
	uploadDownload string
	signalChan     chan os.Signal
}

func (impl *progressServiceImpl) Close() {
	if impl.doneChannel != nil {
		close(impl.doneChannel)
		impl.doneChannel = nil
	}
}

func (impl *progressServiceImpl) ExecAndWait(goroutine func(ProgressService, chan error)) error {
	go goroutine(impl, impl.doneChannel)

	return impl.waitAndShowProgress()
}

func (impl *progressServiceImpl) SetName(name string) {
	impl.name = name
}
func (impl *progressServiceImpl) SetSize(dataSize int64) {
	impl.dataSize = dataSize
}

func (impl *progressServiceImpl) WrapperReader(reader io.Reader) io.Reader {
	return impl.reader.SetReader(reader)
}

func (impl *progressServiceImpl) waitAndShowProgress() error {
	started := time.Now()
	last := time.Now()
	var lastRead int64

	for i := 0; ; i += 1 {
		select {
		case err := <-impl.doneChannel:
			progressSize := impl.reader.N()
			bps := float64(progressSize) / float64(time.Now().Sub(started)/time.Second) / 1024.0
			fmt.Printf("%s %s (%d bytes, %.2f Kbytes/sec)\n",
				impl.name, impl.uploadDownload, progressSize, bps)
			return err

		case <-impl.signalChan:
			fmt.Fprint(os.Stderr, "\n\n**** Break ****\n\n")
			return fmt.Errorf("user break")

		default:
			time.Sleep(time.Millisecond * 100)
			if impl.dataSize > 0 && i%100 == 99 {
				now := time.Now()
				elapsedTotal := now.Sub(started) / time.Second
				elapsedLast := now.Sub(last) / time.Second

				currentRead := impl.reader.N()
				periodRead := currentRead - lastRead

				avgRead := float64(currentRead) / float64(elapsedTotal) / 1024.0
				currRead := float64(periodRead) / float64(elapsedLast) / 1024.0

				fmt.Printf("%s: %s %d (current: %.2f Kbytes/sec, avg: %.2f Kbytes/sec, %.2f%%)\n",
					impl.name, impl.uploadDownload, currentRead, currRead, avgRead,
					float64(currentRead)*100/float64(impl.dataSize))

				lastRead = currentRead
				last = now
			}
		}
	}
}
