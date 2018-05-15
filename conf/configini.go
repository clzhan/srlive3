package conf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	AppName   string
	DebugMode bool
)

type section struct {
	name   string
	fields map[string]string
}

type Config struct {
	//文件io
	File *os.File
	//节点与键值对应
	Sections map[string]section
}

func (cfg *Config) ReadKeyValue(sectionName string, key string) (keyval string, err error) {
	if section, ok := cfg.Sections[sectionName]; ok {
		if keyval, ok := section.fields[key]; ok {
			return keyval, nil
		} else {
			return keyval, errors.New("key value not found")
		}
	} else {
		return keyval, errors.New("section  not found")
	}

	return keyval, nil
}

func (cfg *Config) initParam() (err error) {
	var value string

	if value, err = cfg.ReadKeyValue("App", "Name"); err != nil {
		AppName = "livedefault"
	} else {
		AppName = value
	}
	fmt.Println("value is :", AppName)

	if value, err = cfg.ReadKeyValue("App", "DebugMode"); err != nil {
		DebugMode = false
	} else {
		if value == "on" {
			DebugMode = true
		} else {
			DebugMode = false
		}
	}

	return nil
}

func (cfg *Config) LoadConfig(configurename string) (err error) {
	cfg.File, err = os.OpenFile(configurename, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer cfg.File.Close()

	cfg.Sections = make(map[string]section)

	var ls []string
	// 读取所有行
	bio := bufio.NewReader(cfg.File)
	for {
		line, _, err := bio.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("read line error")
		}
		ls = append(ls, string(line))
	}

	var sectionName string
	for _, v := range ls {
		//去掉左右空格
		v = strings.TrimSpace(v)

		//空行
		if v == "" {
			continue
		}

		//判断注释
		if strings.HasPrefix(v, "#") {
			continue
		}
		if strings.HasPrefix(v, "//") {
			continue
		}

		//判断section 取出sectionName
		if strings.HasPrefix(v, "[") {
			if strings.HasSuffix(v, "]") {
				sectionName = v[1 : len(v)-1]
				continue
			}
		}

		//取key 键
		index := strings.Index(v, "=")
		if index <= 0 {
			return errors.New("load config failed key in section error")
		}

		key := strings.TrimSpace(v[:index])
		value := strings.TrimSpace(v[index+1:])
		if len(value) == 0 {
			return errors.New("The key value is not set")
		}

		if _, ok := cfg.Sections[sectionName]; ok {
			cfg.Sections[sectionName].fields[key] = value
		} else {
			//新new一个
			sec := section{
				name: sectionName,
				fields: map[string]string{
					key: value,
				},
			}

			cfg.Sections[sectionName] = sec
		}

	}
	fmt.Println(cfg.Sections)
	if err = cfg.initParam(); err != nil {
		return err

	}

	return nil
}
