package conf

import (
	"os"
	"path/filepath"
	"strings"
)

/**
 * 加载运行路径中的全部INI文件
 */
func Load() (map[string]string, error) {
	confs := make(map[string]string)
	// 遍历当前程序运行路径，读取所有INI文件
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		lowCaseName := strings.TrimSpace(strings.ToLower(info.Name()))
		if !strings.HasSuffix(lowCaseName, ".ini") {
			return nil
		}
		// 解析当前配置文件
		err = parseFile(path, confs)
		return err
	})
	return confs, nil
}

/**
 * 加载指定路径的配置文件
 */
func LoadFile(path string) (map[string]string, error) {
	confs := make(map[string]string)
	if err := parseFile(path, confs); err != nil {
		return nil, err
	}
	return confs, nil
}

// 解析指定路径的文件
// 将解析结果保存至confs，使用key、value存储
func parseFile(path string, confs map[string]string) error {

	return nil
}
