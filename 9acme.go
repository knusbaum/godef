// +build plan9

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"errors"
)

type acmeFile struct {
	name       string
	body       []byte
	offset     int
	runeOffset int
}

func acmeCurrentFile() (*acmeFile, error) {
	winid := os.Getenv("winid")
	if winid == "" {
		return nil, fmt.Errorf("$winid not set - not running inside acme?")
	}
	path := "/mnt/acme/" + winid + "/"
	addrF, err := os.Open(path + "addr")
	if err != nil {
		return nil, err
	}
	defer addrF.Close()

//	b := make([]byte, 40)
//	addrF.Read(b)
//	fmt.Printf("Start ADDR: %s\n", string(b))

	ctlF, err := os.Create(path + "ctl")
	if err != nil {
		return nil, err
	}
	ctlF.Write([]byte("addr=dot\n"))
	ctlF.Close()

	buf := make([]byte, 40)
	n, err := addrF.Read(buf)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("DOT ADDR: %s\n", string(buf))
	a := strings.Fields(string(buf[0:n]))
	if len(a) < 2 {
		return nil, errors.New("short read from acme addr")
	}

	q0, err := strconv.Atoi(a[0])
	if err != nil {
		return nil, fmt.Errorf("invalid read from acme addr: %s", err)
	}

	bodyF, err := os.Open(path + "body")
	if err != nil {
		return nil, fmt.Errorf("failed to read window body: %s", err)
	}
	defer bodyF.Close()
	body, err := ioutil.ReadAll(bodyF)
	if err != nil {
		return nil, fmt.Errorf("failed to read window body: %s", err)
	}

	tagF, err := os.Open(path + "tag")
	if err != nil {
		return nil, fmt.Errorf("failed to read window tag: %s", err)
	}
	defer tagF.Close()
	tagB, err := ioutil.ReadAll(tagF)
	if err != nil {
		return nil, fmt.Errorf("failed to read window tag: %s", err)
	}
	tag := string(tagB)

	i := strings.Index(tag, " ")
	if i == -1 {
		return nil, fmt.Errorf("strange tag with no spaces")
	}

	w := &acmeFile{
		name:       tag[0:i],
		body:       body,
		offset:     runeOffset2ByteOffset(body, q0),
		runeOffset: q0,
	}
	return w, nil
}

func runeOffset2ByteOffset(b []byte, off int) int {
	r := 0
	for i, _ := range string(b) {
		if r == off {
			return i
		}
		r++
	}
	return len(b)
}