package logs

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// InitLogs init formats for the log files
//
// stdout logs have color and are simplified
//
// file logs are formated in the buildHookForOutputLogs method
func InitLogs() {
	logrus.SetReportCaller(true)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {

			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	})

	rotateFileHook, err := buildHookForOutputLogs()
	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	logrus.AddHook(rotateFileHook)
}

//buildHookForOutputLogs format for the file logs
func buildHookForOutputLogs() (logrus.Hook, error) {
	now := time.Now()
	path := fmt.Sprintf("/var/log/oyster-guaridan/%s", getFileName(now))

	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}

	
	//more information on the rotateFileHook https://github.com/natefinch/lumberjack#type-logger
	return rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: path,
		// MaxSize is the maximum size in megabytes of the log file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize: 20, // megabytes

		// MaxBackups is the maximum number of old log files to retain.  The default
		// is to retain all old log files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups: 10, // number of files

		// MaxAge is the maximum number of days to retain old log files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default is not to remove old log files
		// based on age.
		MaxAge: 2, //days

		Level:     logrus.InfoLevel,
		Formatter: formatter,
	})
}

func getFileName(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d.log",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
