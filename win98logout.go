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
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"unsafe"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//go:generate sh -c "go run vendor/github.com/mattn/go-gtk/tools/make_inline_pixbuf/make_inline_pixbuf.go iconPNG icons/shutdown.png > icon.gen.go"

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not get current user: %v\n", err)
		os.Exit(1)
	}

	options, err := loadConfig(usr.HomeDir + `/.config/win98logout/win98logout.ini`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not read configuration: %v\n", err)
		os.Exit(1)
	}

	gtk.Init(&os.Args)

	pb := gdkpixbuf.NewPixbufFromData(iconPNG)
	if pb.GetWidth() == -1 {
		fmt.Fprintf(os.Stderr, "ERROR: invalid embedded pixbuf\n")
		os.Exit(1)
	}

	dlg := gtk.NewDialog()
	dlg.SetPosition(gtk.WIN_POS_CENTER)
	dlg.SetTitle("Shutdown Linux")
	dlg.SetModal(true)
	dlg.SetResizable(false)
	dlg.SetIcon(pb)
	dlg.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, nil)

	vbox := dlg.GetVBox()

	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	framebox2 := gtk.NewVBox(false, 1)

	vpaned.Pack2(framebox2, false, false)

	label := gtk.NewLabel("What do you want the computer to do?")
	framebox2.PackStart(label, false, true, 18)

	hbox := gtk.NewHBox(false, 1)

	hbox.Add(gtk.NewImageFromPixbuf(pb))

	buttonbox := gtk.NewVBox(false, 1)
	var first *gtk.RadioButton
	var foundActive bool
	for _, lo := range options {
		var g *glib.SList
		if first != nil {
			g = first.GetGroup()
		}
		btn := gtk.NewRadioButtonWithLabel(g, lo.Label)
		if len(lo.Command) == 0 {
			btn.SetSensitive(false)
		} else {
			if !foundActive {
				btn.SetActive(true)
				foundActive = true
			}
		}
		buttonbox.Add(btn)
		if first == nil {
			first = btn
		}
	}

	hbox.Add(buttonbox)

	// add to layout
	framebox2.PackStart(hbox, false, false, 0)

	hbox = gtk.NewHBox(false, 1)

	okButton := gtk.NewButtonWithLabel("OK")
	okButton.Clicked(func() {
		var activeButtonLabel string

		// scan through all buttons to find the active one
		first.GetGroup().ForEach(func(data unsafe.Pointer, user_data interface{}) {
			var rb gtk.RadioButton
			rb.GWidget = gtk.AS_GWIDGET(data)
			if rb.GetActive() {
				activeButtonLabel = rb.GetLabel()
			}
		})

		for _, lo := range options {
			if lo.Label == activeButtonLabel {
				activate(lo.Command)
			}
		}

		dlg.Destroy()
	})
	okButton.SetCanDefault(true)
	dlg.SetDefault(&okButton.Widget)
	hbox.Add(okButton)

	cancelButton := gtk.NewButtonWithLabel("Cancel")
	cancelButton.Clicked(func() {
		dlg.Destroy()
	})
	hbox.Add(cancelButton)

	// add to layout
	framebox2.PackStart(hbox, false, false, 20)

	// initialize & display
	dlg.SetSizeRequest(350, 200)
	dlg.ShowAll()
	gtk.Main()
}

func activate(s string) {
	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not run %q: %v\n", s, err)
		os.Exit(1)
	}
}
