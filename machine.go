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
	<li class="nav-item"><a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_Roms" role="tab">Roms</a></li>
	<li class="nav-item"><a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_History" role="tab">History</a></li>
	<li class="nav-item"><a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_Command" role="tab">Command</a></li>
	<li class="nav-item"><a class="nav-link" data-toggle="tab" href="#` + machine.Name + `_None" role="tab">None</a></li>
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

	// info += `<div class="card-block">`
	info += `	<div class="tab-content">`
	info += "\n"

	info += `		<div class="tab-pane" id="` + machine.Name + `_Roms" role="tabpanel">`
	info += machine.GetRomInfo()
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

func (machine Machine) GetRomInfo() (info string) {
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
			"<td>" + rom.Crc + "</td>" +
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
			"<td> Sha1 </td>" +
			// "<td>" + disk.Sha1 + "</td>" +
			"<td>" + fmt.Sprintf("%b", disk.DiskStatus) + "</td>" +
			"<td>" + disk.Status + "</td>" +
			"<td>" + ancesterMachineName + "</td>" +
			"</tr>"
		info += "\n"
	}

	info += `			</table>`
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
	s := GetGeneralInfo(machine.Name, kind)

	info += `		<div id="accordion" role="tablist" aria-multiselectable="true">`
	info += "\n"
	var id string
	reader := bufio.NewReader(strings.NewReader(s))
	k := 0
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(lineBytes)

		switch {
		case line == "$cmd":
			k++
			titleBytes, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			title := string(titleBytes)

			id = fmt.Sprintf(machine.Name+"%d", k)
			// info += `
			// <div class="panel panel-default">
			// <div class="panel-heading" role="tab" id="heading_` + id + `">
			// <h6 class="panel-title">
			// <a data-toggle="collapse" data-parent="#accordion" href="#` + id + `" aria-expanded="true" aria-controls="` + id + `">
			// ` + title + `
			// </a>
			// </h6>
			// </div>
			// <div id="` + id + `" class="panel-collapse collapse in" role="tabpanel" aria-labelledby="heading_` + id + `">
			// `
			info += `
				<div class="card">
					<div class="card-header" role="tab" id="heading_` + id + `">
						<h6 class="mb-0">
							<a data-toggle="collapse" data-parent="#accordion" href="#` + id + `" aria-expanded="true" aria-controls="` + id + `">
								` + title + `
							</a>
						</h6>
					</div>
					<div id="` + id + `" class="collapse" role="tabpanel" aria-labelledby="heading_` + id + `">
						<div class="card-block">
							<dl class="row">
				`

			_, _, err = reader.ReadLine()
			if err == io.EOF {
				info += `</dl></div></div></div>`
				break
			}

		case line == "$end":
			info += `</dl>`
			info += `</div></div></div>`

		default:
			// info += `<p>` + line + `</p>`
			// info += line + "<br/>"
			// line = strings.TrimSpace(line)
			n := strings.Index(line, "     ")
			// if strings.Contains(line, "     ") {
			if n > 1 {
				dt := strings.TrimSpace(line[0:n])
				dd := strings.TrimSpace(line[n:len(line)])
				info += `<dt class="col-sm-5"><small>` + dt + `</small></dt>` +
					`<dd class="col-sm-7 text-right"><small>` + dd + `</small></dd>`
			} else {
				info += `<dt class="col-sm-12"><small>` + line + `</small></dt>`
			}
			info += "\n"
		}
	}
	info += `		</div>`

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
				if record == true {
					buffer.WriteString(line + "\n")
				}
				continue
			}
			// reach another entry, stop recording
			if record == true {
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
				buffer.WriteString(line + "\n")
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
			// - *** -
			{`>\s*(-[^<>-]+-)\s*<`, `><font color='red'><b>$1</b></font><`},

			{`\s*\$\w+\s*\n`, ``},
			{`\n`, `<br/>`},
		}

	case "mameinfo":
		rgxTable = []RegexpTable{
			// [***]
			{`>\s*([^<>:]+:)\s*<`, `><font color='red'><b>$1</b></font><`},

			{`\s*\$\w+\s*\n`, ``},
			{`\n`, `<br/>`},
		}

	case "command":
		rgxTable = []RegexpTable{
			{`_6_3_2_1_4_1_2_3_6`, `<img width='32' height='32' src='data/icons/63214.svg'/><img width='32' height='32' src='data/icons/41236.svg'/>`},
			{`_4_1_2_3_6_3_2_1_4`, `<img width='32' height='32' src='data/icons/41236.svg'/><img width='32' height='32' src='data/icons/63214.svg'/>`},

			// directions, generate duplicated symbols
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
			{`_([a-zA-DGKPS])`, `<font><kbd>$1</kbd></font>`},
			{`_\+`, `<font color='red'>✚</font>`},
			//  ------  ───
			// {`<br>[─]{8,}\s*<br>`, `<hr>`},
			// [***]
			{`>\s*(\[[^\]<>]*\])`, `><font color='red'><b>$1</b></font>`},
			{`(^\s*\[[^\]<>]*\])`, `<font color='red'><b>$1</b></font>`},
			// special moves
			{`★`, `<mark>★</mark>`},
			{`☆`, `<mark>☆</mark>`},
			{`●`, `<mark>●</mark>`},
			{`○`, `<mark>○</mark>`},
			{`◎`, `<mark>◎</mark>`},

			// {`\n`, `<br/>`},
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
