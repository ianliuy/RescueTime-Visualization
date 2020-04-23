# RescueTime-Github
**Automatically** Push `RescueTime Daily Data` to your `Github Repository`

Personal Demo: [yiyangiliu/RescueTime-Record](https://github.com/yiyangiliu/RescueTime-Record)

## Quick Start

1. [Create a repository](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-new-repository) and [clone it](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository) to your local

2. [set system environment variables](https://www.google.com/search?q=set+system+environment+variables)

|key|value|
|-|-|
|RESCUETIME_API_KEY|[Your RescueTime API Key](https://www.rescuetime.com/anapi/manage)|
|GITHUB_USERNAME|yiyangiliu|
|GITHUB_PASSWORD|abC123!@#|

3. `go build *.go` and run `main.exe`



```golang
Variables:
	 rtapi: string, Your RescueTime API Key,
			https://www.rescuetime.com/anapi/manage
		un: string, Username of your Github account
		pw: string, Password of your Github account
	 token:, string, "personal access token" of your Github account:
			https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line
		repo: string, "HTTPS URLs" of your repository
			https://help.github.com/en/github/using-git/which-remote-url-should-i-use#cloning-with-https-urls-recommended
		dir: string, Directory path that your repository cloned into
		fpath: string, Path of "README.md" file of your repository
	 auth: http.BasicAuth, the "auth" Type contains your Github username & password
			or username & your "personal access token"
	 nrt, rescuetime.RescueTime, basic RescueTime object
	 data, rescuetime.AnalyticData, a json-like object contains your todays detailed data
	 today, []string, transformed by "data" to a slice of string
	 				 that printed like a markdown table,
	 	example:
			[[|Rank|Activity|Time|Category|Label|],
			 [|-|-|-|-|-|],
			 [|1|goland64|4h37m|Dev|2|],
			 [|2|github.com|1h14m|Dev|2|],
			 ...
			 [|15|dllhost|4m30s|Utils|1|]]
	 history: []string, read old "README.md" by lines
	 hd: string, the latest date of your old "README.md", like "2020-04-21"
	 td: string, todays date, like "2020-04-22"
			the "README.md" file will update only when hd < td
	 cont: []string, the new content of "README.md", mixed by "today" and "history"
```
