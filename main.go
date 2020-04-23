package main

import (
	"fmt"
	"github.com/go-git/go-git/plumbing/transport/http"
	"os"
	"time"
)

func main() {
	//Variables:
	//  rtapi: string, Your RescueTime API Key,
	//		https://www.rescuetime.com/anapi/manage
	//	un: string, Username of your Github account
	//	pw: string, Password of your Github account
	//  token:, string, "personal access token" of your Github account:
	//		https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line
	//	repo: string, "HTTPS URLs" of your repository
	//		https://help.github.com/en/github/using-git/which-remote-url-should-i-use#cloning-with-https-urls-recommended
	//	dir: string, Directory path that your repository cloned into
	//	fpath: string, Path of "README.md" file of your repository
	//  auth: http.BasicAuth, the "auth" Type contains your Github username & password
	//		or username & your "personal access token"
	//  nrt, rescuetime.RescueTime, basic RescueTime object
	//  data, rescuetime.AnalyticData, a json-like object contains your today's detailed data
	//  today, []string, transformed by "data" to a slice of string
	//  				 that printed like a markdown table,
	//  	example:
	//		[[|Rank|Activity|Time|Category|Label|],
	//		 [|-|-|-|-|-|],
	//		 [|1|goland64|4h37m|Dev|2|],
	//		 [|2|github.com|1h14m|Dev|2|],
	//		 ...
	//		 [|15|dllhost|4m30s|Utils|1|]]
	//  history: []string, read old "README.md" by lines
	//  hd: string, the latest date of your old "README.md", like "2020-04-21"
	//  td: string, today's date, like "2020-04-22"
	//		the "README.md" file will update only when hd < td
	//  cont: []string, the new content of "README.md", mixed by "today" and "history"

	err := os.Setenv("http_proxy", "http://127.0.0.1:1080")
	if err != nil {}
	err = os.Setenv("https_proxy", "http://127.0.0.1:1080")
	if err != nil {}


	rtapi := os.Getenv("RESCUETIME_API_KEY")  // like: "B63IavC02qsRZ4QZjl7lURlX6wiV_D_m9Z4ReXvR"

	un := os.Getenv("GITHUB_USERNAME") // like: "yiyangiliu"
	pw := os.Getenv("GITHUB_PASSWORD") // like: "abC123!@#"
	//token := os.Getenv("GITHUB_AUTH_TOKEN")
	repo := "https://github.com/yiyangiliu/RescueTime-Record.git"
	dir := "C:/SakilaGithub/RescueTime-Record"
	fpath := "C:/SakilaGithub/RescueTime-Record/README.md"
	fmt.Printf("rtapi: %s\nun: %s\npw: %s\n", rtapi, un, pw)


	auth := &http.BasicAuth{un, pw}
	// use "personal access token" instead of password
	//auth := &http.BasicAuth{un, token}

	// if you haven't clone your repo to local your need to run following line of code
	//cloneYourRepository(repo, dir, auth)

	nrt := NewRescueTime(rtapi)
	a := NewAnalyticDataQueryParameters("",
		"",
		"",
		"", //2020-04-21
		"", //2020-04-21
		"",
		"",
		"")
	data, _ := nrt.GetAnalyticData("", &a)
	//j,_ := json.Marshal(data)
	//fmt.Println(string(j))

	today :=  getToday(&data) // change name, which default is "yiyangiliu"
	//for _, row := range today { fmt.Println(row)}
	if len(today) < 15 {fmt.Printf("len(today): %#v\n", len(today))}
	history := getHistory(fpath)
	//for _, row := range history { fmt.Println(row)}

	hd := history[3][14:24]
	td := time.Now().Format("2006-01-02")
	if td < hd {fmt.Println("Today is the latest day, nothing need to update")
	} else {
		cont := getContent(today, history)
		//for _, row := range cont { fmt.Println(row)}

		err := writef(cont, fpath)
		if err == nil {fmt.Println("Update success\n")}

		commitAndPush(repo, dir, auth)
		fmt.Println("Commit & Push success")
	}
}
