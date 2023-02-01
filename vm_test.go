///*
//@Author   : longtao.wu@zstack.io
//*/
package govirsh

//
//import (
//	"fmt"
//	"testing"
//	"time"
//)
//
//func TestNewVm(t *testing.T) {
//	vm, err := NewVm("name", "", "", "/home/longtao/下载/image-52450c73d20c8c678a0f39d4ef98dea34dec6568.qcow2", "", "/home/longtao/temp/vm.xml", "", 2, 2)
//	if err != nil {
//		panic(err)
//	}
//	s, si, err := vm.Domain.GetState()
//	if err != nil {
//		panic(err)
//	}
//	time.Sleep(20)
//	fmt.Println(s, si)
//	vm.DestroyVm()
//	time.Sleep(20)
//	s, si, err = vm.Domain.GetState()
//	if err == nil {
//		panic("destroy vm false")
//	}
//	fmt.Println(s, si)
//}
