package models

type WeaponRefine struct {
	Orig    string `json:"orig"`
	Refined string `json:"refined"`
	Use     []Give `json:"use"`
	Give    Give   `json:"give"`
}

type Give struct {
	ResType int64 `json:"res_type"`
	Count   int64 `json:"count"`
}
