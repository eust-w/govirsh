package govirsh

import (
	"errors"
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
)

type LibNetList struct {
	netList []libvirt.Network
}

func (l *LibNetList) GetNetListByName(name string) (libvirt.Network, error) {
	for _, net := range l.netList {
		currentName, err := net.GetName()
		if currentName == name {
			if err != nil {
				return net, err
			}
			return net, nil
		}
	}
	return libvirt.Network{}, errors.New(fmt.Sprintf("net[%s] not found", name))
}

func AddLibvirtNet(xmlContent string) (libvirt.Network, error) {
	conn, _ := libvirt.NewConnect("qemu:///system")
	conn.NetworkDefineXML(xmlContent)
	xml, err := conn.NetworkCreateXML(xmlContent)
	if err != nil {
		return libvirt.Network{}, err
	}
	return *xml, nil
}

func (l LibNetList) GetNetListName() ([]string, error) {
	return getNetListName(l.netList)
}

func getNetListName(l []libvirt.Network) ([]string, error) {
	var nameList = make([]string, 0, len(l))
	for _, net := range l {
		currentName, err := net.GetName()
		if err != nil {
			return nameList, err
		}
		nameList = append(nameList, currentName)
	}
	return nameList, nil
}

func NewLibNetList() LibNetList {
	conn, _ := libvirt.NewConnect("qemu:///system")
	activeNetwork := make([]libvirt.Network, 0, 0)
	networks, err := conn.ListAllNetworks(0)
	for _, v := range networks {
		b, _ := v.IsActive()
		if b {
			activeNetwork = append(activeNetwork, v)
		}
	}
	if err != nil {
		panic(err)
	}
	return LibNetList{netList: activeNetwork}
}
