/*json test*/

package main

import (
	"fmt"

	"encoding/json"
	"flag"
	"log"

	//"strings"
	"os"

	"github.com/mytest/xmrtools"
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
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.IntVar(&rpcport, "p", 9401, "monerod rpc port")
	flag.BoolVar(&t, "t", false, "test")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if t {

		return
	}

	if h {
		flag.Usage()
		return
	}
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
