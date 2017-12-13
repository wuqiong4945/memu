package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// NewMame factory Mame struct
func NewMame() (mame *Mame) {
	mame = new(Mame)

	cacheFile, err := os.OpenFile("cache.gob", os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer cacheFile.Close()
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in open gob file : %s\n", err)
	}
	if err == nil {
		dec := gob.NewDecoder(cacheFile)
		err = dec.Decode(mame)
		CheckError(err)
		if err != nil || mame.Build == "" {
			fmt.Printf("Error in decoding gob : %s\n", err)
		}
	}

	mame.Update()
	return
}

// Audit roms
func (mame *Mame) Audit() {
	t := time.Now()
	// defer fmt.Printf("Audit time: %s\n", time.Now().Sub(t).String())

	// reset status
	for k := range mame.Machines {
		mame.Machines[k].MachineStatus = 0
		for m := range mame.Machines[k].Roms {
			mame.Machines[k].Roms[m].RomStatus = 0
		}
		for m := range mame.Machines[k].Disks {
			mame.Machines[k].Disks[m].DiskStatus = 0
		}
	}

	romDefaultPath := "roms," + mamePath[0:strings.LastIndex(mamePath, "/")] + "/roms"
	romPath := cfg.Section("general").Key("rompath").MustString(romDefaultPath)
	//fullPath, _ := filepath.Abs(path)
	// idleFiles := []string{}
	dirs := strings.Split(romPath, ",")
	for _, dir := range dirs {
		list, err := ioutil.ReadDir(dir)
		CheckError(err)
		fmt.Println("\n===== " + dir + " =====")

		// Iterate through the files in the directory
		for _, info := range list {
			switch {
			// dispose machine in chd file
			case info.IsDir() == true:
				// just chd game has sub dir
				mame.AuditCHDFolder(dir, info.Name())

			case strings.HasSuffix(info.Name(), ".zip"):
				mame.AuditZipFile(dir, info.Name())

			case strings.HasSuffix(info.Name(), ".7z"):
				mame.Audit7zFile(dir, info.Name())

			default:
				fmt.Println("\n--- " + info.Name() + " ---")
				fmt.Println(dir + "/" + info.Name() + " is not a vaild file.")
			}

		}
		fmt.Println("\n==========")
	}

	mame.UpdateAllMachineStatus()
	isFlush = true

	fmt.Printf("Audit time: %s\n", time.Now().Sub(t).String())
	return
}

// Audit7zFile audit rom in 7z format
// TODO: implement it
func (mame *Mame) Audit7zFile(dir, fileName string) {
	if !strings.HasSuffix(fileName, ".7z") {
		return
	}

	fmt.Println("\n--- " + fileName + " ---")
	machineName := strings.TrimSuffix(fileName, ".7z")
	machineFilePath := dir + "/" + fileName

	machine := mame.Machine(machineName)
	if machine == nil {
		fmt.Println(machineFilePath + " is not a vaild file.")
		return
	}
}

// AuditCHDFolder audit CHD rom
func (mame *Mame) AuditCHDFolder(dir, folderName string) {
	fmt.Println("\n--- " + folderName + " ---")
	CHDMachineName := folderName
	CHDFolderPath := dir + "/" + CHDMachineName

	machine := mame.Machine(CHDMachineName)
	if machine == nil {
		fmt.Println(CHDFolderPath + " is not a vaild directory.")
		return
	}
	machine.MachineStatus |= MACHINE_EXIST

	// Iterate through the files in the chd directory
	list, err := ioutil.ReadDir(CHDFolderPath)
	CheckError(err)
	for _, info := range list {
		CHDFileName := info.Name()
		CHDFilePath := CHDFolderPath + "/" + CHDFileName
		if !strings.HasSuffix(CHDFileName, ".chd") {
			fmt.Println(CHDFilePath + " is not a CHD file.")
			continue
		}

		// check sha1 of chd file
		data, err := ioutil.ReadFile(CHDFilePath)
		CheckError(err)
		sha1 := fmt.Sprintf("%x", sha1.Sum(data))
		disk := machine.Disk(sha1)
		if disk == nil {
			fmt.Printf("%s(sha1:%s) is redundant file.\n", CHDFileName, sha1)
			machine.MachineStatus |= MACHINE_EXIST_R
			continue
		}

		machine.MachineStatus |= MACHINE_EXIST_P
		disk.DiskStatus |= DISK_EXIST
		diskName := CHDFileName[0:strings.LastIndex(CHDFileName, ".")]
		if diskName != disk.Name {
			disk.DiskStatus |= DISK_EXIST_WN
		}
	}
}

