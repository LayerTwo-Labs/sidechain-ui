package main

import (
	"fmt"
	"os"

	"github.com/biter777/processex"
)

type ChainData struct {
	ParentChain      bool   `json:"parentchain,omitempty"`
	BinName          string `json:"binname,omitempty"`
	Regtest          int    `json:"regtest"`
	Port             int    `json:"rpcport"`
	RPCUser          string `json:"rpcuser"`
	RPCPass          string `json:"rpcpassword"`
	Dir              string `json:"dir,omitempty"`
	ConfDir          string `json:"confdir,omitempty"`
	DataDir          string `json:"datadir,omitempty"`
	Slot             *int   `json:"slot,omitempty"`
	MinerBreakForBMM *int   `json:"minerbreakforbmm,omitempty"`
}

type ChainState struct {
	ID               string `json:"id"`
	State            State  `json:"state"`
	RefreshBMM       bool   `json:"refresh_bmm"`
	AvailableBalance float64
	PendingBalance   float64
}

type State uint

const (
	Unknown State = iota
	Waiting
	Running
)

func getChainProcess(name string) (*os.Process, error) {
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

// func LaunchChain(chainData ChainData) {
// 	p, err := getChainProcess(chainData.BinName)
// 	if p != nil && err == nil {
// 		// We are already running...
// 		println(chainData.BinName + " already running...")
// 		return
// 	}
// 	chainDataDir := switchboardDir + "/data/" + chain.ID
// 	if _, err := os.Stat(chainDataDir); os.IsNotExist(err) {
// 		os.MkdirAll(chainDataDir, 0o755)
// 	}
// 	var regtest string = "0"
// 	if chain.Regtest {
// 		regtest = "1"
// 	}
// 	args := []string{"-regtest=" + regtest, "-datadir=" + chainDataDir, "-rpcport=" + chain.Port, "-rpcuser=" + chain.RPCUser, "-rpcpassword=" + chain.RPCPass, "-server=1"}
// 	cmd := exec.Command(switchboardDir+"/"+chain.Bin, args...)
// 	err := cmd.Start()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	chainState[chain.ID] = ChainState{ID: chain.ID, State: Waiting, RefreshBMM: true, CMD: cmd}

// 	setMainContentUI(selectedChainDataIndex)
// 	// if err != nil {
// 	// 	return
// 	// }
// }

func StopChain(chainData ChainData) {
	// _, err := getChainProcess()
	// if err != nil {
	// 	return
	// }
	// Shutdown chain gracefully via rpc
}
