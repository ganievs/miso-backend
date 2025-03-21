package provider

type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type Metadata struct {
	Versions []Version `json:"versions"`
}

type Platform struct {
	Hashes []string `json:"hashes"`
	URL    string   `json:"url"`
}

type Archives struct {
	Archives map[string]Platform `json:"archives"`
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
