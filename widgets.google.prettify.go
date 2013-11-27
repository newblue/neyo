package neyo

import (
	"fmt"
	"github.com/wendal/mustache"
)

const (
	GOOGLE_PRETTIFY = `
    <script src="http://cdnjs.cloudflare.com/ajax/libs/prettify/188.0.0/prettify.js"></script>
    <script>
    var pres = document.getElementsByTagName("pre");
    for (var i=0; i < pres.length; ++i) {
        pres[i].className = "prettyprint %s";
    }
    prettyPrint();
    </script>
    `
)

type google_prettify Mapper

func (self google_prettify) Prepare(mapper Mapper, topCtx mustache.Context) Mapper {
	if mapper["google_prettify"] != nil && !mapper["google_prettify"].(bool) {
		return nil
	}
	return Mapper(self)
}

func BuildGoogle_prettify(cnf Mapper, topCtx mustache.Context) (Widget, error) {
	if enable, ok := cnf["linenums"].(bool); ok && enable { //是否显示行号
		self := make(google_prettify)
		self["google_prettify"] = fmt.Sprintf(GOOGLE_PRETTIFY, "linenums")
		return self, nil
	}
	self := make(google_prettify)
	self["google_prettify"] = fmt.Sprintf(GOOGLE_PRETTIFY, "")
	return self, nil
}

func init() {
	WidgetBuilders["google_prettify"] = BuildGoogle_prettify // 代码高亮
}
