package main

import (
	"bufio"
	"fmt"
	"github.com/AbhiAgarwal/go-rescuetime"
	git "github.com/go-git/go-git"
	. "github.com/go-git/go-git/_examples"
	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/plumbing/transport/http"
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
		adqp.ResolutionTime = ""
	} else {
		adqp.ResolutionTime = ResolutionTime
	}
	if RestrictGroup == "" {
		adqp.RestrictGroup = ""
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
	if RestrictThing == "" {
		adqp.RestrictThing = ""
	} else {
		adqp.RestrictThing = RestrictThing
	}
	if RestrictThingy == "" {
		adqp.RestrictThingy = ""
	} else {
		adqp.RestrictThingy = RestrictThingy
	}
	return adqp
}

func sec2hour(sec int) string {
	s := ""
	h := sec / 3600
	if h != 0 {
		s = s + strconv.Itoa(h)
		if len(s) == 1 {s = "0" + s}
		s = s + "h"
	}
	m := (sec % 3600) / 60
	s2 := ""
	if m != 0 {
		s2 = s2 + strconv.Itoa(m)
		if len(s2) == 1 {s2 = "0" + s2}
		s2 = s2 + "m"
	}
	second := sec % 60
	s3 := ""
	if !strings.Contains(s, "h") {
		s3 = s3 + strconv.Itoa(second)
		if len(s3) == 1 {s3 = "0" + s3}
		s3 = s3 + "s"
	}
	return s + s2 + s3
}

func sec2asterisk(sec int) string {
	if sec >= 21600 {return "lllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll"}
	l := sec / 240
	if l == 0 {return ""}
	return strings.Repeat("l", l)
}

func getToday(data *rescuetime.AnalyticData) []string {
	var cont []string
	t := time.Now().Format("2006-01-02 15:04")
	hrow := "## yiyangiliu " + t[:10] + " Detailed Activaties, "
	shrow := "Update at " + t[11:]
	frow := "|Rank|Activity|Len|Time|Category|Label|"
	srow := "|-|-|-|-|-|-|"
	var trow []string
	for i, row := range data.Rows{
		if i == 15 {break}
		t := sec2hour(row.TimeSpentSeconds)
		l := sec2asterisk(row.TimeSpentSeconds)
		act := row.Activity
		switch act {
		case "mobile - tv.danmaku.bili": 			act = "(m)bilibili"
		case "YouTube for Android":					act = "(m)Youtube"
		case "mobile - com.reddit.frontpage": 		act = "(m)Reddit"
		case "mobile - air.tv.douyu.android": 		act = "(m)douyu"
		case "WeChat / Weixin": 					act = "Wechat"
		case "Google Chrome for Android": 			act = "(m)Chrome"
		case "Visual Studio Code": 					act = "VS Code"
		case "Windows Explorer": 					act = "Win Explorer"
		case "mobile - com.hengye.share": 			act = "(m)share"
		case "wechat":								act = "(m)Wechat"
		case "Google Chrome":						act = "Chrome"
		}
		if strings.Contains(act, ".github.io") {
			act = "*.github.io"
		}
		if len(act) > 14 {
			act = act[:14]
		}
		categ := row.Category
		switch categ {
		case "Editing & IDEs":						categ = "IDE"
		case "General Software Development":		categ = "Dev"
		case "General Social Networking":			categ = "SNS"
		case "General Reference & Learning":		categ = "Ref&Learn"
		case "Internet Utilities":					categ = "Utils"
		case "General Utilities":					categ = "Utils"
		case "Uncategorized":						categ = "Unknown"
		case "Presentation":						categ = "Pre"
		}
		p := strconv.Itoa(row.Productivity)
		switch p {
		case "2": p = "💖"
		case "1": p = "❤"
		case "0": p = "🙂"
		case "-1": p = "😥"
		case "-2": p = "💚"
		}
		if len(categ) > 12 {
			categ = categ[:12]
		}
		r :="|" + strconv.Itoa(row.Rank) + "|" + act + "|" + l +"|" + t + "|" + categ + "|" + p + "|"
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
	cont = append(cont, h[:5]...)
	cont = append(cont, t...)
	cont = append(cont, "")
	cont = append(cont, h[5:]...)
	return cont
}

func coverContent(t, h []string) []string {
	var cont []string
	cont = append(cont, h[:5]...)
	cont = append(cont, t...)
	cont = append(cont, h[26:]...)
	return cont
}

func cloneYourRepository(repo, dir string, auth *http.BasicAuth) {
	// Clone the given repository to the given directory
	Info("git clone %s %s", repo, dir)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: 	  auth,
		URL:      repo,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

func commitAndPush(repo, dir string, auth *http.BasicAuth) {
	// reference:
	// commit - Commit changes to the current branch to an existent repository
	// 		https://github.com/go-git/go-git/blob/master/_examples/commit/main.go
	// push - Push repository to default remote (origin)
	// 		https://github.com/go-git/go-git/blob/master/_examples/push/main.go
	r, err := git.PlainOpen(dir)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// Adds the new file to the staging area.
	Info("git add README.md")
	_, err = w.Add("README.md")
	CheckIfError(err)

	// We can verify the current status of the worktree using the method Status.
	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	Info("git commit -m \"Update README.md " + time.Now().Format("2006-01-02") + "\"")
	commit, err := w.Commit("Update README.md " + time.Now().Format("2006-01-02"), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "yiyangiliu",
			Email: "i@yiyangliu.me",
			When:  time.Now(),
		},
	})
	CheckIfError(err)

	//Prints the current HEAD to verify that all worked well.
	Info("git show -s")
	obj, err := r.CommitObject(commit)
	CheckIfError(err)

	fmt.Println(obj)

	Info("git push")
	// push using default options
	err = r.Push(&git.PushOptions{
		RemoteName: "",
		RefSpecs:   nil,
		Auth:       auth,
		Progress:   nil,
		Prune:      false,
	})
	CheckIfError(err)
}