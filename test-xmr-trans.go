/*
一个比特币结点，两个xmr结点，测试正常的pop挖矿，pop挖矿的分叉合并
*/

package main

import (
	//"strconv"
	"log"
	"os"
	"os/exec"

	//"os/signal"
	"fmt"
	"time"

	//"strings"
	//"bufio"
	//"regexp"
	//"syscall"
	//"net"
	//"net/http"
	//"bytes"
	//"io/ioutil"
	"flag"

	"github.com/mytest/btctools"
	"github.com/mytest/common"
	"github.com/mytest/xmrtools"
)

var (
	BTC_dir string = "C:/magnachain/btc0.18/btc/regtest"
	c       chan os.Signal
)

//func SysSignal() {
//	select {
//    case <-c:
//		fmt.Println("exit signal")
//		os.Exit(3)
//		break
//    default:
//    }
//}

var (
	h   bool
	cmd string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&cmd, "cmd", "", "Command: stop")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `test-xmr-trans version: test-xmr-trans/1.0.1
Usage: test-xmr-trans [-h -cmd=command]
Options:
`)
	flag.PrintDefaults()
}

func CloseAll(btc_cli, btcdatadir string, walletrpcport1, walletrpcport2, rpcPort1, rpcPort2 int) {
	btctools.StopBtc(btc_cli, btcdatadir)

	xmrtools.XMRRpc(walletrpcport1, "stop_wallet", ``)
	xmrtools.XMRRpc(walletrpcport2, "stop_wallet", ``)
	xmrtools.XMRUrlCall(rpcPort1, "stop_daemon", "")
	xmrtools.XMRUrlCall(rpcPort2, "stop_daemon", "")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//go SysSignal()

	localip := xmrtools.GetLocalIp()
	//fmt.Println("localip", localip)

	btcd := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoind.exe"
	btc_cli := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoin-cli.exe"
	btcdatadir := "-datadir=C:/magnachain/btc0.18/btc"

	xmrd := "C:/dev/bitcoin-0.18/monero-v0.14/build/release/bin/monerod.exe"
	xmrWalletRPC := "C:/dev/bitcoin-0.18/monero-v0.14/build/release/bin/monero-wallet-rpc.exe"

	xmrdirroot := "C:/magnachain/btc0.18/xmr/"
	dir1, dir2 := xmrdirroot+"d1", xmrdirroot+"d2"
	p2pPort1, p2pPort2 := 8401, 8402
	rpcPort1, rpcPort2 := 9401, 9402
	zmqRpcPort1, zmqRpcPort2 := 9501, 9502

	walletrpcport1, walletrpcport2 := 9601, 9602

	CloseAll(btc_cli, btcdatadir, walletrpcport1, walletrpcport2, rpcPort1, rpcPort2)
	if cmd == "stop" {
		fmt.Println("just stop")
		return
	}

	pdir1, pdir2 := fmt.Sprintf("--data-dir=%s", dir1), fmt.Sprintf("--data-dir=%s", dir2)
	pppPort1, pppPort2 := fmt.Sprintf("--p2p-bind-port=%d", p2pPort1), fmt.Sprintf("--p2p-bind-port=%d", p2pPort2)
	prPort1, prPort2 := fmt.Sprintf("--rpc-bind-port=%d", rpcPort1), fmt.Sprintf("--rpc-bind-port=%d", rpcPort2)
	pzPort1, pzPort2 := fmt.Sprintf("--zmq-rpc-bind-port=%d", zmqRpcPort1), fmt.Sprintf("--zmq-rpc-bind-port=%d", zmqRpcPort2)

	//limitrateup := "--limit-rate-up=819200"
	//limitratedown := "--limit-rate-down=819200"
	//limitrate := "--limit-rate=819200"

	testAddr1 := "4AZ4HFjsRw8ZMttnREegMZ23qXYMfUPyka8cNz18vMCjH3b1JWL5fV9cWuCWMANmusHS21Z23kiaheYztq4wJoZCCciDXvb"
	testAddr2 := "49Wq1rBbUJMTbcHrsYaXNf91bjMmJi9bVUvYZEtUwjJc6QYVU4EsQ8Scn4sFwM5Boy1wwYU4sm5tVVtwMkoovMvdBvAMC68"
	//testPrivateVK1 := "5a22b1c029e7405c374c141260d7744253d7b8b270e3d9aec811b3e6c16b9e03"
	//testPrivateVK2 := "897b712fc1f16d50fa11ee0ef644dbb3082b138042e51f65d79c9ddd48df8808"
	//testPubVK1 := "eca1cf1f88ede8d45559c9502ab442623c8d496eb449dee081c43e700b992b66"
	//testPubVK2 := "51e11c6fdfce2bd5e5931a849f712c05a9c931a1ed0d5aacc135753cef8b9e60"

	xmrParam1 := []string{pdir1, pppPort1, prPort1, pzPort1}
	xmrParam2 := []string{pdir2, pppPort2, prPort2, pzPort2}

	commonParams := []string{}
	commonParams = append(commonParams, "--btc-rpc-ip=127.0.0.1")
	commonParams = append(commonParams, "--btc-rpc-port=9001")
	commonParams = append(commonParams, "--btc-rpc-login=user:pwd")
	commonParams = append(commonParams, "--regtest")
	commonParams = append(commonParams, "--fixed-difficulty=1")
	commonParams = append(commonParams, "--non-interactive")
	//commonParams = append(commonParams, "--rpc-login=user:pwd") // Digest Authentication, i don't know how to write
	commonParams = append(commonParams, "--log-level=0")
	commonParams = append(commonParams, "--allow-local-ip")
	commonParams = append(commonParams, "--btcbidstart=200")
	commonParams = append(commonParams, "--popforkheight=1000")
	commonParams = append(commonParams, "--xmrbidstart=424")

	common.ClearDataDir(dir1)
	common.ClearDataDir(dir2)
	common.ClearDataDir(BTC_dir)

	//---------------------------
	addpeer := fmt.Sprintf("--add-peer=%s:%d", localip, p2pPort2)
	//xmrParam1 = append(xmrParam1, addpeer)
	xmrParam1 = append(xmrParam1, commonParams...)
	go xmrtools.StartXMRD(xmrd, xmrParam1...)
	xmrtools.WaitToXMRLoadFinish(rpcPort1)

	//---------------------------
	addpeer = fmt.Sprintf("--add-peer=%s:%d", localip, p2pPort1)
	xmrParam2 = append(xmrParam2, addpeer)
	xmrParam2 = append(xmrParam2, commonParams...)
	go xmrtools.StartXMRD(xmrd, xmrParam2...)
	xmrtools.WaitToXMRLoadFinish(rpcPort2)
	//----------------------------

	xmrtools.SetLogCategories(rpcPort1)
	xmrtools.SetLogCategories(rpcPort2)

	//---------------------------
	time.Sleep(0 * time.Second)
	xmrtools.WaitXMRGetPeer(rpcPort1)
	xmrtools.WaitXMRGetPeer(rpcPort2)

	//start wallet
	fmt.Println("start wallet rpc")
	w1dir, w2dir := xmrdirroot+"w1", xmrdirroot+"w2"

	xmrtools.StartWalletRPC(xmrWalletRPC, rpcPort1, walletrpcport1, w1dir)
	xmrtools.StartWalletRPC(xmrWalletRPC, rpcPort2, walletrpcport2, w2dir)

	fmt.Println("create wallet1")
	xmrtools.XMRRpc(walletrpcport1, "create_wallet", `{"filename":"ww1","password":"","language":"English"}`)
	waddr1, vk1, sec1 := xmrtools.GetMinerAddress(walletrpcport1)
	fmt.Println("create wallet2")
	xmrtools.XMRRpc(walletrpcport2, "create_wallet", `{"filename":"ww1","password":"","language":"English"}`)
	waddr2, vk2, sec2 := xmrtools.GetMinerAddress(walletrpcport2)

	//---------------------------
	go func() {
		var btccmd *exec.Cmd
		fmt.Println("start BTC")
		btccmd = exec.Command(btcd, btcdatadir)
		btccmd.Start()
		btccmd.Wait()
	}()
	btctools.WaitToLoadFinish(btc_cli, btcdatadir)
	fmt.Println("start BTC finish")

	btcAddress, _ := btctools.CliCommand(btc_cli, btcdatadir, "btc getnewaddress:", "getnewaddress")

	fmt.Println("BTC generate 200")
	btctools.CliCommand(btc_cli, btcdatadir, "generate blocks", "generatetoaddress", "200", string(btcAddress))

	var xmrBlockCount int64 = 1
	for xmrBlockCount < 424+1 { //before bid
		xmrtools.XMRGenBlock(rpcPort1, 106, waddr1, sec1, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)

		xmrtools.XMRGenBlock(rpcPort2, 106, waddr2, sec2, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)
	}

	recipients := []xmrtools.Recipient{}
	recipients = append(recipients, xmrtools.Recipient{int64(100000000000), testAddr1})
	recipients = append(recipients, xmrtools.Recipient{int64(200000000000), testAddr2})

	log.Println("refresh wallet", xmrtools.XMRRpc(walletrpcport1, "refresh", `{"start_height":1}`))

	log.Println(xmrtools.SendTo(walletrpcport1, recipients))

	log.Println("start 1 vs 5 gen loop", xmrBlockCount)
	var lastsendblockheight int64 = 0
	switchop := true
	mocktime := time.Now().Unix()
	xmrmocktime := mocktime
	for {
		xmrtools.XMRBid(rpcPort1, "1", 1, vk1)
		xmrtools.XMRBid(rpcPort2, "1", 1, vk2)

		btctools.SetMockTime(btc_cli, btcdatadir, mocktime)
		btctools.CliCommand(btc_cli, btcdatadir, "btc gen 1 blocks", "generatetoaddress", "1", string(btcAddress))

		for i := 0; i < 5; i++ {
			mocktime = mocktime + 120
			xmrmocktime = xmrmocktime + 120 + 2
			xmrtools.SetMockTime(rpcPort1, xmrmocktime)
			xmrtools.SetMockTime(rpcPort2, xmrmocktime)
			if switchop {
				xmrtools.XMRGenBlock(rpcPort1, 1, waddr1, sec1, &xmrBlockCount)
				xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)
			} else {
				xmrtools.XMRGenBlock(rpcPort2, 1, waddr2, sec2, &xmrBlockCount)
				xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)
			}
		}

		switchop = !switchop

		if xmrBlockCount-lastsendblockheight > 4 || xmrBlockCount == 1001 {
			//log.Println(xmrtools.SendTo(walletrpcport1, recipients))
			lastsendblockheight = xmrBlockCount
		}
		//if lastsendblockheight >= 1030 {
		//	break
		//}
	}

	fmt.Println("main end")
}
