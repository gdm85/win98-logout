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

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: make-inline-pixbuf resourceName input > output\n")
		os.Exit(1)
	}

	image := gtk.NewImageFromFile(os.Args[2])
	pb := image.GetPixbuf()
	if pb.GetWidth() == -1 {
		fmt.Fprintf(os.Stderr, "ERROR: invalid pixbuf image\n")
		os.Exit(2)
	}

	var pbd gdkpixbuf.PixbufData
	pbd.Data = pb.GetPixelsWithLength()
	pbd.Width, pbd.Height, pbd.RowStride, pbd.HasAlpha = pb.GetWidth(), pb.GetHeight(), pb.GetRowstride(), pb.GetHasAlpha()
	pbd.Colorspace, pbd.BitsPerSample = pb.GetColorspace(), pb.GetBitsPerSample()

	fmt.Printf("package main \n\nimport \"github.com/mattn/go-gtk/gdkpixbuf\"\n\nvar (\n")
	fmt.Printf("\t%s = %#v\n", os.Args[1], pbd)

	fmt.Println(")\n")
}
