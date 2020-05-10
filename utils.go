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

func fromSummaryGetsuString(suList []rescuetime.DailySummary, tod string) string {
	suString := ""
	for _, su := range suList {
		if su.Date == tod {
			fmt.Println(suString)
			var m float64
			tAct := ""
			tT := ""
			tL := ""
			m = 0
			if su.BusinessHours > m {
				m = su.BusinessHours
				tAct = "Bussiness"
				tT = strings.Replace(su.BusinessDurationFormatted, " ", "", -1)
				tL = "🧡"
			}
			if su.CommunicationAndSchedulingHours > m {
				m = su.CommunicationAndSchedulingHours
				tAct = "Communication"
				tT = strings.Replace(su.CommunicationAndSchedulingDurationFormatted, " ", "", -1)
				tL = "🧡"
			}
			if su.DesignAndCompositionHours > m {
				m = su.DesignAndCompositionHours
				tAct = "Composition"
				tT = strings.Replace(su.DesignAndCompositionDurationFormatted, " ", "", -1)
				tL = "💖"
			}
			if su.EntertainmentHours > m {
				m = su.EntertainmentHours
				tAct = "Entertainment"
				tT = strings.Replace(su.EntertainmentDurationFormatted, " ", "", -1)
				tL = "💚"
			}
			if su.NewsHours > m {
				m = su.NewsHours
				tAct = "News"
				tT = strings.Replace(su.NewsDurationFormatted, " ", "", -1)
				tL = "💜"
			}
			if su.ReferenceAndLearningHours > m {
				m = su.ReferenceAndLearningHours
				tAct = "Reference"
				tT = strings.Replace(su.ReferenceAndLearningDurationFormatted, " ", "", -1)
				tL = "🧡"
			}
			if su.ShoppingHours > m {
				m = su.ShoppingHours
				tAct = "Shopping"
				tT = strings.Replace(su.ShoppingDurationFormatted, " ", "", -1)
				tL = "💚"
			}
			if su.SocialNetworkingHours > m {
				m = su.SocialNetworkingHours
				tAct = "SNS"
				tT = strings.Replace(su.SocialNetworkingDurationFormatted, " ", "", -1)
				tL = "💚"
			}
			if su.SoftwareDevelopmentHours > m {
				m = su.SoftwareDevelopmentHours
				tAct = "Dev"
				tT = strings.Replace(su.SoftwareDevelopmentDurationFormatted, " ", "", -1)
				tL = "💖"
			}
			if su.UncategorizedHours > m {
				m = su.UncategorizedHours
				tAct = "Unknown"
				tT = strings.Replace(su.UncategorizedDurationFormatted, " ", "", -1)
				tL = "💛"
			}
			if su.UtilitiesHours > m {
				m = su.UtilitiesHours
				tAct = "Utilities"
				tT = strings.Replace(su.UtilitiesDurationFormatted, " ", "", -1)
				tL = "💛"
			}
			suString = "Total logged time: " + "" + strings.Replace(su.TotalDurationFormatted, " ", "", -1) + ";"
			suString = suString + " Today's theme is: " + "**" + tAct + "**" + tL + " (" + tT + ")"
		}
	}
	return suString
}

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

