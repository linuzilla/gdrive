package execution_service

import (
	"bytes"
	"io"
)

type CommandExecutor interface {
	Exec(dict map[string]string, reader io.Reader) (*bytes.Buffer, error)
	Pipe(dict map[string]string, reader io.Reader, writer io.Writer) error
}

type CommandExecutionService interface {
	Executor(command string) CommandExecutor
}
