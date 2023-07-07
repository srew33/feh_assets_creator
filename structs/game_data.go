package structs

import (
	"errors"
	"feh_assets_creator/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"feh_assets_creator/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/mozillazg/go-pinyin"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "embed"
)

var SkillField = []string{"Weapon", "Assist", "Special", "PassiveA", "PassiveB", "PassiveC"}

//go:embed passive.yaml
var passiveRegex []byte

type GameData struct {
	Person              []models.Person
	Skill               []models.Skill
	Weapon              []models.Weapon
	Move                []models.Move
	WeaponRefine        []models.WeaponRefine
	SkillAccessory      []models.SkillAccessory
	SubscriptionCostume []models.SubscriptionCostumeElement
	HolyGrail           []models.HolyGrail
	TransCn             map[string]*string
	TransJp             map[string]*string
	TransEn             map[string]*string
	SkillSeries         [][]string
}

func (l *GameData) setTranslations() {

	zhExt := map[string]string{
		"CUSTOM_STATS_HP":     "血量",
		"CUSTOM_STATS_ATK":    "攻擊",
		"CUSTOM_STATS_SPD":    "速度",
		"CUSTOM_STATS_DEF":    "防守",
		"CUSTOM_STATS_RES":    "魔防",
		"CUSTOM_STATS_NULL":   "中性",
		"CUSTOM_STATS_TRAITS": "+{}-{}",
		"CUSTOM_LANG_ZH":      "繁体中文",
		"CUSTOM_LANG_JA":      "日本語",
		"CUSTOM_LANG_EN":      "English",
		"ATK":                 "攻",
		"AGI":                 "速",
		"DEF":                 "防",
		"RES":                 "抗",
		"":                    "",
		"SID_CUSTOM_無し":       "错误",
		"MSID_CUSTOM_無し":      "错误",
	}
	otherExt := map[string]string{
		"CUSTOM_STATS_HP":     "HP",
		"CUSTOM_STATS_ATK":    "ATK",
		"CUSTOM_STATS_SPD":    "AGI",
		"CUSTOM_STATS_DEF":    "DEF",
		"CUSTOM_STATS_RES":    "RES",
		"CUSTOM_STATS_NULL":   "",
		"CUSTOM_STATS_TRAITS": "+{}-{}",
		"CUSTOM_LANG_ZH":      "繁体中文",
		"CUSTOM_LANG_JA":      "日本語",
		"CUSTOM_LANG_EN":      "English",
		"ATK":                 "ATK",
		"AGI":                 "AGI",
		"DEF":                 "DEF",
		"RES":                 "RES",
		"":                    "",
		"SID_CUSTOM_無し":       "error",
		"MSID_CUSTOM_無し":      "error",
	}

	if l.TransCn == nil {
		l.TransCn = make(map[string]*string)
	}
	if l.TransEn == nil {
		l.TransEn = make(map[string]*string)
	}
	if l.TransJp == nil {
		l.TransJp = make(map[string]*string)
	}
	// 因为值复制的问题，这里要先使用一个变量固定住V的值
	for k, v := range zhExt {
		s := v
		p := &s
		l.TransCn[k] = p
	}
	for k, v := range otherExt {
		s := v
		p := &s
		l.TransEn[k] = p
		l.TransJp[k] = p
	}
}

func (l *GameData) setPinyin() {
	arg := pinyin.NewArgs()
	for i := range l.Person {
		p := &l.Person[i]
		p.TranslatedNames = make(map[string]string)

		if t, ok := l.TransCn[fmt.Sprintf("M%s", p.IDTag)]; ok {
			p.TranslatedNames["zh_TW"] = strings.Join(pinyin.LazyPinyin(*t, arg), "_")
		} else {
			p.TranslatedNames["zh_TW"] = p.Roman
		}
		p.TranslatedNames["en_US"] = p.Roman
		p.TranslatedNames["ja_JP"] = p.Roman
	}
}

func (l *GameData) setSubscription() {

	for j := range l.SubscriptionCostume {
		find := utils.IndexFunc(len(l.Person), func(i int) bool { return l.Person[i].IDTag == l.SubscriptionCostume[j].HeroID })
		// find := slices.utils.IndexFunc(l.Person, func(p Person) bool { return p.IDTag == v.HeroID })
		if find != -1 {
			l.Person[find].ResplendentHero = true
		}
	}
}

