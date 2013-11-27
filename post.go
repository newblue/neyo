package neyo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TPL_NEW_POST = `---
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

// 创建一个新post
// TODO 移到到其他地方?
func CreateNewPost(title string) (path string) {
	if !IsGorDir(".") {
		Log(ERROR, "Not Gor Dir, need config.yml")
	}
	path = filepath.Join("posts/", strings.Replace(title, " ", "-", -1)+".md")
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		Log(ERROR, "Post file(%s) exist?", path)
	}
	err = ioutil.WriteFile(path, []byte(fmt.Sprintf(TPL_NEW_POST, title, time.Now().Format("2006-01-02"))), DEFAULT_FILE_MODE)
	if err != nil {
		Log(ERROR, "%s", err)
	}
	Log(INFO, "Create Post at %s", path)
	return
}

func CreateNewPostWithImgs(title, imgsrc string) (path string) {

	cfg := loadConfig(".")
	for k, v := range cfg {
		Log(INFO, "%s = %s", k, v)
	}
	path = CreateNewPost(title)

	start := strings.LastIndex(path, "/") + 1
	end := strings.LastIndex(path, ".")
	if start < 0 || end < 0 {
		Log(ERROR, "path not complate? %s", path)
	}
	post := path[start:end]

	// 如果创建失败直接exit，所以不用检查
	imgs := cpPostImgs(post, imgsrc, cfg)
	tags := generateImgLinks(imgs, cfg)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, DEFAULT_FILE_MODE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, tag := range tags {
		if _, err = f.WriteString("\n" + tag + "\n"); err != nil {
			panic(err)
		}
	}
	return
}

func cpPostImgs(post string, imgsrc string, cfg Mapper) (imgtag []string) {
	files, err := ioutil.ReadDir(imgsrc)
	if files == nil || err != nil {
		Log(INFO, "no img file exists.")
		return nil
	}

	if !strings.HasSuffix(imgsrc, "/") {
		imgsrc += "/"
	}

	imgdst := filepath.Join(cfg.GetString("localdir"), post)
	_, err = os.Stat(imgdst)
	if os.IsNotExist(err) {
		os.MkdirAll(imgdst, DEFAULT_DIR_MODE)
	}

	imgtag = make([]string, len(files))
	i := 0
	for idx, f := range files {
		err := cp(filepath.Join(imgdst, f.Name()), filepath.Join(imgsrc, f.Name()))
		if err != nil {
			Log(ERROR, "%5d resouce file copy %s error", idx, f.Name())
			continue
		}
		imgtag[i] = filepath.Join(post, f.Name())
		i++
	}

	imgtag = imgtag[:i]
	return
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func generateImgLinks(files []string, cfg Mapper) (links []string) {
	links = make([]string, len(files))
	for i, f := range files {
		//tmp := strings.TrimLeft(f, "rc/")
		links[i] = fmt.Sprintf(cfg.GetString("imgtag"), cfg.GetString("urlperfix")+f)
		println(i, links[i])
	}

	return
}

func loadConfig(root string) (imgs_cfg Mapper) {
	var cfg Mapper
	var err error

	if root == "" {
		root = "."
	}
	root, err = filepath.Abs(root)
	root += "/"
	Log(INFO, "root = %s", root)

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
