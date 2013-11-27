package neyo

import (
	"errors"
	"fmt"
	"github.com/wendal/mustache"
)

const (
	COMMENTS_UYAN = ` <!-- UY BEGIN -->
    <div id="uyan_frame"></div>
    <script type="text/javascript" src="http://v2.uyan.cc/code/uyan.js?uid=%d"></script>
    <!-- UY END -->
    `
	COMMENTS_DUOSHUO = `
	<!-- Duoshuo Comment BEGIN -->
	<div class="ds-thread"></div>
	<script type="text/javascript">
	var duoshuoQuery = {short_name:"%s"};//require,replace your short_name
	(function() {
					var ds = document.createElement('script');
					ds.type = 'text/javascript';ds.async = true;
					ds.src = 'http://static.duoshuo.com/embed.js';
					ds.charset = 'UTF-8';
					(document.getElementsByTagName('head')[0]
					|| document.getElementsByTagName('body')[0]).appendChild(ds);
	})();
	</script>
	<!-- Duoshuo Comment END -->
	`
	COMMENTS_DISQUS = `
<div id="disqus_thread"></div>
<script>
    var disqus_developer = 1;
    var disqus_shortname = '%s'; // required: replace example with your forum shortname
    /* * * DON'T EDIT BELOW THIS LINE * * */
    (function() {
        var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
        dsq.src = 'http://' + disqus_shortname + '.disqus.com/embed.js';
        (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
    })();
</script>
<noscript>Please enable JavaScript to view the <a href="http://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
<a href="http://disqus.com" class="dsq-brlink">blog comments powered by <span class="logo-disqus">Disqus</span></a>
`
)

type CommentsWidget Mapper

func (self CommentsWidget) Prepare(mapper Mapper, topCtx mustache.Context) Mapper {
	if mapper["comments"] != nil && !mapper["comments"].(bool) {
		Log(INFO, "Disable comments")
		return nil
	}
	return Mapper(self)
}

func BuildCommentsWidget(cnf Mapper, topCtx mustache.Context) (Widget, error) {
	Log(DEBUG, "Build comment widget %s", cnf.Layout())
	switch cnf.Layout() {
	case "disqus":
		disqus := cnf[cnf.Layout()].(map[string]interface{})
		short_name := disqus["short_name"]
		if short_name == nil {
			return nil, errors.New("CommentsWidget Of disqus need short_name")
		}
		self := make(CommentsWidget)
		self["comments"] = fmt.Sprintf(COMMENTS_DISQUS, short_name)
		return self, nil
	case "uyan":
		uyan := cnf[cnf.Layout()].(map[string]interface{})
		uid := uyan["uid"]
		self := make(CommentsWidget)
		self["comments"] = fmt.Sprintf(COMMENTS_UYAN, uid)
		return self, nil
	case "duoshuo":
		duoshuo := cnf[cnf.Layout()].(map[string]interface{})
		short_name := duoshuo["short_name"]
		if short_name == nil {
			return nil, errors.New("CommentsWidget Of duoshuo need short_name")
		}
		self := make(CommentsWidget)
		self["comments"] = fmt.Sprintf(COMMENTS_DUOSHUO, short_name)
		return self, nil
	}
	return nil, errors.New("CommentsWidget Only for disqus yet")
}

func init() {
	WidgetBuilders["comments"] = BuildCommentsWidget //社会化评论
}
