///*
//@Author   : longtao.wu@zstack.io
//*/
package govirsh

//
//import (
//	"fmt"
//	"testing"
//)
//
//func TestVirtNet(t *testing.T) {
//	ll := NewLibNetList()
//	nameList, err := ll.GetNetListName()
//	if err != nil {
//		panic(err)
//	}
//	virtNet1, err := ll.GetNetListByName(nameList[1])
//	if err != nil {
//		panic(err)
//	}
//	vt := NewVirtNet(virtNet1)
//	list, err := vt.GetDhcpHostListByIpKey()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(list)
//	_ = vt.DelDhcpHostByIp("192.168.100.141")
//	err = vt.AddDhcpHost("70:af:e7:2f:3f:8a", "name141", "192.168.100.141")
//	if err != nil {
//		panic(err)
//	}
//	list, err = vt.GetDhcpHostListByMacKey()
//	if err != nil {
//		panic(err)
//	}
//	if !StringSliceEqual(list["70:af:e7:2f:3f:8a"], []string{"name141", "192.168.100.141"}) {
//		panic("add host failed")
//	}
//	err = vt.DelDhcpHostByIp("192.168.100.141")
//	if err != nil {
//		panic(err)
//	}
//	list, err = vt.GetDhcpHostListByNameKey()
//	if err != nil {
//		panic(err)
//	}
//	if list["name141"] != nil {
//		panic("del failed")
//	}
//}
//
//func StringSliceEqual(a, b []string) bool {
//	if len(a) != len(b) {
//		return false
//	}
//
//	if (a == nil) != (b == nil) {
//		return false
//	}
//
//	for i, v := range a {
//		if v != b[i] {
//			return false
//		}
//	}
//	return true
//}
