package gg

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CloneResult struct {
	Repo Repo  `json:"repo"`
	Err  error `json:"err"`
}
