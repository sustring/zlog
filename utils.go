package zlog

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

type exitFunc func()

// SecurityExitProcess 监听系统信号,执行exitFunc释放程序资源,优雅退出.
func SecurityExitProcess(exitFunc exitFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Printf("\n[ INFO ] (system) - security exit by %s signal.\n", s)
			exitFunc()
		default:
			fmt.Printf("\n[ INFO ] (system) - unknow exit by %s signal.\n", s)
			exitFunc()
		}
	}
}
