package logging

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/mushroomyuan/gorder/common/tracing"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
	logrus.AddHook(&traceHook{})
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{

		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		logger.SetFormatter(&prefixed.TextFormatter{
			ForceFormatting: true,
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
		})

	}
}

func InfoWithCost(ctx context.Context, fields logrus.Fields, start time.Time, format string, args ...any) {
	fields["Cost"] = time.Since(start).Milliseconds()
	Infof(ctx, fields, format, args...)

}

func Infof(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithTime(time.Now()).WithFields(fields).Infof(format, args...)
}

func Errorf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithTime(time.Now()).WithFields(fields).Errorf(format, args...)
}

func Panicf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithTime(time.Now()).WithFields(fields).Panicf(format, args...)
}

func Warnf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithTime(time.Now()).WithFields(fields).Warnf(format, args...)
}

func logf(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Logf(level, format, args...)
}

type traceHook struct {
}

func (t traceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t traceHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data["trace"] = tracing.TraceID(entry.Context)
		entry = entry.WithTime(time.Now())
	}
	return nil
}
