package utils

import (
	"path"
	"runtime"
)

type CallerInfo struct {
	funcName   string
	fileName   string
	lineNumber int
}

func (caller *CallerInfo) FuncName() string {
	return caller.funcName
}

func (caller *CallerInfo) FileName() string {
	return caller.fileName
}

func (caller *CallerInfo) LineNumber() int {
	return caller.lineNumber
}

func GetCallerInfo(skip int) *CallerInfo {
	var funcName string
	var fileName string
	var lineNumber int

	fpcs := make([]uintptr, 1)

	if n := runtime.Callers(2+skip, fpcs); n == 0 {
		funcName = "n/a"
	} else if fcn := runtime.FuncForPC(fpcs[0] - 1); fcn == nil {
		funcName = "n/a"
	} else {
		funcName = fcn.Name()
	}

	if _, file, no, ok := runtime.Caller(1 + skip); ok {
		fileName = path.Base(path.Dir(file)) + "/" + path.Base(file)
		lineNumber = no
	}

	return &CallerInfo{
		funcName:   funcName,
		fileName:   fileName,
		lineNumber: lineNumber,
	}
}
