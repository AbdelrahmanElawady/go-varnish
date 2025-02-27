package main

/*
#cgo CFLAGS: -I/usr/include/varnish
#include <stdlib.h>
#include <vdef.h>
#include <vrt.h>
*/
import "C"
import (
	"fmt"
	"os"
)

//export vmod_hello
func vmod_hello(ctx *C.struct_vrt_ctx, name C.VCL_STRING) C.VCL_STRING {
	goName := C.GoString(name)
	greeting := fmt.Sprintf("Hello %s from Go!", goName)
	return C.VCL_STRING(C.CString(greeting))
}

//export vmod_add
func vmod_add(ctx *C.struct_vrt_ctx, a C.VCL_INT, b C.VCL_INT) C.VCL_INT {
	return C.VCL_INT(a + b)
}

//export vmod_write
func vmod_write(ctx *C.struct_vrt_ctx, h *C.struct_gethdr_s) {
	value := C.VRT_GetHdr(ctx, C.VCL_HEADER(h))
	if value != nil {
		goValue := C.GoString(value)
		os.WriteFile("/tmp/test.txt", []byte(goValue), 0644)
	}
}

func main() {}
