package dmi

import (
	"bytes"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/lunixbochs/struc"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

type Svc struct {
	printer   *printer.Svc
	RawDMISvc *RawSvc
}

func NewSvc(printer *printer.Svc, rawDMISvc *RawSvc) *Svc {
	return &Svc{
		printer:   printer,
		RawDMISvc: rawDMISvc,
	}
}

func (s *Svc) GetData() (*DMI, error) {
	rawDmi, err := s.RawDMISvc.GetRaw()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get SMBIOS stream")
	}

	defer rawDmi.Stream.Close()

	decoder := smbios.NewDecoder(rawDmi.Stream)
	structures, err := decoder.Decode()
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode SMBIOS stream")
	}

	version := NewSMBIOSVersion(rawDmi.EntryPoint.Version())

	dmi := &DMI{
		Version: version,
	}

	for _, structure := range structures {
		switch structure.Header.Type {
		case CSystemInformationHeaderType:
			systemInfo, err := s.parseSystemInformation(structure, version)
			if err != nil {
				s.printer.VErr(errors.Wrap(err, "unable to parse system info"))
			}
			dmi.SystemInformation = systemInfo
		}
	}

	return dmi, nil
}

func (s *Svc) parseSystemInformation(structure *smbios.Structure, version *SMBIOSVersion) (*SystemInformation, error) {
	// Spec contains info only for 2.0+
	if version.Lesser(&SMBIOSVersion{2, 0, 0}) {
		return &SystemInformation{}, nil
	}

	// 2.4+
	if version.GreaterOrEqual(&SMBIOSVersion{2, 4, 0}) {
		ref := &SystemInformationRefSpec24{}
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return SystemInformationFromSpec24(ref, structure.Strings), nil
	}

	// 2.1+
	if version.GreaterOrEqual(&SMBIOSVersion{2, 1, 0}) {
		ref := &SystemInformationRefSpec21{}
		if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
			return nil, errors.Wrap(err, "unable to unpack structure")
		}
		return SystemInformationFromSpec21(ref, structure.Strings), nil
	}

	// 2.0+
	ref := &SystemInformationRefSpec20{}
	if err := struc.Unpack(bytes.NewReader(structure.Formatted), ref); err != nil {
		return nil, errors.Wrap(err, "unable to unpack structure")
	}
	return SystemInformationFromSpec20(ref, structure.Strings), nil
}
