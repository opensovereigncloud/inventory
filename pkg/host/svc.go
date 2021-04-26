package host

import (
	"encoding/json"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	CVersionFilePath = "/etc/sonic/sonic_version.yml"
	CMachineType     = "Machine"
	CSwitchType      = "Switch"
)

type Info struct {
	Type     string
	Hostname string
}

type Distro struct {
	BuildVersion  string
	DebianVersion string
	KernelVersion string
	AsicType      string
	CommitId      string
	BuildDate     string
	BuildNumber   uint32
	BuildBy       string
}

type Svc struct {
	printer       *printer.Svc
	sonicInfoPath string
}

func NewSvc(printer *printer.Svc, basePath string) *Svc {
	return &Svc{
		printer:       printer,
		sonicInfoPath: path.Join(basePath, CVersionFilePath),
	}
}

func (s *Svc) GetData() (*Info, *Distro, error) {
	info := Info{}
	distro := Distro{}
	name, err := os.Hostname()
	if err != nil {
		s.printer.VErr(errors.Errorf("failed to get hostname"))
	}
	info.Hostname = name

	rawInfo := make(map[string]interface{})
	sonicConfigFileExists, err := sonicConfigExists(s.sonicInfoPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to determine host type")
	}
	if sonicConfigFileExists {
		sonicInfo, err := ioutil.ReadFile(s.sonicInfoPath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to read SONiC version file")
		}
		err = yaml.Unmarshal(sonicInfo, &rawInfo)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to collect SONiC version")
		}
		err = convertMapStruct(&distro, rawInfo)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to process SONiC version")
		}
		info.Type = CSwitchType
	} else {
		info.Type = CMachineType
		//todo: collect distro info from regular machine
	}
	return &info, &distro, nil
}

func convertMapStruct(obj *Distro, m map[string]interface{}) error {
	for k, v := range m {
		m[strings.Replace(k, "_", "", 1)] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

func sonicConfigExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}
