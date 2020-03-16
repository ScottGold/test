package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

//---------------------------------
var (
	h        bool
	tfile    string
	rootpath string
	copyto   string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&tfile, "tfile", "", "file to be search.")
	flag.StringVar(&rootpath, "rootpath", "", "root path, need with / end")
	flag.StringVar(&copyto, "copyto", "", "copy to path, need with / end")
}

func usage() {
	fmt.Fprintf(os.Stderr, `mygo version: mygo/1.10.0
Usage: get-require-files [-h] [-tfile=filename]

Options:
`)
	flag.PrintDefaults()
}

//---------------------------------

func ReadFileLines(filename string) []string {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	return lines
}

//查找文件中的关键词
func SearchFile(filename string, waitToSearch []string, allfiles map[string]int) []string {
	lines := ReadFileLines(filename)
	for _, l := range lines {
		//re := regexp.MustCompile(`import\s*"[A-Za-z0-9/.]+"`)
		re := regexp.MustCompile(`import(\s*|\s*\{[ A-Za-z0-9_]*\}\s*from\s*)"[A-Za-z0-9/.]+"`)
		idex := re.FindIndex([]byte(l))

		if len(idex) == 2 {
			//fmt.Println(l)
			re2 := regexp.MustCompile(`"`)
			idex2 := re2.FindAllIndex([]byte(l), -1)
			//fmt.Println(len(idex2))
			if len(idex2) == 2 {
				fpath := l[idex2[0][1]:idex2[1][0]]

				//fpath = rootpath + fpath
				lastindex := strings.LastIndex(filename, "/")
				if lastindex != -1 {
					relatepath := filename[:lastindex+1]
					fpath = relatepath + fpath
					fpath = path.Clean(fpath)
					//fmt.Println(fpath)
				}

				if allfiles[fpath] == 0 {
					allfiles[fpath] = 1
					waitToSearch = append(waitToSearch, fpath)
				}
			}
		}
	}
	return waitToSearch
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	if tfile == "" {
		fmt.Println("error, -tfile cat not empty string")
		return
	}

	allfiles := make(map[string]int)
	waitToSearch := []string{}

	tfile = rootpath + tfile
	allfiles[tfile] = 1
	waitToSearch = append(waitToSearch, tfile)

	for {
		if len(waitToSearch) == 0 {
			break
		}

		waitToSearch = SearchFile(waitToSearch[0], waitToSearch, allfiles)

		waitToSearch = waitToSearch[1:]
	}

	//copy files
	for k, v := range allfiles {
		if v != 0 {
			fmt.Println("cp", k, copyto)
			//generate a sh copy file
			//TODO: file duplicate
		}
	}
	fmt.Println(len(allfiles))
}
