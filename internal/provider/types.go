package provider

type Versions struct {
	Versions map[string]struct{} `json:"versions"`
}

type Platform struct {
	Hashes []string `json:"hashes"`
	URL    string   `json:"url"`
}

type Archives struct {
	Archives map[string]Platform `json:"archives"`
}
