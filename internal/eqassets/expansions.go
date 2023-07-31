package eqassets

import _ "embed"

type Expansion struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Mask int    `json:"mask"`
	Icon string `json:"icon"`
}

//go:embed expansions.json
var expansionJson []byte
