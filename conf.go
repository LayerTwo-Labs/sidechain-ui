package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
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

func ConfInit(as *AppState) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Look for drivechain and bail if not found
	drivechainDir := homeDir + string(os.PathSeparator) + drivechainDirName
	if _, err := os.Stat(drivechainDir); os.IsNotExist(err) {
		// Drive chain not found at default location
		log.Fatal(err)
	}

	// Find drivechain.conf if not there write default
	// TODO: Copy old write ours?
	drivechainConfDir := drivechainDir + string(os.PathSeparator) + drivehainConfName
	if _, err := os.Stat(drivechainConfDir); os.IsNotExist(err) {
		err = os.WriteFile(drivechainConfDir, drivechainConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Look for sidechain dir and create if not found
	sidechainDir := homeDir + string(os.PathSeparator) + sidechainDirName
	if _, err := os.Stat(sidechainDir); os.IsNotExist(err) {
		os.MkdirAll(sidechainDir, 0o755)
	}

	// Find sidechains conf and if not found write default
	sidechainConfDir := sidechainDir + string(os.PathSeparator) + sidechainConfName
	if _, err := os.Stat(sidechainConfDir); os.IsNotExist(err) {
		err = os.WriteFile(sidechainConfDir, sidechainConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	drivechainChainData := ChainData{}
	drivechainChainData.ParentChain = true
	drivechainChainData.Dir = drivechainDir
	drivechainChainData.ConfDir = drivechainConfDir
	drivechainChainData.BinName = drivehainBinName

	// Load in drivechain conf
	loadConf(&drivechainChainData)
	as.pcd = drivechainChainData

	sidechainChainData := ChainData{}
	sidechainChainData.ParentChain = false
	sidechainChainData.Dir = sidechainDir
	sidechainChainData.ConfDir = sidechainConfDir
	sidechainChainData.BinName = sidechainBinName

	// Load in sidechain conf
	loadConf(&sidechainChainData)
	as.scd = sidechainChainData

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

func loadConf(chainData *ChainData) {
	readFile, err := os.Open(chainData.ConfDir)
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
