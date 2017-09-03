package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/patterns/cefcapi/cefgo"
)

func main() {
	// TODO keep comments
	fmt.Printf("\nProcess args: ")
	if len(os.Args) == 1 {
		fmt.Printf("none (Main process)")
	} else {
		for _, v := range os.Args[1:] {
			if len(v) > 128 {
				fmt.Printf("... ")
			} else {
				fmt.Printf("%s ", v)
			}
		}
	}
	fmt.Println("\n")

	wb := cef.NewEmbed(os.Args)
	wb.InitializeApp()

	code := wb.ExecuteProcess(os.Args)
	if code >= 0 {
		os.Exit(code)
	}

	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exdir := filepath.Dir(ex)
	wb.Initialize(exdir)

////	wb.InitializeGtk()
	gtk.Init(nil)
	x := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	x.SetPosition(gtk.WIN_POS_CENTER)
	x.SetTitle("cef2go + cefcapi test")
	x.Connect("destroy", func(ctx *glib.CallbackContext) {
			println("end window", ctx.Data().(string))
			cef.QuitMessageLoop()
			gtk.MainQuit()
	}, "testing")
	x.ShowAll()
	xid := gdk.WindowFromUnsafe(
			unsafe.Pointer(x)).GetNativeWindowID()

	wb.InitializeClient()

	wb.CreateBrowser(xid, "https://www.google.com/ncr")

	cef.RunMessageLoop()

	cef.Shutdown()
}

