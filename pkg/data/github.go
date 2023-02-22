package data

type GithubCreds struct {
	Token string
}

type GithubFileTree struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Sha           string `json:"sha"`
	Size          int32  `json:"size"`
	Url           string `json:"url"`
	DownloadUrl   string `json:"download_url"`
	HtmlUrl       string `json:"html_url"`
	Type          string `json:"type"`
	DefaultBranch string `json:"default_branch"`
}

type GithubFolderTree struct {
	FolderHash string `json:"sha"`
}

type RawUrl struct {
	Username   string `json:"username"`
	Repository string `json:"repo"`
	Url        string `json:"url"`
}

type commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type RepoTree struct {
	Name   string `json:"name"`
	Commit commit `json:"commit"`
}
