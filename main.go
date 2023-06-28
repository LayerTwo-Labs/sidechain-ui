package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/biter777/processex"
)

//go:embed binaries/linux-testchaind
var linuxBytes []byte

//go:embed default.conf
var defaultConfBytes []byte

const (
	DATA_DIR  = ".testchain"
	BIN_NAME  = "testchaind"
	CONF_NAME = "testchain.conf"
)

var (
	as  *AppState
	mui *MainUI
)

var (
	chainDir  string
	confDir   string
	chainData ChainData
)

func main() {
	dirSetup()
	as = NewAppState("com.ismyhc.sidechain-ui", "Sidechain UI")
	mui = NewMainUi(as)

	launchChain()
	mui.as.w.ShowAndRun()
}

func dirSetup() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	chainDir = homeDir + string(os.PathSeparator) + DATA_DIR
	if _, err := os.Stat(chainDir); os.IsNotExist(err) {
		os.MkdirAll(chainDir, 0o755)
	}

	confDir = chainDir + string(os.PathSeparator) + CONF_NAME
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err = os.WriteFile(confDir, defaultConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// load
	readFile, err := os.Open(confDir)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	confMap := make(map[string]interface{})

	for _, line := range fileLines {
		a := strings.Split(line, "=")
		if len(a) == 2 {
			k := strings.TrimSpace(a[0])
			v := strings.TrimSpace(a[1])
			if k != "" {
				iv, err := (strconv.ParseInt(v, 0, 64))
				if err != nil {
					confMap[k] = v
				} else {
					confMap[k] = int(iv)
				}
			}
		}
	}

	jsonData, _ := json.Marshal(confMap)

	var chainData ChainData
	err = json.Unmarshal(jsonData, &chainData)
	if err != nil {
		log.Fatal(err)
	}

	target := runtime.GOOS
	switch target {
	case "darwin":
		break
	case "linux":
		if _, err := os.Stat(chainDir + BIN_NAME); os.IsNotExist(err) {
			err = os.WriteFile(chainDir+BIN_NAME, linuxBytes, 0o755)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "windows":
		break
	}
}

func launchChain() {
	_, err := getChainProcess()
	if err != nil {
		return
	}
}

func stopChain() {
	_, err := getChainProcess()
	if err != nil {
		return
	}
	// Shutdown chain gracefully via rpc
}

func killChain() {
	process, err := getChainProcess()
	if err != nil {
		return
	}
	process.Kill()
}

func getChainProcess() (*os.Process, error) {
	process, _, err := processex.FindByName(BIN_NAME)
	if err == processex.ErrNotFound {
		fmt.Printf("Process %v not running", BIN_NAME)
		return nil, err
	}
	if err != nil {
		fmt.Printf("Process %v find error: %v", BIN_NAME, err)
		return nil, err
	}
	fmt.Printf("Process %v PID: %v", BIN_NAME, process[0].Pid)
	if len(process) > 0 {
		return process[0], nil
	}
	return nil, fmt.Errorf("something went wrong finding process")
}
