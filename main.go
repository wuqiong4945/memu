package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-ini/ini"
)

var cfg *ini.File
var iniFileName = "memu.ini"
var mamePath = "mame/mame64"
var isFlush = false

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
	mamePath = cfg.Section("general").Key("mame").MustString("mame/mame64")

	fmt.Println("system infomation : " + runtime.GOOS + "/" + runtime.GOARCH)
	t := time.Now()
	fmt.Println("start memu!")

	mame := NewMame()
	fmt.Printf("starting took amount of time: %s\n", time.Now().Sub(t).String())
	mame.Update()
	mame.Audit()

	fmt.Println("current mame version is : " + mame.Build)
	fmt.Printf("%#v\n", mame.Debug)
	// fmt.Println(mame.Machine("qsound"))

	out := mame.VerifyRoms("pgm")
	fmt.Println(string(out))

	// out = mame.Machine("sfa3").Start()
	// fmt.Println(string(out))

	// info := GetInfo("aoh", "history")
	// fmt.Println(info)
	var info string
	for _, machine := range mame.Machines {
		if machine.MachineStatus&MACHINE_EXIST == MACHINE_EXIST {
			info += machine.GetStatusInfo()
		}
	}
	info += mame.Machine("mslug3").GetStatusInfo()
	info += mame.Machine("kov2p").GetStatusInfo()
	info += mame.Machine("pgm").GetStatusInfo()

	/*  ms, err := mame.Search("chess") */
	// CheckError(err)
	// if err == nil {
	// info += "<table>"
	// for _, m := range ms {
	// info += "<tr><td>" + m.Name + "</td><td>" + m.Description + "</td><td>" + m.Year + "</td></tr>"
	// }
	// info += "</table>"
	/* } */

	html, _ := os.OpenFile("info.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer html.Close()
	html.WriteString(info)

	// out, _ = exec.Command(mamePath, "-lx").Output()
	// html.WriteString(string(out))

	fmt.Println(mame.Version())
	if isFlush == true {
		mame.Flush()
	}
	fmt.Printf("Start memu took amount of time: %s\n", time.Now().Sub(t).String())
}

func CheckError(err error) {
	if err == nil {
		return
	}
	log.Println(err)
}
