package govirsh

type createVmOption struct {
	Name       string
	MemorySize string
	Xml        string
	Qcow2      string
	CpuNum     string
	Ip         string
}

type CreateVmOption interface {
	apply(*createVmOption)
}

type funcOption struct {
	f func(*createVmOption)
}

func (fdo *funcOption) apply(do *createVmOption) {
	fdo.f(do)
}

func newFuncOption(f func(*createVmOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func defaultCreateVmOptions() createVmOption {
	return createVmOption{
		Name:       "",
		MemorySize: "",
		Xml:        "",
		Qcow2:      "",
		CpuNum:     "",
		Ip:         "",
	}
}

func SetName(name string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.Name = name
	})
}

func SetMemorySize(memorySize string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.MemorySize = memorySize
	})
}

func SetXml(xml string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.Xml = xml
	})
}

func SetQcow2(qcow2 string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.Qcow2 = qcow2
	})
}

func SetCpuNum(cpuNum string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.CpuNum = cpuNum
	})
}

func SetIp(ip string) CreateVmOption {
	return newFuncOption(func(o *createVmOption) {
		o.Ip = ip
	})
}
