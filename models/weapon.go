package models

type Weapon struct {
	IDTag      string    `json:"id_tag"`
	SpriteBase []*string `json:"sprite_base"`
	BaseWeapon *string   `json:"base_weapon"`
	Index      int64     `json:"index"`
	Color      int64     `json:"color"`
	Range      int64     `json:"range"`
	Unknown1   int64     `json:"_unknown1"`
	SortID     int64     `json:"sort_id"`
	EquipGroup int64     `json:"equip_group"`
	ResDamage  bool      `json:"res_damage"`
	IsStaff    bool      `json:"is_staff"`
	IsDagger   bool      `json:"is_dagger"`
	IsBreath   bool      `json:"is_breath"`
	IsBeast    bool      `json:"is_beast"`
}
