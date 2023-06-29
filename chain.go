package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

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
	ID               string  `json:"id"`
	State            State   `json:"state"`
	RefreshBMM       bool    `json:"refreshbmm"`
	AvailableBalance float64 `json:"availablebalance"`
	PendingBalance   float64 `json:"pendingbalance"`
	Height           int     `json:"height,omitempty"`
}

type State uint

const (
	Unknown State = iota
	Waiting
	Running
)

var (
	sidechainChainStateUpdate     *time.Ticker
	quitsidechainChainStateUpdate chan struct{}
)

func getChainProcess(name string) (*os.Process, error) {
	process, _, err := processex.FindByName(sidechainBinName)
	if err == processex.ErrNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(process) > 0 {
		return process[0], nil
	}
	return nil, fmt.Errorf("something went wrong finding process")
}

func LaunchChain(cd *ChainData, cs *ChainState) {
	p, err := getChainProcess(cd.BinName)
	if p != nil && err == nil {
		// We are already running...
		println(cd.BinName + " already running...")
		return
	}

	args := []string{"-conf=" + cd.ConfDir}
	cmd := exec.Command(cd.Dir+string(os.PathSeparator)+cd.BinName, args...)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cs.State = Waiting
	println(cd.BinName + " Started...")
}

func StopChain(chainData ChainData) {
	// _, err := getChainProcess()
	// if err != nil {
	// 	return
	// }
	// Shutdown chain gracefully via rpc
}

func StartSidechainStateUpdate(as *AppState, mui *MainUI) {
	sidechainChainStateUpdate = time.NewTicker(1 * time.Second)
	quitsidechainChainStateUpdate = make(chan struct{})
	go func() {
		for {
			select {
			case <-sidechainChainStateUpdate.C:
				updateUI := false
				// getblockcount
				bcr, err := MakeRpcRequest(&as.scd, "getblockcount", []interface{}{})
				if err != nil {
					if as.scs.State != Waiting {
						as.scs.State = Waiting
						updateUI = true
					}
					fmt.Printf(err.Error())
				} else {
					defer bcr.Body.Close()
					if bcr.StatusCode == 200 {
						var res RPCGetBlockCountResponse
						err := json.NewDecoder(bcr.Body).Decode(&res)
						if err == nil {
							println(res.Result)
							if as.scs.Height != res.Result {
								as.scs.Height = res.Result
								updateUI = true
							}
						}
					}
				}

				if updateUI {
					mui.Refresh()
				}

			case <-quitsidechainChainStateUpdate:
				sidechainChainStateUpdate.Stop()
				return
			}
		}
	}()
}
