//This file is part of xtagger. ©2023-2026 Jörg Walter.
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

package config

import (
	"github.com/jwdev42/xtagger/internal/hashes"
	"log/slog"
)

// Return program Preferences filled with default values
func DefaultPreferences() *Preferences {
	return &Preferences{
		UseHash:      hashes.SHA256,
		Paths:        []string{},
		Names:        []string{},
		LogLevel:     defaultLogLevel(),
		Threads:      1,
		UseRecursion: true,
	}
}

// Preferences hold's the program configuration during runtime
type Preferences struct {
	Command         Command // Specified command
	UseHash         hashes.Algo
	Paths           []string       // User-specified directories and files
	Names           []string       // User-specified record identifiers
	LogLevel        *slog.LevelVar // Currently active log level
	Quota           int64
	Threads         int
	FollowSymlinks  bool
	PrintRecords    bool
	UsePrint0       bool
	UseRecursion    bool
	TagConstraint   TagConstraint
	UntagConstraint UntagConstraint
	PrintConstraint PrintConstraint
}

// defaultLogLevel returns a LevelVar that hold's the program's
// default log level.
func defaultLogLevel() *slog.LevelVar {
	lvl := &slog.LevelVar{}
	lvl.Set(slog.LevelError)
	return lvl
}
