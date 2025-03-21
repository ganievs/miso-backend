package module

type Version struct {
	Version string `json:"version"`
}

type Metadata struct {
	Versions []Version `json:"versions"`
}

type Module struct {
	Namespace    string
	Name         string
	TargetSystem string
}
