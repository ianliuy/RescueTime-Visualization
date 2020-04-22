package main

import (
"bufio"
"fmt"
"github.com/AbhiAgarwal/go-rescuetime"
"log"
"os"
"strconv"
"strings"
"time"
)

func NewAnalyticDataQueryParameters(Perspective, ResolutionTime, RestrictGroup, RestrictBegin, RestrictEnd, RestrictKind, RestrictThing, RestrictThingy string) rescuetime.AnalyticDataQueryParameters {
	var adqp rescuetime.AnalyticDataQueryParameters
	if Perspective == ""{
		adqp.Perspective = "rank"
	} else {
		adqp.Perspective = Perspective
	}
	if ResolutionTime == ""{
		adqp.ResolutionTime = "day"
	} else {
		adqp.ResolutionTime = ResolutionTime
	}
	if RestrictBegin == "" {
		adqp.RestrictBegin = time.Now().Format("2006-01-02") // Today
	} else {
		adqp.RestrictBegin = RestrictBegin
	}
	if RestrictEnd == "" {
		adqp.RestrictEnd = time.Now().Format("2006-01-02")
	} else {
		adqp.RestrictEnd = RestrictEnd
	}
	if RestrictKind == "" {
		adqp.RestrictKind = "activity"
	} else {
		adqp.RestrictKind = RestrictKind
	}
	return adqp
}

func sec2hour(sec int) string {
	s := ""
	h := sec / 3600
	if h != 0 {
		s = s + strconv.Itoa(h) + "h"
	}
	m := (sec % 3600) / 60
	if m != 0 {
		s = s + strconv.Itoa(m) + "m"
	}
	second := sec % 60
	if !strings.Contains(s, "h") {
		s = s + strconv.Itoa(second) + "s"
	}
	return s
}

func getToday(data *rescuetime.AnalyticData) []string {
	var cont []string
	t := time.Now().Format("2006-01-02 15:04")
	hrow := "## yiyangiliu " + t[:10] + " Detailed Activaties, "
	shrow := "Updated at " + t[11:]
	frow := "|Rank|Activity|Time|Category|Label|"
	srow := "|-|-|-|-|-|"
	var trow []string
	for i, row := range data.Rows{
		if i == 15 {break}
		t := sec2hour(row.TimeSpentSeconds)
		act := row.Activity
		if len(act) > 14 {
			act = act[:14]
		}
		categ := row.Category
		switch categ {
		case "Editing & IDEs":
			categ = "IDE"
		case "General Software Development":
			categ = "Dev"
		case "General Social Networking":
			categ = "SNS"
		case "General Reference & Learning":
			categ = "Ref&Learn"
		case "Internet Utilities":
			categ = "Utils"
		case "General Utilities":
			categ = "Utils"
		case "Uncategorized":
			categ = "Unknown"
		case "Presentation":
			categ = "Pre"
		}
		if len(categ) > 12 {
			categ = categ[:12]
		}
		r :="|" + strconv.Itoa(row.Rank) + "|" + act + "|" + t + "|" + categ + "|" + strconv.Itoa(row.Productivity)+ "|"
		trow = append(trow, r)
		//fmt.Println(string()
	}
	cont = append(cont, hrow, "", shrow, "", frow, srow)
	cont = append(cont, trow...)
	return cont
}

func NewRescueTime(api_key string) *rescuetime.RescueTime {
	var rt rescuetime.RescueTime
	rt.APIKey = api_key
	return &rt
}

func getHistory(fpath string) []string {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	Scanner := bufio.NewScanner(f)
	Scanner.Split(bufio.ScanLines)
	var lines []string

	for Scanner.Scan() {
		lines = append(lines, Scanner.Text())
	}

	err = f.Close()

	if err != nil {
		log.Fatalf("failed to close file: %s", err)
	}
	return lines
}

// writeLines writes the lines to the given file.
func writef(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func getContent(t, h []string) []string {
	var cont []string
	cont = append(cont, h[:3]...)
	cont = append(cont, t...)
	cont = append(cont, "")
	cont = append(cont, h[3:]...)
	return cont
}

func main() {
	nrt := NewRescueTime("B63IavC02qsRZ4QZjl7lURlX6wiV_D_m9Z4ReXvR")
	a := NewAnalyticDataQueryParameters("",
		"",
		"",
		"", //2020-04-21
		"", //2020-04-21
		"",
		"",
		"")
	data, _ := nrt.GetAnalyticData("local", &a)
	today :=  getToday(&data)
	history := getHistory("C:/SakilaGithub/RescueTime-Record/README.md")

	hd := history[3][14:24]
	td := time.Now().Format("2006-01-02")
	if td == hd {fmt.Println("No need to update")
	} else {}

	cont := getContent(today, history)
	//for _, row := range cont {fmt.Println(row)}
	err := writef(cont, "C:/SakilaGithub/RescueTime-Record/README.md")
	if err == nil {fmt.Println("Update success")}

	//os.Args("")

}
