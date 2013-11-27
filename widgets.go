package neyo

import (
	"errors"
	"github.com/wendal/mustache"
	"os"
	"path/filepath"
)

var (
	WidgetBuilders = make(map[string]WidgetBuilder)
)

type WidgetBuilder func(Mapper, mustache.Context) (Widget, error)

type Widget interface {
	Prepare(mapper Mapper, ctx mustache.Context) Mapper
}

func LoadWidgets(topCtx mustache.Context) ([]Widget, string, error) {
	widgets := make([]Widget, 0)
	assets := ""

	err := filepath.Walk("widgets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		cnf_path := path + "/config.yml"
		fst, err := os.Stat(cnf_path)
		if err != nil || fst.IsDir() {
			return nil //ignore
		}
		cnf, err := ReadYml(cnf_path)
		if err != nil {
			return errors.New(cnf_path + ":" + err.Error())
		}
		if cnf["layout"] != nil {
			widget_enable, ok := cnf["layout"].(bool)
			if ok && !widget_enable {
				Log(INFO, "Disable >", cnf_path)
			}
		}
		builderFunc := WidgetBuilders[info.Name()]
		if builderFunc == nil { // 看看是否符合自定义挂件的格式
			_widget, _assets, _err := BuildCustomWidget(info.Name(), path, cnf)
			if _err != nil {
				Log(ERROR, "No WidgetBuilder %s %s", cnf_path, _err)
			}
			if _widget != nil {
				widgets = append(widgets, _widget)
				if _assets != nil {
					for _, asset := range _assets {
						assets += asset + "\n"
					}
				}
			}
			return nil
		}
		widget, err := builderFunc(cnf, topCtx)
		if err != nil {
			return err
		}
		widgets = append(widgets, widget)
		Log(DEBUG, "Load widget %s", cnf_path)
		return nil
	})
	return widgets, assets, err
}

func PrapareWidgets(widgets []Widget, mapper Mapper, topCtx mustache.Context) mustache.Context {
	mappers := make([]interface{}, 0)
	for _, widget := range widgets {
		mr := widget.Prepare(mapper, topCtx)
		if mr != nil {
			for k, v := range mr {
				mapper[k] = v
			}
			mappers = append(mappers, mr)
		}
	}
	return mustache.MakeContexts(mappers...)
}
