/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * All rights reserved. See LICENSE file for details.                         *
 ******************************************************************************/

package tlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TLog struct {
	fileHandle *os.File
	logfile    string
	logStdout  bool
	debug      bool
	logLevel   bool
	prefix     string
	dateFormat string
}

// This package implements interfaces.Logger
var _ Logger = (*TLog)(nil)

// Option is a function that configures a TLog
type Option func(*TLog) error

// New creates a new instance of TLog with the provided options
func New(options ...Option) (Logger, error) {
	t := &TLog{
		logLevel:   true,
		dateFormat: "2006-01-02 15:04:05",
	}

	for _, option := range options {
		if err := option(t); err != nil {
			return nil, err
		}
	}

	// Call the OS-specific constructor
	return t.open()
}

// WithPrefix sets a process name or similar short identifier
//
//goland:noinspection GoUnusedExportedFunction
func WithPrefix(prefix string) Option {
	return func(u *TLog) error {
		if prefix == "" {
			u.prefix = ""
		} else {
			u.prefix = " " + strings.TrimSpace(prefix)
		}
		return nil
	}
}

// WithDateFormat sets the date format for the TLog
//
//goland:noinspection GoUnusedExportedFunction
func WithDateFormat(dateFormat string) Option {
	return func(u *TLog) error {
		u.dateFormat = dateFormat
		return nil
	}
}

// WithLogFile sets the log file for the TLog
//
//goland:noinspection GoUnusedExportedFunction
func WithLogFile(logfile string) Option {
	return func(u *TLog) error {
		u.logfile = logfile
		return nil
	}
}

// WithLogStdout enables or disables logging to stdout
//
//goland:noinspection GoUnusedExportedFunction
func WithLogStdout(logStdout bool) Option {
	return func(u *TLog) error {
		u.logStdout = logStdout
		return nil
	}
}

// WithLevel enables or disables logging the level
//
//goland:noinspection GoUnusedExportedFunction
func WithLevel(logLevel bool) Option {
	return func(u *TLog) error {
		u.logLevel = logLevel
		return nil
	}
}

// WithDebug enables or disables debug logging
//
//goland:noinspection GoUnusedExportedFunction
func WithDebug(debug bool) Option {
	return func(u *TLog) error {
		u.debug = debug
		return nil
	}
}

// open sets up the logger. This function is not exported, it is called by New
func (t *TLog) open() (*TLog, error) {
	var err error
	var fh *os.File

	if t.logfile != "" {

		// Sanitize the file path
		t.logfile = filepath.Clean(t.logfile)

		// Create the directory if it doesn't exist
		dir := filepath.Dir(t.logfile)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open the log file
		fh, err = os.OpenFile(t.logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			t.fileHandle = nil
			// If unable to log to file, force stdout logging
			t.logStdout = true
		} else {
			t.fileHandle = fh

			// Attempt to set the file mode to 0644 on a best-effort basis
			_ = os.Chmod(t.logfile, 0644)
		}
	} else {
		// If no log file is specified, force stdout logging
		t.logStdout = true
	}
	return t, nil
}

// Close closes the logger
func (t *TLog) Close() {
	if t.fileHandle != nil {
		_ = t.fileHandle.Sync()
		_ = t.fileHandle.Close()
	}
}

// formatMessage formats the log message with a timestamp.
func (t *TLog) formatMessage(level string, message string) string {
	var levelStr string
	if t.logLevel {
		levelStr = " [" + level + "]"
	} else {
		levelStr = ""
	}
	return fmt.Sprintf("%s%s%s %s",
		time.Now().Format(t.dateFormat),
		t.prefix, levelStr, message)
}

// writeLog writes a log message
func (t *TLog) writeLog(level string, message string) {

	tmp := t.formatMessage(level, message) + "\n"

	//  Write and flush
	if t.fileHandle != nil {
		_, _ = t.fileHandle.WriteString(tmp)
		_ = t.fileHandle.Sync()
	}

	if t.logStdout {
		_, _ = os.Stdout.Write([]byte(tmp))
	}
}

// Debug logs a debug message.
func (t *TLog) Debug(message string) {
	if t.debug {
		t.writeLog("DEBUG", message)
	}
}

// Info logs an informational message.
func (t *TLog) Info(message string) {
	t.writeLog("INFO", message)
}

// Notice logs a notice message.
func (t *TLog) Notice(message string) {
	t.writeLog("NOTICE", message)
}

// Warning logs a warning message.
func (t *TLog) Warning(message string) {
	t.writeLog("WARNING", message)
}

// Error logs an error message.
func (t *TLog) Error(message string) {
	t.writeLog("ERROR", message)
}

// Fatal logs a fatal error message.
func (t *TLog) Fatal(message string) {
	t.writeLog("FATAL", message)
	t.FatalExit()
}

// Debugf logs a formatted debug message.
func (t *TLog) Debugf(format string, v ...any) {
	if t.debug {
		t.writeLog("DEBUG", fmt.Sprintf(format, v...))
	}
}

// Infof logs a formatted informational message.
func (t *TLog) Infof(format string, v ...any) {
	t.writeLog("INFO", fmt.Sprintf(format, v...))
}

// Noticef logs a formatted notice message.
func (t *TLog) Noticef(format string, v ...any) {
	t.writeLog("NOTICE", fmt.Sprintf(format, v...))
}

// Warningf logs a formatted warning message.
func (t *TLog) Warningf(format string, v ...any) {
	t.writeLog("WARNING", fmt.Sprintf(format, v...))
}

// Errorf logs a formatted error message.
func (t *TLog) Errorf(format string, v ...any) {
	t.writeLog("ERROR", fmt.Sprintf(format, v...))
}

// Fatalf logs a formatted fatal message.
func (t *TLog) Fatalf(format string, v ...any) {
	t.writeLog("FATAL", fmt.Sprintf(format, v...))
	t.FatalExit()
}

// FatalExit attempts to close the log and exits with a status code of 1
func (t *TLog) FatalExit() {
	t.writeLog("FATAL", "Exiting with code 1 on fatal error")
	t.Close()
	os.Exit(1)
}
