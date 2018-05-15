package httpserver

import (
	"net"
	"net/http"
	"runtime"
	"strings"

	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/clzhan/srlive3/log"
	"github.com/clzhan/srlive3/utils"
	"github.com/clzhan/srlive3/stream"
)

var ostype = runtime.GOOS

var crossdomainxml = []byte(`<?xml version="1.0" ?>
<cross-domain-policy>
	<allow-access-from domain="*" />
	<allow-http-request-headers-from domain="*" headers="*"/>
</cross-domain-policy>`)

type HttpServer struct {
	listener net.Listener
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

func (server *HttpServer) Serve(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.handleConn(w, r)
	})
	http.Serve(l, mux)
	return nil
}
func (server *HttpServer) GetListener() net.Listener {
	return server.listener
}

func (server *HttpServer) handleConn(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("http flv handleConn panic: ", r)
		}
	}()

	//跨域
	if path.Base(r.URL.Path) == "crossdomain.xml" {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(crossdomainxml)
		return
	}

	switch path.Ext(r.URL.Path) {
	case ".m3u8":
		url := r.URL.String()
		u := r.URL.Path
		path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".m3u8")
		paths := strings.SplitN(path, "/", 2)
		log.Info("url:", url, "path:", path, "paths:", paths)

		_, found := stream.FindObject(path)
		if !found {
			log.Error("object not find key :", path)
			http.NotFound(w, r)
			return
		}
		var m3u8 string


		if ostype == "windows"{
			m3u8 = util.GetProjectPath() + "\\" + "hlsstream"+"\\"+ paths[0] + "\\" + paths[1] + ".m3u8"
		}else{
			m3u8 = util.GetProjectPath() + "/" + "hlsstream"+"/"+ paths[0] + "/" + paths[1] + ".m3u8"
		}

		if data, err := loadFile(m3u8); nil == err {
			log.Info("Send m3u8 path：",m3u8);

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Content-Type", "application/x-mpegURL")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			if _, err = w.Write(data); err != nil {
				log.Error("write m3u8 file err=", err)
			}
		}else{
			log.Info("m3u8  path.......not found...");
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	case ".ts":
		app, ts := parseTsFile(r.URL.Path)
		log.Info("app.......ts.....", app, ts)
		if ostype == "windows"{
			ts = util.GetProjectPath() + "\\" + "hlsstream"+"\\"+ app + "\\" + ts
		}else{
			ts = util.GetProjectPath() + "/" + "hlsstream"+"/"+ app + "/" + ts
		}
		if data, err := loadFile(ts); nil == err {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "video/mp2ts")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			if _, err = w.Write(data); err != nil {
				log.Error("write ts file err=", err)
			}
		}else{
			log.Info("app.......ts..... not found")
		}
	case ".flv":
		url := r.URL.String()
		u := r.URL.Path
		path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
		paths := strings.SplitN(path, "/", 2)
		log.Info("url:", url, "path:", path, "paths:", paths)
		var obj *stream.StreamObject
		obj, found := stream.FindObject(path)
		if !found {
			log.Error("object not find key :", path)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		stream := NewHttpFlvStream(path)
		stream.SetObj(obj)
		obj.HttpAttach(stream)
		log.Debug("HttpFlvStream BeginHanle")

		stream.WriteLoop(w, r)
		log.Debug("HttpFlvStream mid")

		stream.Close()
		log.Debug("HttpFlvStream EndHanle")

	default:
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

}
func parseTsFile(p string) (app string, tsFile string) {
	if i := strings.Index(p, "/"); i >= 0 {
		if j := strings.LastIndex(p, "/"); j > 0 {
			app = p[i+1 : j]
		}
	}

	if i := strings.LastIndex(p, "/"); i > 0 {
		tsFile = p[i+1:]
	}

	return
}


func loadFile(filename string) (data []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(util.PanicTrace())
		}
	}()

	var f *os.File
	if f, err = os.Open(filename); err != nil {
		log.Error("Open file ", filename, " failed, err is", err)
		return
	}
	defer f.Close()

	if data, err = ioutil.ReadAll(f); err != nil {
		log.Error("read file ", filename, " failed, err is", err)
		return
	}

	return
}
