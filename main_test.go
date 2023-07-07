package main

import (
	"feh_assets_creator/structs"
	"feh_assets_creator/utils"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestM(t *testing.T) {
	config := structs.TaskConfig{
		NewVersionPath: `H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0707a_summer.zip`,
		OldVersionPath: `H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0706a_fuin.zip`,

		UseCache:    false,
		ToSql:       true,
		ParseRarity: true,
		// BasePath:    `H:\GitProject\flutter\feh\feh_assets\baseline`,
		// SrcPath:     `I:`,
	}

	// err := structs.CreateBin(`H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0609b_abc.zip`,
	// 	`H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0608c_legend.zip`, option)
	if utils.Exist("output") {
		err := os.RemoveAll("output")
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 5)
		}
	}

	ti := time.Now()

	// if newPath == "" || oldPath == "" {
	// 	fmt.Println("新旧版本文件地址不能为空")
	// 	time.Sleep(time.Second * 5)
	// 	os.Exit(-1)
	// }
	o := structs.Output{
		MinSupVersion: MIN_SUPPORT_VERSION,
		Version:       ti.UnixMilli(),
	}
	// o := structs.Output{MinSupVersion: 60}
	err := o.Create(config, ti)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 3)
		return
	}

	elapsed := time.Since(ti)
	fmt.Printf("数据处理完成，输出版本：%d，运行用时：%s\n", ti.UnixMilli(), elapsed)

	if config.SrcPath != "" {
		fmt.Println("开始制作资源文件")
		err = utils.CreateAssets(config.SrcPath, "output", config.BasePath)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 3)
			return
		}
	}

	zipFile := fmt.Sprintf("update_%s.zip", ti.Format("20060102"))

	err = structs.Zip("output", zipFile, "output/data.bin")
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 3)
		return
	}

	elapsed = time.Since(ti)
	fmt.Println("任务完成，总共用时:", elapsed)
	// time.Sleep(time.Second * 20)
}
