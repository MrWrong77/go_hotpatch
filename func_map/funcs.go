package funcmap

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var FuncMap = map[string]uintptr{}

func init() {
	generateFuncName2PtrDict()
}

func generateFuncName2PtrDict() {
	fileFullPath := os.Args[0]

	cmd := exec.Command("nm", fileFullPath)
	contentBytes, err := cmd.Output()
	if err != nil {
		println(err)
		return
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		arr := strings.Split(line, " ")
		if len(arr) < 3 {
			continue
		}
		funcSymbol, addr := arr[2], arr[0]
		addrUint, _ := strconv.ParseUint(addr, 16, 64)
		FuncMap[funcSymbol] = uintptr(addrUint)
	}
}
