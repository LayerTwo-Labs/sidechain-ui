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

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/biter777/processex"
)

//go:embed binaries/linux-testchaind
var linuxBytes []byte

//go:embed sidechain.conf
var sidechainConfBytes []byte

//go:embed drivechain.conf
var drivechainConfBytes []byte

// TODO: Make these configurable in UI
const (
	sidechainDirName  = ".testchain"
	sidechainBinName  = "testchaind"
	sidechainConfName = "testchain.conf"
	drivechainDirName = ".drivechain"
	drivehainBinName  = "drivechaind"
	drivehainConfName = "drivechain.conf"
)

var (
	as  *AppState
	mui *MainUI
)

var (
	sidechainDir        string
	sidechainConfDir    string
	sidechainChainData  ChainData
	drivechainDir       string
	drivechainConfDir   string
	drivechainChainData ChainData
)

func main() {
	dirSetup()
	name := strings.ReplaceAll(sidechainDirName, ".", "")
	caser := cases.Title(language.English)
	as = NewAppState("com.ismyhc.sidechain-ui", caser.String(name))
	mui = NewMainUi(as)

	launchChain()
	mui.as.w.ShowAndRun()
}

func dirSetup() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Look for drivechain and bail if not found
	drivechainDir = homeDir + string(os.PathSeparator) + drivechainDirName
	if _, err := os.Stat(drivechainDir); os.IsNotExist(err) {
		// Drive chain not found at default location
		log.Fatal(err)
	}

	// Find drivechain.conf if not there write default
	// TODO: Copy old write ours?
	drivechainConfDir = drivechainDir + string(os.PathSeparator) + drivehainConfName
	if _, err := os.Stat(drivechainConfDir); os.IsNotExist(err) {
		err = os.WriteFile(drivechainConfDir, drivechainConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Look for sidechain dir and create if not found
	sidechainDir = homeDir + string(os.PathSeparator) + sidechainDirName
	if _, err := os.Stat(sidechainDir); os.IsNotExist(err) {
		os.MkdirAll(sidechainDir, 0o755)
	}

	// Find sidechains conf and if not found write default
	sidechainConfDir = sidechainDir + string(os.PathSeparator) + sidechainConfName
	if _, err := os.Stat(sidechainConfDir); os.IsNotExist(err) {
		err = os.WriteFile(sidechainConfDir, sidechainConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Load in drivechain conf
	loadConf(drivechainConfDir, &drivechainChainData)

	// Load in sidechain conf
	loadConf(sidechainConfDir, &sidechainChainData)

	// Write sidechain binary
	target := runtime.GOOS
	switch target {
	case "darwin":
		break
	case "linux":
		if _, err := os.Stat(sidechainDir + sidechainBinName); os.IsNotExist(err) {
			err = os.WriteFile(sidechainDir+sidechainBinName, linuxBytes, 0o755)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "windows":
		break
	}
}

func loadConf(path string, chainData *ChainData) {
	readFile, err := os.Open(path)
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
	err = json.Unmarshal(jsonData, &chainData)
	if err != nil {
		log.Fatal(err)
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
	process, _, err := processex.FindByName(sidechainBinName)
	if err == processex.ErrNotFound {
		fmt.Printf("Process %v not running", sidechainBinName)
		return nil, err
	}
	if err != nil {
		fmt.Printf("Process %v find error: %v", sidechainBinName, err)
		return nil, err
	}
	fmt.Printf("Process %v PID: %v", sidechainBinName, process[0].Pid)
	if len(process) > 0 {
		return process[0], nil
	}
	return nil, fmt.Errorf("something went wrong finding process")
}
