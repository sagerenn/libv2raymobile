package libv2raymobile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	// core "github.com/v2fly/v2ray-core/v5"
	core "github.com/xtls/xray-core/core"
	// serial "github.com/v2fly/v2ray-core/v5/infra/conf/serial"
	serial "github.com/xtls/xray-core/infra/conf/serial"
	// _ "github.com/v2fly/v2ray-core/v5/main/distro/all"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type CoreManager struct {
	inst      *core.Instance
	shouldOff chan int
}

func (m *CoreManager) runConfigSync(confPath string) {
	bs := readFileAsBytes(confPath)

	r := bytes.NewReader(bs)
	config, err := serial.LoadJSONConfig(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	if m.inst != nil {
		fmt.Println("m.inst != nil")
		return
	}
	m.inst, err = core.New(config)
	if err != nil {
		log.Println(err)
		return
	}

	err = m.inst.Start()
	if err != nil {
		log.Println(err)
		return
	}

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()

	{
		m.shouldOff = make(chan int, 1)
		<-m.shouldOff
	}
}

func (m *CoreManager) RunConfig(confPath string) {
	go m.runConfigSync(confPath)
}

func readFileAsBytes(filePath string) (bs []byte) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the file into a byte slice
	bs = make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	return
}

func (m *CoreManager) Stop() {
	m.shouldOff <- 1
	m.inst.Close()
	m.inst = nil
}

func SetEnv(key string, val string) {
	os.Setenv(key, val)
}
