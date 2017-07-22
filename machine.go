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

	cardType := "card"
	switch {
	case machine.Isbios == "yes" ||
		machine.Isdevice == "yes":
		cardType += " card-warning"
	case machine.MachineStatus&MACHINE_EXIST_V == MACHINE_EXIST_V:
		cardType += " card-info"
	default:
		cardType += " card-danger"
	}

	info += "\n"
	// info += `<div class="col-sm-3">`
	info += `<div class="` + cardType + `">`
	// info += `<div class="card">`
	info += "\n"

	// header
	machineStatus := fmt.Sprintf("%b", machine.MachineStatus)
	info += `	<div class="card-header">` + machine.Name + " (" + machineStatus + ")" + ``
	info += `		<ul class="nav nav-tabs card-header-tabs" role="tablist">
      <li class="nav-item">
        <a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_Roms" role="tab">Roms</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_History" role="tab">History</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_Command" role="tab">Command</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_None" role="tab">None</a>
      </li>
		</ul>`
	info += "\n"
	info += `	</div>`
	info += "\n"

	// image
	picName := "snap/" + machine.Name + ".png"
	_, err := os.Stat(picName)
	switch {
	case !os.IsNotExist(err):
		info += `	<img class="card-img-top img-fluid" src="` + picName + `" alt="` + machine.Name + `">`
	case os.IsNotExist(err) && machine.Cloneof != "":
		picName := "snap/" + machine.Cloneof + ".png"
		if _, err := os.Stat(picName); !os.IsNotExist(err) {
			info += `	<img class="card-img-top img-fluid" src="` + picName + `" alt="` + machine.Name + `">`
		}
	}
	info += "\n"

	info += `	<div class="tab-content">`
	info += "\n"
	// info += `<div class="card-block">`
	info += `		<div class="tab-pane" id="` + machine.Name + `_Roms" role="tabpanel">`
	info += "\n"
	// block
	info += `			<table class="table table-striped table-sm">`
	info += "\n"
	for _, rom := range machine.Roms {
		m := ancesterMachineByRom(machine, rom)
		ancesterMachineName := m.Name
		if ancesterMachineName == machine.Name {
			ancesterMachineName = ""
		}
		info += "				<tr>" +
			"<td>rom</td>" +
			"<td>" + rom.Name + "</td>" +
			"<td>" + fmt.Sprintf("%b", rom.RomStatus) + "</td>" +
			"<td>" + rom.Status + "</td>" +
			"<td>" + ancesterMachineName + "</td>" +
			"</tr>"
		info += "\n"
	}

	for _, disk := range machine.Disks {
		m := ancesterMachineByDisk(machine, disk)
		ancesterMachineName := m.Name
		if ancesterMachineName == machine.Name {
			ancesterMachineName = ""
		}
		info += "				<tr>" +
			"<td>disk</td>" +
			"<td>" + disk.Name + "</td>" +
			"<td>" + fmt.Sprintf("%b", disk.DiskStatus) + "</td>" +
			"<td>" + disk.Status + "</td>" +
			"<td>" + ancesterMachineName + "</td>" +
			"</tr>"
		info += "\n"
	}

	info += `			</table>`
	info += "\n"
	info += `		</div>`
	info += "\n"

	info += `		<div class="tab-pane" id="` + machine.Name + `_History" role="tabpanel">`
	info += machine.GetHistoryInfo()
	info += `		</div>`
	info += "\n"

	info += `		<div class="tab-pane" id="` + machine.Name + `_Command" role="tabpanel">`
	info += machine.GetCommandInfo()
	info += `		</div>`
	info += "\n"

	info += `		<div class="tab-pane" id="` + machine.Name + `_None" role="tabpanel">`
	info += `		</div>`
	info += "\n"

	// info += `</div>` // block
	info += `	</div>`
	info += "\n"
	// info += `	</div>`

	// footer
	info += `	<div class="card-footer text-muted text-right">` + machineStatus + `</div>`
	info += "\n"

	// info += `</div>`
	info += `</div>`
	info += "\n"
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
			// {"_2_1_4_1_2_3_6", "<img width='32' height='32' src='data/icons/bl.svg'/><img width='32' height='32' src='data/icons/lbr.svg'/>"},
			// {"_2_3_6_3_2_1_4", "<img width='32' height='32' src='data/icons/br.svg'/><img width='32' height='32' src='data/icons/rbl.svg'/>"},

			{"_4_1_2_3_6", "<img width='32' height='32' src='data/icons/lbr.svg'/>"},
			{"_6_3_2_1_4", "<img width='32' height='32' src='data/icons/rbl.svg'/>"},
			{"_6_3_2_3", "<img width='32' height='32' src='data/icons/rb.svg'/><img width='32' height='32' src='data/icons/br.svg'/>"}, // obscure
			{"_4_1_2_1", "<img width='32' height='32' src='data/icons/lb.svg'/><img width='32' height='32' src='data/icons/bl.svg'/>"}, // obscure

			{"_2_3_6", "<img width='32' height='32' src='data/icons/br.svg'/>"},
			{"_6_3_2", "<img width='32' height='32' src='data/icons/rb.svg'/>"},
			{"_2_1_4", "<img width='32' height='32' src='data/icons/bl.svg'/>"},
			{"_4_1_2", "<img width='32' height='32' src='data/icons/lb.svg'/>"},

			// partly command
			{"_4_2_6", "<img width='32' height='32' src='data/icons/lbr.svg'/>"},
			{"_6_2_4", "<img width='32' height='32' src='data/icons/rbl.svg'/>"},
			{"_6_2_3", "<img width='32' height='32' src='data/icons/rb.svg'/><img width='32' height='32' src='data/icons/br.svg'/>"}, // obscure
			{"_4_2_1", "<img width='32' height='32' src='data/icons/lb.svg'/><img width='32' height='32' src='data/icons/bl.svg'/>"}, // obscure

			{"_1_2_3", "<img width='32' height='32' src='data/icons/bol.svg'/><img width='32' height='32' src='data/icons/b.svg'/><img width='32' height='32' src='data/icons/rob.svg'/>"},
			{"_2_3", "<img width='32' height='32' src='data/icons/br.svg'/>"}, // obscure
			{"_2_1", "<img width='32' height='32' src='data/icons/bl.svg'/>"}, // obscure

			{"_2_8", "<img width='32' height='32' src='data/icons/bu.svg'/>"},
			{"_8_2", "<img width='32' height='32' src='data/icons/ub.svg'/>"},
			{"_6_4", "<img width='32' height='32' src='data/icons/rl.svg'/>"},
			{"_4_6", "<img width='32' height='32' src='data/icons/lr.svg'/>"},
			{"_6_6", "<img width='32' height='32' src='data/icons/rr.svg'/>"},
			{"_4_4", "<img width='32' height='32' src='data/icons/ll.svg'/>"},
			{"_2_2", "<img width='32' height='32' src='data/icons/bb.svg'/>"},
			{"_8_8", "<img width='32' height='32' src='data/icons/uu.svg'/>"},

			// partly command
			{"_2_6", "<img width='32' height='32' src='data/icons/br.svg'/>"},
			{"_6_2", "<img width='32' height='32' src='data/icons/rb.svg'/>"},
			{"_2_4", "<img width='32' height='32' src='data/icons/bl.svg'/>"},
			{"_4_2", "<img width='32' height='32' src='data/icons/lb.svg'/>"},

			{"_1", "<img width='32' height='32' src='data/icons/bol.svg'/>"},
			{"_2", "<img width='32' height='32' src='data/icons/b.svg'/>"},
			{"_3", "<img width='32' height='32' src='data/icons/rob.svg'/>"},
			{"_4", "<img width='32' height='32' src='data/icons/l.svg'/>"},
			// {"_5",  " " },
			{"_6", "<img width='32' height='32' src='data/icons/r.svg'/>"},
			{"_7", "<img width='32' height='32' src='data/icons/lou.svg'/>"},
			{"_8", "<img width='32' height='32' src='data/icons/u.svg'/>"},
			{"_9", "<img width='32' height='32' src='data/icons/uor.svg'/>"},
			{"_N", "<font color='green'><b>ℕ</b></font>"},
			// {R"(_(\d))",          "dir-$1.png" },
			// buttons
			{"_A", "<font color='green'><kbd>A</kbd></font>"},
			{"_B", "<font color='green'><kbd>B</kbd></font>"},
			{"_C", "<font color='green'><kbd>C</kbd></font>"},
			{"_D", "<font color='green'><kbd>D</kbd></font>"},
			{"_E", "<font color='green'><kbd>E</kbd></font>"},
			{"_F", "<font color='green'><kbd>F</kbd></font>"},
			{"_\\+", "<font color='red'>✚</font>"},
			{"_K", "<font color='green'><kbd>K</kbd></font>"},
			{"_P", "<font color='green'><kbd>P</kbd></font>"},
			{"_S", "<font color='green'><kbd>S</kbd></font>"},
			// {`_([A-DGKNPS\+])`, "btn-$1.png" },
			// {`_([a-f])`,         "btn-n$1.png" },
			//  ------  ───
			{`<br>[─]{8,}\s*<br>`, "<hr>"},
			// [***]
			{`>\s*(\[[^\]<>]*\])`, "><font color='red'><b>$1</b></font>"},
			{`(^\s*\[[^\]<>]*\])`, "<font color='red'><b>$1</b></font>"},
			// special moves
			{"★", "<font color='gold'><kbd>★</kbd></font>"},
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