func (l *GameData) setRedeemable() {

	for j := range l.HolyGrail {
		find := utils.IndexFunc(len(l.Person), func(i int) bool { return l.Person[i].IDTag == l.HolyGrail[j].Reward[0].IDTag })
		// find := slices.utils.IndexFunc(l.Person, func(p Person) bool { return p.IDTag == v.HeroID })
		if find != -1 {
			l.Person[find].Redeemable = true
		}
	}
}

func (l *GameData) setBstAndDefault() error {
	for i := 0; i < len(l.Person); i++ {
		p := &l.Person[i]

		err := p.SetStats(&p.DefaultStats, 1, 40, 5, models.PersonArgs{})
		if err != nil {
			return err
		}
		skillBst := 0
		for _, skills := range p.Skills {
			for _, skill := range skills {
				if skill != nil {
					// find := utils.IndexFunc2(l.Skill, func(i int) bool { return l.Skill[i].IDTag == *skill })
					find := utils.IndexFunc(len(l.Skill), func(i int) bool { return l.Skill[i].IDTag == *skill })
					if find == -1 {
						return fmt.Errorf("加载技能出错: 未找到 %s 技能", *skill)
					}

					s := l.Skill[find]
					if s.Category == 3 && s.TimingID == 18 {
						if p.Legendary != nil && p.Legendary.Kind == 1 && s.SkillParams.Atk != 0 {
							skillBst = int(s.SkillParams.Atk)
						} else {
							skillBst = int(s.SkillParams.HP)
						}
						// if p.Legendary != nil || p.Legendary.Kind != 1 {
						// 	skillBst = int(s.SkillParams.HP)
						// } else if s.SkillParams.Atk != 0 {
						// 	skillBst = int(s.SkillParams.Atk)
						// } else {
						// 	skillBst = int(s.SkillParams.HP)
						// }
					}
				}

			}
		}

		legendaryBst := 0
		if p.Legendary != nil {
			legendaryBst = int(p.Legendary.Bst)
		}

		bsts := []int{legendaryBst, skillBst, int(p.DefaultStats.Sum())}
		// slices.SortStableFunc(bsts, func(a int, b int) bool { return b > a })
		sort.Ints(bsts)
		l.Person[i].Bst = int64(bsts[2])
	}

	return nil
}

func (l *GameData) setSkillRarity() error {
	for i := range l.Person {
		p := &l.Person[i]
		pRarity := 0
		if p.MinRarity == 0 {
			if p.MaxRarity == 0 {
				pRarity = 5
			} else {
				pRarity = int(p.MaxRarity)
			}
		} else {
			pRarity = int(p.MinRarity)
		}

		for _, skills := range p.Skills {
			for _, tag := range skills {
				if tag != nil {
					find := utils.IndexFunc(len(l.Skill), func(i int) bool { return l.Skill[i].IDTag == *tag })
					// find := slices.utils.IndexFunc(l.Skill, func(s Skill) bool { return s.IDTag == *tag })
					if find == -1 {
						return fmt.Errorf("加载技能出错: 未找%s技能", *tag)
					}
					s := &l.Skill[find]
					switch pRarity {
					case 1:
						s.Rarity1.Add(p.IDTag)
					case 2:
						s.Rarity2.Add(p.IDTag)
					case 3:
						s.Rarity3.Add(p.IDTag)
					case 4:
						s.Rarity4.Add(p.IDTag)
					case 5:
						s.Rarity5.Add(p.IDTag)
					}
				}
			}
		}
	}
	return nil
}

func (l *GameData) setAccessorySkill() error {
	for _, sa := range l.SkillAccessory {
		find := utils.IndexFunc(len(l.Skill), func(i int) bool { return sa.IDTag == l.Skill[i].IDTag })
		if find == -1 {
			return fmt.Errorf("加载技能出错: 未找%s技能", sa.IDTag)
		}

		l.Skill[find].IsSkillAccessory = true
	}

	return nil
}
func (l *GameData) setRefineList() error {
	for i := range l.WeaponRefine {
		wr := &l.WeaponRefine[i]
		find := utils.IndexFunc(len(l.Skill), func(i int) bool { return wr.Refined == l.Skill[i].IDTag })
		if find == -1 {
			return fmt.Errorf("加载技能出错: 未找%s技能", wr.Refined)
		}

		l.Skill[find].OrigSkill = &wr.Orig
	}
	return nil
}

