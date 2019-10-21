/*

 */

package common

import (
	"fmt"
	"os"
)

func ClearDataDir(dir1 string) {
	errRM := os.RemoveAll(dir1)
	if errRM != nil {
		fmt.Println("rm", dir1, "data fail", errRM.Error())
		panic("ClearDataDir fail")
	}
	os.Mkdir(dir1, os.ModeDir)
}
