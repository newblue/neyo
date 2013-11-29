package neyo

import (
	//	"archive/zip"
	//	"bytes"
	//	"encoding/base64"
	//	"io"
	//	"io/ioutil"
	"os"
	"path/filepath"
)

func New(path string, init_data map[string][]byte) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		Log(ERROR, "Path Exist?!")
	}

	err = os.MkdirAll(path, DEFAULT_DIR_MODE)
	if err != nil {
		Log(ERROR, "Make diretory error %s", err)
	}

	//decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(zipgo))
	//b, _ := ioutil.ReadAll(decoder)

	//z, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	//if err != nil {
	//	Log(ERROR, "Read zip data error %s", err)
	//}

	//Log(INFO, "Unpack init content zip")
	Log(INFO, "Init project data.")

	for name, data := range init_data {
		dst_path := filepath.Join(path, name)
		os.MkdirAll(filepath.Dir(dst_path), DEFAULT_DIR_MODE)
		file, err := os.OpenFile(dst_path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DEFAULT_FILE_MODE)
		if err != nil {
			Log(ERROR, "Open %s >> %s", dst_path, err)
		}

		if _, err := file.Write(data); err != nil {
			Log(ERROR, "Write data to %s >> %s.", dst_path, err)
		} else {
			Log(INFO, name)
		}

		file.Sync()
		file.Close()
	}
	/*
		for _, zf := range z.File {
			if zf.FileInfo().IsDir() {
				continue
			}
			dst := filepath.Join(path, zf.FileInfo().Name())
			os.MkdirAll(filepath.Dir(dst), DEFAULT_DIR_MODE)
			f, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DEFAULT_FILE_MODE)
			if err != nil {
				Log(ERROR, "Open %s error %s", dst, err)
			}

			defer f.Sync()
			defer f.Close()

			rc, err := zf.Open()
			if err != nil {
				Log(ERROR, "Open %s error %s", zf.FileInfo().Name(), err)
			}
			defer rc.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				Log(ERROR, "Copy %s to %s error %s.", zf.FileInfo().Name(), dst, err)
			}
		}
	*/
	Log(INFO, "All done.")
}
