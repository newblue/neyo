package neyo

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// 存在核心配置文件的路径,才可能是Gor的目录
func IsYoProjectDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, CONFIG_YAML))
	return err == nil
}

// 以Json方式打印对象,方便调试
func PrintJson(v interface{}) {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		Log(ERROR, "Json marshal: %s", err)
	} else {
		Log(INFO, "PrintJson <<END\n%s\nEND\n", string(buf))
	}
}

func Copy(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		return err
	}
	return d.Close()
}
