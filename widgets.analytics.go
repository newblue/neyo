package neyo

import (
	"errors"
	"fmt"
	"github.com/wendal/mustache"
)

const (
	GOOGLE_ANALYTICS = `
<script type="text/javascript">

  var _gaq = _gaq || [];
  var pluginUrl = '//www.google-analytics.com/plugins/ga/inpage_linkid.js';
  _gaq.push(['_require', 'inpage_linkid', pluginUrl]);
  _gaq.push(['_setAccount', '%s']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();

</script>`

	CNZZ_ANALYTICS = `<script src="http://s25.cnzz.com/stat.php?id=%d&web_id=%d" language="JavaScript"></script>`
)

func init() {
	WidgetBuilders["analytics"] = BuildAnalyticsWidget //访问统计
}

type AnalyticsWidget Mapper

func (self AnalyticsWidget) Prepare(mapper Mapper, topCtx mustache.Context) Mapper {
	if mapper["analytics"] != nil && !mapper["analytics"].(bool) {
		return nil
	}
	return Mapper(self)
}

func BuildAnalyticsWidget(cnf Mapper, topCtx mustache.Context) (Widget, error) {
	switch cnf.Layout() {
	case "google": // 鼎鼎大名的免费,但有点拖慢加载速度,原因你懂的
		google := cnf[cnf.Layout()].(map[string]interface{})
		tracking_id := google["tracking_id"]
		if tracking_id == nil {
			return nil, errors.New("AnalyticsWidget Of Google need tracking_id")
		}
		self := make(AnalyticsWidget)
		self["analytics"] = fmt.Sprintf(GOOGLE_ANALYTICS, tracking_id)
		return self, nil
	case "cnzz": //免费,而且很快,但强制嵌入一个反向链接,靠!
		cnzz := cnf[cnf.Layout()].(map[string]interface{})
		tracking_id := cnzz["tracking_id"]
		if tracking_id == nil {
			return nil, errors.New("AnalyticsWidget Of CNZZ need tracking_id")
		}
		self := make(AnalyticsWidget)
		self["analytics"] = fmt.Sprintf(CNZZ_ANALYTICS, tracking_id, tracking_id)
		return self, nil
	}
	return nil, errors.New("AnalyticsWidget Only for Goolge/CNZZ yet")
}
