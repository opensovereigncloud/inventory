package utils

import "os"

const (
	CVersionFilePath = "/etc/sonic/sonic_version.yml"
	CMachineType     = "Machine"
	CSwitchType      = "Switch"
)

func GetHostType() (string, error) {
	//todo: determining how to check host type without checking files
	if _, err := os.Stat(CVersionFilePath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		} else {
			return CMachineType, nil
		}
	}
	return CSwitchType, nil
}
