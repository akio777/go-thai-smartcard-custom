package logger

import "github.com/sirupsen/logrus"

var (
	_logger *logrus.Logger
)

func InitLogger() {
	_logger = logrus.New()
	_logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
		DisableQuote:    true,
	})
	_logger.SetReportCaller(true)
}

func LOGGER() *logrus.Logger {
	return _logger
}
