package dot

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// GContext 全局上下文
type GContext struct {
	Logger *logrus.Entry
}

// NewGContext 新建一个全局上下文
func NewGContext() *GContext {
	return &GContext{
		Logger: logrus.WithFields(nil),
	}
}

// SetupLogrus 配置logrus
func SetupLogrus() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			file = f.File
			function = ""
			if len(f.File) > 0 {
				last := strings.LastIndex(f.File, "/")
				file = fmt.Sprintf("%s:%d", f.File[last+1:], f.Line)
			}
			return
		},
	})
	logrus.SetReportCaller(true)
}
