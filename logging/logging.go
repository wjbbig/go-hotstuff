package logging

import "github.com/sirupsen/logrus"

var logging *logrus.Logger

func init() {
	logging = logrus.New()
}

func GetLogger() *logrus.Logger {
	return logging
}
