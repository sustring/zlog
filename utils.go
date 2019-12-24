package zlog

import (
	"fmt"
	"os"
	"runtime"
)

// 操作系统和架构类型
const (
	Linux   = "linux"
	Windows = "windows"
)

// ThirdCallerInfo 返回祖爷调用函数的信息.
func ThirdCallerInfo() (file string, line int, funcName string) {
	pc, file, line, _ := runtime.Caller(3)
	funcName = runtime.FuncForPC(pc).Name()

	return file, line, funcName
}

// GetOS 获取设备操作系统.
func GetOS() (os string) {
	return runtime.GOOS
}

// IsExist 检查文件是否存在,如果存在,返回文件的具体信息.
func IsExist(path string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return fileInfo, nil
	}

	if os.IsNotExist(err) {
		return fileInfo, fmt.Errorf("%s not exist", path)
	}

	return fileInfo, err
}
