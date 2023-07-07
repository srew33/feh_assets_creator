package models

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goccy/go-json"
)

type Skill struct {
	JsonSkill
	AdditionalSkill
}

func DefaultSkill() Skill {
	return Skill{
		JsonSkill: JsonSkill{},
		AdditionalSkill: AdditionalSkill{
			Rarity1: Set[string]{mapset.NewSet[string]()},
			Rarity2: Set[string]{mapset.NewSet[string]()},
			Rarity3: Set[string]{mapset.NewSet[string]()},
			Rarity4: Set[string]{mapset.NewSet[string]()},
			Rarity5: Set[string]{mapset.NewSet[string]()},
		},
	}

}

func (t *Skill) UnmarshalJSON(data []byte) error {

	t1 := AdditionalSkill{}
	if err := json.Unmarshal(data, &t1); err != nil {
		return err
	}
	if t1.Rarity1.Set == nil {
		t1.Rarity1.Set = mapset.NewSet[string]()
	}
	if t1.Rarity2.Set == nil {
		t1.Rarity2.Set = mapset.NewSet[string]()
	}
	if t1.Rarity3.Set == nil {
		t1.Rarity3.Set = mapset.NewSet[string]()
	}
	if t1.Rarity4.Set == nil {
		t1.Rarity4.Set = mapset.NewSet[string]()
	}
	if t1.Rarity5.Set == nil {
		t1.Rarity5.Set = mapset.NewSet[string]()
	}
	t.AdditionalSkill = t1

	t2 := JsonSkill{}
	if err := json.Unmarshal(data, &t2); err != nil {
		return err
	}
	t.JsonSkill = t2

	return nil
}

type JsonSkill struct {
	IDTag           string     `json:"id_tag"`
	RefineBase      *string    `json:"refine_base"`
	NameID          string     `json:"name_id"`
	DescID          string     `json:"desc_id"`
	RefineID        *string    `json:"refine_id"`
	BeastEffectID   *string    `json:"beast_effect_id"`
	Prerequisites   SliceStr   `json:"prerequisites"`
	NextSkill       *string    `json:"next_skill"`
	Sprites         SliceStr   `json:"sprites"`
	Stats           Stats      `json:"stats"`
	ClassParams     Stats      `json:"class_params"`
	CombatBuffs     Stats      `json:"combat_buffs"`
	SkillParams     Stats      `json:"skill_params"`
	SkillParams2    Stats      `json:"skill_params2"`
	RefineStats     Stats      `json:"refine_stats"`
	IDNum           int64      `json:"id_num"`
	SortID          int64      `json:"sort_id"`
	IconID          int64      `json:"icon_id"`
	WEPEquip        int64      `json:"wep_equip"`
	MOVEquip        int64      `json:"mov_equip"`
	SPCost          int64      `json:"sp_cost"`
	Category        int64      `json:"category"`
	TomeClass       int64      `json:"tome_class"`
	Exclusive       StringBool `json:"exclusive"`
	EnemyOnly       StringBool `json:"enemy_only"`
	Range           int64      `json:"range"`
	Might           int64      `json:"might"`
	CooldownCount   int64      `json:"cooldown_count"`
	AssistCD        StringBool `json:"assist_cd"`
	Healing         StringBool `json:"healing"`
	SkillRange      int64      `json:"skill_range"`
	Score           int64      `json:"score"`
	PromotionTier   int64      `json:"promotion_tier"`
	PromotionRarity int64      `json:"promotion_rarity"`
	Refined         StringBool `json:"refined"`
	RefineSortID    int64      `json:"refine_sort_id"`
	WEPEffective    int64      `json:"wep_effective"`
	MOVEffective    int64      `json:"mov_effective"`
	WEPShield       int64      `json:"wep_shield"`
	MOVShield       int64      `json:"mov_shield"`
	WEPEFFWeakness  int64      `json:"wep_eff_weakness"`
	MOVEFFWeakness  int64      `json:"mov_eff_weakness"`
	WEPWeakness     int64      `json:"wep_weakness"`
	MOVWeakness     int64      `json:"mov_weakness"`
	WEPAdaptive     int64      `json:"wep_adaptive"`
	MOVAdaptive     int64      `json:"mov_adaptive"`
	TimingID        int64      `json:"timing_id"`
	AbilityID       int64      `json:"ability_id"`
	Limit1ID        int64      `json:"limit1_id"`
	Limit1Params    SliceI64   `json:"limit1_params"`
	Limit2ID        int64      `json:"limit2_id"`
	Limit2Params    SliceI64   `json:"limit2_params"`
	TargetWEP       int64      `json:"target_wep"`
	TargetMOV       int64      `json:"target_mov"`
	PassiveNext     *string    `json:"passive_next"`
	Timestamp       *string    `json:"timestamp"`
	RandomAllowed   int64      `json:"random_allowed"`
	MinLV           int64      `json:"min_lv"`
	MaxLV           int64      `json:"max_lv"`
	TtInheritBase   StringBool `json:"tt_inherit_base"`
	RandomMode      int64      `json:"random_mode"`
	Limit3ID        int64      `json:"limit3_id"`
	Limit3Params    SliceI64   `json:"limit3_params"`
	RangeShape      int64      `json:"range_shape"`
	TargetEither    StringBool `json:"target_either"`
	DistantCounter  StringBool `json:"distant_counter"`
	CantoRange      int64      `json:"canto_range"`
	PathfinderRange int64      `json:"pathfinder_range"`
	ArcaneWeapon    StringBool `json:"arcane_weapon"`
}

type AdditionalSkill struct {
	IsSkillAccessory StringBool  `json:"is_skill_accessory"`
	Rarity1          Set[string] `json:"rarity1"`
	Rarity2          Set[string] `json:"rarity2"`
	Rarity3          Set[string] `json:"rarity3"`
	Rarity4          Set[string] `json:"rarity4"`
	Rarity5          Set[string] `json:"rarity5"`
	OrigSkill        *string     `json:"orig_skill"`
	Series           string      `json:"series"`
}
