package logger

import (
	"io/ioutil"
	"os"

	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

// Init initialize the logging system
func Init(discard bool) {
	if discard {
		logrus.SetOutput(ioutil.Discard)
		return
	}
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = time.RFC3339
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

// Info logging with INFO level
func Info(msg string, a ...interface{}) {
	logrus.Info(fmt.Sprintf(msg, a...))
}

// Error logging with ERROR level
func Error(msg string, a ...interface{}) {
	logrus.Error(fmt.Sprintf(msg, a...))
}

// Errorf logging with ERROR level and returns an error struct
func Errorf(msg string, a ...interface{}) error {
	err := fmt.Errorf(msg, a...)
	logrus.Error(err.Error())
	return err
}

// Debug logging with DEBUG level
func Debug(msg string, a ...interface{}) {
	logrus.Debug(fmt.Sprintf(msg, a...))
}

// APILogger provide Logrus integration for Gin APIs
func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.String()
		userAgent := c.Request.UserAgent()

		timeFormatted := end.Format("2006-01-02 15:04:05")

		msg := fmt.Sprintf(
			"%s %s \"%s %s\" %d %s %s",
			clientIP,
			timeFormatted,
			method,
			path,
			statusCode,
			latency,
			userAgent,
		)

		logrus.StandardLogger().WithFields(logrus.Fields{
			"time":       timeFormatted,
			"method":     method,
			"path":       path,
			"latency":    latency,
			"ip":         clientIP,
			"comment":    comment,
			"status":     statusCode,
			"user-agent": userAgent,
		}).Info(msg)
	}
}
