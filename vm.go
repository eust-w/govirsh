package govirsh

import (
	"errors"
	"fmt"
	"github.com/libvirt/libvirt-go"
	"os"
	"strings"
	"sync"
	"time"
)

type Vm struct {
	name       string
	uuid       string
	cpu        uint
	memory     uint
	ip         string
	mac        string
	baseQcow   string
	emulator   string
	baseXml    string
	qcow       string
	libvirtNet string
	xml        string
	govirtNet  VirtNet
	Domain     *libvirt.Domain
}

func (v *Vm) buildXml() {
	libXml := NewLibvirtXmlByPath(v.baseXml)
	libXml.setQcow(v.qcow)
	libXml.setNameAndUuid(v.name, v.uuid)
	govirtName, _ := v.govirtNet.n.GetName()
	libXml.setFirstNicName(govirtName)
	libXml.setFirstNicMac(v.mac)
	libXml.setEmulator(v.emulator)
	libXml.setMemory(v.memory)
	libXml.setCPU(v.cpu)
	v.xml = libXml.domainToString()
}

func (v *Vm) buildGovirtNet() error {
	netList := NewLibNetList()
	if v.libvirtNet == "" {
		v.libvirtNet = "govirt"
	}
	govirt, err := netList.GetNetListByName(v.libvirtNet)
	if err != nil {
		data := "<network>\n    <name>govirt</name>\n    <uuid>98361b46-1581-acb7-1643-85a412626e70</uuid>\n    <forward mode='open'/>\n    <bridge name='govirt' stp='off' delay='0'/>\n    <mac address='52:54:00:b9:e7:1e'/>/home/longtao/workspace/projects/govirsh\n    <ip address='192.168.100.1' netmask='255.255.255.0'>\n        <dhcp>\n            <range start='192.168.100.128' end='192.168.100.254'/>\n        </dhcp>\n    </ip>\n</network>"
		govirt, err = AddLibvirtNet(data)
		_ = openIpForward()
		_ = setIptables()
	}
	v.govirtNet = NewVirtNet(govirt)
	return nil
}

func openIpForward() error {
	cmdArgs := fmt.Sprintf("sudo echo %s >/proc/sys/net/ipv4/ip_forward", "1")
	_, err := outputStringCmd(cmdArgs)
	if err != nil {
		return err
	}
	_, err = outputStringCmd("sudo sysctl -p")
	return err
}
func setIptables() error {
	cmd := "sudo  iptables -t nat -C POSTROUTING -s 192.168.100.0/24 -j MASQUERADE||iptables -t nat -A POSTROUTING -s 192.168.100.0/24 -j MASQUERADE -w 10"
	_, err := outputStringCmd(cmd)
	return err
}

func (v *Vm) SetIp(ip string) {
	v.ip = ip
}

func (v *Vm) buildIp() error {
	if v.mac == "" {
		v.mac, _ = generateUnicastMacAddress()
	}
	ip, _, err := v.govirtNet.AddDhcpHostWithAutoIpAndName(v.mac)
	if err != nil {
		return err
	}
	v.ip = ip
	return nil
}

func (v *Vm) createVm() error {
	var err error
	conn, _ := libvirt.NewConnect("qemu:///system")
	v.Domain, err = conn.DomainCreateXML(v.xml, 0)
	return err
}

func (v *Vm) RemoveQcow() error {
	var err error
	err = os.Remove(v.qcow)
	return err
}

func (v *Vm) DestroyVm() error {
	var err error
	err = v.Domain.Destroy()
	if err != nil {
		return err
	}
	// todo: 临时的解决办法，先去掉dhcp的清理工作
	//err = v.govirtNet.DelDhcpHostByIp(v.ip)
	//if err != nil {
	//	return err
	//}
	//err = v.govirtNet.n.Free()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (v Vm) GetName() string {
	return v.name
}
func (v Vm) GetUuid() string {
	return v.uuid
}
func (v Vm) GetIp() string {
	return v.ip
}

func (v Vm) GetMac() string {
	return v.mac
}

//func rebootVirtNet() {
//	utils.RunCmd("virsh net-destroy govirt")
//	time.Sleep(time.Second * 2)
//	utils.RunCmd(" virsh net-start govirt")
//
//}

func (v *Vm) CheckIp() (string, error) {
	mac := v.GetMac()
	var ip string
	var err error
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Second)
		ip, err = checkIP(mac)
		if err == nil {
			v.SetIp(ip)
			return ip, nil
		}
	}
	return "", errors.New("check ip failed")
}

func checkIP(mac string) (string, error) {
	out, err := outputStringCmd("virsh net-dhcp-leases govirt |grep " + mac + "|awk '{print $5}'")
	if err != nil || out == "" {
		return "", errors.New("check ip failed")
	}
	p := strings.Split(out, "/")
	if len(p) <= 0 {
		return "", errors.New("check ip failed")
	}
	return p[0], nil
}

//加锁粒度较大，后期修改
var newVmLocker sync.Mutex

func NewVm(name, uuid, mac, baseQcow, emulato, baseXml, qcow string, cpu, memory uint) (Vm, error) {
	newVmLocker.Lock()
	defer newVmLocker.Unlock()
	var err error
	defer func() {
		if err != nil {
		}
	}()
	vm := Vm{name: name, uuid: uuid, mac: mac, baseQcow: baseQcow, emulator: emulato, baseXml: baseXml, qcow: qcow, cpu: cpu, memory: memory}
	err = vm.buildGovirtNet()
	if err != nil {
		return Vm{}, err
	}
	err = vm.buildIp()
	if err != nil {
		return Vm{}, err
	}
	vm.buildXml()
	err = vm.createVm()
	//可以在创建出vm后通过polling `virsh net-dhcp-leases govirt|grep mac地址`来获取dhcp动态分配的ip地址，然后写入vm
	if err != nil {
		return Vm{}, err
	}
	return vm, err
}
