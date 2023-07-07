package models

type SkillAccessory struct {
	IDTag        string  `json:"id_tag"`
	NextSeal     *string `json:"next_seal"`
	PrevSeal     *string `json:"prev_seal"`
	SsCoin       int64   `json:"ss_coin"`
	SsBadgeType  int64   `json:"ss_badge_type"`
	SsBadge      int64   `json:"ss_badge"`
	SsGreatBadge int64   `json:"ss_great_badge"`
}
