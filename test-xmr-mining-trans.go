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

	"github.com/ScottGold/test/btctools"
	"github.com/ScottGold/test/xmrtools"
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

func main() {
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//go SysSignal()

	localip := xmrtools.GetLocalIp()

	btcd := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoind.exe"
	btc_cli := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoin-cli.exe"
	btcdatadir := "-datadir=C:/magnachain/btc0.18/btc"

	xmrbuildbin := 'C:/dev/bitcoin-0.18/monero-v0.15/build/MINGW64_NT-10.0-17763/master/release/bin'
	xmrd := xmrbuildbin + "/monerod.exe"
	xmrWalletRPC := xmrbuildbin+ "/monero-wallet-rpc.exe"

	xmrdirroot := "C:/magnachain/btc0.18/xmr/"
	dir1, dir2 := xmrdirroot+"d1", xmrdirroot+"d2"
	p2pPort1, p2pPort2 := 8401, 8402
	rpcPort1, rpcPort2 := 9401, 9402
	zmqRpcPort1, zmqRpcPort2 := 9501, 9502

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
	commonParams = append(commonParams, "--log-level=4")
	commonParams = append(commonParams, "--allow-local-ip")
	commonParams = append(commonParams, "--btcbidstart=200")
	commonParams = append(commonParams, "--popforkheight=1000")

	xmrtools.ClearDataDir(dir1)
	xmrtools.ClearDataDir(dir2)
	xmrtools.ClearDataDir(BTC_dir)

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

	//xmrtools.SetLogCategories(rpcPort1)
	//xmrtools.SetLogCategories(rpcPort2)

	//---------------------------
	time.Sleep(0 * time.Second)
	xmrtools.WaitXMRGetPeer(rpcPort1)
	xmrtools.WaitXMRGetPeer(rpcPort2)

	//start wallet
	fmt.Println("start wallet rpc")
	walletrpcport1, walletrpcport2 := 9601, 9602
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
	for xmrBlockCount < 280+1 { //before bid
		xmrtools.XMRGenBlock(rpcPort1, 70, waddr1, sec1, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)

		xmrtools.XMRGenBlock(rpcPort2, 70, waddr2, sec2, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)
	}

	recipients := []xmrtools.Recipient{}
	recipients = append(recipients, xmrtools.Recipient{int64(100000000000), testAddr1})
	recipients = append(recipients, xmrtools.Recipient{int64(200000000000), testAddr2})

	log.Println(xmrtools.SendTo(walletrpcport1, recipients))

	log.Println("start one four loop", xmrBlockCount)
	var lastsendblockheight int64 = 0
	for {
		xmrtools.XMRBid(rpcPort1, "1", 1, vk1)
		xmrtools.XMRBid(rpcPort2, "1", 1, vk2)

		xmrtools.XMRGenBlock(rpcPort1, 2, waddr1, sec1, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)

		xmrtools.XMRGenBlock(rpcPort2, 2, waddr2, sec2, &xmrBlockCount)
		xmrtools.WaitXMRSyncBlock(rpcPort1, rpcPort2, xmrBlockCount)
		btctools.CliCommand(btc_cli, btcdatadir, "btc gen 1 blocks", "generatetoaddress", "1", string(btcAddress))

		if xmrBlockCount-lastsendblockheight > 20 {
			log.Println(xmrtools.SendTo(walletrpcport1, recipients))
			lastsendblockheight = xmrBlockCount
		}
	}

	fmt.Println("main end")
}
