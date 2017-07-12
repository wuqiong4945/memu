package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func (machine Machine) Start() (result []byte) {
	f, err := os.Open(mamePath)
	CheckError(err)
	if os.IsNotExist(err) {
		return
	}
	f.Close()

	result, _ = exec.Command(mamePath, machine.Name).Output()
	return
}

func (machine Machine) GetStatusInfo() (info string) {
	machineStatus := fmt.Sprintf("%b", machine.MachineStatus)
	info += `<table border="1">`
	info += `<tr style="color:red">` +
		"<th>" + machine.Name + "</th>" +
		"<th>" + machineStatus + "</th>" +
		"</tr>"

	var ancesterMachineByRom func(Machine, Rom) Machine
	ancesterMachineByRom = func(machine Machine, rom Rom) Machine {
		upperMachine := machine.UpperMachine()
		if rom.Merge == "" ||
			upperMachine == nil ||
			upperMachine.Rom(rom.Crc) == nil {
			return machine
		}
		return ancesterMachineByRom(*upperMachine, rom)
	}
	for _, rom := range machine.Roms {
		m := ancesterMachineByRom(machine, rom)
		ancesterMachineName := m.Name
		if ancesterMachineName == machine.Name {
			ancesterMachineName = ""
		}
		info += "<tr>" +
			"<td>rom</td>" +
			"<td>" + rom.Name + "</td>" +
			"<td>" + fmt.Sprintf("%b", rom.RomStatus) + "</td>" +
			"<td>" + rom.Status + "</td>" +
			"<td>" + ancesterMachineName + "</td>" +
			"</tr>"
	}

	var ancesterMachineByDisk func(Machine, Disk) Machine
	ancesterMachineByDisk = func(machine Machine, disk Disk) Machine {
		upperMachine := machine.UpperMachine()
		if disk.Merge == "" ||
			upperMachine == nil ||
			upperMachine.Disk(disk.Sha1) == nil {
			return machine
		}
		return ancesterMachineByDisk(*upperMachine, disk)
	}
	for _, disk := range machine.Disks {
		m := ancesterMachineByDisk(machine, disk)
		ancesterMachineName := m.Name
		if ancesterMachineName == machine.Name {
			ancesterMachineName = ""
		}
		info += "<tr>" +
			"<td>disk</td>" +
			"<td>" + disk.Name + "</td>" +
			"<td>" + fmt.Sprintf("%b", disk.DiskStatus) + "</td>" +
			"<td>" + disk.Status + "</td>" +
			"<td>" + ancesterMachineName + "</td>" +
			"</tr>"
	}

	info += "</table>"
	return
}

func (machine Machine) GetHistoryInfo() (info string) {
	kind := "history"
	info = GetGeneralInfo(machine.Name, kind)
	info = Convert(info, kind)

	return
}
func (machine Machine) GetCommandInfo() (info string) {
	kind := "command"
	info = GetGeneralInfo(machine.Name, kind)
	info = Convert(info, kind)

	return
}
func (machine Machine) GetMameinfoInfo() (info string) {
	kind := "mameinfo"
	info = GetGeneralInfo(machine.Name, kind)
	info = Convert(info, kind)

	return
}

func GetGeneralInfo(machineName, kind string) (info string) {
	path := cfg.Section("general").Key(kind).MustString(kind + ".dat")
	file, err := os.Open(path)
	CheckError(err)
	if err != nil {
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	reader := bufio.NewReader(file)
	record := false

mainLoop:
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		line := string(lineBytes)
		switch {
		case strings.HasPrefix(line, "#"):
			continue

		case strings.HasPrefix(line, "$"):
			if !strings.Contains(line, "=") {
				continue
			}
			// reach another entry, stop recording
			if record {
				record = false
				break mainLoop // finished
			}

			if strings.HasPrefix(line, "$info=") {
				line = strings.TrimPrefix(line, "$info=")
				keys := strings.Split(line, ",")
				for _, key := range keys {
					// found the entry, start recording
					if key == machineName {
						record = true
						break
					}
				}
			}

		default:
			if record {
				buffer.WriteString(line)
				buffer.WriteString("<br>")
			}
		}

	}

	info = buffer.String()
	return
}

