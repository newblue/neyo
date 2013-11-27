package neyo

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

func MakePayLoad(root string) (site *WebSite, err error) {
	site = &WebSite{}
	err = nil

	if root == "" {
		root = "."
	}
	root, err = filepath.Abs(root)
	if err != nil {
		Log(ERROR, "%s", err)
	}
	site.Root = root
	Log(DEBUG, "Root: %s", site.Root)

	site.LoadMainConfig()
	site.CheckMainConfig()
	site.MakeBasicURLs()

	site.FixPostPageConfigs()

	site.LoadPages()
	site.LoadPosts()

	return
}

func (self *WebSite) LoadMainConfig() {

	cnf, err := ReadYml(filepath.Join(self.Root, CONFIG_YAML))
	if err != nil {
		Log(ERROR, "Read %s error %s", CONFIG_YAML, err)
	}
	ToStruct(cnf, reflect.ValueOf(&self.TopCnf))
	Log(DEBUG, "config.yml\n%s\n", self.TopCnf)

	siteCnf, err := ReadYml(self.Root + SITE_YAML)
	if err != nil {
		Log(ERROR, "Read %s error %s", SITE_YAML, err)
	}
	ToStruct(siteCnf, reflect.ValueOf(&self.SiteCnf))
	Log(DEBUG, "site.yml\n%s\n", self.SiteCnf)
}

func (self *WebSite) CheckMainConfig() {

	themeName := self.TopCnf.Theme
	if themeName == "" {
		Log(ERROR, "Miss theme config!")
	}
	// 载入theme的设置
	themeCnf, err := ReadYml(fmt.Sprintf("%s/themes/%s/theme.yml", self.Root, themeName))
	if err != nil {
		Log(ERROR, "No %s theme? %s ", themeName, err.Error())
	}
	ToStruct(themeCnf, reflect.ValueOf(&self.ThemeCnf))

	self.Layouts = LoadLayouts(self.Root, themeName)
	if self.Layouts == nil || len(self.Layouts) == 0 {
		Log(ERROR, "Theme without any layout!!")
	}

	production_url := self.TopCnf.Production_url
	if production_url == "" {
		Log(ERROR, "Miss production_url")
	}
	if !strings.HasPrefix(production_url, "http://") &&
		!strings.HasPrefix(production_url, "https://") {
		Log(ERROR, "production_url must start with https:// or http://")
	}

	rootUrl := production_url
	pos := strings.Index(rootUrl[len("https://"):], "/")
	basePath := ""
	if pos == -1 {
		basePath = "/"
	} else {
		basePath = rootUrl[len("https://")+pos:]
		if !strings.HasSuffix(basePath, "/") {
			basePath += "/"
		}
	}
	self.RootURL = rootUrl
	self.BasePath = basePath
}

func (self *WebSite) MakeBasicURLs() {
	urls := make(map[string]string)
	urls["media"] = filepath.Join(self.BasePath, "assets/media")
	urls["theme"] = filepath.Join(self.BasePath, "assets/", self.TopCnf.Theme)
	urls["theme_media"] = filepath.Join(urls["theme"], "/media")
	urls["theme_javascripts"] = filepath.Join(urls["theme"], "/javascripts")
	urls["theme_stylesheets"] = filepath.Join(urls["theme"], "/stylesheets")
	urls["base_path"] = self.BasePath

	//需要重写?
	/*
		if site["urls"] != nil { //允许用户自定义基础URL,实现CDN等功能
			var site_url Mapper
			site_url = site["urls"].(map[string]interface{})
			for k, v := range site_url {
				urls[k] = v.(string)
			}
		}
	*/
	self.BaiseURLs = urls
}

func (self *WebSite) FixPostPageConfigs() {

	postsCnf := self.TopCnf.Posts
	if postsCnf.Permalink == "" {
		postsCnf.Permalink = "/:categories/:title/"
	}
	if postsCnf.Summary_lines < 5 {
		postsCnf.Summary_lines = 20
	}
	if postsCnf.Latest < 5 {
		postsCnf.Latest = 5
	}
	if postsCnf.Layout == "" {
		postsCnf.Layout = "post"
	}

	pagesCnf := self.TopCnf.Pages
	if pagesCnf.Layout == "" {
		pagesCnf.Layout = "page"
	}
	if pagesCnf.Permalink == "" {
		pagesCnf.Permalink = "pretty"
	}
}

