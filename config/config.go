package config

import (
	"fmt"
	"errors"
	"strings"
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
		var url string
		if false == strings.HasPrefix(guid, "00000000_") {
			url = "http://social.ssoor.com/issued/settings/20160521/" + account + "/" + guid + ".settings"
		} else {
			url = "http://api.ieceo.cn/" + buildVer + "/Init/Default/GUID/" + guid
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
