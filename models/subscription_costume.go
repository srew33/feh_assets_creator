package models

type SubscriptionCostumeElement struct {
	IDNum       int64  `json:"id_num"`
	AvailStart  string `json:"avail_start"`
	AvailFinish string `json:"avail_finish"`
	HeroID      string `json:"hero_id"`
}