func (l *GameData) appendBlessing() {
	blessingDict := map[int64]string{
		1: "火",
		2: "水",
		3: "风",
		4: "地",
		5: "光",
		6: "暗",
		7: "天",
		8: "理",
	}

	newBlessing := make(map[string]models.Skill)
	blessingId := 9000
	find := utils.IndexFunc(len(l.Skill), func(i int) bool { return l.Skill[i].IDTag == "SID_ダメージ強化R差分" })
	// find := slices.utils.IndexFunc(l.Skill, func(s Skill) bool { return s.IDTag == "SID_ダメージ強化R差分" })
	base := l.Skill[find]

	for _, p := range l.Person {
		if p.Legendary != nil {
			legendary := p.Legendary
			if !legendary.BonusEffect.IsDefault() {
				key := fmt.Sprintf("%s %s", blessingDict[legendary.Element], legendary.BonusEffect.String())
				if _, ok := newBlessing[key]; ok {
					s := newBlessing[key]
					s.Rarity5.Add(p.IDTag)
				} else {
					blessing := base
					blessing.IDTag = key
					blessing.NameID = legendary.BonusEffect.String()
					blessing.DescID = legendary.BonusEffect.String()
					blessing.SortID = legendary.Element
					blessing.IDNum = int64(blessingId)
					blessing.IconID = legendary.Element
					// Rarityx本质是一个引用类型，浅拷贝不能产生独立对象，这里所有稀有度字段都指向base的对应字段的地址
					// 由于只用到Rarity5，这里只克隆这一个字段
					blessing.Rarity5.Set = blessing.Rarity5.Clone()
					blessing.Rarity5.Add(p.IDTag)
					blessing.Category = 15
					blessing.Stats = legendary.BonusEffect
					newBlessing[key] = blessing
					blessingId += 1
				}
			}
		}
	}

	for _, v := range newBlessing {
		l.Skill = append(l.Skill, v)
	}

	base.IDTag = "SID_CUSTOM_無し"
	base.NameID = "SID_CUSTOM_無し"
	base.DescID = "SID_CUSTOM_無し"
	l.Skill = append(l.Skill, base)
}

func (t *GameData) setSkillSeries() error {

	keys, res, err := initRegex()
	if err != nil {
		return err
	}

	for i := range t.Skill {
		s := &t.Skill[i]

		if s.Category < 6 {
			cRes := res[s.Category]
			matched := false

		LOOP1:
			for i, v := range cRes {
				for _, r := range v {
					if r.MatchString(s.IDTag) {
						s.Series = keys[s.Category][i]
						matched = true
						break LOOP1
					}
				}
			}

			if !matched && len(cRes) != 0 {
				s.Series = "其他"
			}
		}

	}

	t.SkillSeries = keys

	return nil

}

func (l *GameData) setRecently(old GameData) {
	// 新增人物
	for j := range l.Person {
		if utils.IndexFunc(len(old.Person), func(i int) bool { return old.Person[i].IDTag == l.Person[j].IDTag }) == -1 {
			// if slices.utils.IndexFunc(old.Person, func(o Person) bool { return o.IDTag == p.IDTag }) == -1 {
			l.Person[j].RecentlyUpdate = true
		}
	}

	// 新增了技能的人物
	for _, skill := range l.Skill {
		// skill := l.Skill[j]
		if utils.IndexFunc(len(old.Skill), func(i int) bool { return old.Skill[i].IDTag == skill.IDTag }) == -1 {
			// if slices.utils.IndexFunc(old.Skill, func(o Skill) bool { return o.IDTag == skill.IDTag }) == -1 {
			tag := ""
			if skill.RefineBase == nil {
				tag = skill.IDTag
			} else {
				tag = *skill.RefineBase
			}

			find := make([]int, 0)
			for i := range l.Person {
				p := l.Person[i]
				if utils.IndexFunc(len(p.Skills), func(i int) bool {
					return utils.IndexFunc(len(p.Skills), func(n int) bool {
						return p.Skills[i][n] != nil && *p.Skills[i][n] == tag
					}) != -1
				}) != -1 {
					find = append(find, i)
				}
			}

			for _, i := range find {
				l.Person[i].RecentlyUpdate = true
			}
		}
	}

	// 新增了神装的人物

	for _, sub := range l.SubscriptionCostume {
		if utils.IndexFunc(len(old.SubscriptionCostume), func(i int) bool { return old.SubscriptionCostume[i].HeroID == sub.HeroID }) == -1 {
			f := utils.IndexFunc(len(l.Person), func(i int) bool { return l.Person[i].IDTag == sub.HeroID })
			l.Person[f].RecentlyUpdate = true
		}
		// if slices.utils.IndexFunc(old.SubscriptionCostume, func(sce SubscriptionCostumeElement) bool { return sub.HeroID == sce.HeroID }) == -1 {
		// 	f := slices.utils.IndexFunc(l.Person, func(p Person) bool { return p.IDTag == sub.HeroID })
		// 	l.Person[f].RecentlyUpdate = true
		// }
	}
}

