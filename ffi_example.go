package ghp

/*
#include "lib/ph7.h"
#include "ghp.h"

// 注意注册给 php vm go 函数要 export 并且要写 C 的声明
int SayHello(ph7_context *, int, ph7_value **);
int ReturnArr(ph7_context *, int, ph7_value **);
*/
import "C"
import (
	"fmt"
)

//export SayHello
func SayHello(pCtx *C.ph7_context, argc C.int, argv **C.ph7_value) C.int {
	fmt.Printf("HELLO\n")
	return C.int(0)
}

//export ReturnArr
func ReturnArr(pCtx *C.ph7_context, argc C.int, argv **C.ph7_value) C.int {
	// 这里不用处理 array 的内存释放, add_elem 内部会 copy,
	// 函数返回, array 之类也会释放
	pArray := C.ph7_context_new_array(pCtx)
	pValue := C.ph7_context_new_scalar(pCtx)

	C.ph7_value_int(pValue, C.int(1))
	C.ph7_array_add_elem(pArray, nil, pValue)

	C.ph7_value_int(pValue, C.int(2))
	C.ph7_array_add_elem(pArray, nil, pValue)

	C.ph7_value_int(pValue, C.int(3))
	C.ph7_array_add_elem(pArray, nil, pValue)

	// 返回值
	C.ph7_result_value(pCtx, pArray)

	return C.PH7_OK
}

var InjectFuns = []C.gph_foreign_func{
	{
		zName: C.CString("SayHello"), // C.free(unsafe.Pointer(...))
		xProc: C.gph_xProc(C.SayHello),
	},
	{
		zName: C.CString("ReturnArr"), // C.free(unsafe.Pointer(...))
		xProc: C.gph_xProc(C.ReturnArr),
	},
}
