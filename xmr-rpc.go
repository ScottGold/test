/*json test*/

package main

import (
	"fmt"

	"encoding/json"
	"flag"
	"log"

	//"strings"
	"os"
	"time"

	"github.com/ScottGold/test/xmrtools"
)

type peer struct {
	id           uint64
	ip           uint32
	port         uint16
	rpc_port     uint16
	last_seen    uint64
	pruning_seed uint32
}

type RespGetPeerList struct {
	gray_list  json.RawMessage //[]peer
	status     string
	white_list json.RawMessage //[]peer

	//_Gray []peer
	//_White []peer
}

var (
	h                           bool
	rpcport                     int
	t                           bool
	CURRENT_TRANSACTION_VERSION int = 1002
	createwallet                bool
	generate                    int
	bid                         bool
	mining                      bool
	setlog                      bool
	get_balance                 bool
	cmd                         string
	waitformining               bool
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.IntVar(&rpcport, "p", 0, "monerod rpc port")
	flag.BoolVar(&t, "t", false, "test")
	flag.BoolVar(&createwallet, "createwallet", false, "create wallet, other else open the wallet")

	flag.IntVar(&generate, "generate", 0, "generate a block")

	// to monerod
	flag.BoolVar(&bid, "bid", false, "bid")
	flag.BoolVar(&mining, "mining", false, "start mining loop")
	flag.BoolVar(&setlog, "setlog", false, "set monerod log catalog")

	// to wallet
	flag.BoolVar(&get_balance, "get_balance", false, "get_balance")
	flag.StringVar(&cmd, "cmd", "", "command without parameters like: get_height, stop_wallet, get_version")

	flag.BoolVar(&waitformining, "waitformining", false, "waiting for popheight to mining.")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if t {
		fmt.Println("test")
		return
	}

	if h {
		flag.Usage()
		return
	}
	/* json test
	fmt.Println("rpcport", rpcport)

	//fmt.Println(XMRRpc(rpcport, "get_block_count", ""))
	peerList := xmrtools.XMRUrlCall(rpcport, "get_peer_list", "")
	fmt.Println(peerList)

	//dec := json.NewDecoder(strings.NewReader(peerList))
	//_, err := dec.Token()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	var kPeerList RespGetPeerList
	//err = dec.Decode(&kPeerList)
	err := json.Unmarshal([]byte(peerList), &kPeerList)
	if err == nil {
		var g, w []peer
		if len(kPeerList.gray_list) > 0 {
			err = json.Unmarshal(kPeerList.gray_list, &g)
			if err != nil {
				log.Printf("decode gray:", err)
			}
		}
		if len(kPeerList.white_list) > 0 {
			err = json.Unmarshal(kPeerList.white_list, &w)
			if err != nil {
				log.Printf("decode white:", err)
			}
		}
	} else {
		log.Printf("decode root:", err)
	}
	//
	//_, err = dec.Token()
	//if err != nil {
	//	log.Fatal(err)
	//}
	fmt.Println("---------------------")
	fmt.Println(kPeerList, kPeerList.status)
	*/

	//
	//xmrbuildbin := "/home/user1/pop-blockchain/monero-v0.15-no-new-pull/build/release/bin"
	xmrbuildbin := "c:/dev/bitcoin-0.18/poporodev/poporo/build/MINGW64_NT-10.0-17763/master/release/bin"
	//xmrd := xmrbuildbin + "/poporod"
	xmrWalletRPC := xmrbuildbin + "/poporo-wallet-rpc.exe"

	//18080 P2P_DEFAULT_PORT
	rpcPort1 := 18081 //note rpc port RPC_DEFAULT_PORT
	if rpcport != 0 {
		rpcPort1 = rpcport
	}
	walletrpcport1 := 9601

	if setlog {
		xmrtools.SetLogCategories(rpcPort1,
			`{"categories": "*:ERROR,net:ERROR,net.throttle:ERROR,net.p2p:FATAL,blockchain:TRACE,blockchain.db.lmdb:ERROR,cn:DEBUG,miner:INFO"}`)
	}

	//check if running
	retGetVersion := xmrtools.XMRRpc(walletrpcport1, "get_version", ``)
	if retGetVersion != "" {
		walletversion := xmrtools.ParseFieldInt(retGetVersion, "version")
		if walletversion == 0 {
			fmt.Println("get version error")
			retGetVersion = ""
		} else {
			fmt.Println("Wallet is running! version", walletversion)
		}
	}
	if retGetVersion == "" {
		//w1dir := "/home/user1/data/popwallet"
		w1dir := "C:/magnachain/btc0.18/xmr/xmrwallet"
		cleardir := false
		xmrtools.StartWalletRPC(xmrWalletRPC, rpcPort1, walletrpcport1, w1dir, cleardir)
	}

	fmt.Println("")
	var waddr1, vk1, sec1 string
	if createwallet {
		log.Println("create wallet1")
		xmrtools.XMRRpc(walletrpcport1, "create_wallet", `{"filename":"wallet","password":"","language":"English"}`)
	} else {
		log.Println("open wallet")
		xmrtools.XMRRpc(walletrpcport1, "open_wallet", `{"filename":"wallet","password":""}`)
	}

	//下面两个需要用到minder address
	waddr1, vk1, sec1 = xmrtools.GetMinerAddress(walletrpcport1)
	if bid {
		fmt.Println("bid")
		amount := "0.0005"
		blockHeightOfBidAddress := int64(1) // only one address now, so any value will available
		log.Println(xmrtools.XMRBid(rpcPort1, amount, blockHeightOfBidAddress, vk1))
	}
	//注意 XMRGenBlock 不能用在主网上
	if generate > 0 { //gen one block
		var xmrBlockCount int64 = 0
		xmrtools.XMRGenBlock(rpcPort1, int64(generate), waddr1, sec1, &xmrBlockCount)
	}
	if mining { //loop gen blocks
		ret := xmrtools.WalletStartMining(walletrpcport1, waddr1, sec1)
		fmt.Println(ret)
	}

	if get_balance {
		ret := xmrtools.XMRRpc(walletrpcport1, "get_balance", `{"account_index":0,"address_indices":[0,1]}}'`)
		fmt.Println(ret)
	}
	if cmd != "" {
		if cmd == "get_version" && retGetVersion != "" {
			log.Println(retGetVersion)
		} else {
			ret := xmrtools.XMRRpc(walletrpcport1, cmd, ``)
			log.Println(ret)
		}
	}
	if waitformining {
		for {
			h := xmrtools.GetXMRHeight(rpcPort1)
			if h >= int64(2051693) {
				log.Println("start mining")
				ret := xmrtools.WalletStartMining(walletrpcport1, waddr1, sec1)
				fmt.Println(ret)
				return
			} else {
				log.Println("height ", h, "/2051693")
				time.Sleep(30 * time.Second)
			}
		}
	}
	//stop wallet default
	//fmt.Println("stop wallet")
	//xmrtools.XMRRpc(walletrpcport1, "stop_wallet", ``)
}

func usage() {
	fmt.Fprintf(os.Stderr, `mygo version: mygo/1.10.0
Usage: mygo [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)
	flag.PrintDefaults()
}

/*
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func main() {
	const jsonStream = `
	[
		{"Name": "Ed", "Text": "Knock knock."},
		{"Name": "Sam", "Text": "Who's there?"},
		{"Name": "Ed", "Text": "Go fmt."},
		{"Name": "Sam", "Text": "Go fmt who?"},
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	]
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("111%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var m Message
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("----%v: %v\n", m.Name, m.Text)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("222%T: %v\n", t, t)

}

*/
