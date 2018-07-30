package log

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)



// NewBeegolog new object beego log
func NewBeegoLog(path string,debug,console,separate bool)(Hlog,error){

	bLog:=logs.NewLogger()
	if path!=""{
		err := os.MkdirAll(filepath.Dir(path), 0664)
		if err != nil {
			logs.Error("fail to create log dir")
			return nil,err
		}
		if separate{
			err = bLog.SetLogger(logs.AdapterMultiFile, fmt.Sprintf(`{"filename":"%s","separate":["error"]}`,path))
			if err != nil {
				logs.Error("fail to config logrus")
			}
		}else{
			err =bLog.SetLogger(logs.AdapterFile,fmt.Sprintf(`{"filename":"%s"}`,path))
			if err != nil {
				logs.Error("fail to config logrus")
			}
		}
		
	}

	bLog.EnableFuncCallDepth(true)
	bLog.SetLogFuncCallDepth(4)
	if debug{
		bLog.SetLevel(logs.LevelDebug)
	}else{
		logs.SetLevel(logs.LevelInfo)
	}
	if console{
		bLog.SetLogger("console")
	}
	return &BeegoLog{logger:bLog},nil
}

// TranBeeLog transform beego Hlog
func TranBeeLog(b *logs.BeeLogger)Hlog{
	return &BeegoLog{logger:b}
}

// BeegoLog beego log
type BeegoLog struct {
	logger *logs.BeeLogger
}

// Debug log debug
func (b *BeegoLog)Debug(f interface{},v ...interface{}){
	b.logger.Debug(formatLog(f,v...))
}

// Info log info
func (b *BeegoLog)Info(f interface{},v ...interface{}){
	b.logger.Info(formatLog(f,v...))

}

func (b *BeegoLog)Warn(f interface{},v ...interface{}){
	b.logger.Warning(formatLog(f,v...))
}

// Error log error
func (b *BeegoLog)Error(f interface{},v ...interface{}){
	b.logger.Error(formatLog(f,v...))
}


func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}

	return fmt.Sprintf(msg, v...)
}