//This file is part of xtagger. ©2023 Jörg Walter.
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see <https://www.gnu.org/licenses/>.

package logging

import (
	"fmt"
	"log/slog"
	"strings"
)

// Custom log levels
const (
	LevelTrace  slog.Level = -8
	LevelNotice slog.Level = 2
	LevelFatal  slog.Level = 12
)

// Log level names
const (
	LevelNameTrace  = "TRACE"
	LevelNameDebug  = "DEBUG"
	LevelNameInfo   = "INFO"
	LevelNameNotice = "NOTICE"
	LevelNameWarn   = "WARN"
	LevelNameError  = "ERROR"
	LevelNameFatal  = "FATAL"
)

// Map log level numbers to strings
var logLevelNames = map[slog.Level]string{
	LevelTrace:      LevelNameTrace,
	slog.LevelDebug: LevelNameDebug,
	slog.LevelInfo:  LevelNameInfo,
	LevelNotice:     LevelNameNotice,
	slog.LevelWarn:  LevelNameWarn,
	slog.LevelError: LevelNameError,
	LevelFatal:      LevelNameFatal,
}

// Apply custom log level names to a handler.
// Use this for HandlerOptions.ReplaceAttr.
func ReplaceLogLevelNames(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		if level, ok := a.Value.Any().(slog.Level); ok {
			// Set log level string
			a.Value = slog.StringValue(LogLevelName(level))
		}
	}
	return a
}

// Return the string representation of a log level, has support for
// the application-defined additional log levels.
func LogLevelName(level slog.Level) string {
	if customName, exists := logLevelNames[level]; exists {
		return customName
	}
	return level.String()
}

// Return the log level for the given name. If no name is found,
// the function will return a non-nil error.
// Will always return slog.LevelInfo on error.
func ParseLogLevelName(name string) (slog.Level, error) {
	switch strings.ToUpper(name) {
	case LevelNameTrace:
		return LevelTrace, nil
	case LevelNameDebug:
		return slog.LevelDebug, nil
	case LevelNameInfo:
		return slog.LevelInfo, nil
	case LevelNameNotice:
		return LevelNotice, nil
	case LevelNameWarn:
		return slog.LevelWarn, nil
	case LevelNameError:
		return slog.LevelError, nil
	case LevelNameFatal:
		return LevelFatal, nil
	}
	return slog.LevelInfo, fmt.Errorf("No log level found for name %q", name)
}
