package govirsh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func generateUnicastMacAddress() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0]&0xFC, buf[1], buf[2], buf[3], buf[4], buf[5]), nil
}

func generateIp(start, end string, exitList []string) (string, error) {
	dot := "."
	startNum, _ := strconv.Atoi(strings.Split(start, dot)[3])
	endNum, _ := strconv.Atoi(strings.Split(end, dot)[3])
	var exit []int
	for _, e := range exitList {
		exitE, _ := strconv.Atoi(strings.Split(e, dot)[3])
		exit = append(exit, exitE)
	}
	sort.Ints(exit)
	for i := startNum + 1; i < endNum; i++ {
		index := sort.SearchInts(exit, i)
		if index != len(exit) {
			continue
		}
		ip := strings.Join(append(strings.Split(start, dot)[:3], strconv.Itoa(i)), dot)
		return ip, nil
	}
	return "", errors.New("unable to assign ip address")
}

func generateName(exitList []string) (name string, err error) {
	for i := 1; i <= 10000; i++ {
		name = "govirt-" + strconv.Itoa(i)
		for _, v := range exitList {
			if name == v {
				goto end
			}
		}
		return name, nil
	end:
	}
	return "", errors.New("ip not enough")
}

func createQcow(qcowPath string, size int) error {
	cmd := fmt.Sprintf("qemu-img create -f qcow2  -o preallocation=metadata %s %vG", qcowPath, size)
	_, err := outputStringCmd(cmd)
	return err
}

func createQcowWith50G(qcowPath string) error {
	return createQcow(qcowPath, 50)
}

func outputStringCmd(cmd string) (string, error) {
	c := exec.Command("/bin/bash", "-c", cmd)
	out, err := c.Output()
	if err != nil {
		log.Fatalln("OutputStringCmd:", cmd, "; error is:", err, "output is:", string(out))
	}
	return string(out), err
}

func getListMapKey(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
