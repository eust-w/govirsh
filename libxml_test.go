/*
@Author   : longtao.wu@zstack.io
*/

package govirsh

import (
	"fmt"
	"testing"
)

func TestNewLibvirtXmlByPath(t *testing.T) {
	name := "ztest_name"
	uuid := "ztest_uuid"
	mac := "70:af:e7:2f:3f:8a"
	nicName := "govirt"
	var cpuNum uint = 15
	var memory uint = 15
	emulator := "/test/qemu-kvm"
	qcow := "../fsdf/"
	kk := NewLibvirtXmlByPath("../static/defaultVm.xml")
	kk.setNameAndUuid(name, uuid)
	setedName, setedUuid := kk.getNameAndUuid()
	if !(setedName == name && setedUuid == uuid) {
		panic("setNameAndUuid error")
	}
	kk.setCPU(cpuNum)
	if kk.getCPU() != cpuNum {
		panic("setCpu error")
	}
	err := kk.setEmulator(emulator)
	if err != nil {
		panic(err)
	}
	if kk.getEmulator() != emulator {
		panic("setEmulator error")
	}
	kk.setMemory(memory)
	err = kk.setFirstNicMac(mac)
	if err != nil {
		panic(err)
	}
	if kk.getFirstNicMac() != mac {
		panic("firstNicMac")
	}
	kk.setFirstNicName(nicName)
	if kk.getFirstNicName() != nicName {
		panic("firstNicMac")
	}
	err = kk.setQcow(qcow)
	if err != nil {
		panic(err)
	}
	if kk.getQcow() != qcow {
		panic("qcow2")
	}
	fmt.Println(kk.domainToString())
}
