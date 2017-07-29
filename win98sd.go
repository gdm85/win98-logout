package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
)

//go:generate sh -c "go run embedder/make_inline_pixbuf.go iconPNG icons/shut_down_normal.png > icon.gen.go"

func main() {
	gtk.Init(&os.Args)

	pb := gdkpixbuf.NewPixbufFromData(iconPNG)
	if pb.GetWidth() == -1 {
		panic("invalid embedded pixbuf")
	}

	fmt.Printf("size %vx%v\n", pb.GetWidth(), pb.GetHeight())

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
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	framebox2 := gtk.NewVBox(false, 1)

	vpaned.Pack2(framebox2, false, false)

	label := gtk.NewLabel("What do you want the computer to do?")
	framebox2.PackStart(label, false, true, 18)

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
	standBy := gtk.NewRadioButtonWithLabel(nil, "Stand by")
	buttonbox.Add(standBy)
	shutdown := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Shutdown")
	buttonbox.Add(shutdown)
	restart := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Restart")
	buttonbox.Add(restart)
	vtSwitch := gtk.NewRadioButtonWithLabel(standBy.GetGroup(), "Switch to a VT")
	buttonbox.Add(vtSwitch)

	buttons.Add(buttonbox)
	shutdown.SetActive(true)

	framebox2.PackStart(buttons, false, false, 0)

	buttons = gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	okButton := gtk.NewButtonWithLabel("OK")
	okButton.Clicked(func() {
		var obj gtk.ILabel
		if standBy.GetActive() {
			obj = standBy
		} else if shutdown.GetActive() {
			obj = shutdown
		} else if restart.GetActive() {
			obj = restart
		} else if vtSwitch.GetActive() {
			obj = vtSwitch
		}
		fmt.Println("selected:", obj.GetLabel())
	})
	okButton.SetCanDefault(true)
	window.SetDefault(&okButton.Widget)
	buttons.Add(okButton)

	cancelButton := gtk.NewButtonWithLabel("Cancel")
	cancelButton.Clicked(func() {
		window.Destroy()
	})
	buttons.Add(cancelButton)

	framebox2.PackStart(buttons, false, false, 20)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.SetSizeRequest(350, 200)
	window.ShowAll()
	gtk.Main()
}
