package model

type HashAlgorithm string
type ExternalReference string

const (
	SHA1 HashAlgorithm     = "SHA-1"
	VCS  ExternalReference = "vcs"
)

type Dependency struct {
	Name               string
	Version            string
	PackageURL         string
	Hashes             map[HashAlgorithm]string
	ExternalReferences map[ExternalReference]string
}
