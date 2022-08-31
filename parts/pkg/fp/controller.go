package fp

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type OSXFp struct {
	Chrome []string `json:"chrome"`
	Safari []string `json:"safari"`
}

type LinuxFP struct {
	Chrome []string `json:"chrome"`
	Opera  []string `json:"opera"`
}

type WindowsFP struct {
	Chrome []string `json:"chrome"`
	Edge   []string `json:"edge"`
	Opera  []string `json:"opera"`
}

type FpData struct {
	Linux   LinuxFP   `json:"linux"`
	Windows WindowsFP `json:"windows"`
	Osx     OSXFp     `json:"osx"`
}

func Check(c []string, ua string) bool {
	for _, s := range c {
		if s == ua {
			return true
		}
	}
	return false
}

func CheckLin(f LinuxFP, ua string) bool {
	var res bool
	if res = Check(f.Opera, ua); res {
		return true
	}
	if res = Check(f.Chrome, ua); res {
		return true
	}
	return false
}

func CheckWin(f WindowsFP, ua string) bool {
	var res bool
	if res = Check(f.Opera, ua); res {
		return true
	}
	if res = Check(f.Chrome, ua); res {
		return true
	}
	if res = Check(f.Edge, ua); res {
		return true
	}
	return false
}

func CheckOsx(f OSXFp, ua string) bool {
	var res bool
	if res = Check(f.Chrome, ua); res {
		return true
	}
	if res = Check(f.Safari, ua); res {
		return true
	}
	return false
}

func CreateController() FpData {
	var data FpData
	jsonFile, err := os.Open("./parts/pkg/fp/fp.json")

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
	if res = CheckLin(d.Linux, ua); res {
		return true
	}
	if res = CheckOsx(d.Osx, ua); res {
		return true
	}
	if res = CheckWin(d.Windows, ua); res {
		return true
	}
	return false
}
