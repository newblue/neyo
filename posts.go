package neyo

import (
	"fmt"
)

// 列出全部post -- 纯属无聊?

func ListPosts() {
	var payload Mapper
	payload, err := BuildPayload("./")
	if err != nil {
		Log(ERROR, "%s", err)
	}
	posts := payload["db"].(map[string]interface{})["posts"].(map[string]interface{})["chronological"].([]string)
	fmt.Printf("Posts total(%d):\n", len(posts))
	for _, id := range posts {
		fmt.Printf("   - %s\n", id)
	}
}
