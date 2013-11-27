package neyo

import (
	"errors"
	"fmt"
	"github.com/wendal/mustache"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CustomWidget struct {
	name   string
	layout *DocContent
	mapper Mapper
}

func (c *CustomWidget) Prepare(mapper Mapper, ctx mustache.Context) Mapper {
	return Mapper(map[string]interface{}{c.name: c.layout.Source})
}

func BuildCustomWidget(name string, dir string, cnf Mapper) (Widget, []string, error) {
	layoutName, ok := cnf["layout"]
	if !ok || layoutName == "" {
		Log(WARN, "Skip Widget: %s ", dir)
		return nil, nil, nil
	}

	layoutFilePath := filepath.Join(dir, "/layouts/", layoutName.(string)+".html")
	f, err := os.Open(layoutFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to load widget layout %s %s", dir, err.Error())
		return nil, nil, errors.New(errMsg)
	}
	defer f.Close()
	cont, err := ioutil.ReadAll(f)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to load widget layout %s %s", dir, err.Error())
		return nil, nil, errors.New(errMsg)
	}

	assets := []string{}
	for _, js := range cnf.GetStrings("javascripts") {
		path := filepath.Join("/assets/", dir, "/javascripts/", js)
		assets = append(assets, fmt.Sprintf("<script type=\"text/javascript\" src=\"%s\"></script>", path))
	}
	for _, css := range cnf.GetStrings("stylesheets") {
		path2 := filepath.Join("/assets/", dir, "/stylesheets/", css)
		assets = append(assets, fmt.Sprintf("<link href=\"%s\" type=\"text/css\" rel=\"stylesheet\" media=\"all\">", path2))
	}

	return &CustomWidget{name, &DocContent{string(cont), string(cont), nil}, cnf}, assets, nil
}