func (self *WebSite) LoadPages() {
	pagesCnf := self.TopCnf.Pages
	pages, err := LoadPages(self.Root, pagesCnf.Exclude)
	if err != nil {
		return
	}
	// 构建导航信息(page列表),及整理page的配置信息
	navigation := make([]string, 0)
	self.Pages = make(map[string]PageBean)
	for page_id, page := range pages {
		pageBean := PageBean{}
		ToStruct(page, reflect.ValueOf(&pageBean))
		self.Pages[page_id] = pageBean

		if pageBean.Layout == "" {
			pageBean.Layout = pagesCnf.Layout
		}
		if pageBean.Permalink == "" {
			pageBean.Permalink = pagesCnf.Permalink
		}

		page_url := ""
		switch {
		case strings.HasSuffix(page_id, "index.html"):
			page_url = page_id[0 : len(page_id)-len("index.html")]
		case strings.HasSuffix(page_id, "index.md"):
			page_url = page_id[0 : len(page_id)-len("index.md")]
		default:
			page_url = page_id[0 : len(page_id)-len(filepath.Ext(page_id))]
			if pageBean.Title == "" && !strings.HasSuffix(page_url, "/") {
				pageBean.Title = strings.Title(filepath.Base(page_url))
			}
		}
		if strings.HasPrefix(page_url, "/") {
			pageBean.Url = filepath.Join(self.BasePath, page_url[1:])
		} else {
			pageBean.Url = filepath.Join(self.BasePath, page_url)
		}

		if page_id != "index.html" && page_id != "index.md" {
			navigation = append(navigation, page_id)
		}
	}
	if self.SiteCnf.Navigation == nil || len(self.SiteCnf.Navigation) == 0 {
		self.SiteCnf.Navigation = navigation
	}
}

func (self *WebSite) LoadPosts() {
	postsCnf := self.TopCnf.Posts
	posts, err := LoadPosts(self.Root, postsCnf.Exclude)
	if err != nil {
		return
	}
	self.Posts = make(map[string]PostBean)
	for post_id, _post := range posts {
		postBean := PostBean{}
		ToStruct(_post, reflect.ValueOf(&postBean))
		self.Posts[post_id] = postBean

		if postBean.Layout == "" {
			postBean.Layout = postsCnf.Layout
		}
		if postBean.Permalink == "" {
			postBean.Permalink = postsCnf.Permalink
		}

		if postBean.Tags == nil {
			postBean.Tags = []string{}
		}
		if postBean.Categories == nil {
			postBean.Categories = []string{}
		}
	}

	// 整理post
	tags := make(map[string]*Tag)
	catalogs := make(map[string]*Catalog)
	chronological := make([]string, 0)
	collated := make(CollatedYears, 0)

	_collated := make(map[string]*CollatedYear)

	for id, post := range self.Posts {
		chronological = append(chronological, id)

		for _, _tag := range post.Tags {
			tag := tags[_tag]
			if tag == nil {
				tag = &Tag{0, _tag, make([]string, 0), "/tags/#" + EncodePathInfo(_tag) + "-ref"}
				tags[_tag] = tag
			}
			tag.Count += 1
			tag.Posts = append(tag.Posts, id)
		}

		for _, _catalog := range post.Categories {
			catalog := catalogs[_catalog]
			if catalog == nil {
				catalog = &Catalog{0, _catalog, make([]string, 0), "/categories/#" + EncodePathInfo(_catalog) + "-ref"}
				catalogs[_catalog] = catalog
			}
			catalog.Count += 1
			catalog.Posts = append(catalog.Posts, id)
		}

		_year, _month, _ := post._Date.Date()
		year := fmt.Sprintf("%v", _year)
		month := _month.String()

		_yearc := _collated[year]
		if _yearc == nil {
			_yearc = &CollatedYear{year, make([]*CollatedMonth, 0), make(map[string]*CollatedMonth)}
			_collated[year] = _yearc
		}
		_monthc := _yearc.months[month]
		if _monthc == nil {
			_monthc = &CollatedMonth{month, _month, []string{}}
			_yearc.months[month] = _monthc
			//log.Println("Add>>", year, month, post["id"])
		}
		_monthc.Posts = append(_monthc.Posts, id)

		post_map := make(Mapper)
		post_map["title"] = post.Title
		post_map["_date"] = post._Date
		post_map["id"] = post.Id
		post_map["categories"] = post.Categories
		post_map["permalink"] = post.Permalink
		CreatePostURL(nil, self.BasePath, post_map)
		post.Url = post_map["url"].(string)
	}
	_ = collated
	// TODO 需要重写排序方法
	/*

		sort.Sort(collated)

		for _, catalog := range catalogs {
			catalog.Posts = SortPosts(webSite.Posts, catalog.Posts)
		}
		for _, tag := range tags {
			tag.Posts = SortPosts(webSite.Posts, tag.Posts)
		}

		webSite.Tags = tags
		webSite.Catalogs = catalogs
		webSite.Chronological = SortPosts(webSite.Posts, chronological)
		webSite.Collated = collated
	*/
}
