package btctools

import (
	"log"
	"strconv"

	//"os"
	"fmt"
	"os/exec"
	"strings"
	"time"

	//"bufio"
	"regexp"
)

func WaitToLoadFinish(cli_exe string, datadir string) {
	fmt.Println("wait to ", datadir, " finish")
	for { //wait to init finish
		time.Sleep(3 * time.Second)

		cmd := exec.Command(cli_exe, datadir, "getblockcount")
		out, err := cmd.CombinedOutput()
		if err == nil {
			return
		}
		strOut := string(out)
		match, _ := regexp.MatchString("Loading wallet", strOut)
		if match {
			fmt.Println(datadir, "Loading wallet...")
		} else {
			fmt.Println(strOut)
		}
	}
}

func PrintAllChainBlockCount(bch_cli string, bchdatadir1 string, bchdatadir2 string, btc_cli string, btcdatadir string) {
	t := time.Now()
	fmt.Println(t.Format(time.RFC3339))

	cmd := exec.Command(bch_cli, bchdatadir1, "getblockcount")
	blockcount1, _ := cmd.CombinedOutput()
	fmt.Printf("bch1 getblockcount %s", blockcount1)
	cmd = exec.Command(bch_cli, bchdatadir2, "getblockcount")
	blockcount2, _ := cmd.CombinedOutput()
	fmt.Printf("bch2 getblockcount %s", blockcount2)
	cmd = exec.Command(btc_cli, btcdatadir, "getblockcount")
	blockcount, _ := cmd.CombinedOutput()
	fmt.Printf("BTC getblockcount %s\n", blockcount)
}

func CliCommand(cli string, datadir string, errLog string, rpc string, params ...string) (string, error) {
	var arg []string
	arg = append(arg, "-rpcclienttimeout=0", datadir, rpc)
	arg = append(arg, params...)
	cmd := exec.Command(cli, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(errLog, string(output), err.Error())
		panic("CliCommand fail")
	}
	return string(output), err
}

func GenerateBlocks(cli, datadir, genCount, mineAddress string) {
	defer RecoveFunc()
	maxtry := "10"
	CliCommand(cli, datadir, datadir+" gen blocks "+genCount, "generatetoaddress", genCount, mineAddress, maxtry)
}

func WaitToConnectPeer(bch_cli string, bchdatadir1 string) {
	for {
		strnum, errGCC := CliCommand(bch_cli, bchdatadir1, "getConnections", "getconnectioncount")
		if errGCC != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		strnum = strings.Trim(strnum, " \r\n")
		num, _ := strconv.Atoi(strnum)
		if num > 0 {
			return
		}
		fmt.Println("getconnectioncount:", strnum, num)
	}
}

func WaitToSyncBlock(bch_cli string, bchdatadir1 string, bchdatadir2 string) {
	for {
		getblockcount1, _ := CliCommand(bch_cli, bchdatadir1, "get blocks count", "getblockcount")
		getblockcount2, _ := CliCommand(bch_cli, bchdatadir2, "get blocks count", "getblockcount")

		getblockcount1 = strings.Trim(getblockcount1, " \r\n")
		getblockcount2 = strings.Trim(getblockcount2, " \r\n")
		c1, _ := strconv.Atoi(getblockcount1)
		c2, _ := strconv.Atoi(getblockcount2)
		if c1 == c2 {
			return
		}
		time.Sleep(1 * time.Second)
		fmt.Println("getblockcount:", c1, c2)
	}
}

func GetBlockCount(btc_cli string, btcdatadir string) (blockcount int) {
	strret, _ := CliCommand(btc_cli, btcdatadir, "get blocks count", "getblockcount")
	strret = strings.Trim(strret, " \r\n")
	blockcount, _ = strconv.Atoi(strret)
	return
}

func RecoveFunc( /*finish chan<- int*/ ) {
	if err := recover(); err != nil {
		log.Println("RecoveFunc catch error:", err)
	}
	//finish <- 1
}

func StopBtc(btc_cli, btcdatadir string) {
	defer RecoveFunc() //recover defer must write before the panic
	CliCommand(btc_cli, btcdatadir, "stop btc", "stop")
}

func SetMockTime(cli string, datadir string, timestamp int64) {
	CliCommand(cli, datadir, "", "setmocktime", strconv.FormatInt(timestamp, 10))
}
