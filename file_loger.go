package zlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FileLogImpl 日志文件执行对象
type FileLogImpl struct {
	log             *log.Logger
	mutex           sync.Mutex
	curFileTime     int64
	logFileFullPath string

	logQueue chan oneLog

	logSwitch     int
	modules       []string
	logFilePath   string
	logFilePrefix string
	isTimeSplit   bool
	splitPeriod   int64
	isSizeSplit   bool
	splitSize     int64
	isClear       bool
	savePeriod    int64
}

func newFileLogImpl() *FileLogImpl {
	var modulesTemp = make([]string, 0)

	modulesTemp = append(modulesTemp, ModulesAll)

	impl := &FileLogImpl{logSwitch: Debug | Info | Notice | Warn | Error | Critical,
		modules:     modulesTemp,
		logFilePath: DefaultLogFilePath,
		isTimeSplit: true,
		splitPeriod: DefaultSplitPeriod,
		isSizeSplit: true,
		splitSize:   DefaultSplitSize,
		logQueue:    make(chan oneLog, QueueMaxNumber),
		isClear:     true,
		savePeriod:  DefaultLogFileSavePeriod}

	go impl.listenLogQueue()

	go impl.clearAndRecycle()

	return impl
}

// Printf 输出i一条日志到文件
func (f *FileLogImpl) Printf(level int, module string, format string, v ...interface{}) {
	if f.checkLevel(level) && f.checkModule(module) {
		log := oneLog{level: level, module: module, format: format, args: v}

		log.callerFile, log.callerLine, log.callerName = ThirdCallerInfo()

		filePathSplit := strings.Split(log.callerFile, "/")
		log.callerFile = filePathSplit[len(filePathSplit)-1]

		f.logQueue <- log
	}
}

func (f *FileLogImpl) writeLog(log oneLog) {
	var prefixStr string

	if Debug == log.level {
		prefixStr = "[ DEBUG  ] - " + "(" + log.module + ") "
	} else if Info == log.level {
		prefixStr = "[ INFO   ] - " + "(" + log.module + ") "
	} else if Notice == log.level {
		prefixStr = "[ NOTICE ] - " + "(" + log.module + ") "
	} else if Warn == log.level {
		prefixStr = "[ WARN   ] - " + "(" + log.module + ") "
	} else if Error == log.level {
		prefixStr = "[ ERROR  ] - " + "(" + log.module + ") "
	} else if Critical == log.level {
		prefixStr = "[CRITICAL] - " + "(" + log.module + ") "
	} else {
		prefixStr = "[UNKNOWN ] - " + "(" + log.module + ") "
	}

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var timeTmp = string(timestamp[10:19])
	timeNow := t.Format("2006-01-02 15:04:05")
	timeStr := timeNow + " @ " + timeTmp + " ▶  "

	prefixStr += fmt.Sprintf("%s:%d %s() %s ", log.callerFile, log.callerLine, log.callerName, timeStr)

	f.log.SetPrefix(prefixStr)
	f.log.Printf(log.format, log.args...)
}

func (f *FileLogImpl) listenLogQueue() {
	for l := range f.logQueue {
		if f.isSplitLogFile() {
			err := f.initFileLogImpl()
			if err != nil {
				f.logQueue <- l
				Prints(Warn, "log", "init fileLog impl error : %s", err)
				continue
			}

			f.writeLog(l)
		} else {
			f.writeLog(l)
		}

		time.Sleep(time.Duration(5) * time.Millisecond)
	}
}

func (f *FileLogImpl) isSplitLogFile() bool {
	t := time.Now()
	//  s - 1565095766
	// ms - 1565095766957
	// us - 1565095766957123
	// ns - 1565095766957123123
	timeNow := t.UTC().Unix()

	if (timeNow - f.curFileTime) >= f.splitPeriod {
		return true
	}

	fileInfo, err := IsExist(f.logFileFullPath)
	if err != nil {
		Prints(Warn, "zlog", "checkout log file error : %s", err)
		return false
	}

	if fileInfo.Size() >= f.splitSize {
		return true
	}

	return false
}

func (f *FileLogImpl) initFileLogImpl() error {
	t := time.Now()
	f.curFileTime = t.UTC().Unix()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var timeTmp = string(timestamp[10:19])
	timeNow := t.Format("2006-01-02 15:04:05")
	timeNow = strings.Replace(timeNow, " ", "@", -1)
	filenameTmp := strings.Replace(timeNow, ":", "-", -1)

	var fileName string
	if GetOS() == Windows {
		fileName = f.logFilePath + "\\" + f.logFilePrefix + "=" + filenameTmp + "@" + timeTmp + ".txt"
	} else if GetOS() == Linux {
		fileName = f.logFilePath + "/" + f.logFilePrefix + "=" + filenameTmp + "@" + timeTmp + ".txt"
	} else {
		return fmt.Errorf("error system type = %s", GetOS())
	}

	f.logFileFullPath = fileName

	var logFile io.Writer
	_, err := os.Stat(f.logFilePath)
	if err == nil {
		logFile, err = os.Create(fileName)
		if err != nil {
			Prints(Warn, "zlog", "create log file error : %s", err)
			return err
		}
	} else {
		if os.IsNotExist(err) {
			Prints(Notice, "zlog", "log dir is not exist, begin to create")
			err = os.MkdirAll(f.logFilePath, 0777)
			if err != nil {
				Prints(Warn, "zlog", "create log dir error : %s", err)
				return err
			}

			logFile, err = os.Create(fileName)
			if err != nil {
				Prints(Warn, "zlog", "create log file error : %s", err)
			}
		} else {
			Prints(Warn, "zlog", "check log dir error : %s", err)
			return err
		}
	}

	f.log = log.New(logFile, "[INIT]-", 0)

	return err
}

