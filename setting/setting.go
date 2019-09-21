package setting

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Cfg *ini.File

	// RunMode  string
	ConfPath string
)


func init() {

	ConfPath = ConfFile()
	Cfg = LoadConf(ConfPath)

	// LoadRunMode(Cfg)

	// fmt.Printf("\033[0;36m--------success to parse 'conf/env.ini'------\n\033[0m")
}


// 设置配置文件路径（绝对路径）
func SetConfFile(absPath string) string {
	ConfPath = absPath
	return ConfPath
}

// 获取配置文件路径
func ConfFile() string {
	// 可执行文件绝对路径
	dir := GetCurrentDirectory()
	// log.Printf(dir)
	// /home/fky/workspace/gitrep/storage/nooiecloud/bin
	// absPath := ConfPath + "/../conf/env.ini"

	absPath := dir + "/conf/rtspserver.ini"
	// absPath := GetParentDirectory(dir) + "/conf/rtspserver.ini"

	ConfPath = SetConfFile(absPath)
	return ConfPath
}

// 加载配置文件
func LoadConf(confFile string) *ini.File {
	Cfg, err := ini.Load(confFile)
	if err != nil {
		log.Fatalf("\033[1;31m--------Fail to parse %s, error: %v------\033[0m", confFile, err)
	}
	return Cfg
}

func Conf() *ini.File {
	return Cfg
}
/*
func LoadRunMode(Cfg *ini.File) {
	sec, err := Cfg.GetSection("runmode")
	if err != nil {
		log.Fatalf("Fail to get section 'runmode': %v", err)
	}

	key, err := sec.GetKey("mode")
	if err != nil {
		log.Fatalf("\033[1;31mFail to get key 'mode': %v\033[0m", err)
	}
	RunMode = key.String()

	fmt.Printf("load runmode setting, mode:%v\n", RunMode)
}
*/
func GetCurrentDirectory() string {

	/*
		// 0 表示调用runtime.Caller()所在的位置，1表示runtime.Caller()所在函数的调用位置，依此类推
		// 所以写死1则始终得到的是调用CurrentFile()所在的位置，此函数能在任意调用
		_, file, _, ok := runtime.Caller(1)
		if !ok {
			panic(errors.New("Can not get current file info"))
		}
		log.Printf(file)
		// 0 --> /home/fky/workspace/gitrep/storage/nooiecloud/src/setting/setting.go
		// 1 --> /usr/local/go/src/runtime/proc.go

		// 可执行文件绝对路径
		file = GetCurrentDirectory()
		log.Printf(file)
		// /home/fky/workspace/gitrep/storage/nooiecloud/bin

		// 可执行文件绝对路径文件名
		file, _ = exec.LookPath(os.Args[0])
		log.Printf(file)
		// /home/fky/workspace/gitrep/storage/nooiecloud/bin/nooiecloud
	*/

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetParentDirectory(directory string) string {
	return Substr(directory, 0, strings.LastIndex(directory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}
