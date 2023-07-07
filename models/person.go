package models

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/slices"
)

type Person struct {
	JsonPerson
	TranslatedNames TransName  `json:"translated_names"`
	ResplendentHero StringBool `json:"resplendent_hero"`
	RecentlyUpdate  StringBool `json:"recently_update"`
	Bst             int64      `json:"bst"`
	MaxRarity       int64      `json:"max_rarity"`
	MinRarity       int64      `json:"min_rarity"`
	Type_           int64      `json:"type"`
	DefaultStats    Stats      `json:"default_stats"`
	Redeemable      bool       `json:"redeemable"`
}

type JsonPerson struct {
	IDTag         string        `json:"id_tag"`
	Roman         string        `json:"roman"`
	FaceName      string        `json:"face_name"`
	FaceName2     string        `json:"face_name2"`
	Legendary     *Legendary    `json:"legendary"`
	Dragonflowers Dragonflowers `json:"dragonflowers"`
	Timestamp     *string       `json:"timestamp"`
	IDNum         int64         `json:"id_num"`
	VersionNum    int64         `json:"version_num"`
	SortValue     int64         `json:"sort_value"`
	Origins       int64         `json:"origins"`
	WeaponType    int64         `json:"weapon_type"`
	TomeClass     int64         `json:"tome_class"`
	MoveType      int64         `json:"move_type"`
	Series        int64         `json:"series"`
	RandomPool    int64         `json:"random_pool"`
	PermanentHero StringBool    `json:"permanent_hero"`
	BaseVectorID  int64         `json:"base_vector_id"`
	Refresher     StringBool    `json:"refresher"`
	BaseStats     Stats         `json:"base_stats"`
	GrowthRates   Stats         `json:"growth_rates"`
	Skills        Skills        `json:"skills"`
}

