package msgplog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	DEBUG = 1
	INFO = 2
	WARRING = 3
	ERROR = 4
)

var Logger = new(MSGPLog)

type MSGPLog struct {
	logLevel int
	logger *log.Logger
}

//level 记录日志级别 4：error 3：warring及以上 2：info及以上 1：debug及以上
func (ml *MSGPLog) InitMSGPLog(isfile bool,level int){
	if isfile {
		logFile, _ := os.Create("."+ string(filepath.Separator) + time.Now().Format("20060102_150405") + ".txt")
		ml.logger = log.New(logFile, "MSGP ", log.Lshortfile | log.Ldate | log.Ltime)
	} else {
		ml.logger = log.New(os.Stdout, "MSGP ", log.Lshortfile | log.Ldate | log.Ltime)
	}
	ml.logLevel = level
}

func (ml *MSGPLog) Debug(content interface{}){
	if ml.logLevel <= DEBUG {
		ml.logger.Output(2, fmt.Sprintf(" DEBUG: %v", content))
	}
}

//打印Info
func (ml *MSGPLog) Info(content interface{}){
	if ml.logLevel <= INFO {
		ml.logger.Output(2, fmt.Sprintf(" INFO: %v", content))
	}
}

//打印Info,格式化打印
func (ml *MSGPLog) Infof(fmtstr string, content ...interface{}){
	if ml.logLevel <= INFO {
		ml.logger.Output(2, fmt.Sprintf(" INFO: " + fmtstr, content))
	}
}

//打印Warring
func (ml *MSGPLog) Warring(content interface{}){
	if ml.logLevel <= WARRING {
		ml.logger.Output(2, fmt.Sprintf(" WARRING: %v", content))
	}
}

//打印Error,退出程序
func (ml *MSGPLog) Error(err interface{}){
	if ml.logLevel <= ERROR {
		s := fmt.Sprintf(" ERROR: %v", err)
		ml.logger.Output(2, s)
		panic(s)
	}
}

//打印Error,不退出程序
func (ml *MSGPLog) NError(err interface{}){
	if ml.logLevel <= ERROR {
		s := fmt.Sprintf(" ERROR: %v", err)
		ml.logger.Output(2, s)
	}
}

//打印Error,不退出程序
func (ml *MSGPLog) NErrorf(fmtstr string, err ...interface{}){
	if ml.logLevel <= ERROR {
		s := fmt.Sprintf(" ERROR: "+fmtstr, err)
		ml.logger.Output(2, s)
	}
}

//打印处理结果
func (ml *MSGPLog) Result(title string, content interface{}){
	ml.logger.Output(2, fmt.Sprintf(" RESULT>>%v\n%v", title, content))
}