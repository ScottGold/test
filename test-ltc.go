package main

import (
	//"strconv"
	//"log"
	"fmt"
	//"os"
	"os/exec"
	"strings"

	"time"

	//"bufio"
	//"regexp"

	"github.com/ScottGold/test/btctools"
	"github.com/ScottGold/test/common"
)

func main() {
	btcd := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoind.exe"
	btc_cli := "c:/dev/bitcoin-0.18/bitcoin-0.18/build_msvc/x64/Debug/bitcoin-cli.exe"
	btcdir := "C:/magnachain/btc0.18/btc"
	btcdatadir := "-datadir=" + btcdir

	bchd := "C:/dev/bitcoin-0.18/Litcoin-test/build_msvc/x64/Debug/Litcoind.exe"
	bch_cli := "C:/dev/bitcoin-0.18/Litcoin-test/build_msvc/x64/Debug/Litcoin-cli.exe"
	ltcdir1 := "C:/magnachain/btc0.18/ltc/d1"
	ltcdir2 := "C:/magnachain/btc0.18/ltc/d2"
	bchdatadir1 := "-datadir=" + ltcdir1
	bchdatadir2 := "-datadir=" + ltcdir2

	//try stop first
	btctools.StopBtc(btc_cli, btcdatadir)
	btctools.StopBtc(bch_cli, bchdatadir1)
	btctools.StopBtc(bch_cli, bchdatadir2)
	//
	common.ClearDataDir(btcdir + "/regtest")
	common.ClearDataDir(ltcdir1 + "/regtest")
	common.ClearDataDir(ltcdir2 + "/regtest")

	fmt.Println("start BTC")
	btccmd := exec.Command(btcd, btcdatadir)
	btccmd.Start()
	btctools.WaitToLoadFinish(btc_cli, btcdatadir)
	fmt.Println("start BTC finish")

	btcAddress, _ := btctools.CliCommand(btc_cli, btcdatadir, "btc getnewaddress:", "getnewaddress")

	fmt.Println("BTC generate 200")
	btctools.CliCommand(btc_cli, btcdatadir, "generate blocks", "generatetoaddress", "200", string(btcAddress))

	fmt.Println("start BCH 1")
	bch1cmd := exec.Command(bchd, bchdatadir1)
	bch1cmd.Start()
	btctools.WaitToLoadFinish(bch_cli, bchdatadir1)
	fmt.Println("start BCH 1 finish")

	fmt.Println("start BCH 2")
	bch2cmd := exec.Command(bchd, bchdatadir2)
	bch2cmd.Start()
	btctools.WaitToLoadFinish(bch_cli, bchdatadir2)
	fmt.Println("start BCH 2 finish")

	btctools.CliCommand(bch_cli, bchdatadir1, "addNode", "addnode", "127.0.0.1:8202", "onetry")

	btctools.WaitToConnectPeer(bch_cli, bchdatadir1)

	bchAddress1, _ := btctools.CliCommand(bch_cli, bchdatadir1, "bch1 getNewAddress", "getnewaddress")
	strBchAddress1 := strings.Trim(string(bchAddress1), " \r\n")
	fmt.Println("bch1 getnewaddress ", strBchAddress1)

	bchAddress2, _ := btctools.CliCommand(bch_cli, bchdatadir2, "bch2 getNewAddress", "getnewaddress")
	strBchAddress2 := strings.Trim(string(bchAddress2), " \r\n")
	fmt.Println("bch2 getnewaddress ", strBchAddress2)

	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter to continue: ")
	//reader.ReadString('\n')

	//first 424
	btctools.CliCommand(bch_cli, bchdatadir1, "cli1 generate 12 blocks", "generatetoaddress", "12", strBchAddress1)
	btctools.WaitToSyncBlock(bch_cli, bchdatadir1, bchdatadir2)
	btctools.CliCommand(bch_cli, bchdatadir2, "cli2 generate 12 blocks", "generatetoaddress", "12", strBchAddress2)
	btctools.WaitToSyncBlock(bch_cli, bchdatadir1, bchdatadir2)
	for i := 0; i < 4; i++ {
		btctools.CliCommand(bch_cli, bchdatadir1, "cli1 generate 50 blocks", "generatetoaddress", "50", strBchAddress1)
		btctools.WaitToSyncBlock(bch_cli, bchdatadir1, bchdatadir2)
		btctools.CliCommand(bch_cli, bchdatadir2, "cli2 generate 50 blocks", "generatetoaddress", "50", strBchAddress2)
		btctools.WaitToSyncBlock(bch_cli, bchdatadir1, bchdatadir2)
	}

	//bid start
	mocktime := time.Now().Unix()
	bchComplianceTime := int64(150)
	bchMineTime := int64(30)
	btcMineCount := int64(0)

	btctools.PrintAllChainBlockCount(bch_cli, bchdatadir1, bchdatadir2, btc_cli, btcdatadir)

	for {
		btctools.CliCommand(bch_cli, bchdatadir1, "cli1 bid", "bid", "1", "1")
		btctools.CliCommand(bch_cli, bchdatadir2, "cli2 bid", "bid", "1", "1")

		if btcMineCount >= 600 {
			fmt.Println("r", btcMineCount)
			btcMineCount = 0
		}
		if btcMineCount == 0 {
			btctools.SetMockTime(btc_cli, btcdatadir, mocktime)
			btctools.CliCommand(btc_cli, btcdatadir, "btc gen 1 blocks", "generatetoaddress", "1", string(btcAddress))
		}
		for i := 0; i < 4; i++ {
			timespace := bchComplianceTime
			blockcount := btctools.GetBlockCount(bch_cli, bchdatadir1)
			fmt.Println("blockcount", blockcount)
			if blockcount >= 1000 {
				timespace = bchMineTime
			}
			mocktime = mocktime + timespace
			btcMineCount = btcMineCount + timespace

			btctools.SetMockTime(bch_cli, bchdatadir1, mocktime)
			btctools.SetMockTime(bch_cli, bchdatadir2, mocktime)

			btctools.GenerateBlocks(bch_cli, bchdatadir1, "1", strBchAddress1)
			newblockcount := btctools.GetBlockCount(bch_cli, bchdatadir1)
			if newblockcount == blockcount { // chain 1 gen fail
				btctools.GenerateBlocks(bch_cli, bchdatadir2, "1", strBchAddress2)
				newblockcount = btctools.GetBlockCount(bch_cli, bchdatadir2)
			}

			btctools.WaitToSyncBlock(bch_cli, bchdatadir1, bchdatadir2)
		}
	}

	btccmd.Wait()
	bch1cmd.Wait()
	bch2cmd.Wait()
}
