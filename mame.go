package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha1"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func NewMame() (mame *Mame) {
	mame = new(Mame)

	cacheFile, err := os.OpenFile("cache.gob", os.O_CREATE|os.O_RDWR, os.ModePerm)
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in open gob file : %s\n", err)
		return
	}
	defer cacheFile.Close()

	dec := gob.NewDecoder(cacheFile)
	err = dec.Decode(mame)
	CheckError(err)
	if err != nil || mame.Build == "" {
		fmt.Printf("Error in decoding gob : %s\n", err)
	}

	return
}

func (mame *Mame) Fresh() {
	*mame = Mame{}

	f, err := os.Open("a.xml")
	CheckError(err)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	CheckError(err)
	err = xml.Unmarshal(data, mame)
	CheckError(err)

	mame.Flush()
	return
}

func (mame Mame) Flush() {
	cacheFile, err := os.OpenFile("cache.gob", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in opening gob file : %s\n", err)
		return
	}
	defer cacheFile.Close()

	enc := gob.NewEncoder(cacheFile)
	err = enc.Encode(mame)
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in encoding gob : %s\n", err)
	}

	return
}

func (mame *Mame) Update() {
	if mame.Build == "" {
		mame.Fresh()
		return
	}

	f, err := os.Open("a.xml")
	CheckError(err)
	defer f.Close()

	reader := bufio.NewReader(f)
	var version string
	reg := regexp.MustCompile(`^\s*(\w+)\s*=\s*"([^\"]*)"\s+`)
	for {
		lineBytes, _, err := reader.ReadLine()
		line := strings.TrimSpace(string(lineBytes))
		if strings.HasPrefix(line, "<mame ") {
			line = strings.TrimPrefix(line, "<mame ")
			line = strings.TrimSuffix(line, ">") + " "

			for {
				attr := reg.FindString(line)
				if attr == "" {
					break
				}

				key := reg.ReplaceAllString(attr, `$1`)
				value := reg.ReplaceAllString(attr, `$2`)
				if key == "build" {
					version = value
					break
				}

				line = strings.TrimPrefix(line, attr)
			}

			break
		}

		if err == io.EOF {
			break
		}
	}

	if version > mame.Build {
		mame.Fresh()
	}

	return
}

func (mame *Mame) Machine(machineName string) (machine *Machine) {
	for k, _ := range mame.Machines {
		if mame.Machines[k].Name == machineName {
			machine = &mame.Machines[k]
			return
		}
	}

	return
}

func (mame Mame) VerifyRoms(machineName string) (result []byte) {
	mamePath := cfg.Section("general").Key("mame").MustString("mame/mame64")
	result, _ = exec.Command(mamePath, "-verifyroms", machineName).Output()
	return
}

func (mame *Mame) Audit() {
	t := time.Now()
	defer fmt.Printf("Audit time: %s\n", time.Now().Sub(t).String())

	romDirs := "roms"
	//fullPath, _ := filepath.Abs(path)

	// idleFiles := []string{}
	dirs := strings.Split(romDirs, ";")
	for _, dir := range dirs {
		list, err := ioutil.ReadDir(dir)
		CheckError(err)
		fmt.Println("\n===== " + dir + "是目录 =====")

		// Iterate through the files in the directory
		for _, info := range list {
			switch {
			// dispose machine in chd file
			case info.IsDir() == true:
				// just chd game has sub dir
				fmt.Println("\n--- " + info.Name() + "是目录 ---")
				CHDMachineName := info.Name()
				CHDDirectory := dir + "/" + CHDMachineName
				machine := mame.Machine(CHDMachineName)
				if machine == nil {
					fmt.Println(CHDDirectory + " is not a vaild directory.")
					continue
				}

				// Iterate through the files in the sub directory, only for chd machine
				sublist, err := ioutil.ReadDir(CHDDirectory)
				CheckError(err)
				for _, subinfo := range sublist {
					CHDFileName := subinfo.Name()
					CHDFilePath := CHDDirectory + "/" + CHDFileName
					if !strings.HasSuffix(CHDFileName, ".chd") {
						fmt.Println(CHDFilePath + " is not a CHD file.")
						// idleFiles = append(idleFiles, directoryPath)
						continue
					}
					// check sha1 of chd file
					data, err := ioutil.ReadFile(CHDFilePath)
					CheckError(err)
					sha1 := fmt.Sprintf("%x", sha1.Sum(data))
					disk := machine.Disk(sha1)
					if disk == nil {
						continue
					}

					fmt.Println(disk.Name)
					fmt.Println(sha1)
					// this.addMachine(machine)
					// machineStatus = this.MachineStatus[directoryName]
				}

			case strings.HasSuffix(info.Name(), ".zip"):
				machineName := strings.TrimSuffix(info.Name(), ".zip")
				machineFilePath := dir + "/" + info.Name()
				machine := mame.Machine(machineName)
				CheckError(err)
				if machine == nil {
					continue
				}
				fmt.Println("\n--- " + machine.Name + ".zip ---")
				// Open a zip archive for reading.
				z, err := zip.OpenReader(machineFilePath)
				CheckError(err)

				// Iterate through the files in the archive
				for _, f := range z.File {
					crc := fmt.Sprintf("%x", f.CRC32)
					for i := 8 - len(crc); i > 0; i = i - 1 {
						crc = "0" + crc
					}
					// fmt.Println(crc)
					rom := machine.Rom(crc)
					if rom == nil {
						fmt.Printf("Contents of %s(crc:%s) is redundant file.\n", f.Name, crc)
						continue
					}
					rom.Availabl = true
				}

				z.Close()

				fmt.Println(machine.Roms)
			case strings.HasSuffix(info.Name(), ".7z"):
			default:
			}

		}
	}

	mame.Flush()
	return
}

func (machine *Machine) Rom(crc string) (rom *Rom) {
	for k, _ := range machine.Roms {
		if machine.Roms[k].Crc == crc {
			rom = &machine.Roms[k]
			return
		}
	}

	return
}

func (machine *Machine) Disk(sha1 string) (disk *Disk) {
	for k, _ := range machine.Disks {
		if machine.Disks[k].Sha1 == sha1 {
			disk = &machine.Disks[k]
			return
		}
	}

	return
}
