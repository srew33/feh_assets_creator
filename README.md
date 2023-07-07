# feh_assets_creator

与[srew33/feh_rebuilder: 使用FLUTTER对FEH BUILDER的重置 (github.com)](https://github.com/srew33/feh_rebuilder) 配套的项目，用来生成所需的数据和附件。

# 开发环境要求

1. go语言
   
   1.18+

2. cgo
   
   由于资源文件的生成需要进行webp文件的读写，而现有webp库全部依赖CGO

# 参数说明

```
-new, -n
  新版本zip文件地址（必填）
-old, -o
  旧版本zip文件地址
-parseRarity, -p
  是否从网络抓取稀有度，默认开启，关闭-parseRarity=false | -p=false
-useCache, -u
  是否使用已缓存文件解析稀有度，默认关闭
-S, -sql
  是否在当前文件夹同时输出一份人物和技能的sqlite数据库， 默认关闭  
-src, -s
  挂载的资源文件所在的盘符,eg: i:\ 
-base, -b
  旧版本资源文件文件夹
```

# 使用说明

1. 通过实机或模拟器确认当前有效的卡池

2. 从[HertzDevil/feh-assets-json: JSON dumps of Fire Emblem Heroes asset files (github.com)](https://github.com/HertzDevil/feh-assets-json)下载对应的数据，关于旧版本的数据选择，见FAQ

3. 按照 [参数说明](#参数说明 ) 设置好所需参数运行，比如：

```bash
feh_assets_creator.exe -n xxx\feh-assets-json-0703c_legend.zip -o xxx\feh-assets-json-0701c_legend.zip -p -useCache -S -s i:
```

生成的文件会在当前文件夹下，文件名为update_xxx(当前日期).zip

4. 如果需要生成资源文件，需要提前挂载好模拟器的数据文件，详见FAQ

5. 如果想查看或分析数据，可以使用-S或-sql参数，程序运行完成后会在当前文件夹下生成一个sqlite3的数据库，文件名为game_data_xxxxx.db，内含人物和技能两个数据表

# FAQ

1. 网络不好时，抓取稀有度时会一直等待或报错
   
   使用各种方式通过浏览器访问 [List_of_Heroes](https://feheroes.fandom.com/wiki/List_of_Heroes) ，并将网页保存到运行路径的同一文件夹，文件名为cache.html，并在运行程序时添加-u 参数。
   在首次运行成功后程序也会将本次抓取的数据写入到这个文件中。

2. 如何生成资源文件
   
   2.1 先更新好游戏版本，然后按照[Fire Emblem Heroes Wiki:Extracting game assets - Fire Emblem Heroes Wiki (fandom.com)](https://feheroes.fandom.com/wiki/Fire_Emblem_Heroes_Wiki:Extracting_game_assets)的说明挂载好模拟器的数据盘
   
   2.2 使用-src或-s 参数设置对应的盘符，比如i:\
   
   2.3 -base或-b参数可以设置基线版本文件夹路径，程序运行时会比对新版本和基准版本的数据差异，只输出最新版本的资源文件，第一次运行时可以不设置基线版本文件夹，此时会输出全部的附件，作为初始化的资源文件和以后的基准版本

3. 如何选择旧版本文件
   
   程序通过比对新旧版本数据来确定哪些技能和人物是最新的，旧版本数据为当前开放的新池的上一个版本数据。
   例： 假设现有A、B、C三个卡池，其中BC两个池子正在开放，A池在C池开放时关闭，则最新的数据为B、C两个池子的数据，旧版本数据下载A池的数据文件，新版本数据下载C池的数据文件。

4. 技能系列字段的设置
   
   目前通过正则表达式的方式来设置每个技能所属的系列，正则表达式文件使用yaml格式，预置的规则会在第一次运行后写入当前文件夹下的passive.yaml文件，可以根据自己的需求自行定制