func (l *GameData) delNone() {
	j := 0
	for _, val := range l.Person {
		if val.IDTag != "PID_無し" {
			l.Person[j] = val
			j++
		}
	}
	l.Person = l.Person[:j]

	j = 0
	for _, val := range l.Skill {
		if val.IDTag != "PID_無し" {
			l.Skill[j] = val
			j++
		}
	}
	l.Skill = l.Skill[:j]

	j = 0
	for _, val := range l.SkillAccessory {
		if val.IDTag != "PID_無し" {
			l.SkillAccessory[j] = val
			j++
		}
	}
	l.SkillAccessory = l.SkillAccessory[:j]
}

func (l *GameData) setRarity(useCache bool) error {

	defer func() {

		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	cache := "cache.html"
	if _, e := os.Stat(cache); os.IsNotExist(e) || !useCache {
		url := "https://feheroes.fandom.com/wiki/List_of_Heroes"

		req, _ := http.NewRequest("GET", url, nil)

		param := req.URL.Query()
		param.Add("Accept", "*/*")
		param.Add("Accept-Encoding", "gzip, deflate")
		param.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.33")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		f, err := os.Create(cache)
		if err != nil {
			return err
		}

		_, err = f.WriteString(string(data))
		if err != nil {
			return err
		}
	}
	f, err := os.Open(cache)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return err
	}
	heroes := doc.Find("#mw-content-text").Find("table").Find("tbody").Find(".hero-filter-element")
	nameRegex := regexp.MustCompile("^(.{1,}): (.{1,})$")
	rarityRegex := regexp.MustCompile("^([0-9])( – ){0,1}([0-9]){0,1}(.{0,})$")

	heroes.Each(func(i int, s *goquery.Selection) {

		allInfo := s.Find("td")
		capN := nameRegex.FindStringSubmatch(allInfo.Slice(1, 2).Text())
		name := capN[1]
		honour := capN[2]
		capR := rarityRegex.FindStringSubmatch(allInfo.Slice(6, 7).Text())
		minR := "0"
		maxR := "0"
		type_ := "0"
		if len(capR) >= 5 {
			minR = capR[1]
			maxR = capR[3]
			type_ = capR[4]
		}

		foundTag := ""
		for k, v := range l.TransEn {
			if v != nil && *v == honour && strings.HasPrefix(k, "MPID_HONOR_") {
				tag := strings.Split(k, "MPID_HONOR_")[1]
				if *l.TransEn[fmt.Sprintf("MPID_%s", tag)] == name {
					foundTag = tag
					break
				}
			}
		}
		if foundTag != "" {
			findP := utils.IndexFunc(len(l.Person), func(i int) bool { return l.Person[i].IDTag == fmt.Sprintf("PID_%s", foundTag) })
			// findP := slices.utils.IndexFunc(l.Person, func(p Person) bool { return p.IDTag == fmt.Sprintf("PID_%s", foundTag) })
			if findP != -1 {
				p := &l.Person[findP]
				min, err := strconv.ParseInt(minR, 10, 64)
				if err != nil {
					min = 0
				}
				max, err := strconv.ParseInt(maxR, 10, 64)
				if err != nil {
					max = 0
				}
				if min < max {
					p.MinRarity, p.MaxRarity = min, max
				} else {
					p.MinRarity, p.MaxRarity = max, min
				}

				if strings.Contains(type_, "Grand Hero Battle") {
					p.Type_ = 1
				} else if strings.Contains(type_, "Special") {
					p.Type_ = 2
				} else if strings.Contains(type_, "Story") {
					p.Type_ = 3
				} else if strings.Contains(type_, "Legendary") {
					p.Type_ = 4
				} else if strings.Contains(type_, "Mythic") {
					p.Type_ = 5
				} else if strings.Contains(type_, "Tempest Trials") {
					p.Type_ = 6
				} else {
					p.Type_ = 0
				}

			}

		}

	})

	return nil
}

