/*
一个比特币结点，两个xmr结点，测试正常的pop挖矿，pop挖矿的分叉合并
*/

package xmrtools

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	//"os/signal"
	"fmt"
	"strings"
	"time"

	//"bufio"
	"regexp"
	//"syscall"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
)

var (
	//BTC_dir string = "C:/magnachain/btc0.18/btc/regtest"
	defaultLogType string = "TRACE"
)

func GetLocalIp() string {
	var localip string
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ip := strings.Split(addr.String(), "/")[0]
		fmt.Println("ip", ip)
		match, _ := regexp.MatchString("192.168.", ip) //TODO:
		if match {
			localip = ip
		}
	}
	_ = localip
	return localip
}

func XMRRpc(rpcPort int, method, params string) string {
	url := fmt.Sprintf("http://127.0.0.1:%d/json_rpc", rpcPort)
	var jsonStr []byte
	if params == "" {
		jsonStr = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"0","method":"%s"}`, method))
	} else {
		jsonStr = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"0","method":"%s","params":%s}`, method, params))
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.SetBasicAuth("user", "pwd") // Use HTTP Digest Authentication
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("RPC error", err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.Status != "200 Ok" {
		fmt.Println("http status", resp.Status)
		return ""
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func XMRUrlCall(rpcPort int, method, params string) string {
	url := fmt.Sprintf("http://127.0.0.1:%d/%s", rpcPort, method)
	reader := strings.NewReader(params)
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("RPC error", err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.Status != "200 Ok" {
		fmt.Println("http status", resp.Status)
		return ""
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func WaitToXMRLoadFinish(rpcPort int) {
	for {
		time.Sleep(5 * time.Second)
		ret := XMRRpc(rpcPort, "get_version", "")
		if ret != "" {
			fmt.Printf("XMR %d rpc connect ok!\n", rpcPort)
			return
		}
		fmt.Printf("Wait XMR %d to finish...\n", rpcPort)
	}
}

func ParseFieldString(json, field string) string {
	re := regexp.MustCompile(fmt.Sprintf(`"%s": "[A-Za-z0-9]+"`, field))
	clen := len(fmt.Sprintf(`"%s": "`, field))
	idex := re.FindIndex([]byte(json))
	if idex == nil || len(idex) < 2 {
		fmt.Println("find index fail", idex)
		return ""
	}
	strSub := json[idex[0]+clen : idex[1]-1]
	return strSub
}

//json: get_height return val
func ParseFieldInt(json, field string) (count int64) {
	re := regexp.MustCompile(fmt.Sprintf(`"%s": [0-9]+`, field))
	clen := len(fmt.Sprintf(`"%s": `, field))
	idex := re.FindIndex([]byte(json))
	if idex == nil || len(idex) < 2 {
		fmt.Println("find index fail", idex)
		return 0
	}
	strSub := json[idex[0]+clen : idex[1]]
	count, _ = strconv.ParseInt(strSub, 10, 64)
	return
}

func GetBlockCount(rpcPort int) (count int64) {
	bc := XMRRpc(rpcPort, "get_block_count", "")
	count = ParseFieldInt(bc, "count")
	return
}

func GetXMRHeight(rpcPort int) int64 {
	bh := XMRUrlCall(rpcPort, "get_height", "")
	return ParseFieldInt(bh, "height")
}

func GetMinerAddress(walletRpcPort int) (address, pub_view_key, view_key string) {
	getaddress := XMRRpc(walletRpcPort, "get_address", `{"account_index":0,"address_index":[]}`)
	//fmt.Println(getaddress)
	address = ParseFieldString(getaddress, "address")
	pub_view_key = ParseFieldString(getaddress, "pub_view_key")
	fmt.Println("address", address)
	fmt.Println("pub_view_key", pub_view_key)

	view_key = XMRRpc(walletRpcPort, "query_key", `{"key_type":"view_key"}`)
	view_key = ParseFieldString(view_key, "key")
	fmt.Println("view_key", view_key) //view_key spend_key mnemonic
	return
}

func WaitXMRSyncBlock(rpcPort1, rpcPort2 int, expectheight int64) {
	usetime := 0
	waitCount := 0
	for {
		var sleeptime time.Duration = 1 + time.Duration(waitCount*5)
		time.Sleep(sleeptime * time.Second)

		count1 := GetBlockCount(rpcPort1)
		count2 := GetBlockCount(rpcPort2)
		if count1 == count2 {
			if expectheight == 0 || count1 == expectheight {
				break
			}
		}

		log.Println("WaitXMRSyncBlock count", count1, count2, "expect", expectheight)
		usetime++
		if usetime > 30 {
			usetime = 0
			waitCount++
		}
	}
}

func StartXMRD(xmrd string, params ...string) {
	name := params[0]
	fi := strings.LastIndex(xmrd, "/")
	exename := xmrd[fi+1:]
	fmt.Println("start", exename, name)
	fmt.Println(params)
	cmd := exec.Command(xmrd, params...)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Command start with error: %v", err)
		panic("StartXMRD fail")
	}
	err = cmd.Wait()
	fmt.Printf("xmr %s exit error: %v\n", exename, name, err)
}

func ClearDataDir(dir1 string) {
	errRM := os.RemoveAll(dir1)
	if errRM != nil {
		fmt.Println("rm", dir1, "data fail", errRM.Error())
		panic("ClearDataDir fail")
	}
	os.Mkdir(dir1, os.ModeDir)
}

func WaitXMRGetPeer(rpcPort int) {
	fmt.Println("wait peer to connect")
	for {
		ret := XMRUrlCall(rpcPort, "get_peer_list", "")
		match, _ := regexp.MatchString("white_list", ret)
		if match {
			fmt.Println(ret)
			fmt.Println("node", rpcPort, "get peer ok")
			break
		} else {
			fmt.Println(ret)
			time.Sleep(3 * time.Second)
		}
	}
}

func XMRGenBlock(rpcPort int, amount_of_blocks int64, wallet_address string, miner_sec_key string, xmrBlockCount *int64) bool {
	log.Println(rpcPort, "generating", amount_of_blocks, "block, blockcount", *xmrBlockCount+amount_of_blocks)

	prev_block := ""
	starting_nonce := 0
	params := fmt.Sprintf("{\"amount_of_blocks\":%d,\"wallet_address\":\"%s\", \"prev_block\":\"%s\", \"starting_nonce\":%d, \"miner_sec_key\":\"%s\"}",
		amount_of_blocks, wallet_address, prev_block, starting_nonce, miner_sec_key)

	ret := XMRRpc(rpcPort, "generateblocks", params)
	match, _ := regexp.MatchString(`"status": "OK"`, ret)
	if match { //TODO: may some of them success
		*xmrBlockCount = *xmrBlockCount + amount_of_blocks
	}
	return match
}

func XMRBid(rpcPort int, amount string, block_height int64, pub_view_key string) string {
	params := fmt.Sprintf("{\"amount\":\"%s\",\"block_height\":%d,\"pub_view_key\":\"%s\"}", amount, block_height, pub_view_key)
	return XMRUrlCall(rpcPort, "bid", params)
}

/*

-----------Level is one of the following-------------
FATAL - higher level
ERROR
WARNING
INFO
DEBUG
TRACE - lower level A level automatically includes higher level. By default,
*/
func SetLogCategories(rpcPort int) {
	params := fmt.Sprintf(`{"categories": "*:%s,net:ERROR,net.throttle:ERROR,net.p2p:FATAL,blockchain.db.lmdb:ERROR"}`, defaultLogType)
	XMRUrlCall(rpcPort, "set_log_categories", params)
}

func SetBan(rpcPort1, banTime int) {
	LocalIP := GetLocalIp()
	var params, localip string
	if banTime > 0 {
		params = fmt.Sprintf(`{"bans":[{"host":"%s","ban":true,"seconds":%d}]}`, LocalIP, banTime)
		localip = fmt.Sprintf(`{"bans":[{"host":"%s","ban":true,"seconds":%d}]}`, "127.0.0.1", banTime)
	} else {
		params = fmt.Sprintf(`{"bans":[{"host":"%s","ban":false}]}`, LocalIP)
		localip = fmt.Sprintf(`{"bans":[{"host":"%s","ban":false}]}`, "127.0.0.1")
	}

	XMRRpc(rpcPort1, "set_bans", params)
	XMRRpc(rpcPort1, "set_bans", localip)
}

func StartWalletRPC(xmrWalletRPC string, dhostRpcPort, rpcPort int, walletDir string, clearData bool) {
	walletParam := []string{}
	walletParam = append(walletParam, fmt.Sprintf("--daemon-port=%d", dhostRpcPort))
	walletParam = append(walletParam, fmt.Sprintf("--rpc-bind-port=%d", rpcPort))
	walletParam = append(walletParam, fmt.Sprintf("--wallet-dir=%s", walletDir))
	walletParam = append(walletParam, "--non-interactive")
	walletParam = append(walletParam, "--disable-rpc-login")
	walletParam = append(walletParam, fmt.Sprintf("--log-file=%s/wallet_rpc%d.log", walletDir, rpcPort))
	walletParam = append(walletParam, "--log-level=4")

	if clearData {
		ClearDataDir(walletDir)
	}

	go StartXMRD(xmrWalletRPC, walletParam...)

	WaitToXMRLoadFinish(rpcPort)
}

type Recipient struct {
	Amount  int64
	Address string
}

func (v *Recipient) String() string {
	return fmt.Sprintf(`
    {"amount":%d,"address":"%s"}`, v.Amount, v.Address)
}

func SendTo(rpcPort int, recipient []Recipient) string {
	var dest string
	for i, v := range recipient {
		if i > 0 {
			dest = dest + ","
		}
		dest = dest + fmt.Sprintf(`
    {"amount":%d,"address":"%s"}`, v.Amount, v.Address)
	}
	params := fmt.Sprintf(
		`{
  "destinations":[%s],
  "account_index":0,
  "subaddr_indices":[0],
  "priority":0,
  "ring_size":7,
  "get_tx_key": true
}`, dest)

	log.Println("transfer ", params)
	return XMRRpc(rpcPort, "transfer", params)
}

func SetMockTime(rpcPort int, timestamp int64) string {
	params := fmt.Sprintf(`{"timestamp":%d}`, timestamp)
	return XMRRpc(rpcPort, "setmocktime", params)
}
