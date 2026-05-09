package provider

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type PlatformArtifact struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Filename string `json:"filename"`
	Shasum   string `json:"shasum"`
}

type VersionMetadata struct {
	Protocols   []string           `json:"protocols"`
	Platforms   []PlatformArtifact `json:"platforms"`
	SigningKeys SigningKeys        `json:"signing_keys"`
}

type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type Metadata struct {
	Versions []Version `json:"versions"`
}

type SigningKeys struct {
	GPGPublicKeys []GpgPublicKeys `json:"gpg_public_keys"`
}

type GpgPublicKeys struct {
	KeyID          string `json:"key_id"`
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceURL      string `json:"source_url"`
}
