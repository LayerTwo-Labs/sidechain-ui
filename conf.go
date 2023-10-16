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

////go:embed binaries/linux-testchain-qt
//var linuxBytes []byte

//go:embed binaries/linux-testchaind
var linuxBytes []byte

//go:embed sidechain.conf
var sidechainConfBytes []byte

// TODO: Make these configurable in UI
const (
	sidechainDirName = ".testchain"
	// sidechainBinName = "testchain-qt"
	sidechainBinName   = "testchaind"
	sidechainConfName  = "testchain.conf"
	drivechainConfName = "drivechain.conf"
)

func ConfInit(as *AppState) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	drivechainDir := homeDir + string(os.PathSeparator) + ".drivechain"
	if _, err := os.Stat(drivechainDir); os.IsNotExist(err) {
		println(err.Error())
		return err
	}

	drivechainConfDir := drivechainDir + string(os.PathSeparator) + drivechainConfName
	if _, err := os.Stat(drivechainConfDir); os.IsNotExist(err) {
		// drivechain.conf not found.
		// TODO: Ask user for location of drivechain.conf
		println(err.Error())
		return err
	}
	// Look for sidechain dir and create if not found
	sidechainDir := homeDir + string(os.PathSeparator) + sidechainDirName
	if _, err := os.Stat(sidechainDir); os.IsNotExist(err) {
		os.MkdirAll(sidechainDir, 0o755)
	}

	// Find sidechains conf and if not found write default
	sidechainConfDir := sidechainDir + string(os.PathSeparator) + sidechainConfName
	if _, err := os.Stat(sidechainConfDir); os.IsNotExist(err) {
		// append datadir
		sidechainConfBytes = append(sidechainConfBytes, "\ndatadir="+sidechainDir...)
		err = os.WriteFile(sidechainConfDir, sidechainConfBytes, 0o755)
		if err != nil {
			println(err.Error())
			return err
		}
	}

	drivechainChainData := ChainData{}
	drivechainChainData.ParentChain = true
	drivechainChainData.Dir = drivechainDir
	drivechainChainData.ConfDir = drivechainConfDir
	drivechainChainData.MinimumFee = 0.001 // TODO: Figure out how to estimate this

	// Load in drivechain conf
	loadConf(&drivechainChainData)
	as.pcd = drivechainChainData
	as.pcs = ChainState{}

	sidechainChainData := ChainData{}
	sidechainChainData.ParentChain = false
	sidechainChainData.Dir = sidechainDir
	sidechainChainData.ConfDir = sidechainConfDir
	sidechainChainData.BinName = sidechainBinName
	sidechainChainData.MinimumFee = 0.001 // TODO: Figure out how to estimate this
	sidechainChainData.BMMFee = 0.001
	sidechainChainData.RefreshBMM = true

	// Load in sidechain conf
	loadConf(&sidechainChainData)
	as.scd = sidechainChainData
	as.scs = ChainState{}
	as.scs.Slot = as.scd.Slot

	// Init data
	as.scbmmtd = []BMMTableItem{}

	// Write sidechain binary
	target := runtime.GOOS
	switch target {
	case "darwin":
		break
	case "linux":
		binDr := sidechainDir + string(os.PathSeparator) + sidechainBinName
		if _, err := os.Stat(binDr); os.IsNotExist(err) {
			err = os.WriteFile(binDr, linuxBytes, 0o755)
			if err != nil {
				log.Fatal(err)
				println(err.Error())
				return err
			}
		}
	case "windows":
		break
	}
	return nil
}

func loadConf(chainData *ChainData) error {
	readFile, err := os.Open(chainData.ConfDir)
	if err != nil {
		println(err.Error())
		return err
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
		println(err.Error())
		return err
	}
	return nil
}
