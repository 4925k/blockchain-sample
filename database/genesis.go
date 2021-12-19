package database

import (
	"encoding/json"
	"io/ioutil"
)

var genesisJSON = `{
    "genesis_time": "2021-12-17T00:00:00.000000000Z",
    "chain_id": "nefoli",
    "balances": {
        "dibek": 1000000
    }
  }`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJSON), 0644)
}
