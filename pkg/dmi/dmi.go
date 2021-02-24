package dmi

type DMI struct {
	Version           *SMBIOSVersion
	BIOSInformation   *BIOSInformation
	SystemInformation *SystemInformation
}
