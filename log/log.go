package log

import "github.com/gobuffalo/logger"

// injected by app.go
var Log logger.FieldLogger

func Error(i interface{}) {
	Log.Error(i)
}

func Errorf(s string, i ...interface{}) {
	Log.Errorf(s, i)
}

func Info(i interface{}) {
	Log.Info(i)
}

func Infof(s string, i ...interface{}) {
	Log.Infof(s, i)
}

func Fatal(i interface{}) {
	Log.Fatal(i)
}
