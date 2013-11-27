package neyo

import (
	"bytes"
	"encoding/json"
	"github.com/wendal/goyaml2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 读取配置文件()
func ReadConfig(root string) (cnf map[string]interface{}, err error) {
	cnf, err = ReadYmlCnf(root)
	return
}

// 读取YAML格式的配置文件(兼容JSON)
func ReadYmlCnf(root string) (map[string]interface{}, error) {
	path := filepath.Join(root, CONFIG_YAML)
	return ReadYml(path)
}

// 从文件读取YAML
func ReadYml(path string) (cnf map[string]interface{}, err error) {
	Log(DEBUG, "Read %s", path)
	err = nil
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	cnf, err = ReadYmlReader(f)
	return
}

// 从Reader读取YAML
func ReadYmlReader(r io.Reader) (cnf map[string]interface{}, err error) {

	err = nil
	buf, err := ioutil.ReadAll(r)
	if err != nil || len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" {
		Log(DEBUG, "\tLook lile a json, try it")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			Log(DEBUG, "\tIt is json map")
			return
		}
	}

	_map, _err := goyaml2.Read(bytes.NewBuffer(buf))
	if _err != nil {
		Log(ERROR, "goyaml2 ", string(buf), _err)
		//err = goyaml.Unmarshal(buf, &cnf)
		err = _err
		return
	}
	if _map == nil {
		Log(INFO, "goyaml2 output nil? Pls report bug\n%s", string(buf))
	}
	cnf, ok := _map.(map[string]interface{})
	if !ok {
		Log(INFO, "Not a Map? >> %s %s", string(buf), _map)
		cnf = nil
	}
	return
}
