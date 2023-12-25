package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Level int

var (
	FileHandle         *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func init() {
	// 创建日志文件
	FileHandle = openLogFile()
	// 重写日志
	logger = log.New(FileHandle, DefaultPrefix, log.LstdFlags)
}

func setPrefix(level Level) {
	// 跳了2层才达到真正的栈，所以DefaultCallerDepth=2
	// 比如Debug，那么Debug一层，setPrefix一层
	pc, file, line, ok := runtime.Caller(DefaultCallerDepth)

	//-test start
	pcs := make([]uintptr, DefaultCallerDepth)
	_ = runtime.Callers(0, pcs)
	pcNames := make([]string, DefaultCallerDepth)
	// dynamic get the package name and the minimum caller depth
	println("+++++++++++")
	for i := 0; i < DefaultCallerDepth; i++ {
		funcName := runtime.FuncForPC(pcs[i]).Name()
		pcNames = append(pcNames, funcName)
		fmt.Printf("%v \n", getPackage(funcName))
	}
	fmt.Printf("%v \n", pcNames)

	println("+++++++++++")
	//-test end
	if ok {
		logPrefix = fmt.Sprintf("[%s]-[%s]-[%s:%d]", levelFlags[level], runtime.FuncForPC(pc).Name(), filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}

// Get the package where the current method is located
func getPackage(str string) string {
	for {
		lastPeriod := strings.LastIndex(str, ".")
		lastSlash := strings.LastIndex(str, "/")
		if lastPeriod > lastSlash {
			str = str[:lastPeriod]
		} else {
			break
		}
	}
	return str
}

func Debug(v ...any) {
	setPrefix(DEBUG)
	logger.Println(v)
}
func Info(v ...any) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...any) {
	setPrefix(WARN)
	logger.Println(v)
}

func Error(v ...any) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...any) {
	setPrefix(FATAL)
	logger.Println(v)
}
