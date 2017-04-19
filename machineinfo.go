package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func GetInfo(machineName, kind string) (info string) {
	path := cfg.Section("general").Key(kind).MustString(kind + ".dat")
	info = GetGeneralInfo(machineName, path)
	info = Convert(info, kind)

	html, _ := os.OpenFile("info.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer html.Close()
	html.WriteString(info)

	return
}

func GetGeneralInfo(machineName, path string) (info string) {
	file, err := os.Open(path)
	CheckError(err)
	if err != nil {
		fmt.Printf("Error in open file : %s\n", err)
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	reader := bufio.NewReader(file)
	record := false
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		line := string(lineBytes)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "$") {
			if strings.Contains(line, "=") {
				// reach another entry, stop recording
				if record {
					record = false
					break // finished
				}

				if strings.HasPrefix(line, "$info=") {
					line = strings.TrimPrefix(line, "$info=")
					keys := strings.Split(line, ",")
					for _, key := range keys {
						if key == machineName { // found the entry, start recording
							record = true
							break
						}
					}
				}
			} else {
				continue
			}
		} else if record {
			buffer.WriteString(line)
			buffer.WriteString("<br>")
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
