package zlog

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	backgroundColour = 0
)

// ScreenLogImpl 日志打印执行对象
type ScreenLogImpl struct {
	logSwitch int
	modules   []string
}

func newScreenLogImpl() *ScreenLogImpl {
	var modulesTemp = make([]string, 0)
	modulesTemp = append(modulesTemp, ModulesAll)

	impl := &ScreenLogImpl{logSwitch: Debug | Info | Notice | Warn | Error | Critical, modules: modulesTemp}

	return impl
}

func (s *ScreenLogImpl) showOneLog(log oneLog) {
	var isHighLight int
	var background int
	var foreground int

	if log.level == Debug {
		background = backgroundColour
		foreground = 37
	} else if log.level == Info {
		background = backgroundColour
		foreground = 32
	} else if log.level == Notice {
		background = backgroundColour
		foreground = 34
	} else if log.level == Warn {
		background = backgroundColour
		foreground = 33
	} else if log.level == Error {
		background = backgroundColour
		foreground = 35
	} else if log.level == Critical {
		background = backgroundColour
		foreground = 31
	} else {
		background = backgroundColour
		foreground = 36
	}

	var colourBegin = ""
	var colourEnd = ""
	if GetOS() == "windows" {
		colourBegin = ""
		colourEnd = fmt.Sprintf("\n")
	} else if GetOS() == Linux {
		colourBegin = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, isHighLight, background, foreground)
		colourEnd = fmt.Sprintf("%c[0m\n", 0x1B)
	} else {
		fmt.Printf("error system type = %s\n", GetOS())
	}

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var timeTmp = string(timestamp[10:19])
	timeNow := t.Format("2006-01-02 15:04:05")
	timeStr := timeNow + " @ " + timeTmp + " ▶  "

	var formatLog = fmt.Sprintf("%s%s(%s) %s:%d %s() %s %s%s",
		colourBegin,
		LevelMap[log.level],
		log.module,
		log.callerFile,
		log.callerLine,
		log.callerName,
		timeStr,
		log.format,
		colourEnd)

	fmt.Printf(formatLog, log.args...)
}

func (s *ScreenLogImpl) checkLevel(level int) bool {
	var retCheck bool

	ret := s.logSwitch & level
	if ret == 0 {
		retCheck = false
	} else {
		retCheck = true
	}

	return retCheck
}

func (s *ScreenLogImpl) checkModule(module string) bool {
	var retCheck = false

	if len(s.modules) == 0 {
		return false
	}

	for _, value := range s.modules {
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

// Printf 日志打印
func (s *ScreenLogImpl) Printf(level int, module string, format string, v ...interface{}) {
	if s.checkLevel(level) && s.checkModule(module) {
		log := oneLog{level: level, module: module, format: format, args: v}

		log.callerFile, log.callerLine, log.callerName = ThirdCallerInfo()

		filePathSplit := strings.Split(log.callerFile, "/")
		log.callerFile = filePathSplit[len(filePathSplit)-1]

		s.showOneLog(log)
	}
}

// SetScreenLogConfig 设置配置
func (s *ScreenLogImpl) SetScreenLogConfig(c Config) {
	s.logSwitch = c.LogLevels

	if len(c.Modules) == 0 {
		s.modules = append(s.modules, ModulesAll)
	} else {
		s.modules = c.Modules
	}
}

var instanceScreenLogImpl *ScreenLogImpl
var instanceScreenLogImplOnce sync.Once

func getScreenLogImplInstance() *ScreenLogImpl {
	instanceScreenLogImplOnce.Do(func() {
		instanceScreenLogImpl = newScreenLogImpl()
	})

	return instanceScreenLogImpl
}

func screenAddOneLog() *ScreenLogImpl {
	return getScreenLogImplInstance()
}

func initScreenLoger() {
	getScreenLogImplInstance()
}
