
// Wrap the CEF C API behind Go interface

package cef

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>


#include "capi/cef_base.h"
#include "capi/cef_app.h"
#include "capi/cef_client.h"
#include "capi/cef_life_span_handler.h"
#cgo CFLAGS: -I. -I..
#cgo LDFLAGS: -L/go/src/github.com/patterns/cefcapi/Release -lX11 -lcef
#cgo pkg-config: --libs --cflags gtk+-2.0

// Global expected by capi/cef_client.h
cef_life_span_handler_t g_life_span_handler = {};
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Embed struct {
	app *C.cef_app_t
	args *C.cef_main_args_t
	client *C.cef_client_t
	winf *C.cef_window_info_t
	sbinf unsafe.Pointer
}

const (
    LOGSEVERITY_DEFAULT = C.LOGSEVERITY_DEFAULT
    LOGSEVERITY_VERBOSE = C.LOGSEVERITY_VERBOSE
    LOGSEVERITY_INFO = C.LOGSEVERITY_INFO
    LOGSEVERITY_WARNING = C.LOGSEVERITY_WARNING
    LOGSEVERITY_ERROR = C.LOGSEVERITY_ERROR
//    LOGSEVERITY_ERROR_REPORT = C.LOGSEVERITY_ERROR_REPORT
    LOGSEVERITY_DISABLE = C.LOGSEVERITY_DISABLE
)

func NewEmbed(args []string) Embed {
	b := Embed{}
	return b
}

// initialize app handler
func (b *Embed) InitializeApp() {
	var ah C.cef_app_t
	pah := (*C.cef_app_t)(
            C.calloc(1, C.size_t(unsafe.Sizeof(ah))))
	b.app = pah
	C.initialize_cef_app(b.app)
}

// set main args
func (b *Embed) setMainArgs(args []string) {
	var ma C.cef_main_args_t
	pma := (*C.cef_main_args_t)(
            C.calloc(1, C.size_t(unsafe.Sizeof(ma))))

	argv := make([]*C.char, len(args))

	for k, v := range args {
		argv[C.int(k)] = C.CString(v)
////		defer C.free(unsafe.Pointer(argv[C.int(k)]))
	}

    pma.argc = C.int(len(args))
	pma.argv = &argv[0]
	b.args = pma
}

// sub-processes
func (b *Embed) ExecuteProcess(args []string) int {
	fmt.Println("cef_execute_process, args=", args)
	b.setMainArgs(args)
	var code C.int = C.cef_execute_process(b.args, b.app, b.sbinf)
	return int(code)
}

// initialize CEF
func (b *Embed) Initialize(exdir string) int {
	fmt.Println("cef_initialize")
	if b.args == nil {
		fmt.Println("ERROR:  ExecuteProcess call got skipped!")
		return 0
	}

	var st C.cef_settings_t
	pst := (*C.cef_settings_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(st))))

	pst.size = C.size_t(unsafe.Sizeof(st))

	cpath := exdir + "/webcache"
	fmt.Println("CachePath=", cpath)
	var cachePath *C.char = C.CString(cpath)
	defer C.free(unsafe.Pointer(cachePath))
	C.cef_string_from_utf8(cachePath, C.strlen(cachePath),
		&pst.cache_path)

	pst.log_severity =
		(C.cef_log_severity_t)(C.int(LOGSEVERITY_DEFAULT))

	lpath := exdir + "/debug.log"
	fmt.Println("LogPath=", lpath)
	var logfile *C.char = C.CString(lpath)
	defer C.free(unsafe.Pointer(logfile))
	C.cef_string_from_utf8(logfile, C.strlen(logfile),
		&pst.log_file)

	rpath := exdir
	fmt.Println("ResourcesDirPath=", rpath)
	var resPath *C.char = C.CString(rpath)
	defer C.free(unsafe.Pointer(resPath))
	C.cef_string_from_utf8(resPath, C.strlen(resPath),
		&pst.resources_dir_path)

	gpath := exdir + "/locales"
	fmt.Println("LocalesDirPath=", gpath)
	var locPath *C.char = C.CString(gpath)
	defer C.free(unsafe.Pointer(locPath))
	C.cef_string_from_utf8(locPath, C.strlen(locPath),
		&pst.locales_dir_path)

	pst.no_sandbox = C.int(1)

	ret := C.cef_initialize(b.args, pst, b.app, b.sbinf)
	return int(ret)
}
/*
func (b *Embed) InitializeGtk() {
	C.initialize_gtk()

	cs := C.CString("cefcapi cgo example")
	defer C.free(unsafe.Pointer(cs))
	gw := C.create_gtk_window(cs, C.int(800), C.int(600))
	wi := C.cef_window_info_t{}
	xid := C.gdk_x11_drawable_get_xid(C.gtk_widget_get_window(gw))
	wi.parent_window = C.ulong(xid)
	b.winf = &wi

	//TODO figure out X11
////	C.XSetErrorHandler(C.x11_error_handler)
////	C.XSetIOErrorHandler(C.x11_io_error_handler)
}

*/

// initialize client handler
func (b *Embed) InitializeClient() {
	var cl C.cef_client_t
	pcl := (*C.cef_client_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(cl))))
	b.client = pcl
	C.initialize_cef_client(b.client)
}

func (b *Embed) InitializeLifespanHandler() {
	pls := (*C.cef_life_span_handler_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(C.g_life_span_handler))))
	C.initialize_cef_life_span_handler(pls)
}

func (b *Embed) CreateBrowser(xid int32, url string) {
	fmt.Println("cef_browser_host_create_browser")

	var wi C.cef_window_info_t
	pwi := (*C.cef_window_info_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(wi))))
	FillWindowInfo(pwi, xid)

	var ur C.cef_string_t
	pur := (*C.cef_string_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(ur))))
	var chur *C.char = C.CString(url)
	defer C.free(unsafe.Pointer(chur))
	C.cef_string_from_utf8(chur, C.strlen(chur), pur)

	var st C.cef_browser_settings_t
	pst := (*C.cef_browser_settings_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(st))))
	pst.size = C.size_t(unsafe.Sizeof(st))

	C.cef_browser_host_create_browser(pwi, b.client, pur, pst, nil)
}

func RunMessageLoop() {
	fmt.Println("cef_run_message_loop")
	C.cef_run_message_loop()
}

func QuitMessageLoop() {
	fmt.Println("cef_quit_message_loop")
	C.cef_quit_message_loop()
}

func Shutdown() {
	fmt.Println("cef_shutdown")
	C.cef_shutdown()
}

// Linux platform
func FillWindowInfo(windowInfo *C.cef_window_info_t, xid int32) {
    fmt.Println("FillWindowInfo")

    windowInfo.parent_window = C.ulong(xid)
}
