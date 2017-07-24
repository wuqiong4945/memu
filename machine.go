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
			{`<br>\s+`, `<br>`},
			// - *** -
			{`>\s*(-[^<>-]+-)\s*<`, `><font color='red'><b>$1</b></font><`},
		}

	case "mameinfo":
		rgxTable = []RegexpTable{
			{`<br>\s+`, `<br>`},
			// [***]
			{`>\s*([^<>:]+:)\s*<`, `><font color='red'><b>$1</b></font><`},
		}

	case "command":
		rgxTable = []RegexpTable{
			{`<br>\s+`, `<br>`},
			// directions, generate duplicated symbols
			// {`_2_1_4_1_2_3_6`, `<img width='32' height='32' src='data/icons/bl.svg'/><img width='32' height='32' src='data/icons/lbr.svg'/>`},
			// {`_2_3_6_3_2_1_4`, `<img width='32' height='32' src='data/icons/br.svg'/><img width='32' height='32' src='data/icons/rbl.svg'/>`},

			{`_4_1_2_3_6`, `<img width='32' height='32' src='data/icons/41236.svg'/>`},
			{`_6_3_2_1_4`, `<img width='32' height='32' src='data/icons/63214.svg'/>`},
			{`_4_7_8_9_6`, `<img width='32' height='32' src='data/icons/47896.svg'/>`},
			{`_6_9_8_7_4`, `<img width='32' height='32' src='data/icons/69874.svg'/>`},
			{`_6_3_2_3`, `<img width='32' height='32' src='data/icons/632.svg'/><img width='32' height='32' src='data/icons/236.svg'/>`}, // obscure
			{`_2_3_6_3`, `<img width='32' height='32' src='data/icons/236.svg'/><img width='32' height='32' src='data/icons/632.svg'/>`}, // obscure
			{`_4_1_2_1`, `<img width='32' height='32' src='data/icons/412.svg'/><img width='32' height='32' src='data/icons/214.svg'/>`}, // obscure
			{`_2_1_4_1`, `<img width='32' height='32' src='data/icons/214.svg'/><img width='32' height='32' src='data/icons/412.svg'/>`}, // obscure

			// for good looking
			{`_6_4_6_4`, `<img width='32' height='32' src='data/icons/64.svg'/><img width='32' height='32' src='data/icons/64.svg'/>`},
			{`_8_2_8_2`, `<img width='32' height='32' src='data/icons/82.svg'/><img width='32' height='32' src='data/icons/82.svg'/>`},

			{`_2_3_6`, `<img width='32' height='32' src='data/icons/236.svg'/>`},
			{`_6_3_2`, `<img width='32' height='32' src='data/icons/632.svg'/>`},
			{`_2_1_4`, `<img width='32' height='32' src='data/icons/214.svg'/>`},
			{`_4_1_2`, `<img width='32' height='32' src='data/icons/412.svg'/>`},

			{`_2_2_2`, `<img width='32' height='32' src='data/icons/222.svg'/>`},
			{`_4_4_4`, `<img width='32' height='32' src='data/icons/444.svg'/>`},
			{`_6_6_6`, `<img width='32' height='32' src='data/icons/666.svg'/>`},
			{`_8_8_8`, `<img width='32' height='32' src='data/icons/888.svg'/>`},

			// partly command
			{`_4_2_6`, `<img width='32' height='32' src='data/icons/41236.svg'/>`},
			{`_6_2_4`, `<img width='32' height='32' src='data/icons/63214.svg'/>`},
			{`_6_2_3`, `<img width='32' height='32' src='data/icons/632.svg'/><img width='32' height='32' src='data/icons/236.svg'/>`}, // obscure
			{`_2_6_3`, `<img width='32' height='32' src='data/icons/236.svg'/><img width='32' height='32' src='data/icons/632.svg'/>`}, // obscure
			{`_4_2_1`, `<img width='32' height='32' src='data/icons/412.svg'/><img width='32' height='32' src='data/icons/214.svg'/>`}, // obscure
			{`_2_4_1`, `<img width='32' height='32' src='data/icons/214.svg'/><img width='32' height='32' src='data/icons/412.svg'/>`}, // obscure

			{`_1_2_3`, `<img width='32' height='32' src='data/icons/1.svg'/><img width='32' height='32' src='data/icons/2.svg'/><img width='32' height='32' src='data/icons/3.svg'/>`},
			{`_4_6_6`, `<img width='32' height='32' src='data/icons/4.svg'/><img width='32' height='32' src='data/icons/66.svg'/>'/>`},

			{`_2_3`, `<img width='32' height='32' src='data/icons/236.svg'/>`}, // obscure
			{`_6_3`, `<img width='32' height='32' src='data/icons/632.svg'/>`}, // obscure
			{`_2_1`, `<img width='32' height='32' src='data/icons/214.svg'/>`}, // obscure
			{`_4_1`, `<img width='32' height='32' src='data/icons/412.svg'/>`}, // obscure

			{`_2_2`, `<img width='32' height='32' src='data/icons/22.svg'/>`},
			{`_4_4`, `<img width='32' height='32' src='data/icons/44.svg'/>`},
			{`_6_6`, `<img width='32' height='32' src='data/icons/66.svg'/>`},
			{`_8_8`, `<img width='32' height='32' src='data/icons/88.svg'/>`},

			{`_4_6`, `<img width='32' height='32' src='data/icons/46.svg'/>`},
			{`_6_4`, `<img width='32' height='32' src='data/icons/64.svg'/>`},
			{`_2_8`, `<img width='32' height='32' src='data/icons/28.svg'/>`},
			{`_8_2`, `<img width='32' height='32' src='data/icons/82.svg'/>`},

			// partly command
			{`_2_6`, `<img width='32' height='32' src='data/icons/236.svg'/>`},
			{`_6_2`, `<img width='32' height='32' src='data/icons/632.svg'/>`},
			{`_2_4`, `<img width='32' height='32' src='data/icons/214.svg'/>`},
			{`_4_2`, `<img width='32' height='32' src='data/icons/412.svg'/>`},

			{`_([1-9N])`, `<img width='32' height='32' src='data/icons/$1.svg'/>`},
			// buttons
			{`_([a-fA-DGKPS])`, `<font><kbd>$1</kbd></font>`},
			{`_\+`, `<font color='red'>✚</font>`},
			//  ------  ───
			{`<br>[─]{8,}\s*<br>`, `<hr>`},
			// [***]
			{`>\s*(\[[^\]<>]*\])`, `><font color='red'><b>$1</b></font>`},
			{`(^\s*\[[^\]<>]*\])`, `<font color='red'><b>$1</b></font>`},
			// special moves
			{`★`, `<font style='color:white;background-color:red'>★</font>`},
			{`☆`, `<font style='color:white;background-color:silver'>☆</font>`},
			{`●`, `<font style='color:white;background-color:yellwo'>●</font>`},
			{`○`, `<font style='color:white;background-color:orange'>○</font>`},
			{`◎`, `<font style='color:white;background-color:red'>◎</font>`},
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
