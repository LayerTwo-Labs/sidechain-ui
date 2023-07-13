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
	ParentChain bool    `json:"parentchain,omitempty"`
	BinName     string  `json:"binname,omitempty"`
	Regtest     int     `json:"regtest"`
	Port        int     `json:"rpcport"`
	RPCUser     string  `json:"rpcuser"`
	RPCPass     string  `json:"rpcpassword"`
	Dir         string  `json:"dir,omitempty"`
	ConfDir     string  `json:"confdir,omitempty"`
	DataDir     string  `json:"datadir,omitempty"`
	Slot        *int    `json:"slot,omitempty"`
	RefreshBMM  bool    `json:"refreshbmm,omitempty"`
	MinimumFee  float64 `json:"minimumfee,omitempty"`
	BMMFee      float64 `json:"bmmfee,omitempty"`
}

type ChainState struct {
	ID               string  `json:"id"`
	State            State   `json:"state"`
	AvailableBalance float64 `json:"availablebalance"`
	PendingBalance   float64 `json:"pendingbalance"`
	Height           int     `json:"height,omitempty"`
	Slot             *int    `json:"slot,omitempty"`
}

func (cs *ChainState) FormatedAvailableBalance(withSymbol bool) string {
	if cs.Slot != nil && withSymbol {
		return fmt.Sprintf("%.8f SC%d", cs.AvailableBalance, *cs.Slot)
	} else {
		return fmt.Sprintf("%.8f", cs.AvailableBalance)
	}
}

func (cs *ChainState) FormatedPendingBalance(withSymbol bool) string {
	if cs.Slot != nil && withSymbol {
		return fmt.Sprintf("%.8f SC%d", cs.PendingBalance, *cs.Slot)
	} else {
		return fmt.Sprintf("%.8f", cs.PendingBalance)
	}
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
	println(cd.ConfDir)
	cmd := exec.Command(cd.Dir+string(os.PathSeparator)+cd.BinName, args...)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cs.State = Waiting
	println(cd.BinName + " Started...")
}

func StopChain(cd *ChainData, cs *ChainState) error {
	req, err := MakeRpcRequest(cd, "stop", []interface{}{})
	if err != nil {
		return err
	} else {
		defer req.Body.Close()
		return nil
	}
}

func StartSidechainStateUpdate(as *AppState, mui *MainUI) {
	sidechainChainStateUpdate = time.NewTicker(1 * time.Second)
	quitsidechainChainStateUpdate = make(chan struct{})
	go func() {
		for {
			select {
			case <-sidechainChainStateUpdate.C:
				updateUI := false
				if GetBlockHeight(&as.scd, &as.scs) && !updateUI {
					updateUI = true
				}
				if GetBalance(&as.scd, &as.scs) && !updateUI {
					updateUI = true
				}
				RefreshBMM(&as.scd, &as.scs)
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

func GetBlockHeight(cd *ChainData, cs *ChainState) bool {
	currentHeight := cs.Height
	bcr, err := MakeRpcRequest(cd, "getblockcount", []interface{}{})
	if err != nil {
		println(err.Error())
	} else {
		defer bcr.Body.Close()
		if bcr.StatusCode == 200 {
			var res RPCGetBlockCountResponse
			err := json.NewDecoder(bcr.Body).Decode(&res)
			if err == nil {
				cs.Height = res.Result
				if currentHeight != cs.Height {
					return true
				}
			}
		}
	}
	return false
}

func GetBalance(cd *ChainData, cs *ChainState) bool {
	currentBalance := cs.AvailableBalance
	bcr, err := MakeRpcRequest(cd, "getbalance", []interface{}{})
	if err != nil {
		println(err.Error())
	} else {
		defer bcr.Body.Close()
		if bcr.StatusCode == 200 {
			var res RPCGetBalanceResponse
			err := json.NewDecoder(bcr.Body).Decode(&res)
			if err == nil {
				cs.AvailableBalance = res.Result
				if currentBalance != cs.AvailableBalance {
					return true
				}
			}
		}
	}
	return false
}

func RefreshBMM(cd *ChainData, cs *ChainState) {
	fee := cd.BMMFee
	if fee == 0 {
		fee = 0.001
	}
	req, err := MakeRpcRequest(cd, "refreshbmm", []interface{}{fee})
	if err != nil {
		println(err.Error())
	} else {
		defer req.Body.Close()
		if req.StatusCode == 200 {
			var res RPCRefreshBMMResponse
			err := json.NewDecoder(req.Body).Decode(&res)
			if err != nil {
				println(err.Error())
			} else {
				// fmt.Printf("%+v\n", res) TODO:
			}
		} else {
			PrintNonSuccessRPCResponse(req)
		}
	}
}

func WithdrawFromSidechain(cd *ChainData, mc *ChainData, mcAddress string, rfAddress string, amount float64, scfee float64, mcfee float64) error {
	// TODO: Validate
	args := []interface{}{
		mcAddress,
		rfAddress,
		amount,
		scfee,
		mcfee,
	}
	req, err := MakeRpcRequest(cd, "createwithdrawal", args)
	if err != nil {
		return err
	} else {
		defer req.Body.Close()
		if req.StatusCode == 200 {
			return nil
		} else {
			PrintNonSuccessRPCResponse(req)
		}
		return fmt.Errorf("withdraw request unsuccessful, rpc status code: %v", req.StatusCode)
	}
}

func GetSidechainDepositAddress(cd *ChainData, formated bool) (*string, error) {
	method := "getnewaddress"
	args := []interface{}{"", "legacy"}
	if formated {
		method = "getdepositaddress"
		args = []interface{}{}
	}
	req, err := MakeRpcRequest(cd, method, args)
	if err != nil {
		return nil, err
	} else {
		defer req.Body.Close()
		if req.StatusCode == 200 {
			var res RPCGetDepositAddressResponse
			err := json.NewDecoder(req.Body).Decode(&res)
			if err != nil {
				return nil, err
			} else {
				return &res.Result, nil
			}
		} else {
			return nil, fmt.Errorf("cannot get deposit address, rpc status code: %v", req.StatusCode)
		}
	}
}

func GetDrivechainDepositAddress(mc *ChainData) (*string, error) {
	req, err := MakeRpcRequest(mc, "getnewaddress", []interface{}{"", "legacy"})
	if err != nil {
		return nil, err
	} else {
		defer req.Body.Close()
		if req.StatusCode == 200 {
			var res RPCGetDepositAddressResponse
			err := json.NewDecoder(req.Body).Decode(&res)
			if err != nil {
				return nil, err
			} else {
				return &res.Result, nil
			}
		} else {
			return nil, fmt.Errorf("cannot get deposit address, rpc status code: %v", req.StatusCode)
		}
	}
}
