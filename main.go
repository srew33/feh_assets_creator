package main

import (
	"feh_assets_creator/structs"
	"feh_assets_creator/utils"
	"flag"
	"fmt"
	"os"
	"time"
)

var config = structs.TaskConfig{}

var MIN_SUPPORT_VERSION = 50

func registFlags() {
	flag.StringVar(&config.NewVersionPath, "new", "", "新版本zip文件地址（必填）")
	flag.StringVar(&config.NewVersionPath, "n", "", "新版本zip文件地址（必填）")
	flag.StringVar(&config.OldVersionPath, "old", "", "旧版本zip文件地址")
	flag.StringVar(&config.OldVersionPath, "o", "", "旧版本zip文件地址")
	flag.BoolVar(&config.ParseRarity, "parseRarity", true, "是否从网络抓取稀有度，默认开启，关闭-parseRarity=false || -p=false")
	flag.BoolVar(&config.UseCache, "useCache", false, "是否使用已缓存文件解析稀有度，默认关闭")
	flag.BoolVar(&config.ParseRarity, "p", true, "是否从网络抓取稀有度，默认开启")
	flag.BoolVar(&config.UseCache, "u", false, "是否使用已缓存文件解析稀有度，默认关闭")
	flag.BoolVar(&config.ToSql, "S", false, "是否在当前文件夹同时输出一份人物和技能的sqlite数据库， 默认关闭")
	flag.BoolVar(&config.ToSql, "sql", false, "是否在当前文件夹同时输出一份人物和技能的sqlite数据库， 默认关闭")

	flag.StringVar(&config.SrcPath, "src", "", "挂载的资源文件所在的盘符,eg: i: ")
	flag.StringVar(&config.BasePath, "base", "", "基准资源文件文件夹")
	flag.StringVar(&config.SrcPath, "s", "", "挂载的资源文件所在的盘符,eg: i: ")
	flag.StringVar(&config.BasePath, "b", "", "基准资源文件文件夹")

	flag.Usage = func() {
		flagSet := flag.CommandLine
		order := []string{"new", "n", "old", "o", "parseRarity", "p", "useCache", "u", "S", "sql", "src", "s", "base", "b"}
		l := len(order)
		for i := 0; i < l/2; i++ {
			flag := flagSet.Lookup(order[i*2])
			fmt.Printf("-%s, -%s\n", order[i*2], order[i*2+1])
			fmt.Printf("  %s\n", flag.Usage)
		}
		fmt.Printf("!!!!!!!!!!!生成文件时会删除当前目录下的output文件夹，请注意!!!!!!!!!!!")
		// for _, name := range order {
		// 	flag := flagSet.Lookup(name)
		// 	fmt.Printf("-%s\n", flag.Name)
		// 	fmt.Printf("  %s\n", flag.Usage)
		// }
	}

}

func main() {
	registFlags()
	flag.Parse()

	// err := structs.CreateBin(`H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0609b_abc.zip`,
	// 	`H:\GitProject\flutter\feh\feh_update_creator\data\sources\feh-assets-json-0608c_legend.zip`, option)
	if utils.Exist("output") {
		err := os.RemoveAll("output")
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 5)
		}
	}

	t := time.Now()

	// if newPath == "" || oldPath == "" {
	// 	fmt.Println("新旧版本文件地址不能为空")
	// 	time.Sleep(time.Second * 5)
	// 	os.Exit(-1)
	// }
	o := structs.Output{
		MinSupVersion: MIN_SUPPORT_VERSION,
		Version:       t.UnixMilli(),
	}
	err := o.Create(config, t)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 3)
		return
	}

	elapsed := time.Since(t)
	fmt.Printf("数据处理完成，输出版本：%d，运行用时：%s\n", t.UnixMilli(), elapsed)

	if config.SrcPath != "" {
		fmt.Println("开始制作图标")
		err = utils.CreateAssets(config.SrcPath, "output", config.BasePath)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 3)
			return
		}
	}

	zipFile := fmt.Sprintf("update_%s.zip", t.Format("20060102"))

	err = structs.Zip("output", zipFile, "output/data.bin")
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 3)
		return
	}

	elapsed = time.Since(t)
	fmt.Println("任务完成，总共用时:", elapsed)

}
