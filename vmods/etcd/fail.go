package main

/*
#cgo CFLAGS: -I/usr/include/varnish
#include <vdef.h>
#include <vrt.h>
#include <stdlib.h>

void fail(VRT_CTX, const char *fmt) {
	VRT_fail(ctx, "%s", fmt);
}
*/
import "C"
import "unsafe"

func fail(ctx *C.struct_vrt_ctx, msg string) {
	cMsg := C.CString(msg)
	C.fail(ctx, cMsg)
	C.free(unsafe.Pointer(cMsg))
}
