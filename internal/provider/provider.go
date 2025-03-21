package provider

type Provider struct {
	Protocols           []string    `json:"protocols"`
	OS                  string      `json:"os"`
	Arch                string      `json:"arch"`
	Filename            string      `json:"filename"`
	DownloadURL         string      `json:"download_url"`
	SHASumsURL          string      `json:"shasums_url"`
	SHASumsSignatureURL string      `json:"shasums_signature_url"`
	SHASum              string      `json:"shasum"`
	SigningKeys         SigningKeys `json:"signing_keys"`
}

func NewProvider(p Provider) *Provider {
	return &Provider{
		Protocols:           p.Protocols,
		OS:                  p.OS,
		Arch:                p.Arch,
		Filename:            p.Filename,
		DownloadURL:         p.DownloadURL,
		SHASumsURL:          p.SHASumsURL,
		SHASumsSignatureURL: p.SHASumsSignatureURL,
		SHASum:              p.SHASum,
		SigningKeys:         p.SigningKeys,
	}
}