// AuditZipFile audit rom in zip format
func (mame *Mame) AuditZipFile(dir, fileName string) {
	if !strings.HasSuffix(fileName, ".zip") {
		return
	}

	fmt.Println("\n--- " + fileName + " ---")
	machineName := strings.TrimSuffix(fileName, ".zip")
	machineFilePath := dir + "/" + fileName

	machine := mame.Machine(machineName)
	if machine == nil {
		fmt.Println(machineFilePath + " is not a vaild file.")
		return
	}
	machine.MachineStatus |= MACHINE_EXIST

	// Open zip file for reading.
	z, err := zip.OpenReader(machineFilePath)
	CheckError(err)
	// Iterate through the files in zip file
	for _, f := range z.File {
		crc := fmt.Sprintf("%x", f.CRC32)
		for i := 8 - len(crc); i > 0; i = i - 1 {
			crc = "0" + crc
		}
		// fmt.Println(crc)
		rom := machine.Rom(crc)
		if rom == nil {
			fmt.Printf("%s(crc:%s) is redundant file.\n", f.Name, crc)
			machine.MachineStatus |= MACHINE_EXIST_R
			continue
		}

		machine.MachineStatus |= MACHINE_EXIST_P
		rom.RomStatus |= ROM_EXIST
		if f.Name != rom.Name {
			rom.RomStatus |= ROM_EXIST_WN
		}
	}

	z.Close()
}

// Machine gets machine by name
func (mame Mame) Machine(machineName string) (machine *Machine) {
	for k := range mame.Machines {
		if mame.Machines[k].Name == machineName {
			machine = &mame.Machines[k]
			return
		}
	}

	return
}

// Fresh exports mame xml info
func (mame *Mame) Fresh() {
	out, err := exec.Command(mamePath, "-listxml").Output()
	CheckError(err)
	err = xml.Unmarshal(out, mame)
	CheckError(err)

	// delete repeat rom
	for k, machine := range mame.Machines {
		var roms []Rom
		for _, rom := range machine.Roms {
			isRepeatRom := false
			for _, r := range roms {
				if rom.Name == r.Name ||
					rom.Crc == r.Crc {
					isRepeatRom = true
					break
				}
			}
			if isRepeatRom == false {
				roms = append(roms, rom)
			}
		}
		mame.Machines[k].Roms = roms
	}

	isFlush = true
	return
}

// Flush exports mame info to gob file
func (mame Mame) Flush() {
	cacheFile, err := os.OpenFile("cache.gob", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer cacheFile.Close()
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in opening gob file : %s\n", err)
		return
	}

	enc := gob.NewEncoder(cacheFile)
	err = enc.Encode(mame)
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in encoding gob : %s\n", err)
	}

	return
}

// Search machine in name and description
func (mame Mame) Search(key string) (machines []Machine, err error) {
	reg, err := regexp.Compile(key)
	if err != nil {
		CheckError(err)
		return
	}

	for _, machine := range mame.Machines {
		if reg.FindStringIndex(machine.Description) != nil ||
			reg.FindStringIndex(machine.Name) != nil {
			machines = append(machines, machine)
			continue
		}
	}

	return
}

// Update checks mame version, and fresh mame info
func (mame *Mame) Update() {
	if mame.Build == "" {
		mame.Fresh()
		return
	}

	version := mame.Version()
	l := len(mame.Build)
	if len(version) > l {
		version = version[len(version)-l : len(version)]
	}
	if version > mame.Build {
		mame.Fresh()
	}

	return
}

func (mame *Mame) UpdateAllMachineStatus() {
	// first : deal with bios/device machine
	for k := range mame.Machines {
		machine := &mame.Machines[k]
		switch {
		case machine.MachineStatus&MACHINE_EXIST != MACHINE_EXIST:
			continue
		case machine.Isbios == "yes" || machine.Isdevice == "yes":
			machine.UpdateStatus()
		}
	}
	// second : deal with major machine
	for k := range mame.Machines {
		machine := &mame.Machines[k]
		if machine.Cloneof == "" &&
			machine.Isbios != "yes" &&
			machine.Isdevice != "yes" {
			machine.UpdateStatus()
		}
	}
	// third : deal with clone machine
	for k := range mame.Machines {
		machine := &mame.Machines[k]
		if machine.Cloneof != "" &&
			machine.Isbios != "yes" &&
			machine.Isdevice != "yes" {
			machine.UpdateStatus()
		}
	}

}

// VerifyRoms uses native mame verifyroms
func (mame Mame) VerifyRoms(machineName string) (result []byte) {
	result, _ = exec.Command(mamePath, "-verifyroms", machineName).Output()
	return
}

// Version uses mame help order get version information
func (mame Mame) Version() (version string) {
	out, _ := exec.Command(mamePath, "-help").Output()
	// b := bytes.NewBuffer(out)
	// b.ReadString(' ')
	// version, _ = b.ReadString(' ')
	v, _, _ := bufio.NewReader(bytes.NewBuffer(out)).ReadLine()
	version = string(v)
	return
}