func (f *FileLogImpl) checkLevel(level int) bool {
	var retCheck bool

	ret := f.logSwitch & level
	if ret == 0 {
		retCheck = false
	} else {
		retCheck = true
	}

	return retCheck
}

func (f *FileLogImpl) checkModule(module string) bool {
	var retCheck = false

	if len(f.modules) == 0 {
		return false
	}

	for _, value := range f.modules {
		if value == ModulesAll {
			retCheck = true
			break
		} else if value == ModulesNone {
			retCheck = false
			break
		} else if value == module {
			retCheck = true
		} else {
			continue
		}
	}

	return retCheck
}

func (f *FileLogImpl) clearAndRecycle() {
	for {
		if f.isClear {
			_, err := os.Stat(f.logFilePath)
			if err != nil {
				time.Sleep(time.Duration(30) * time.Second)
				continue
			}

			rd, err := ioutil.ReadDir(f.logFilePath)
			if err == nil {
				for _, fi := range rd {
					if fi.IsDir() {
						continue
					}

					var fileFullPath string
					if GetOS() == Windows {
						fileFullPath = f.logFilePath + "\\" + fi.Name()
					} else if GetOS() == Linux {
						fileFullPath = f.logFilePath + "/" + fi.Name()
					} else {
						Prints(Warn, "zlog", "error system type = %s", GetOS())
						continue
					}

					fileInfo, err := IsExist(fileFullPath)
					if err != nil {
						Prints(Warn, "zlog", "check file exist error : %s", err)
						continue
					}

					if strings.Contains(fileInfo.Name(), f.logFilePrefix) {
						timeNow := time.Now()
						timeFile := fileInfo.ModTime()
						timeCmp := timeFile.Add(time.Duration(f.savePeriod*60*60*24) * time.Second)
						// timeCmp := timeFile.Add(time.Duration(20) * time.Second)
						if timeCmp.Before(timeNow) {
							err = os.Remove(fileFullPath)
							if err != nil {
								Prints(Warn, "zlog", "delete file error : %s", err)
							}
						}
					}
				}
			} else {
				Prints(Warn, "zlog", "read dir = %s error : %s", f.logFileFullPath, err)
			}
		}

		time.Sleep(time.Duration(60*60*12) * time.Second)
		// time.Sleep(time.Duration(5) * time.Second)
	}
}

// SetFileLogConfig 配置
func (f *FileLogImpl) SetFileLogConfig(c Config) {

	if c.LogFilePath == "" {
		f.logFilePath = DefaultLogFilePath
	} else {
		f.logFilePath = c.LogFilePath
	}

	if c.LogLevels == 0 {
		f.logSwitch = Warn | Error | Critical
	} else {
		f.logSwitch = c.LogLevels
	}

	f.isTimeSplit = c.IsTimeSplit
	if c.SplitPeriod <= 0 {
		f.splitPeriod = DefaultSplitPeriod
	} else {
		f.splitPeriod = c.SplitPeriod
	}

	f.isSizeSplit = c.IsSizeSplit
	if c.SplitSize <= 0 {
		f.splitSize = DefaultSplitSize
	} else {
		f.splitSize = c.SplitSize
	}

	if !c.IsClear {
		f.isClear = c.IsClear
	}

	if c.SavePeriod <= 0 {
		f.savePeriod = DefaultLogFileSavePeriod
	} else {
		f.savePeriod = c.SavePeriod
	}

	if c.LogFilePrefix == "" {
		f.logFilePrefix = DefaultLogFilePrefix
	} else {
		f.logFilePrefix = c.LogFilePrefix
	}

	if len(c.Modules) == 0 {
		f.modules = append(f.modules, ModulesAll)
	} else {
		f.modules = c.Modules
	}

	f.log = nil
}

var instanceSysFileLogImpl *FileLogImpl
var instanceFileLogImplOnce sync.Once

func getSysFileLogInstance() *FileLogImpl {
	instanceFileLogImplOnce.Do(func() {
		instanceSysFileLogImpl = newFileLogImpl()
	})

	return instanceSysFileLogImpl
}

func fileAddOneLog() *FileLogImpl {
	if getSysFileLogInstance().log == nil {
		getSysFileLogInstance().mutex.Lock()

		if getSysFileLogInstance().log == nil {
			err := getSysFileLogInstance().initFileLogImpl()
			if err != nil {
				panic(fmt.Sprintf("init fileLog impl error : %s", err))
			}
		}

		getSysFileLogInstance().mutex.Unlock()
	}

	return getSysFileLogInstance()
}

func initFileLoger() {
	getSysFileLogInstance()
}
