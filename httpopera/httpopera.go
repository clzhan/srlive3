package httpopera

import (
	"net"
	"net/http"
	"encoding/json"
	"fmt"
)

var crossdomainxml = []byte(`<?xml version="1.0" ?>
<cross-domain-policy>
	<allow-access-from domain="*" />
	<allow-http-request-headers-from domain="*" headers="*"/>
</cross-domain-policy>`)

type HttpOpera struct {
	listener net.Listener
}

func NewHttpHttpOpera() *HttpOpera {
	return &HttpOpera{}
}

func (s *HttpOpera) Serve(l net.Listener) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/versions", func(w http.ResponseWriter, r *http.Request) {
		s.GetVersions(w, r)
	})
	mux.HandleFunc("/api/v1/summaries", func(w http.ResponseWriter, r *http.Request) {
		s.GetSummaries(w, r)
	})
	mux.HandleFunc("/api/v1/rusages", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/self_proc_stats", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/system_proc_stats", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/meminfos", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/authors", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/features", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/requests", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/vhosts", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("/api/v1/streams", func(w http.ResponseWriter, r *http.Request) {
		s.GetStreams(w, r)
	})
	mux.HandleFunc("/api/v1/streams/", func(w http.ResponseWriter, r *http.Request) {
		s.GetStreams(w, r)
	})
	mux.HandleFunc("/api/v1/clients", func(w http.ResponseWriter, r *http.Request) {

	})

	http.Serve(l, mux)
	return nil
}

func (server *HttpOpera) GetVersions(w http.ResponseWriter, req *http.Request) {
	type Ver struct {
		Major     int `json:"major"` //指定字段的tag，实现json字符串的首字母小写
		Minor     int  `json:"minor"`
		Revision   int `json:"revison"`
		Version   string `json:"version"`
	}

	type Versions struct {
		Code     int `json:"code"`
		Server   int `json:"server"`
		Data     Ver `json:"data"`
	}

	group := Versions{
		Code:     0,
		Server:   112,
		Data: Ver{ Major:2, Minor:0, Revision:243, Version:"2.0.243"},
	}

	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	//跨域问题
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control,X-Proxy-Authorization,X-Requested-With,Content-Type")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func ( server*HttpOpera) GetStreams(w http.ResponseWriter, req *http.Request) {
	s :=`{"code":0,"server":4060,"streams":[]}`
	var b []byte = []byte(s)
	//跨域问题
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control,X-Proxy-Authorization,X-Requested-With,Content-Type")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func (server *HttpOpera) GetSummaries(w http.ResponseWriter, req *http.Request) {
	s := `{"code":0,"data"
           :{"ok":true,
		   "now_ms":1520649632780,
		   "self":
		   {"version":"2.0.243",
		   "pid":1114,
		   "ppid":1,
		   "argv":"./objs/srs -c ./conf/srs.conf",
		   "cwd":"/home/winlin/online.srs",
		   "mem_kbyte":108712,
		   "mem_percent":0.21713,
		   "cpu_percent":0.0166556,
		   "srs_uptime":3684310},
		   "system":
		   {"cpu_percent":0.0334448,
		   "disk_read_KBps":0,
		   "disk_write_KBps":56,
		   "disk_busy_percent":0.00536013,
		   "mem_ram_kbyte":500676,
		   "mem_ram_percent":0.524323,
		   "mem_swap_kbyte":0,
		   "mem_swap_percent":0,
		   "cpus":1,"cpus_online":1,
		   "uptime":39379300,
		   "ilde_time":37788600,
		   "load_1m":0.01,
		   "load_5m":0.02,
		   "load_15m":0,
		   "net_sample_time":1520649626778,
		   "net_recv_bytes":1529161877633,
		   "net_send_bytes":1526858524314,
		   "net_recvi_bytes":1868239326324,
		   "net_sendi_bytes":1868266705477,
		   "srs_sample_time":1520649632780,
		   "srs_recv_bytes":221801772138,
		   "srs_send_bytes":27357536969,
		   "conn_sys":106,
		   "conn_sys_et":75,
		   "conn_sys_tw":9,
		   "conn_sys_udp":7,
		   "conn_srs":51}
		   }
		   }`

	var b []byte = []byte(s)
	//跨域问题
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control,X-Proxy-Authorization,X-Requested-With,Content-Type")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}