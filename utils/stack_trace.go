package utils

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

func StackTrace(skip int, depth int) string {
	strlist := []string{}

	var funcName string
	var fileName string
	var lineNumber int

	fpcs := make([]uintptr, 1)

	for i := 0; i < depth; i++ {
		if n := runtime.Callers(2+skip+i, fpcs); n == 0 {
			funcName = "n/a"
		} else if fcn := runtime.FuncForPC(fpcs[0] - 1); fcn == nil {
			funcName = "n/a"
		} else {
			funcName = fcn.Name()
		}

		if _, file, no, ok := runtime.Caller(1 + skip + i); ok {
			fileName = path.Base(path.Dir(file)) + "/" + path.Base(file)
			lineNumber = no
		}
		strlist = append(strlist, fmt.Sprintf("\tat %s (%s:%d)\n", funcName, fileName, lineNumber))
	}

	return strings.Join(strlist, "")
}

func Print(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(1, 1))
}

func Println(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(1, 1))
}

func Print2(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(2, 1))
}

func Print3(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(3, 1))
}

func Print4(v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(4, 1))
}

func TracePrint(skip int, cnt int, v ...interface{}) {
	fmt.Print(fmt.Sprintln(v...) + StackTrace(skip+1, cnt))
}

func TraceDetail(o interface{}) {
	fmt.Print(Detail(o) + StackTrace(1, 1))
}

func Printf(format string, v ...interface{}) {
	fmt.Print(fmt.Sprintf(format, v...) + "\n" + StackTrace(1, 1))
}