func Convert(info, kind string) string {
	type RegexpTable struct{ rgx, fmt string }
	var rgxTable []RegexpTable
	switch kind {
	case "history":
		rgxTable = []RegexpTable{
			{`<br>\s+`, "<br>"},
			// - *** -
			{`>\s*(-[^<>-]+-)\s*<`, "><font color='red'><b>$1</b></font><"},
		}

	case "mameinfo":
		rgxTable = []RegexpTable{
			{`<br>\s+`, "<br>"},
			// [***]
			{`>\s*([^<>:]+:)\s*<`, "><font color='red'><b>$1</b></font><"},
		}

	case "command":
		rgxTable = []RegexpTable{
			{`<br>\s+`, "<br>"},
			// directions, generate duplicated symbols
			{"_2_1_4_1_2_3_6", "<font color='blue'></font>"},
			{"_2_3_6_3_2_1_4", "<font color='blue'></font>"},
			{"_4_1_2_3_6", "<font color='blue'></font>"},
			{"_6_3_2_1_4", "<font color='blue'></font>"},
			{"_2_3_6", "<font color='blue'></font>"},
			{"_2_1_4", "<font color='blue'></font>"},
			{"_1", "<font color='green'></font>"},
			{"_2", "<font color='green'></font>"},
			{"_3", "<font color='green'></font>"},
			{"_4", "<font color='green'></font>"},
			// {"_5",  " " },
			{"_6", "<font color='green'></font>"},
			{"_7", "<font color='green'></font>"},
			{"_8", "<font color='green'></font>"},
			{"_9", "<font color='green'></font>"},
			{"_N", "<font color='green'><b>ℕ</b></font>"},
			// {R"(_(\d))",          "dir-$1.png" },
			// buttons
			{"_A", "<font color='green'><b>Ⓐ</b></font>"},
			{"_B", "<font color='green'><b>Ⓑ</b></font>"},
			{"_C", "<font color='green'><b>Ⓒ</b></font>"},
			{"_D", "<font color='green'><b>Ⓓ</b></font>"},
			{"_E", "<font color='green'><b>Ⓔ</b></font>"},
			{"_F", "<font color='green'><b>Ⓕ</b></font>"},
			{"_\\+", "<font color='red'>✚</font>"},
			{"_K", "<font color='green'><b>Ⓚ</b></font>"},
			{"_P", "<font color='green'><b>Ⓟ</b></font>"},
			// {`_([A-DGKNPS\+])`, "btn-$1.png" },
			// {`_([a-f])`,         "btn-n$1.png" },
			//  ------  ───
			{`<br>[─]{8,}\s*<br>`, "<hr>"},
			// [***]
			{`>\s*(\[[^\]<>]*\])`, "><font color='red'><b>$1</b></font>"},
			{`(^\s*\[[^\]<>]*\])`, "<font color='red'><b>$1</b></font>"},
			// special moves
			{"★", "<font color='gold'>★</font>"},
			{"☆", "<font color='silver'>☆</font>"},
			{"●", "<font color='yellow'>●</font>"},
			{"○", "<font color='orange'>○</font>"},
			{"◎", "<font color='red'>◎</font>"},
		}
	default:
	}

	var re *regexp.Regexp
	for _, table := range rgxTable {
		re = regexp.MustCompile(table.rgx)
		info = re.ReplaceAllString(info, table.fmt)
	}

	return info
}

func (machine *Machine) UpdateStatus() {
	if machine == nil {
		return
	}

	upperMachine := machine.UpperMachine()
	if machine.Romof != "" && upperMachine != nil {
		for k, rom := range machine.Roms {
			if rom.RomStatus&ROM_EXIST == ROM_EXIST || rom.Merge == "" {
				continue
			}
			upperRom := upperMachine.Rom(rom.Crc)
			if upperRom != nil {
				machine.Roms[k].RomStatus = upperRom.RomStatus
			}
		}
		for k, disk := range machine.Disks {
			if disk.DiskStatus&DISK_EXIST == DISK_EXIST || disk.Merge == "" {
				continue
			}
			upperDisk := upperMachine.Disk(disk.Sha1)
			if upperDisk != nil {
				machine.Disks[k].DiskStatus = upperDisk.DiskStatus
			}
		}
	}
	for _, rom := range machine.Roms {
		if rom.RomStatus&ROM_EXIST != ROM_EXIST && rom.Status != "nodump" {
			machine.MachineStatus &^= MACHINE_EXIST_V
			return
		}
	}
	for _, disk := range machine.Disks {
		if disk.DiskStatus&DISK_EXIST != DISK_EXIST && disk.Status != "nodump" {
			machine.MachineStatus &^= MACHINE_EXIST_V
			return
		}
	}

	if machine.MachineStatus&MACHINE_EXIST == MACHINE_EXIST {
		machine.MachineStatus |= MACHINE_EXIST_V
		machine.MachineStatus &^= MACHINE_EXIST_P
	}
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

func (machine Machine) UpperMachine() (upperMachine *Machine) {
	if machine.Romof == "" {
		return
	}

	upperMachine = mame.Machine(machine.Romof)
	return
}
