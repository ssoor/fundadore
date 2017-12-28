package config

import (
	"fmt"
	"errors"
	"encoding/json"

	"github.com/ssoor/fundadore/api"
)


type Config struct {
	Redirect   Redirect   `json:"redirect"`
	Internest  Internest  `json:"internest"`
	Fundadore  Fundadore  `json:"fundadore"`
	Statistics Statistics `json:"statistics"`
	Youniverse Youniverse `json:"youniverse"`
}



func GetSettings(buildVer string, account string, guid string) (config Config, err error) {
	url := "http://api.ieceo.cn/" + buildVer + "/Init/Default/GUID/" + account + "/dealerSubADID/" + guid
	//url := "http://api.ieceo.lp.com/" + buildVer + "/Init/Default/GUID/" + account + "/dealerSubADID/" + guid
	
	if 0 == len(guid) {
		url = "http://api.ieceo.cn/" + buildVer + "/Init/Default/GUID/" + account
		//url = "http://api.ieceo.lp.com/" + buildVer + "/Init/Default/GUID/" + account
	}

	jsonConfig, err := api.GetURL(url)
	if err != nil {
		return config, errors.New(fmt.Sprint("Query setting interface failed, err: ", err))
	}

	if err = json.Unmarshal([]byte(jsonConfig), &config); err != nil {
		return config, errors.New("Unmarshal setting interface failed.")
	}

	return config, nil
}
