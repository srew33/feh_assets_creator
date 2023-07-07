package utils

import (
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"howett.net/plist"
)

func CreateFaces(src string, to string, base string) error {
	srcP := filepath.Join(src, "data/com.nintendo.zaba/files/assets/Common/Face")
	if _, err := os.Stat(srcP); os.IsNotExist(err) {
		return fmt.Errorf("源文件夹不存在: %s", srcP)
	}
	toP := filepath.Join(to, "update/faces")
	baseP := ""
	if base != "" {
		baseP = filepath.Join(base, `faces`)
	}

	dirCreated := false
	shouldCopy := false

	entries, err := os.ReadDir(srcP)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			dir := e.Name()
			targetFile := filepath.Join(toP, fmt.Sprintf("%s.webp", dir))
			if baseP == "" {
				shouldCopy = true
			} else {
				baseFile := filepath.Join(baseP, fmt.Sprintf("%s.webp", dir))
				if Exist(baseFile) {
					shouldCopy = false
				} else {
					shouldCopy = true
				}
			}

			if shouldCopy {
				if !dirCreated {
					err = os.MkdirAll(toP, 0750)
					if err != nil {
						return err
					}
					dirCreated = true
				}
				srcFile := filepath.Join(srcP, dir, "Face_FC.png")
				if _, err := os.Stat(srcFile); os.IsNotExist(err) {
					// ? 有的文件夹没有Face_FC.png文件，也许最好从加载好的结构体中直接读取?
					// return fmt.Errorf("源文件不存在: %s", srcFile)
					continue
				}
				r, err := os.Open(srcFile)
				if err != nil {
					return err
				}
				w, err := os.Create(targetFile)
				if err != nil {
					return err
				}
				_, err = io.Copy(w, r)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func CreateIcons(src string, to string, base string) error {
	srcP := filepath.Join(src, "data/com.nintendo.zaba/files/assets/Common/UI")
	if !Exist(srcP) {
		return fmt.Errorf("源文件夹不存在: %s", srcP)
	}
	toP := filepath.Join(to, "update/icons")
	baseP := ""
	if base != "" {
		baseP = filepath.Join(base, `icons`)
	}

	dirCreated := false

	compare := image.Rect(19, 19, 57, 57)
	writeConf := webp.Options{Lossless: false, Quality: 80.0}

	for srcId := 1; srcId < 100; srcId++ {
		srcFile := filepath.Join(srcP, fmt.Sprintf("Skill_Passive%d.png", srcId))
		if !Exist(srcFile) {
			break
		}

		m, err := DecodeWebp(srcFile)
		if err != nil {
			return err
		}

		for col := 0; col < 13; col++ {
			for row := 0; row < 13; row++ {
				iconId := (srcId-1)*169 + col*13 + row

				croped := imaging.Crop(*m, image.Rect(row*76, col*76, (row+1)*76, (col+1)*76))
				isNUll, err := IsNullPic(croped, compare)
				if err != nil {
					return err
				}
				if !isNUll {
					// 如果不给基准路径或该路径文件不存在，直接写入
					// 否则比较是否为相同图片后操作
					baseFile := filepath.Join(baseP, fmt.Sprintf("%d.webp", iconId))
					if baseP != "" && Exist(baseFile) {
						basePic, err := DecodeWebp(baseFile)
						if err != nil {
							return err
						}
						isSame, err := IsSamePic(*basePic, croped, compare)
						if err != nil {
							return err
						}
						if !isSame {
							if !dirCreated {
								err = os.MkdirAll(toP, 0755)
								if err != nil {
									return err
								}
								dirCreated = true
							}
							toFile := filepath.Join(toP, fmt.Sprintf("%d.webp", iconId))

							err = WriteWebp(toFile, croped, &writeConf)
							if err != nil {
								return err
							}
						}
					} else {
						if !dirCreated {
							err = os.MkdirAll(toP, 0755)
							if err != nil {
								return err
							}
							dirCreated = true
						}
						toFile := filepath.Join(toP, fmt.Sprintf("%d.webp", iconId))

						err = WriteWebp(toFile, croped, &writeConf)
						if err != nil {
							return err
						}
					}
				}

			}
		}
	}

	return nil
}

func CreateStatus(src string, to string, base string) error {
	srcP := filepath.Join(src, "data/com.nintendo.zaba/files/assets/Common/UI")
	if !Exist(srcP) {
		return fmt.Errorf("源文件夹不存在: %s", srcP)
	}
	toP := filepath.Join(to, "update")
	baseP := ""
	if base != "" {
		baseP = base
	}

	writeConf := webp.Options{Lossless: false, Quality: 80.0}

	f, err := os.Open(filepath.Join(srcP, "Status.plist"))
	if err != nil {
		return err
	}
	by, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	pList := make(map[string]any)

	_, err = plist.Unmarshal(by, &pList)
	if err != nil {
		return err
	}
	tasks := []string{"Move", "Weapon"}
	for _, t := range tasks {
		dirCreated := false
		srcImg, err := DecodeWebp(filepath.Join(srcP, "Status.png"))
		if err != nil {
			return err
		}

		frames, ok := pList["frames"].(map[string]any)
		if !ok {
			return fmt.Errorf("plist 数据解析错误")
		}
		frame, ok := frames[fmt.Sprintf("Icon_%s.png", t)].(map[string]any)
		if !ok {
			return fmt.Errorf("plist 数据解析错误")
		}

		crop, err := CropAsset(frame, *srcImg)
		if err != nil {
			return err
		}
		s := crop.Bounds().Size()
		width := s.X
		height := s.Y
		for i := 0; i < int(math.Round(float64(width)/float64(height))); i++ {
			target := filepath.Join(toP, strings.ToLower(t), fmt.Sprintf("%d.webp", i))
			baseFile := filepath.Join(baseP, strings.ToLower(t), fmt.Sprintf("%d.webp", i))
			if baseP == "" || !Exist(baseFile) {
				if !dirCreated {
					err = os.MkdirAll(filepath.Dir(target), 0755)
					if err != nil {
						return err
					}
					dirCreated = true
				}
				out := imaging.Crop(crop, image.Rect(i*height, 0, i*height+height, height))
				WriteWebp(target, out, &writeConf)
			}

		}

	}
	return nil
}

func CreateBlessing(src string, to string, base string) error {
	srcP := filepath.Join(src, "data/com.nintendo.zaba/files/assets/Common/UI")
	if !Exist(srcP) {
		return fmt.Errorf("源文件夹不存在: %s", srcP)
	}
	toP := filepath.Join(to, "update/blessing")
	baseP := ""
	if base != "" {
		baseP = filepath.Join(base, `blessing`)
	}

	dirCreated := false

	writeConf := webp.Options{Lossless: false, Quality: 80.0}

	tasks := []string{
		"Icon_SeasonNone.png",
		"Icon_BlessingFireS.png",
		"Icon_BlessingWaterS.png",
		"Icon_BlessingWindS.png",
		"Icon_BlessingEarthS.png",
		"Icon_BlessingLightS.png",
		"Icon_BlessingDarkS.png",
		"Icon_BlessingHeavenS.png",
		"Icon_BlessingLogicS.png",
	}

	for i, t := range tasks {
		f, err := os.Open(filepath.Join(srcP, "Blessing.plist"))
		if err != nil {
			return err
		}
		by, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		pList := make(map[string]any)

		_, err = plist.Unmarshal(by, &pList)
		if err != nil {
			return err
		}

		srcImg, err := DecodeWebp(filepath.Join(srcP, "Blessing.png"))
		if err != nil {
			return err
		}

		frames, ok := pList["frames"].(map[string]any)
		if !ok {
			return fmt.Errorf("plist 数据解析错误")
		}
		frame, ok := frames[t].(map[string]any)
		if !ok {
			return fmt.Errorf("plist 数据解析错误")
		}

		cropped, err := CropAsset(frame, *srcImg)
		if err != nil {
			return err
		}
		baseFile := filepath.Join(baseP, fmt.Sprintf("%d.webp", i))
		if !Exist(baseFile) {
			target := filepath.Join(toP, fmt.Sprintf("%d.webp", i))
			if !Exist(target) {
				if !dirCreated {
					err = os.MkdirAll(toP, 0755)
					if err != nil {
						return err
					}
					dirCreated = true
				}

				err = WriteWebp(target, cropped, &writeConf)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func CreateAssets(src string, to string, base string) error {
	err := CreateFaces(src, to, base)
	if err != nil {
		return err
	}
	err = CreateIcons(src, to, base)
	if err != nil {
		return err
	}
	err = CreateStatus(src, to, base)
	if err != nil {
		return err
	}
	err = CreateBlessing(src, to, base)
	if err != nil {
		return err
	}

	return nil
}
