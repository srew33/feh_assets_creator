package structs

import (
	"archive/zip"
	"feh_assets_creator/models"
	"feh_assets_creator/tasks"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
)

func readFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	buffer := make([]byte, f.FileHeader.UncompressedSize64)
	_, err = io.ReadFull(rc, buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func filterSome(iterable []*zip.File, check func(i int) bool) []*zip.File {
	length := len(iterable)

	r := make([]*zip.File, 0, length)

	for i := 0; i < length; i++ {
		if check(i) {
			r = append(r, iterable[i])
		}
	}
	return r
}

func filtOne(iterable []*zip.File, check func(i int) bool) (*zip.File, bool) {
	length := len(iterable)

	for i := 0; i < length; i++ {
		if check(i) {
			return iterable[i], true
		}
	}
	return &zip.File{}, false
}

func setData(g *GameData, taskType int, data []byte) error {
	var err error
	switch taskType {
	case tasks.Move:
		err = json.Unmarshal(data, &g.Move)
	case tasks.Weapon:
		err = json.Unmarshal(data, &g.Weapon)
	case tasks.Person:
		t := []models.Person{}
		err = json.Unmarshal(data, &t)
		g.Person = append(g.Person, t...)
	case tasks.Skill:
		t := []models.Skill{}
		err = json.Unmarshal(data, &t)
		g.Skill = append(g.Skill, t...)
	case tasks.WeaponRefine:
		t := []models.WeaponRefine{}
		err = json.Unmarshal(data, &t)
		g.WeaponRefine = append(g.WeaponRefine, t...)
	case tasks.SkillAccessory:
		t := []models.SkillAccessory{}
		err = json.Unmarshal(data, &t)
		g.SkillAccessory = append(g.SkillAccessory, t...)
	case tasks.SubscriptionCostume:
		t := []models.SubscriptionCostumeElement{}
		err = json.Unmarshal(data, &t)
		g.SubscriptionCostume = append(g.SubscriptionCostume, t...)
	case tasks.HolyGrail:
		t := []models.HolyGrail{}
		err = json.Unmarshal(data, &t)
		g.HolyGrail = append(g.HolyGrail, t...)
	case tasks.TransJp:
		t := []models.TransElement{}
		err = json.Unmarshal(data, &t)
		if g.TransJp == nil {
			g.TransJp = make(map[string]*string)
		}
		for _, e := range t {
			g.TransJp[e.Key] = e.Val
		}
	case tasks.TransEn:
		t := []models.TransElement{}
		err = json.Unmarshal(data, &t)
		if g.TransEn == nil {
			g.TransEn = make(map[string]*string)
		}
		for _, e := range t {
			g.TransEn[e.Key] = e.Val
		}
	case tasks.TransCn:
		t := []models.TransElement{}
		err = json.Unmarshal(data, &t)
		if g.TransCn == nil {
			g.TransCn = make(map[string]*string)
		}
		for _, e := range t {
			g.TransCn[e.Key] = e.Val
		}
	}
	return err
}

func Iter[T any](length int, k func(i int) T) []T {
	// r := []T{}
	r := make([]T, 0, length)
	for i := 0; i < length; i++ {
		r = append(r, k(i))
	}
	return r
}

func Zip(src string, dest string, gamaData string) error {
	os.Remove(dest)

	zipfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			err = compress(e, src, "", archive)
			if err != nil {
				return err
			}
		}

	}
	// 单独压入data.bin
	gameFile, err := os.Open(gamaData)
	if err != nil {
		return err
	}
	defer gameFile.Close()
	stat, err := gameFile.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}
	header.Name = "data.bin"

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, gameFile)

	if err != nil {
		return err
	}

	return nil
}

// 遍历压入src 文件夹，src的第一层只压缩文件夹而忽略文件
func compress(file fs.DirEntry, src string, rel string, zw *zip.Writer) error {
	// 如果要压缩整个文件夹，去掉if rel != ""
	es, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range es {
		if e.IsDir() {
			rel := filepath.Join(rel, e.Name())
			// fmt.Println(rel)
			if err != nil {
				return err
			}
			compress(e, filepath.Join(src, e.Name()), rel, zw)
		} else if rel != "" {
			relName := filepath.Join(rel, e.Name())
			absName := filepath.Join(src, e.Name())
			// fmt.Println(relName)
			info, err := file.Info()
			if err != nil {
				return err
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = relName

			writer, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}
			f, err := os.Open(absName)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, f)
			f.Close()
			if err != nil {
				return err
			}

		}
	}

	return nil

}
func load(path string, gameData *GameData) error {

	c, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer c.Close()

	f := c.File

	root := filepath.SplitList(f[0].Name)[0]

	for taskType, taskPath := range tasks.FileTasks {
		found, ok := filtOne(f, func(i int) bool {
			return filepath.ToSlash(f[i].Name) == filepath.ToSlash(filepath.Join(root, taskPath))
		})
		if !ok {
			return fmt.Errorf("没有找到 %s 文件，请检查源文件", taskPath)
		}
		data, err := readFile(found)
		if err != nil {
			return err
		}

		err = setData(gameData, taskType, data)
		if err != nil {
			return err
		}

	}

	for taskType, taskPath := range tasks.DirectoryTasks {
		found := filterSome(f, func(i int) bool {
			return !f[i].FileInfo().IsDir() && filepath.ToSlash(filepath.Dir(f[i].Name)) == filepath.ToSlash(filepath.Join(root, taskPath))
		})

		for i := range found {
			data, err := readFile(found[i])
			if err != nil {
				return err
			}

			err = setData(gameData, taskType, data)
			if err != nil {
				return err
			}
		}

	}

	return nil

}
