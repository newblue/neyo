package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"github.com/newblue/gor"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func new_init(path string) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		gor.Log(gor.ERROR, "Path Exist?!")
	}

	err = os.MkdirAll(path, 0700)
	if err != nil {
		gor.Log(gor.ERROR, "Make diretory error %s", err)
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(INIT_ZIP))
	b, _ := ioutil.ReadAll(decoder)

	z, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		gor.Log(gor.ERROR, "Read zip data error %s", err)
	}

	gor.Log(gor.INFO, "Unpack init content zip")

	for _, zf := range z.File {
		if zf.FileInfo().IsDir() {
			continue
		}
		dst := filepath.Join(path, zf.FileInfo().Name())
		os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		f, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
		if err != nil {
			gor.Log(gor.ERROR, "Open %s error %s", dst, err)
		}

		defer f.Sync()
		defer f.Close()

		rc, err := zf.Open()
		if err != nil {
			gor.Log(gor.ERROR, "Open %s error %s", zf.FileInfo().Name(), err)
		}

		defer rc.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			gor.Log(gor.ERROR, "Copy %s to %s error %s.", zf.FileInfo().Name(), dst, err)
		}
	}
	gor.Log(gor.ERROR, "All done.")
}
