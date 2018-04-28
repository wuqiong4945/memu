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
var mame *Mame

func main() {
	if runtime.GOOS == "windows" {
		mamePath = "mame/mame64.exe"
	}

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

	mame = NewMame()
	fmt.Printf("starting took amount of time: %s\n", time.Now().Sub(t).String())
	mame.Update()
	mame.Audit()

	fmt.Println("current mame version is : " + mame.Build)
	fmt.Printf("%#v\n", mame.Debug)

	out := mame.VerifyRoms("pgm")
	fmt.Println(string(out))

	// out = mame.Machine("sfa3").Start()
	// fmt.Println(string(out))

	var info string
	info += htmlHead
	// info += `<div class="card-group">`
	info += `<div class="card-columns">`
	for _, machine := range mame.Machines {
		if machine.MachineStatus&MACHINE_EXIST == MACHINE_EXIST {
			info += machine.GetStatusInfo()
		}
	}

	info += `</div>`
	info += htmlEnd
	html, _ := os.OpenFile("info.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer html.Close()
	html.WriteString(info)

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

var htmlHead = `<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="theme/css/bootstrap.min.css">
  </head>
  <body>
  `
var htmlEnd = `        <!-- jQuery first, then Popper.js, then Bootstrap JS. -->
    <script src="theme/js/jquery.slim.min.js"></script>
		<script src="theme/js/tether.min.js">_</script>
    <script src="theme/js/popper.js"></script>
    <script src="theme/js/bootstrap.min.js"></script>
  </body>
</html>`
