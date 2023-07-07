package models

type Move struct {
	IDTag    string  `json:"id_tag"`
	MoveCost []int64 `json:"move_cost"`
	Index    int64   `json:"index"`
	Range    int64   `json:"range"`
}
