package rtmp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/clzhan/srlive3/log"
)

type Infomation map[string]interface{}
type Args interface{}

type NetConnection interface {
	Connect(command string, args ...Args) error
	Call(command string, args ...Args) error
	Connected() bool
	Close()
	URL() string
}

var gstreamid = uint32(64)

func gen_next_stream_id(chunkid uint32) uint32 {
	gstreamid += 1
	return gstreamid
}

func newconn(conn net.Conn, srv *Server) (c *RtmpNetConnection) {
	c = new(RtmpNetConnection)
	c.remoteAddr = conn.RemoteAddr().String()
	c.server = srv
	c.readChunkSize = RTMP_DEFAULT_CHUNK_SIZE
	c.writeChunkSize = RTMP_DEFAULT_CHUNK_SIZE
	c.createTime = time.Now()
	c.bandwidth = 512 << 10
	c.conn = conn
	c.br = bufio.NewReader(conn)
	c.bw = bufio.NewWriter(conn)
	c.buf = bufio.NewReadWriter(c.br, c.bw)
	c.rchunks = make(map[uint32]*RtmpChunk)
	c.wchunks = make(map[uint32]*RtmpChunk)
	c.w_buffer = bytes.NewBuffer(nil)
	c.buffer = bytes.NewBuffer(nil)
	c.nextStreamId = gen_next_stream_id
	c.objectEncoding = 0
	return
}

type RtmpNetConnection struct {
	remoteAddr       string
	url              string
	app              string
	createTime       time.Time
	readChunkSize    int
	writeChunkSize   int
	readTimeout      time.Duration
	writeTimeout     time.Duration
	bandwidth        uint32
	limitType        byte
	wirtesequencenum uint32
	sequencenum      uint32
	totalreadbytes   uint32
	totalwritebytes  uint32
	server           *Server           // the Server on which the connection arrived
	conn             net.Conn          // i/o connection
	buf              *bufio.ReadWriter // buffered(lr,rwc), reading from bufio->limitReader->sr->rwc
	br               *bufio.Reader
	bw               *bufio.Writer
	lock             sync.Mutex // guards the following
	rchunks          map[uint32]*RtmpChunk
	wchunks          map[uint32]*RtmpChunk
	connected        bool
	nextStreamId     func(chunkid uint32) uint32
	streamid         uint32
	objectEncoding   int
	w_buffer         *bytes.Buffer
	buffer           *bytes.Buffer
}

func (c *RtmpNetConnection) Connect(URL string, args ...Args) error {
	//rtmp://host:port/app
	log.Debug(URL)
	p := strings.Split(URL, "/")
	address := p[2]
	log.Debug(address)
	app := p[3]
	log.Debug(app)
	host := strings.Split(address, ":")
	if len(host) == 1 {
		address += ":1935"
	}
	log.Debug(address)
	conn, err := net.DialTimeout("tcp", address, time.Second*30)
	if err != nil {
		return err
	}
	c.conn = conn
	c.app = app
	c.remoteAddr = conn.RemoteAddr().String()
	c.br = bufio.NewReader(conn)
	c.bw = bufio.NewWriter(conn)
	c.buf = bufio.NewReadWriter(c.br, c.bw)
	if !client_simple_handshake(c.buf) {
		return errors.New("Client Handshake Fail")
	}
	err = sendConnect(c, app, "", "", URL)
	if err != nil {
		return err
	}
	for {
		msg, err := readMessage(c)
		if err != nil {
			c.Close()
			return err
		}
		log.Debug(msg)
		if _, ok := msg.(*ReplyMessage); ok {
			dec := newDecoder(msg.Body())
			reply := new(ReplyConnectMessage)
			reply.RtmpHeader = msg.Header()
			reply.Command = readString(dec)
			reply.TransactionId = readNumber(dec)
			//dec.readNull()
			reply.Properties = readObject(dec)
			reply.Infomation = readObject(dec)
			log.Debug(reply)
			if NetConnection_Connect_Success == getString(reply.Infomation, "code") {
				c.connected = true
				return nil
			} else {
				return errors.New(getString(reply.Infomation, "code"))
			}
		}
	}
	return nil
}

func (c *RtmpNetConnection) Call(command string, args ...Args) error {
	return nil
}
func (c *RtmpNetConnection) Connected() bool {
	return c.connected
}
func (c *RtmpNetConnection) Close() {
	if c.conn == nil {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conn.Close()
	c.connected = false
}
func (c *RtmpNetConnection) URL() string {
	return c.url
}
