package govirsh

import (
	"errors"
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type VirtNet struct {
	n libvirt.Network
}

// AddDhcpHost 向dhcp中增加host
func (v VirtNet) AddDhcpHost(mac, name, ip string) error {
	xml := libvirtxml.NetworkDHCPHost{MAC: mac, Name: name, IP: ip}
	xmlString, err := xml.Marshal()
	if err != nil {
		return err
	}
	err = v.n.Update(libvirt.NETWORK_UPDATE_COMMAND_ADD_LAST, libvirt.NETWORK_SECTION_IP_DHCP_HOST, 0, xmlString, libvirt.NETWORK_UPDATE_AFFECT_LIVE|libvirt.NETWORK_UPDATE_AFFECT_CONFIG)
	return err
}

func (v VirtNet) AddDhcpHostWithAutoIp(mac, name string) (ip string, err error) {
	start, end, err := v.GetDhcpInfo()
	if err != nil {
		return "", err
	}
	existList, _ := v.GetDhcpHostIpList()
	ip, err = generateIp(start, end, existList)
	if err != nil {
		return "", err
	}
	err = v.AddDhcpHost(mac, name, ip)
	if err != nil {
		return "", err
	}
	return
}

func (v VirtNet) AddDhcpHostWithAutoIpAndName(mac string) (ip, name string, err error) {
	nameHostMap, err := v.GetDhcpHostListByNameKey()
	if err != nil {
		return "", "", err
	}
	name, err = generateName(getListMapKey(nameHostMap))
	if err != nil {
		return "", "", err
	}
	ip, err = v.AddDhcpHostWithAutoIp(mac, name)
	return
}

// DelDhcpHostByName 向dhcp删除host
func (v VirtNet) DelDhcpHostByName(name string) error {
	nameKeyMap, err := v.GetDhcpHostListByNameKey()
	hostInfoList := nameKeyMap[name]
	if err != nil {
		return err
	}
	if len(hostInfoList) != 2 {
		errors.New("host info list error")
	}
	mac, ip := hostInfoList[0], hostInfoList[1]
	xml := libvirtxml.NetworkDHCPHost{MAC: mac, Name: name, IP: ip}
	xmlString, err := xml.Marshal()
	err = v.n.Update(libvirt.NETWORK_UPDATE_COMMAND_DELETE, libvirt.NETWORK_SECTION_IP_DHCP_HOST, 0, xmlString, libvirt.NETWORK_UPDATE_AFFECT_LIVE|libvirt.NETWORK_UPDATE_AFFECT_CONFIG)
	if err != nil {
		return err
	}
	return nil
}

// DelDhcpHostByIp 向dhcp删除host
func (v VirtNet) DelDhcpHostByIp(ip string) error {
	ipKeyMap, err := v.GetDhcpHostListByIpKey()
	hostInfoList := ipKeyMap[ip]
	if err != nil {
		return err
	}
	if len(hostInfoList) != 2 {
		return err
	}
	name, mac := hostInfoList[0], hostInfoList[1]
	xml := libvirtxml.NetworkDHCPHost{MAC: mac, Name: name, IP: ip}
	xmlString, err := xml.Marshal()
	fmt.Println("*****del:", xmlString)
	err = v.n.Update(libvirt.NETWORK_UPDATE_COMMAND_DELETE, libvirt.NETWORK_SECTION_IP_DHCP_HOST, 0, xmlString, libvirt.NETWORK_UPDATE_AFFECT_LIVE|libvirt.NETWORK_UPDATE_AFFECT_CONFIG)
	if err != nil {
		return err
	}
	return nil
}

// DelDhcpHost 向dhcp删除host
func (v VirtNet) DelDhcpHost(mac, name, ip string) error {
	xml := libvirtxml.NetworkDHCPHost{MAC: mac, Name: name, IP: ip}
	xmlString, err := xml.Marshal()
	err = v.n.Update(libvirt.NETWORK_UPDATE_COMMAND_DELETE, libvirt.NETWORK_SECTION_IP_DHCP_HOST, 0, xmlString, libvirt.NETWORK_UPDATE_AFFECT_LIVE|libvirt.NETWORK_UPDATE_AFFECT_CONFIG)
	if err != nil {
		return err
	}
	return nil
}

// GetDhcpHostListByNameKey 获取第一个DHCP中的host列表，格式为
func (v VirtNet) GetDhcpHostListByNameKey() (k map[string][]string, err error) {
	k = make(map[string][]string)
	desc, err := v.n.GetXMLDesc(0)
	if err != nil {
		return
	}
	var xml libvirtxml.Network
	err = xml.Unmarshal(desc)
	if err != nil {
		return
	}
	for _, v := range xml.IPs[0].DHCP.Hosts {
		k[v.Name] = []string{v.MAC, v.IP}
	}
	return
}

// GetDhcpHostListByIpKey 获取第一个DHCP中的host列表，格式为
func (v VirtNet) GetDhcpHostListByIpKey() (k map[string][]string, err error) {
	k = make(map[string][]string)
	desc, err := v.n.GetXMLDesc(0)
	if err != nil {
		return
	}
	var xml libvirtxml.Network
	err = xml.Unmarshal(desc)
	if err != nil {
		return
	}
	for _, v := range xml.IPs[0].DHCP.Hosts {
		k[v.IP] = []string{v.Name, v.MAC}
	}
	return
}

// GetDhcpHostListByMacKey 获取第一个DHCP中的host列表，格式为
func (v VirtNet) GetDhcpHostListByMacKey() (k map[string][]string, err error) {
	k = make(map[string][]string)
	desc, err := v.n.GetXMLDesc(0)
	if err != nil {
		return
	}
	var xml libvirtxml.Network
	err = xml.Unmarshal(desc)
	if err != nil {
		return
	}
	for _, v := range xml.IPs[0].DHCP.Hosts {
		k[v.MAC] = []string{v.Name, v.IP}
	}
	return
}

// 获取第一个dhcp绑定的host的ip列表
func (v VirtNet) GetDhcpHostIpList() (ip []string, err error) {
	desc, err := v.n.GetXMLDesc(0)
	var xml libvirtxml.Network
	err = xml.Unmarshal(desc)
	for _, v := range xml.IPs[0].DHCP.Hosts {
		ip = append(ip, v.IP)
	}
	return
}

// 获取第一个dhcp的信息
func (v VirtNet) GetDhcpInfo() (startIp, endIp string, err error) {

	desc, err := v.n.GetXMLDesc(0)
	var xml libvirtxml.Network
	err = xml.Unmarshal(desc)
	startIp = xml.IPs[0].DHCP.Ranges[0].Start
	endIp = xml.IPs[0].DHCP.Ranges[0].End
	return
}

func NewVirtNet(l libvirt.Network) VirtNet {
	return VirtNet{
		n: l,
	}
}

//todo: 给libvirt net 加锁，不然会导致并发创建虚拟机 指定的network pointer 失效，目前的做法是给创建vm加锁
