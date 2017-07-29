package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//go:generate sh -c "go run embedder/make_inline_pixbuf.go iconPNG icons/shutdown.png > icon.gen.go"

const (
	standbyAction  = ``
	shutdownAction = `sudo shutdown -h now`
	restartAction  = `sudo reboot`
	switchVTAction = `xdotool key ctrl+alt+F1`
)

func main() {
	gtk.Init(&os.Args)

	pb := gdkpixbuf.NewPixbufFromData(iconPNG)
	if pb.GetWidth() == -1 {
		panic("invalid embedded pixbuf")
	}

	window := gtk.NewDialog()
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Shutdown Linux")
	window.SetModal(true)
	window.SetResizable(false)
	window.SetIcon(pb)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, nil)

	vbox := window.GetVBox()

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons := gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkImage
	//--------------------------------------------------------
	buttons.Add(gtk.NewImageFromPixbuf(pb))

	//--------------------------------------------------------
	// GtkRadioButton
	//--------------------------------------------------------
	buttonbox := gtk.NewVBox(false, 1)
	label := gtk.NewLabel("What do you want the computer to do?")
	buttonbox.PackStart(label, false, true, 18)
	standBy := gtk.NewRadioButtonWithLabel(nil, "Stand by")
	if standbyAction == `` {
		standBy.SetSensitive(false)
	}
	buttonbox.Add(standBy)
	shutdown := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Shutdown")
	buttonbox.Add(shutdown)
	restart := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Restart")
	buttonbox.Add(restart)
	vtSwitch := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Switch to a VT")
	buttonbox.Add(vtSwitch)

	buttons.Add(buttonbox)
	shutdown.SetActive(true)

	vbox.PackStart(buttons, false, false, 0)

	buttons = gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	okButton := gtk.NewButtonWithLabel("OK")
	okButton.Clicked(func() {
		if standBy.GetActive() {
			activate(standbyAction)
		} else if shutdown.GetActive() {
			activate(shutdownAction)
		} else if restart.GetActive() {
			activate(restartAction)
		} else if vtSwitch.GetActive() {
			activate(switchVTAction)
		}
		window.Destroy()
	})
	okButton.SetCanDefault(true)
	window.SetDefault(&okButton.Widget)
	buttons.Add(okButton)

	cancelButton := gtk.NewButtonWithLabel("Cancel")
	cancelButton.Clicked(func() {
		window.Destroy()
	})
	buttons.Add(cancelButton)

	vbox.PackStart(buttons, false, false, 20)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.SetSizeRequest(350, 200)
	window.ShowAll()
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
