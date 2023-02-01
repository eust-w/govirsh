///*
//@Author   : longtao.wu@zstack.io
//*/
package govirsh

//
//import (
//	"fmt"
//	"testing"
//	"ztest/static"
//)
//
//func TestNewLibNetList(t *testing.T) {
//	netList := NewLibNetList()
//	list, err := netList.GetNetListName()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(list)
//	if len(list) != 0 {
//		net, err := netList.GetNetListByName(list[0])
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(net)
//	}
//}
//
//func TestAddLibvirtNet(t *testing.T) {
//	data, err := static.DefaultNet.ReadFile("defaultNet.xml")
//	if err != nil {
//		panic(err)
//	}
//	_, err = AddLibvirtNet(string(data))
//	if err != nil {
//		panic(err)
//	}
//}