func getToday(data *rescuetime.AnalyticData, tod, suString string) []string {
	var cont []string
	t := time.Now().Format("2006-01-02 15:04")
	hrow := "## yiyangiliu " + tod + " Detailed Activaties, "
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
		case "WeChat / Weixin": 					act = "(m)Wechat"
		case "Google Chrome for Android": 			act = "(m)Chrome"
		case "Visual Studio Code": 					act = "VS Code"
		case "Windows Explorer": 					act = "Win Explorer"
		case "mobile - com.hengye.share": 			act = "(m)share"
		case "wechat":								act = "Wechat"
		case "Google Chrome":						act = "Chrome"
		case "en.wikipedia.org":					act = "en.wikipedia"
		case "mobile - com.taobao.taobao":			act = "(m)Taobao"
		case "Command Prompt": 						act = "Win Cmd"
		}
		if strings.Contains(act, ".github.io") {
			act = "*.github.io"
		}
		if strings.Contains(act, "mobile - net.oneplus.launcher") {
			act = "(m)launcher"
		}
		if strings.Contains(act, "mobile - com.rammigsoftware.bluecoins") {
			act = "(m)[bluecoins](https://www.google.com/search?q=bluecoins)"
		}

		if len(act) > 14 {
			idx := len(act) - 4
			if strings.Contains(act, ".com") {act = act[:idx]}
			if strings.Contains(act, ".org") {act = act[:idx]}
			if strings.Contains(act, ".net") {act = act[:idx]}
			if len(act) > 14 {act = act[:14]}
		}
		if strings.Contains(act, "cdonnmf") {
			act = "[Saladict](https://github.com/crimx/ext-saladict#saladict-%E6%B2%99%E6%8B%89%E6%9F%A5%E8%AF%8D) PDF"
		}
		if strings.Contains(act, "(m)Taobao") {
			act = "(m)[Taobao](https://en.wikipedia.org/wiki/Taobao"
		}
		if strings.Contains(act, "zhihu") {
			act = strings.Replace(act, "zhihu", "[zhihu](https://en.wikipedia.org/wiki/Zhihu)", -1)
		}
		if strings.Contains(act, "quicker") {
			act = "[quicker](https://getquicker.net/)"
		}
		if strings.Contains(act, "youtube music") {
			act = "[youtube music](https://github.com/ytmdesktop/ytmdesktop)"
		}
		if strings.Contains(act, "huggingface.co") {
			act = "[huggingface.co](https://huggingface.co/)"
		}
		if strings.Contains(act, "paperswithcode") {
			act = "[paperswithcode](https://paperswithcode.com/area/natural-language-processing/)"
		}
		if strings.Contains(act, "linkedin.com") {
			act = "[linkedin.com](https://www.linkedin.com/in/yiyang-liu-aa56b2192/)"
		}
		if strings.Contains(act, "facebook.com") {
			act = "[facebook.com](https://www.facebook.com/Yiyang.Ian.Liu)"
		}
		if strings.Contains(act, "wikipedia") {
			act = "en.wikipedia"
		}
		if strings.Contains(act, "douyu") {
			act = strings.Replace(act, "douyu", "[douyu](https://www.google.com/search?q=douyu+chinese+twitch)", -1)
		}
		if strings.Contains(act, "bilibili") {
			act = strings.Replace(act, "bilibili", "[bilibili](https://www.youtube.com/watch?v=f-wBecEp6Mk&t=560s)", -1)
		}
		if strings.Contains(act, "stackoverflow") {
			act = "stackoverflow"
		}
		if strings.Contains(act, "programcreek") {
			act = "[programcreek](https://www.programcreek.com/python/)"
		}

		categ := row.Category
		switch categ {
		case "Editing & IDEs":						categ = "IDE"
		case "General Software Development":		categ = "Dev"
		case "General Social Networking":			categ = "SNS"
		case "General Reference & Learning":		categ = "Reference"
		case "Internet Utilities":					categ = "Utils"
		case "General Utilities":					categ = "Utils"
		case "Uncategorized":						categ = "Unknown"
		case "Presentation":						categ = "Pre"
		case "Writing":								categ = "Composing"
		case "Engineering & Technology":			categ = "Tech"
		case "Professional Networking":				categ = "Pro"
		case "General Business":					categ = "Business"
		case "Intelligence":						categ = "Insights"
		case "Instant Message":						categ = "IM"
		case "General News":						categ = "News"
		case "General Shopping":					categ = "Shopping"
		}
		p := strconv.Itoa(row.Productivity)
		switch p {
		case "2": p = "💖"
		case "1": p = "🧡"
		case "0": p = "💛"
		case "-1": p = "💜"
		case "-2": p = "💚"
		}
		if len(categ) > 12 {
			categ = categ[:12]
		}
		r :="|" + strconv.Itoa(row.Rank) + "|" + act + "|" + l +"|" + t + "|" + categ + "|" + p + "|"
		trow = append(trow, r)
		//fmt.Println(string()
	}
	if suString == "" {
		cont = append(cont, hrow, "", shrow, "", frow, srow)
	} else {
		cont = append(cont, hrow, "", shrow, "", suString, "", frow, srow)
	}
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

//func histoyWSummary(history string) string {
//	y := history[5, ]
//	return history
//}

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
	cont = append(cont, h[28:]...)
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