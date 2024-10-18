package logger

import "github.com/sirupsen/logrus"

type LogrusWrapper struct {
	Logger *logrus.Logger
}

func (l *LogrusWrapper) Info(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(toFields(keyvals...)).Info(msg)
}

func (l *LogrusWrapper) Error(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(toFields(keyvals...)).Error(msg)
}

func (l *LogrusWrapper) Warn(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(toFields(keyvals...)).Warn(msg)
}

func toFields(keyvals ...interface{}) logrus.Fields {
	fields := logrus.Fields{}
	for i := 0; i < len(keyvals)-1; i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		fields[key] = keyvals[i+1]
	}
	return fields
}
