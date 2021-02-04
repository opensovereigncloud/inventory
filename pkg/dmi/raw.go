package dmi

import (
	"io"
	_ "unsafe"

	"github.com/digitalocean/go-smbios/smbios"
)

type Raw struct {
	Stream     io.ReadCloser
	EntryPoint smbios.EntryPoint
}
