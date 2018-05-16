package rtmp

import (
	"strings"
	"time"
	"bytes"
	"github.com/clzhan/srlive3/log"
	"fmt"
	"errors"
	"strconv"
)

var clientHandler ClientHandler = new(DefaultClientHandler)

type RtmpURL struct {
	protocol     string
	host         string
	port         uint16
	app          string
	instanceName string
}

// Parse url
//
// To connect to Flash Media Server, pass the URI of the application on the server.
// Use the following syntax (items in brackets are optional):
//
// protocol://host[:port]/[appname[/instanceName]]
func ParseURL(url string) (rtmpURL RtmpURL, err error) {
	s1 := strings.SplitN(url, "://", 2)
	if len(s1) != 2 {
		err = errors.New(fmt.Sprintf("Parse url %s error. url invalid.", url))
		return
	}
	rtmpURL.protocol = strings.ToLower(s1[0])
	s1 = strings.SplitN(s1[1], "/", 2)
	if len(s1) != 2 {
		err = errors.New(fmt.Sprintf("Parse url %s error. no app!", url))
		return
	}
	s2 := strings.SplitN(s1[0], ":", 2)
	if len(s2) == 2 {
		var port int
		port, err = strconv.Atoi(s2[1])
		if err != nil {
			err = errors.New(fmt.Sprintf("Parse url %s error. port error: %s.", url, err.Error()))
			return
		}
		if port > 65535 || port <= 0 {
			err = errors.New(fmt.Sprintf("Parse url %s error. port error: %d.", url, port))
			return
		}
		rtmpURL.port = uint16(port)
	} else {
		rtmpURL.port = 1935
	}
	if len(s2[0]) == 0 {
		err = errors.New(fmt.Sprintf("Parse url %s error. host is empty.", url))
		return
	}
	rtmpURL.host = s2[0]

	s2 = strings.SplitN(s1[1], "/", 2)
	rtmpURL.app = s2[0]
	if len(s2) == 2 {
		rtmpURL.instanceName = s2[1]
	}
	return
}
func ConnectPull(url string) (s *RtmpNetStream, err error) {
	// protocol://host[:port]/[appname[/instanceName]]
	rtmpURL,err := ParseURL(url);
	if err != nil {
		return nil,err
	}
	file := rtmpURL.app + "/" + rtmpURL.instanceName
	addr := rtmpURL.protocol + "://" + rtmpURL.host + ":" + strconv.Itoa(int(rtmpURL.port))  + "/" + rtmpURL.app

	conn := newNetConnection()

	err = conn.Connect(addr)
	if err != nil {
		log.Error("err %v",err)
		return
	}
	s = newNetStream(conn, nil, clientHandler)

	log.Info("play......")
	s.play(file,rtmpURL.instanceName, "live")
	return
}
func newNetConnection() (c *RtmpNetConnection) {
	c = new(RtmpNetConnection)
	c.readChunkSize = RTMP_DEFAULT_CHUNK_SIZE
	c.writeChunkSize = RTMP_DEFAULT_CHUNK_SIZE
	c.createTime = time.Now()
	c.bandwidth = 512 << 10
	c.rchunks = make(map[uint32]*RtmpChunk)
	c.wchunks = make(map[uint32]*RtmpChunk)
	c.buffer = bytes.NewBuffer(nil)
	c.w_buffer = bytes.NewBuffer(nil)
	c.nextStreamId = gen_next_stream_id
	c.objectEncoding = 0
	return
}
