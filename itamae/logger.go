package itamae

import (
	"os"

	"charm.land/log/v2"
)

// Logger is the shared structured logger for user-facing messages on stderr.
var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    false,
	ReportTimestamp: false,
	TimeFormat:      "2006-01-02 15:04:05",
	Level:           log.InfoLevel,
})
