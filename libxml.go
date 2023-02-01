package govirsh

import (
	"errors"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"io/ioutil"
	"os"
	"strconv"
)

type LibvirtXml struct {
	xmlString string
	domcfg    *libvirtxml.Domain
}

func (L *LibvirtXml) setNameAndUuid(newName, uuid string) {
	if newName != "" {
		L.domcfg.Name = newName
	}
	if uuid != "" {
		L.domcfg.UUID = uuid
	}
}

func (L LibvirtXml) getNameAndUuid() (string, string) {
	return L.domcfg.Name, L.domcfg.UUID
}

// Unit is MB
func (L *LibvirtXml) setMemory(memory uint) {
	if memory == 0 {
		return
	}
	L.domcfg.Memory.Unit = "G"
	L.domcfg.Memory.Value = memory
	L.domcfg.CurrentMemory.Value = memory
	L.domcfg.CurrentMemory.Unit = "G"
}

// Unit is GB
func (L LibvirtXml) getMemory() string {
	return strconv.Itoa(int(L.domcfg.CurrentMemory.Value)) + L.domcfg.CurrentMemory.Unit
}

func (L *LibvirtXml) setCPU(cores uint) {
	if cores == 0 {
		return
	}
	L.domcfg.VCPU.Value = cores
}

func (L LibvirtXml) getCPU() uint {
	return L.domcfg.VCPU.Value
}

func (L *LibvirtXml) setFirstNicName(name string) {
	if name == "" {
		return
	}
	InterfaceList := L.domcfg.Devices.Interfaces
	Interface := InterfaceList[0]
	if Interface.Source.Bridge != nil {
		Interface.Source.Bridge.Bridge = name
	} else if Interface.Source.Network != nil {
		Interface.Source.Network.Network = name
	} else {
		panic("unknown nic type")
	}
}

func (L LibvirtXml) getFirstNicName() string {
	InterfaceList := L.domcfg.Devices.Interfaces
	Interface := InterfaceList[0]
	if Interface.Source.Bridge != nil {
		return Interface.Source.Bridge.Bridge
	} else if Interface.Source.Network != nil {
		return Interface.Source.Network.Network
	} else {
		return ""
	}
}

func (L *LibvirtXml) setFirstNicMac(address string) error {
	if address == "" {
		return nil
	}
	InterfaceList := L.domcfg.Devices.Interfaces
	if len(InterfaceList) == 0 {
		return errors.New("no interface")
	}
	Interface := InterfaceList[0]
	Interface.MAC.Address = address
	return nil
}

func (L LibvirtXml) getFirstNicMac() string {
	return L.domcfg.Devices.Interfaces[0].MAC.Address
}

func (L *LibvirtXml) setEmulator(path string) error {
	if path == "" {
		return nil
	}
	L.domcfg.Devices.Emulator = path
	return nil
}

func (L LibvirtXml) getEmulator() string {
	return L.domcfg.Devices.Emulator
}

func (L *LibvirtXml) setQcow(qcowPath string) error {
	if qcowPath == "" {
		return nil
	}
	L.domcfg.Devices.Disks[0].Source.File.File = qcowPath
	count := len(L.domcfg.Devices.Disks)
	if count >= 2 {
		for i := 1; i < count; i++ {
			if L.domcfg.Devices.Disks[i].Source.File != nil {
				_qcowPath := qcowPath[:len(qcowPath)-6] + "_empty_" + strconv.Itoa(i) + qcowPath[len(qcowPath)-6:]
				_ = createQcowWith50G(_qcowPath)
				L.domcfg.Devices.Disks[i].Source.File.File = _qcowPath
			}
		}
	}
	return nil
}

func (L LibvirtXml) getQcow() string {
	return L.domcfg.Devices.Disks[0].Source.File.File
}

func (L LibvirtXml) domainToString() string {
	content, _ := L.domcfg.Marshal()
	return content
}

func fileToString(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	content, err := ioutil.ReadAll(file)
	return string(content)
}

func NewLibvirtXmlByString(xml string) LibvirtXml {
	domcfg := &libvirtxml.Domain{}
	err := domcfg.Unmarshal(xml)
	if err != nil {
		return LibvirtXml{}
	}
	return LibvirtXml{
		xmlString: xml,
		domcfg:    domcfg,
	}

}

func NewLibvirtXmlByPath(path string) LibvirtXml {
	xml := fileToString(path)
	return NewLibvirtXmlByString(xml)
}
