package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/newblue/neyo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/pprof"
	"strings"
)

const (
	NAME = "yo"
	VER  = "20131126"
)

var (
	config_command     = flag.NewFlagSet("config", flag.ExitOnError)
	new_command        = flag.NewFlagSet("new", flag.ExitOnError)
	posts_command      = flag.NewFlagSet("posts", flag.ExitOnError)
	payload_command    = flag.NewFlagSet("payload", flag.ExitOnError)
	compile_command    = flag.NewFlagSet("compile", flag.ExitOnError)
	post_command       = flag.NewFlagSet("post", flag.ExitOnError)
	http_command       = flag.NewFlagSet("http", flag.ExitOnError)
	pprof_command      = flag.NewFlagSet("ppprof", flag.ExitOnError)
	update_zip_command = flag.NewFlagSet("update_zip", flag.ExitOnError)
)

func init() {
	flag.Parse()
	fmt.Printf("%s BETA(%s). a tool fork from gor\n\n", NAME, VER)
}

func main() {
	args := flag.Args()
	if len(args) == 0 || len(args) > 3 {
		PrintUsage()
	}
	switch args[0] {
	default:
		PrintUsage()
	case "config":
		_config()
	case "new":
		_new(args)
	case "posts":
		_posts()
	case "payload":
		_payload()
	case "compile":
		_compile(args)
	case "post":
		_post(args)
	case "http":
		_http(args)
	case "pprof":
		_pprof()
	case "zip.go":
		_update_zip(args)
	}
}

func _posts() {
	neyo.ListPosts()
}

func _post(args []string) {
	if len(args) == 1 {
		neyo.Log(neyo.INFO, "%s post <title> {image diretory}", NAME)
	} else if len(args) == 2 {
		path := neyo.CreateNewPost(args[1])
		edit_new_post(path)
	} else {
		path := neyo.CreateNewPostWithImgs(args[1], args[2])
		edit_new_post(path)
	}
}

func get_editor() (editor string) {
	editor = os.Getenv("EDITOR")
	cnf, err := neyo.ReadConfig(".")
	if err != nil {
		neyo.Log(neyo.ERROR, "Read config error %s", err)
	} else if ed, ok := cnf["editor"].(string); ok {
		editor = ed
	}
	return
}

func edit_new_post(path string) {
	if editor := get_editor(); editor != "" {
		fmt.Printf("Are you edit page? (Yes/No)")
		var ask string
		if _, err := fmt.Scan(&ask); err == nil {
			if _ask := strings.ToLower(ask); _ask == "y" || _ask == "yes" {
				fmt.Printf("Edit %s\n", editor, path)
				cmd := exec.Command(editor, path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err == nil {
					if err := cmd.Wait(); err != nil {
						neyo.Log(neyo.ERROR, "Wait %s", err)
					}
				} else {
					neyo.Log(neyo.ERROR, "Start %s", err)
				}
			}
		}
	}
}

func _http(args []string) {
	address := "127.0.0.1"
	public := neyo.DEFAULT_PUBLIC_DIRETORY

	if la := len(args); la < 2 {
		fmt.Printf("%s http <address:port> [public diretory]\n", NAME)
	} else if la == 2 {
		address = args[1]
	} else if la == 3 {
		address, public = args[1], args[2]
	}
	neyo.Log(neyo.INFO, "Listen at %s", address)
	http.ListenAndServe(address, http.FileServer(http.Dir(public)))
}

func _update_zip(args []string) {
	if len(args) == 2 {
		dir := args[1]

		tmp_file, err := ioutil.TempFile("", "temp-")
		if err != nil {
			neyo.Log(neyo.ERROR, "Open temp file error %s", err)
		} else {
			neyo.Log(neyo.DEBUG, "Create temp file %s", tmp_file.Name())
		}
		defer os.Remove(tmp_file.Name())
		defer EncodeIntoGo(tmp_file.Name(), "zip.go", "INIT_ZIP")
		defer tmp_file.Close()

		z := zip.NewWriter(tmp_file)
		defer func() {
			if err := z.Close(); err != nil {
				neyo.Log(neyo.ERROR, "Close zip file %s", err)
			} else {
				neyo.Log(neyo.INFO, "zip.go updated.\n")
			}
		}()

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			basename := filepath.Base(path)
			is_ignore := strings.HasPrefix(basename, ".")

			if info.IsDir() {
				return nil
			} else if is_ignore {
				neyo.Log(neyo.WARN, "Ignore archive %s", path)
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			zip_path := strings.TrimLeft(path[len(dir):len(path)], "/")
			neyo.Log(neyo.DEBUG, "%s\n\t->zip://%s\n", path, zip_path)
			if sf, err := os.Open(path); err != nil {
				neyo.Log(neyo.ERROR, "Open %s error %s", path, err)
			} else if df, err := z.Create(zip_path); err != nil {
				neyo.Log(neyo.ERROR, "Open zip://%s error %s", zip_path, err)
			} else {
				io.Copy(df, sf)
				if err := sf.Close(); err != nil {
					neyo.Log(neyo.ERROR, "Close %s error %s", path, err)
				}
			}
			return nil
		})
	} else {
		fmt.Printf("\t %s zip.go  <diretory>      Archive project directory, and make zip.go\n", NAME)
	}
}
func EncodeIntoGo(filename, gofilename string, varname string) error {
	d, _ := ioutil.ReadFile(filename)
	_zip, _ := os.OpenFile(gofilename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)

	header := fmt.Sprintf(`package main
    const %s="`, varname)

	_zip.Write([]byte(header))
	encoder := base64.NewEncoder(base64.StdEncoding, _zip)
	encoder.Write(d)
	encoder.Close()
	_zip.Write([]byte("\"\n"))
	_zip.Sync()
	return _zip.Close()
}

func _pprof() {
	f, _ := os.OpenFile("neyo.pprof", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 100; i++ {
		err := neyo.Compile(neyo.DEFAULT_PUBLIC_DIRETORY)
		if err != nil {
			neyo.Log(neyo.ERROR, "%s", err)
		}
	}
}

func _config() {
	cnf, err := neyo.ReadConfig(".")
	if err != nil {
		neyo.Log(neyo.ERROR, "Read config error %s", err)
	}
	neyo.Log(neyo.INFO, "RuhohSpec: %s", cnf["RuhohSpec"])
	buf, err := json.MarshalIndent(cnf, "", "  ")
	if err != nil {
		neyo.Log(neyo.ERROR, "Marshal error %s", err)
	}
	fmt.Printf("Global config\n %s", string(buf))
}

func _new(args []string) {
	if len(args) == 1 {
		fmt.Printf("\t%s new <diertory>\n", NAME)
	} else {
		new_init(args[1])
	}
}

func _payload() {
	payload, err := neyo.BuildPlayload("./")
	if err != nil {
		neyo.Log(neyo.ERROR, "Build paly load %s", err)
	}
	buf, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		neyo.Log(neyo.ERROR, "%s", err)
	}
	neyo.Log(neyo.INFO, string(buf))
}

func _compile(args []string) {
	public := neyo.DEFAULT_PUBLIC_DIRETORY
	if len(args) == 2 {
		public = args[1]
	}
	err := neyo.Compile(public)
	if err != nil {
		neyo.Log(neyo.ERROR, "%s", err)
	}
}