//	func utils.IndexFunc2[E any](s []E, f func(i int) bool) int {
//		l := len(s)
//		for i := 0; i < l; i++ {
//			if f(i) {
//				return i
//			}
//		}
//		return -1
//	}
func (l *GameData) ToSql(now time.Time) error {

	date := now.Local().Format("20060102150405")

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("game_data_%s.db", date)), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return err
	}

	db.AutoMigrate(&models.Skill{}, &models.Person{})
	if err != nil {
		return err
	}

	db.Create(&l.Person)
	// 一次性写入会报错
	m := len(l.Skill) / 300
	p := len(l.Skill) % 300
	for i := 0; i < m; i++ {
		db.Create(l.Skill[i*300 : (i+1)*300])
	}
	if p != 0 {
		db.Create(l.Skill[m*300 : m*300+p])
	}

	return nil
}

func initRegex() ([][]string, [][][]*regexp.Regexp, error) {
	fmt.Println("初始化技能分类表达式")

	var by []byte

	src := make(map[string][][]any)

	keys := make([][]string, 0, len(SkillField))

	res := make([][][]*regexp.Regexp, 0, len(SkillField))

	_, err := os.Stat(`passive.yaml`)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// 文件不存在
			by = passiveRegex
			fmt.Println("passive.yaml不存在，写入预置数据到当前文件夹")
			f, err := os.Create(`passive.yaml`)
			if err != nil {
				return nil, nil, fmt.Errorf("写入预置技能分类数据到passive.yaml失败: %w", err)
			}
			defer f.Close()

			_, err = f.Write(by)
			if err != nil {
				return nil, nil, fmt.Errorf("写入预置技能分类数据到passive.yaml失败: %w", err)
			}

		} else {
			// 其他错误
			return nil, nil, fmt.Errorf("读取技能分类表达式失败: %w", err)
		}
	} else {
		// 文件存在
		by, err = os.ReadFile(`passive.yaml`)
		if err != nil {
			return nil, nil, fmt.Errorf("读取技能分类表达式失败: %w", err)
		}
	}

	err = yaml.Unmarshal(by, &src)
	if err != nil {
		return nil, nil, fmt.Errorf("加载技能分类表达式失败: %w", err)
	}

	for _, v1 := range SkillField {
		// t1     ["堡垒", ['^SID_.*の城塞\d$']],
		// ["一击", ['^SID_.*の一撃\d$', '^SID_.*の瞬撃\d$', '^SID_.*の迫撃\d$']],
		t1, ok := src[v1]

		categoryKeys := make([]string, 0)
		categoryRes := make([][]*regexp.Regexp, 0)

		if !ok {
			keys = append(keys, categoryKeys)
			res = append(res, categoryRes)
			continue
		}

		for _, v2 := range t1 {
			// v2 ["一击", ['^SID_.*の一撃\d$', '^SID_.*の瞬撃\d$', '^SID_.*の迫撃\d$']]
			key, ok := v2[0].(string)
			if !ok {
				return nil, nil, fmt.Errorf("初始化正则表达式失败")
			}

			reStrSliRaw, ok := v2[1].([]any)
			if !ok {
				return nil, nil, fmt.Errorf("初始化正则表达式失败")
			}

			reStrSli := make([]string, 0)

			for _, v := range reStrSliRaw {
				s, ok := v.(string)
				if !ok {
					return nil, nil, fmt.Errorf("初始化正则表达式失败")
				}
				reStrSli = append(reStrSli, s)

			}

			res_ := make([]*regexp.Regexp, 0, len(reStrSli))

			for _, re := range reStrSli {
				r, err := regexp.Compile(re)
				if err != nil {
					return nil, nil, fmt.Errorf("初始化正则表达式失败, %w", err)
				}
				res_ = append(res_, r)
			}

			// if result[v1] == nil {
			// 	result[v1] = make(map[string][]*regexp.Regexp)
			// }

			// result[v1][key] = res

			categoryKeys = append(categoryKeys, key)
			categoryRes = append(categoryRes, res_)
		}
		categoryKeys = append(categoryKeys, "其他")
		keys = append(keys, categoryKeys)
		res = append(res, categoryRes)

	}

	return keys, res, nil
}
