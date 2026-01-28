package capture

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestListNetworkInterfaces(t *testing.T) {
	result, err := ListNetworkInterfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

type consoleListener struct {
}

func (l *consoleListener) OnPkg(pkg *OutputPacket) {
	pkgJsonData, _ := json.Marshal(pkg)
	println(string(pkgJsonData))
}

func (l *consoleListener) OnClose() {
	println("OnClose")
}

func TestStartCapture(t *testing.T) {
	err := StartCapture(FilterParam{
		InterfaceName: "en1",
		Port: &struct {
			Src int
			Dst int
			Any int
		}{Src: 0, Dst: 0, Any: 8080},
	}, &consoleListener{})
	if err != nil {
		println(err.Error())
	}
}
