package neyo

// 列出全部post -- 纯属无聊?
func ListPosts() {
	var payload Mapper
	payload, err := BuildPlayload("./")
	if err != nil {
		Log(ERROR, "%s", err)
	}
	posts := payload["db"].(map[string]interface{})["posts"].(map[string]interface{})["chronological"].([]string)
	Log(INFO, "Posts count: %d", len(posts))
	for _, id := range posts {
		Log(INFO, "\t- %s", id)
	}
}
