package utils

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	gowebp "golang.org/x/image/webp"
)

func Exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func IndexFunc(length int, f func(i int) bool) int {
	for i := 0; i < length; i++ {
		if f(i) {
			return i
		}
	}
	return -1
}

func WriteWebp(path string, img image.Image, opt *webp.Options) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = webp.Encode(f, img, opt)
	if err != nil {
		return err
	}

	return nil

}
func DecodeWebp(path string) (*image.Image, error) {

	if Exist(path) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r, err := gowebp.Decode(f)
		if err != nil {
			return nil, err
		}
		return &r, nil
	}

	return nil, fmt.Errorf("解析webp错误: %s不存在", path)

}

var SAME_PIC_THRESHOLD int = 5

var NULL_PIC_THRESHOLD int = 5
var NULL_HASH *goimagehash.ImageHash

func init() {
	n, err := base64.StdEncoding.DecodeString("If+BAwEBAUQB/4IAAQIBBEhhc2gBBgABBEtpbmQBBAAAAA//ggH4AQAAAAD///8BAgA=")
	if err != nil {
		panic(err)
	}
	null, err := goimagehash.LoadImageHash(bytes.NewBuffer(n))
	if err != nil {
		panic(err)
	}
	NULL_HASH = null
}

func IsSamePic(img1 image.Image, img2 image.Image, compare image.Rectangle) (bool, error) {
	croped1 := imaging.Crop(img1, compare)
	croped2 := imaging.Crop(img2, compare)

	h1, err := goimagehash.AverageHash(croped1)
	if err != nil {
		return false, err
	}
	h2, err := goimagehash.AverageHash(croped2)
	if err != nil {
		return false, err
	}
	dist, err := h1.Distance(h2)
	if err != nil {
		return false, err
	}
	if dist <= SAME_PIC_THRESHOLD {
		return true, nil
	}
	return false, nil

}

func IsNullPic(img image.Image, compare image.Rectangle) (bool, error) {
	croped1 := imaging.Crop(img, compare)

	h, err := goimagehash.AverageHash(croped1)
	if err != nil {
		return false, err
	}

	dist, err := h.Distance(NULL_HASH)
	if err != nil {
		return false, err
	}
	if dist >= SAME_PIC_THRESHOLD {
		return false, nil
	}
	return true, nil
}

func CropAsset(frame map[string]any, src image.Image) (*image.NRGBA, error) {
	rotated, ok := frame["textureRotated"].(bool)
	if !ok {
		return nil, fmt.Errorf("plist 数据解析错误")
	}
	spriteSourceSizeStr, ok := frame["spriteSourceSize"].(string)
	if !ok {
		return nil, fmt.Errorf("plist 数据解析错误")
	}
	spriteSourceSize, err := str2ints(spriteSourceSizeStr)
	if err != nil {
		return nil, fmt.Errorf("plist 数据解析错误")
	}

	textureRectStr, ok := frame["textureRect"].(string)
	if !ok {
		return nil, fmt.Errorf("plist 数据解析错误")
	}
	textureRect, err := str2ints(textureRectStr)
	if err != nil {
		return nil, fmt.Errorf("plist 数据解析错误")
	}

	spriteOffsetStr, ok := frame["spriteOffset"].(string)
	if !ok {
		return nil, fmt.Errorf("plist 数据解析错误")
	}
	spriteOffset, err := str2ints(spriteOffsetStr)
	if err != nil {
		return nil, fmt.Errorf("plist 数据解析错误")
	}
	// spriteSourceSize和spriteSize的关系似乎是裁剪而不是拉伸，
	// 所以这里要计算差值，注意计算的时候插值要除2，保持图像居中
	// ! 这里直接取索引最好用get来判断越界问题
	// ! 注意这里的delta需要用未涉及rotated的原值，要在rotated转换之前计算
	delta_w := spriteSourceSize[0] - textureRect[2]
	delta_h := spriteSourceSize[1] - textureRect[3]

	if rotated {
		textureRect = []int{
			textureRect[0],
			textureRect[1],
			textureRect[3],
			textureRect[2],
		}
	}
	// (texture_rect[0] - sprite_offset[0]) as u32,
	// (texture_rect[1] - sprite_offset[1]) as u32,
	// texture_rect[2] as u32,
	// texture_rect[3] as u32,
	cropped := imaging.Crop(src, image.Rect(textureRect[0]-spriteOffset[0], textureRect[1]-spriteOffset[1],
		textureRect[0]-spriteOffset[0]+textureRect[2], textureRect[1]-spriteOffset[1]+textureRect[3]))

	if rotated {
		cropped = imaging.Rotate90(cropped)
	}

	out := image.NewNRGBA(image.Rect(0, 0, spriteSourceSize[0], spriteSourceSize[1]))
	out = imaging.Paste(out, cropped, image.Point{delta_w / 2, delta_h / 2})
	return out, nil

}

func str2ints(src string) ([]int, error) {
	r := []int{}
	l := strings.Split(strings.ReplaceAll(strings.ReplaceAll(src, "{", ""), "}", ""), ",")
	for _, v := range l {
		p, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return []int{}, err
		}
		r = append(r, int(p))
	}
	if len(r) == 0 {
		return []int{}, fmt.Errorf("plist 数据解析错误")
	}
	return r, nil

}

func read(f *zip.File) ([]byte, error) {

	buf, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer buf.Close()
	data, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return data, nil
}
