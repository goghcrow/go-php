# go-php: Run PHP in Golang

ph7 wrapper

## started

```sh
git clone https://github.com/symisc/PH7.git
cd PH7
cc -c -W -Wall -O6 -o lib_ph7.so ph7.c -D PH7_ENABLE_MATH_FUNC -lm
# 根据不同平台选择
# ar r libph7_darwin_arm64.a lib_ph7.s
# ar r libph7_darwin_amd64.a lib_ph7.s
# ar r libph7_linux_amd64.a lib_ph7.s
mv libph7_darwin_arm64.a ghp/lib
```

## example

```go
package test

/*
#include "lib/ph7.h"
#include "ghp.h"

// 注册给 php-vm 的 go 函数要 export 并且要写 C 的声明
int SayHello(ph7_context *, int, ph7_value **);
int ReturnArr(ph7_context *, int, ph7_value **);
*/
import "C"

import (
	"fmt"
	"github.com/goghcrow/ghp"
)

// C-API  https://ph7.symisc.net/c_api_func.html
// PHP builtin functions  https://ph7.symisc.net/builtin_func.html

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

// InjectFuns : ffi, go 实现 PHP 函数
var InjectFuns = []C.gph_foreign_func{
	{
		zName: C.CString("SayHello"),
		xProc: C.gph_xProc(C.SayHello),
	},
	{
		zName: C.CString("ReturnArr"),
		xProc: C.gph_xProc(C.ReturnArr),
	},
}

func TestEval() {
	ghp.Eval(`<?php
echo "Hello World!\n";
var_dump($argv); // 接收 go 传进来的参数
var_dump(ReturnArr()); // 调用go注册的函数
print_r(implode(" ", array_map(function($it) { return $it * 2; }, range(1, 5))));
`, []string{"arg1", "arg2"}, InjectFuns...)
}
```