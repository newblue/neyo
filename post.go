package neyo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	NEW_POST = `---
title: %s
date: '%s'
description:
categories:

tags:

---

`
	IMG_TAG       = `<img src="%s" alt="img: " width="600">`
	IMG_URLPERFIX = `{{urls.media}}/`
	IMG_LOCALDIR  = `media/`
)

func CreateNewPost(title string) (path string) {
	if wd, err := os.Getwd(); err != nil || !IsYoProjectDir(wd) {
		Log(ERROR, "Not yo project diretory, need config.yml")
	}

	path = filepath.Join("posts/", strings.Replace(title, " ", "-", -1)+".md")

	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		Log(ERROR, "Post file(%s) exist?", path)
	}
	post_default := fmt.Sprintf(NEW_POST, title, time.Now().Format("2006-01-02"))

	if err := ioutil.WriteFile(path, []byte(post_default), DEFAULT_FILE_MODE); err != nil {
		Log(ERROR, "%s", err)
	}
	Log(INFO, "Create Post at %s", path)
	return
}

func CreateNewPostWithImgs(title, imgsrc string) (path string) {
	cfg := loadConfig(".")
	for k, v := range cfg {
		Log(DEBUG, "CFG %s = %s", k, v)
	}
	path = CreateNewPost(title)

	start := strings.LastIndex(path, "/") + 1
	end := strings.LastIndex(path, ".")

	if start < 0 || end < 0 {
		Log(ERROR, "%s path not complate?", path)
	}
	post := path[start:end]

	imgs := copyPostImgs(post, imgsrc, cfg)
	tags := generateImgLinks(imgs, cfg)

	if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, DEFAULT_FILE_MODE); err != nil {
		Log(ERROR, "Open %s >> %s", path, err)
	} else {
		defer file.Close()

		for _, tag := range tags {
			if _, err = file.WriteString("\n" + tag + "\n"); err != nil {
				panic(err)
			}
		}
	}
	return
}

func copyPostImgs(post string, imgsrc string, cfg Mapper) (imgtag []string) {
	files, err := ioutil.ReadDir(imgsrc)
	if files == nil || err != nil {
		Log(INFO, "No img file exists.")
		return nil
	}

	if !strings.HasSuffix(imgsrc, "/") {
		imgsrc += "/"
	}

	imgdst := filepath.Join(cfg.GetString("localdir"), post)

	if _, err := os.Stat(imgdst); os.IsNotExist(err) {
		os.MkdirAll(imgdst, DEFAULT_DIR_MODE)
	}

	imgtag = make([]string, len(files))
	i := 0
	for idx, f := range files {
		dst_path := filepath.Join(imgdst, f.Name())
		src_path := filepath.Join(imgsrc, f.Name())
		if err := Copy(dst_path, src_path); err != nil {
			Log(ERROR, "%5d resouce file copy %s error", idx, f.Name())
			continue
		}
		imgtag[i] = filepath.Join(post, f.Name())
		i++
	}

	imgtag = imgtag[:i]
	return
}

func generateImgLinks(files []string, cfg Mapper) (links []string) {
	links = make([]string, len(files))
	for i, f := range files {
		links[i] = fmt.Sprintf(cfg.GetString("imgtag"), cfg.GetString("urlperfix")+f)
		println(i, links[i])
	}

	return
}

func loadConfig(root string) (imgs_cfg Mapper) {
	var cfg Mapper
	var err error

	if wd, err := os.Getwd(); err == nil && root == "" {
		root = wd
	} else if abs, err := filepath.Abs("."); err == nil && root == "" {
		root = abs
	}

	Log(DEBUG, "ROOT %s", root)

	cfg, err = ReadYml(root + CONFIG_YAML)
	if err != nil {
		Log(ERROR, "Fail to read %s %s", root+CONFIG_YAML, err)
		return
	}

	if cfg["imgs"] == nil {
		imgs_cfg = make(Mapper)
	} else {
		imgs_cfg = cfg["imgs"].(map[string]interface{})
	}

	if imgs_cfg["imgtag"] == nil {
		imgs_cfg["imgtag"] = IMG_TAG
	}
	if imgs_cfg["urlperfix"] == nil {
		imgs_cfg["urlperfix"] = IMG_URLPERFIX
	}
	if imgs_cfg["localdir"] == nil {
		imgs_cfg["localdir"] = IMG_LOCALDIR
	}
	return
}
