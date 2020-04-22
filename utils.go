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