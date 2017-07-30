/* win98-logout - https://github.com/gdm85/win98-logout
Copyright (C) 2017 gdm85

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/vaughan0/go-ini"
)

// LogoutOption describes a logout option.
type LogoutOption struct {
	// Label is the displayed text for the logout option.
	Label string
	// Command is the command to execute when the logout option is activated.
	Command string
}

var (
	defaultOptions = []LogoutOption{
		{Label: "Stand by"},
		{Label: "Shutdown", Command: "sudo shutdown -P now"},
		{Label: "Restart", Command: "sudo shutdown -r now"},
		{Label: "Switch to a VT", Command: "xdotool key ctrl+alt+F1"},
	}
)

func loadConfig(path string) ([]LogoutOption, error) {
	f, err := ini.LoadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// return default
			return defaultOptions, nil
		}
		return nil, err
	}

	var options []LogoutOption

	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("option_%d_label", i)

		label, ok := f.Get("options", name)
		if !ok {
			break
		}

		name = fmt.Sprintf("option_%d_command", i)
		cmd, ok := f.Get("options", name)
		if !ok {
			return nil, errors.New("name specified, but not the option")
		}

		options = append(options, LogoutOption{Label: label, Command: cmd})
	}

	return options, nil
}
