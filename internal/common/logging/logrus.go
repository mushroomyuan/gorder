package logging

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func NewLogrusLogger() *logrus.Logger {
	return logrus.New()
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "serverity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		logger.SetFormatter(&prefixed.TextFormatter{
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}
