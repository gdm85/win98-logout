package main

import (
	"fmt"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path/filepath"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.NewDialog()
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Shutdown Linux")
	window.SetModal(true)
	window.SetResizable(false)
	window.SetIconName("gtk-dialog-info")
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
	framebox2.PackStart(label, false, true, 16)

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons := gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkImage
	//--------------------------------------------------------
	dir, _ := filepath.Split(os.Args[0])
	imagefile := filepath.Join(dir, "icons/shut_down_normal.png")

	image := gtk.NewImageFromFile(imagefile)
	buttons.Add(image)

	//--------------------------------------------------------
	// GtkRadioButton
	//--------------------------------------------------------
	buttonbox := gtk.NewVBox(false, 1)
	radiofirst := gtk.NewRadioButtonWithLabel(nil, "Stand by")
	buttonbox.Add(radiofirst)
	sdButton := gtk.NewRadioButtonWithLabel(radiofirst.GetGroup(), "Shutdown")
	buttonbox.Add(sdButton)
	buttonbox.Add(gtk.NewRadioButtonWithLabel(radiofirst.GetGroup(), "Restart"))
	buttonbox.Add(gtk.NewRadioButtonWithLabel(radiofirst.GetGroup(), "Switch to a VT"))
	buttons.Add(buttonbox)
	sdButton.SetActive(true)

	framebox2.PackStart(buttons, false, false, 0)

	buttons = gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	okButton := gtk.NewButtonWithLabel("OK")
	okButton.Clicked(func() {
		fmt.Println("OK clicked:", okButton.GetLabel())
	})
	buttons.Add(okButton)

	cancelButton := gtk.NewButtonWithLabel("Cancel")
	cancelButton.Clicked(func() {
		fmt.Println("cancel clicked:", cancelButton.GetLabel())
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