func (p *JsonPerson) SetStats(stats *Stats, startLevel int, targetLevel int, rarity int, args PersonArgs) error {
	baseStats := p.BaseStats.ToList()
	growthRates := p.GrowthRates.ToList()

	advantage := args.Advantage
	disadvantage := args.Disadvantage
	ascend := args.AscendedAsset

	if advantage == Null && disadvantage == Null {

	} else if advantage != Null && disadvantage != Null && advantage != disadvantage {
		baseStats[advantage-1].Value += 1
		growthRates[advantage-1].Value += 5
		baseStats[disadvantage-1].Value -= 1
		growthRates[disadvantage-1].Value -= 5

		// base_stats[advantage!] = base_stats[advantage]! + 1;
		// base_stats[disadvantage!] = base_stats[disadvantage]! - 1;
		// growth_rates[advantage] = growth_rates[advantage]! + 5;
		// growth_rates[disadvantage] = growth_rates[disadvantage]! - 5;
	} else {
		return fmt.Errorf("错误的属性参数 %+v", args)
	}

	if ascend != advantage && ascend != Null {
		baseStats[ascend].Value += 1
		growthRates[ascend].Value += 5

	} else if ascend == advantage && ascend != Null {
		return fmt.Errorf("错误的属性参数 %+v", args)
	}

	// 获得降序后的属性字典，数值相同时key按key列表的升序
	// 这里假定了base_stats一定是按顺序排列的
	// https://www.reddit.com/r/FireEmblemHeroes/comments/1160n4d/broken_expectations_in_stat_increases_from
	// 说明了突破和神龙之花在计算属性时将HP固定在第一位，因此这里在计算排序列表时也将HP属性固定在第一位即可
	sortedList := make([]Stat, 4)
	copy(sortedList, baseStats[1:])
	// baseStats本身是按key升序，所以相同value不会影响顺序
	sort.SliceStable(sortedList, func(i, j int) bool { return sortedList[i].Value > sortedList[j].Value })

	sortedList = append([]Stat{baseStats[0]}, sortedList...)

	d := Stats{}
	deltaStats := d.ToList()

	// 计算X星属性逻辑，为以3星为基础，1/5星直接-1/+1，2/4星在1/3星基础上取排序后属性的最大的前两个除HP的属性+1，
	switch rarity {
	case 1:
		for i := range deltaStats {
			deltaStats[i].Value -= 1
		}
	case 2:
		c := 0
		// 倒序遍历sortedList，取除了HP的前两个属性+1（即取HP和最小的两个属性-1）
		deltaStats[0].Value -= 1
		for i := len(sortedList) - 1; i >= 0; i-- {
			stat := sortedList[i].Key
			if stat != HP {
				deltaStats[stat-1].Value -= 1
				c += 1
			}
			if c == 2 {
				break
			}
		}

	case 3:
	case 4:
		c := 0
		for i := 0; i < len(sortedList); i++ {
			stat := sortedList[i].Key
			if stat != HP {
				deltaStats[stat-1].Value += 1
				c += 1
			}
			if c == 2 {
				break
			}
		}
	case 5:
		for i := range deltaStats {
			deltaStats[i].Value += 1
		}
	default:
		return fmt.Errorf("错误的属性参数 rarity:%d", rarity)
	}

	// --------------------------------------------------神龙之花-----------------------
	// 如果dragonflowers大于等于5， 每+5则五维再+1
	// 否则按sorted_keys的顺序将对应属性+1
	if args.Dragonflowers >= 5 {
		for i := range deltaStats {
			deltaStats[i].Value += 1
		}
	}

	for i := 0; i < args.Dragonflowers%5; i++ {
		stat := sortedList[i].Key
		deltaStats[stat-1].Value += 1
	}

	// --------------------------------------------------召唤师的羁绊-----------------------
	// HP+5 四维+2

	if args.SummonerSupport {
		for i := range deltaStats {
			if i == HP {
				deltaStats[i].Value += 5
			} else {
				deltaStats[i].Value += 2
			}
		}
	}

	// --------------------------------------------------神装英雄-----------------------
	// 五维+2
	if args.Resplendent {
		for i := range deltaStats {
			deltaStats[i].Value += 2
		}
	}

	//--------------------------------------------------突破---------------------------

	// 如果中性，前三高属性+1,有优劣属性 突破+1时首先把劣势属性+3或+4
	if args.Merged > 0 {
		// 如果是中性，按排序属性去掉开花属性的列表的前三项+1
		if advantage == Null && disadvantage == Null {
			c := 0
			for i := 0; i < len(sortedList); i++ {
				stat := sortedList[i].Key
				if stat != ascend {
					deltaStats[stat-1].Value += 1
					c += 1
				}
				if c == 3 {
					break
				}

			}

		} else {
			if slices.Contains(DisadvantageRates, int(growthRates[disadvantage-1].Value+5)) {
				deltaStats[disadvantage-1].Value += 4
			} else {
				deltaStats[disadvantage-1].Value += 3
			}
		}
	}

	for i := 0; i < args.Merged%5; i++ {
		s := (i * 2) % 5
		e := (i*2 + 1) % 5
		deltaStats[sortedList[s].Key-1].Value += 1
		deltaStats[sortedList[e].Key-1].Value += 1

	}

	// 如果merge大于等于5， 每+5则五维再+2
	for j := 0; j < args.Merged/5; j++ {
		for i := range deltaStats {
			deltaStats[i].Value += 2
		}
	}

	r := make([]Stat, 5)
	for i := 0; i < 5; i++ {
		growValue := math.Trunc(float64((targetLevel - startLevel)) * math.Trunc(float64(growthRates[i].Value)*(0.79+(0.07*float64(rarity)))) / 100.0)
		r[i].Value = baseStats[i].Value + deltaStats[i].Value + int64(growValue)
	}
	stats.FromList(r)
	return nil
}

type PersonArgs struct {
	Advantage       int
	Disadvantage    int
	Merged          int
	Dragonflowers   int
	Resplendent     bool
	SummonerSupport bool
	AscendedAsset   int
}
