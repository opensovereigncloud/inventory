package dmi

// Starting enumeration intentionally from zero
// since header type min value is zero
const (
	CBIOSInformationHeaderType = iota
	CSystemInformationHeaderType
	CBoardInformationHeaderType
)
