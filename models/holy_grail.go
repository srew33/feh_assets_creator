package models

type HolyGrail struct {
	IDNum       int   `json:"id_num"`
	MaxCount    int   `json:"max_count"`
	PayloadSize int   `json:"payload_size"`
	Costs       []int `json:"costs"`
	Reward      []struct {
		Kind   int    `json:"kind"`
		Type   string `json:"_type"`
		Len    int    `json:"len"`
		IDTag  string `json:"id_tag"`
		Rarity int    `json:"rarity"`
	} `json:"reward"`
}
