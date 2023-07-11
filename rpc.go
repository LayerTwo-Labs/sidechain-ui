package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
)

type RPCRequest struct {
	JSONRpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPCGetBlockCountResponse struct {
	Result int `json:"result"`
}

type RPCGetDepositAddressResponse struct {
	Result string `json:"result"`
}

type RPCGetBalanceResponse struct {
	Result float64 `json:"result"`
}

type RPCGetUnconfirmedBalanceResponse struct {
	Result float64 `json:"result"`
}

func MakeRpcRequest(chainData *ChainData, method string, params []interface{}) (*http.Response, error) {
	auth := chainData.RPCUser + ":" + chainData.RPCPass
	authBytes := []byte(auth)
	authEncoded := base64.StdEncoding.EncodeToString(authBytes)
	rpcRequest := RPCRequest{JSONRpc: "2.0", ID: "", Method: method, Params: params}
	body, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+strconv.Itoa(chainData.Port), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+authEncoded)
	req.Header.Add("content-type", "application/json")
	return client.Do(req)
}

func PrintRPCErrorResponse(r *http.Response) {
	println("RPC Error: " + r.Status)
	println("Response:")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	println(buf.String())
}
