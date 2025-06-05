package oci

import (
	"net/http"

	"gofr.dev/pkg/gofr/logging"
)

type gofrError struct {
	msg        string
	statusCode int
	logLevel   logging.Level
}

func (e gofrError) Error() string {
	return e.msg
}

func (e gofrError) StatusCode() int {
	return e.statusCode
}

func (e gofrError) LogLevel() logging.Level {
	return e.logLevel
}

var (
	errInvalidOCICreds = gofrError{
		msg:        "invalid OCI credentials",
		statusCode: http.StatusUnauthorized,
		logLevel:   logging.ERROR,
	}

	errCreateDBClient = gofrError{
		msg:        "failed to create Database client",
		statusCode: http.StatusInternalServerError,
		logLevel:   logging.ERROR,
	}

	errListDBSystems = gofrError{
		msg:        "failed to list DB systems",
		statusCode: http.StatusInternalServerError,
		logLevel:   logging.ERROR,
	}

	errMonitoringClient = gofrError{
		msg:        "failed to create Monitoring client",
		statusCode: http.StatusInternalServerError,
		logLevel:   logging.ERROR,
	}

	errReadingMetrics = gofrError{
		msg:        "error reading metrics for DB system",
		statusCode: http.StatusInternalServerError,
		logLevel:   logging.ERROR,
	}
)
