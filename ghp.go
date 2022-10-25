package ghp

import (
	"fmt"
	"os"
	"unsafe"
)

/*
#cgo darwin,arm64 LDFLAGS: -L./lib -lph7_darwin_arm64
#cgo darwin,amd64 LDFLAGS: -L./lib -lph7_darwin_amd64
#cgo linux,amd64 LDFLAGS: -L./lib -lph7_linux_amd64
#include <stdio.h>
#include <stdlib.h>
#include "lib/ph7.h"
#include "ghp.h"
*/
import "C"

//export gph_output_consumer
//goland:noinspection GoSnakeCaseUsage
func gph_output_consumer(pOutput unsafe.Pointer, nOutputLen C.ulonglong, pUserData unsafe.Pointer) C.int {
	s := C.GoStringN((*C.char)(pOutput), C.int(nOutputLen))
	_, _ = fmt.Fprint(os.Stdout, s)
	return C.PH7_OK
}

//export gph_error_consumer
//goland:noinspection GoSnakeCaseUsage
func gph_error_consumer(msg *C.char) {
	_, _ = fmt.Fprintf(os.Stderr, C.GoString(msg))
}

func Eval(src string, args []string, fs ...C.gph_foreign_func) int {
	pgrm := C.CString(src)
	defer C.free(unsafe.Pointer(pgrm))

	cArgs := make([]*C.char, len(args))
	for i, arg := range args {
		cArgs[i] = C.CString(arg)
	}
	defer func() {
		for _, arg := range cArgs {
			C.free(unsafe.Pointer(arg))
		}
	}()

	return int(C.gph_eval(pgrm, &cArgs[0], C.int(len(args)), &fs[0], C.int(len(fs))))
}
