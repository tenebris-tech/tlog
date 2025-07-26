/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * All rights reserved. See LICENSE file for details.                         *
 ******************************************************************************/

package tlog

// Logger defines an interface for logging and log messages
type Logger interface {
	Debug(string)
	Info(string)
	Notice(string)
	Warning(string)
	Error(string)
	Fatal(string)
	Debugf(string, ...any)
	Infof(string, ...any)
	Noticef(string, ...any)
	Warningf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Close()
}
