package dot

import "github.com/sirupsen/logrus"

var std = NewGContext()

// Logger 日志打印对象
func Logger() *logrus.Entry { return std.Logger }
