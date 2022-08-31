package fp

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type UAArray []string

type OSXFp struct {
	Chrome UAArray
	Safari UAArray
}

type LinuxFP struct {
	Chrome UAArray
	Opera  UAArray
}

type WindowsFP struct {
	Chrome UAArray
	Edge   UAArray
	Opera  UAArray
}

type FpData struct {
	Linux   LinuxFP
	Windows WindowsFP
	Osx     OSXFp
}

func (c UAArray) Check(ua string) bool {
	for _, s := range c {
		if s == ua {
			return true
		}
	}
	return false
}

func (f *LinuxFP) Check(ua string) bool {
	var res bool
	if res = f.Opera.Check(ua); res {
		return true
	}
	if res = f.Chrome.Check(ua); res {
		return true
	}
	return false
}

func (f WindowsFP) Check(ua string) bool {
	var res bool
	if res = f.Opera.Check(ua); res {
		return true
	}
	if res = f.Chrome.Check(ua); res {
		return true
	}
	if res = f.Edge.Check(ua); res {
		return true
	}
	return false
}

func (f OSXFp) Check(ua string) bool {
	var res bool
	if res = f.Chrome.Check(ua); res {
		return true
	}
	if res = f.Safari.Check(ua); res {
		return true
	}
	return false
}

func CreateController() FpData {
	var data FpData
	jsonFile, err := os.Open("fp.json")

	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("Successfully Opened fp.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatal().Err(err)
	}
	return data
}

func (d *FpData) Check(ua string) bool {
	var res bool
	if res = d.Linux.Check(ua); res {
		return true
	}
	if res = d.Osx.Check(ua); res {
		return true
	}
	if res = d.Windows.Check(ua); res {
		return true
	}
	return false
}
