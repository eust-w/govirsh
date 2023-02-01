package govirsh

import (
	"fmt"
	"testing"
)

func TestGenerateIp(t *testing.T) {
	start := "172.20.1.35"
	end := "172.20.1.90"
	exit := []string{"172.20.1.36"}
	a, err := generateIp(start, end, exit)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)
}
