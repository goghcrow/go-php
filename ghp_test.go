package ghp

import "testing"

func TestEval(t *testing.T) {
	Eval(
		`<?php
var_dump($argv); // 接收 go 传进来的参数

var_dump(ReturnArr()); // 调用 go 函数

var_dump("Hello World!");

print_r(implode(" ", array_map(function($it) { return $it * 2; }, range(1, 5))));
`, []string{"arg1", "arg2"}, InjectFuns...)
}
