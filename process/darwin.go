// +build darwin

package process

import (
	"github.com/mitchellh/go-ps"
)

type process struct{}

func (*process) IsRunning(name string) (bool, error) {
	pList, err := ps.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range pList {
		pName := p.(*ps.DarwinProcess).Executable()
		if name == pName {
			return true, nil
		}
	}

	return false, nil
}
