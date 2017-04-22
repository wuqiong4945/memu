package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-ini/ini"
)

var cfg *ini.File
var iniFileName = "memu.ini"

func main() {
	var err error
	cfg, err = ini.Load(iniFileName)
	if err != nil {
		iniFile, err := os.OpenFile(iniFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
		iniFile.WriteString(``)
		iniFile.Close()

		cfg, err = ini.Load(iniFileName)
		if err != nil {
			fmt.Printf("Error in open file : %s\n", err)
		} else {
			fmt.Println("Memu.ini file is created.")
		}
	}

	fmt.Println("system infomation : " + runtime.GOOS + "/" + runtime.GOARCH)
	t := time.Now()
	fmt.Println("start memu!")
	mame := NewMame()
	fmt.Printf("took amount of time: %s\n", time.Now().Sub(t).String())
	// mame.Fresh()
	mame.Update()
	mame.Audit()

	fmt.Println("current mame version is : " + mame.Build)
	// fmt.Printf("%#v\n", mame.Machine("qsound"))
	fmt.Printf("%#v\n", mame.Debug)
	fmt.Println(mame.Machine("qsound"))

	// out := mame.VerifyRoms("qsound")
	out, err := exec.Command("mame/mame64", "-verifyroms", "qsound").Output()
	//out, err = exec.Command("mame/mame64", "kov2p").Output()
	fmt.Println(out, err)

	// info := GetInfo("aoh", "history")
	// fmt.Println(info)
	// exec.Command("firefox", info).Output()
	var info string
	for _, machine := range mame.Machines {
		if machine.MachineStatus != MACHINE_NEXIST {
			info += machine.Status()
		}
	}
	html, _ := os.OpenFile("info.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer html.Close()
	html.WriteString(info)

	fmt.Printf("Start memu took amount of time: %s\n", time.Now().Sub(t).String())
}

func CheckError(err error) {
	if err == nil {
		return
	}
	log.Println(err)
}
