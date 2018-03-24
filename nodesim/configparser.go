package nodesim

import (
	"encoding/json"
	"io/ioutil"
)

func ExecConfigFile(file string, managerIP string) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.E(err)
	}
	var config interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.E(err)
	}
	for _, r := range (config.(map[string]interface{}))["rows"].([]interface{}) {
		row := r.(map[string]interface{})
		north := int(row["north"].(float64))
		south := int(row["south"].(float64))
		east := int(row["east"].(float64))
		west := int(row["west"].(float64))
		num := uint64(row["num"].(float64))
		supernode := row["supernode"].(bool)
		//group := row["group"].(bool)
		if supernode {
			CreateSupernodes(num, managerIP)
		} else {
			CreateNodes(num, managerIP, north, south, east, west)
		}

	}
}
