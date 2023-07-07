package structs

import (
	"compress/zlib"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
)

type Output struct {
	MinSupVersion int     `json:"minimal_support_version"`
	SembastExport int     `json:"sembast_export"`
	Version       int64   `json:"version"`
	Stores        []Store `json:"stores"`
}

type Store struct {
	Name   string            `json:"name"`
	Keys   []string          `json:"keys"`
	Values []json.RawMessage `json:"values"`
}

type TaskConfig struct {
	OldVersionPath string
	NewVersionPath string
	ParseRarity    bool
	UseCache       bool
	ToSql          bool
	SrcPath        string
	BasePath       string
}

func (t *Output) Create(config TaskConfig, now time.Time) error {
	if config.NewVersionPath == "" {
		return fmt.Errorf("没有设置新版本文件地址")
	}

	newData := GameData{}
	oldData := GameData{}

	err := load(config.NewVersionPath, &newData)
	if err != nil {
		return err
	}
	fmt.Println("数据加载完成")

	newData.delNone()

	newData.setTranslations()

	newData.setPinyin()

	newData.setSubscription()

	newData.setRedeemable()

	err = newData.setBstAndDefault()

	if err != nil {
		return err
	}
	err = newData.setAccessorySkill()

	if err != nil {
		return err
	}
	err = newData.setRefineList()

	if err != nil {
		return err
	}
	newData.appendBlessing()

	if config.OldVersionPath != "" {
		err = load(config.OldVersionPath, &oldData)
		if err != nil {
			return err
		}
		newData.setRecently(oldData)
	}

	if config.ParseRarity {
		err = newData.setRarity(config.UseCache)
		if err != nil {
			return err
		}
	}
	err = newData.setSkillRarity()

	if err != nil {
		return err
	}

	err = newData.setSkillSeries()

	if err != nil {
		return err
	}

	if config.ToSql {
		err = newData.ToSql(now)
		if err != nil {
			return err
		}
	}

	t.SembastExport = 1

	t.Stores = append(t.Stores, Store{
		Name: "person",
		Keys: Iter(len(newData.Person), func(i int) string {
			return newData.Person[i].IDTag
		}),
		Values: Iter(len(newData.Person), func(i int) json.RawMessage {
			by, _ := json.Marshal(&newData.Person[i])
			return by
		}),
	})
	t.Stores = append(t.Stores, Store{
		Name: "skill",
		Keys: Iter(len(newData.Skill), func(i int) string {
			return newData.Skill[i].IDTag
		}),
		Values: Iter(len(newData.Skill), func(i int) json.RawMessage {
			by, _ := json.Marshal(&newData.Skill[i])
			return by
		}),
	})
	t.Stores = append(t.Stores, Store{
		Name: "weaponType",
		Keys: Iter(len(newData.Weapon), func(i int) string {
			return newData.Weapon[i].IDTag
		}),
		Values: Iter(len(newData.Weapon), func(i int) json.RawMessage {
			by, _ := json.Marshal(&newData.Weapon[i])
			return by
		}),
	})
	t.Stores = append(t.Stores, Store{
		Name: "translations",
		Keys: []string{"en", "ja", "zh"},
		Values: Iter(3, func(i int) json.RawMessage {
			by := []byte{}
			switch i {
			case 0:
				by, _ = json.Marshal(&newData.TransEn)
			case 1:
				by, _ = json.Marshal(&newData.TransJp)
			case 2:
				by, _ = json.Marshal(&newData.TransCn)
			}
			return by
		}),
	})

	t.Stores = append(t.Stores, Store{
		Name: "skillSeries",
		Keys: SkillField,
		Values: Iter(len(SkillField), func(i int) json.RawMessage {
			by, _ := json.Marshal(&newData.SkillSeries[i])
			return by
		}),
	})

	os.Mkdir("output", os.ModePerm)

	f, err := os.Create("output/data.bin")
	if err != nil {
		return err
	}
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	// f.WriteString(string(j))

	w := zlib.NewWriter(f)
	w.Write(j)
	w.Close()
	f.Close()

	return nil

}
