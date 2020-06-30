package logging

import "github.com/sirupsen/logrus"

var logging *logrus.Logger

func init() {
	logging = logrus.New()
	logging.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	return logging
}
