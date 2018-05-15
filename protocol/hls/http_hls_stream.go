package hls

import (
	"errors"
	"fmt"

	"github.com/clzhan/srlive3/log"
	"github.com/clzhan/srlive3/conf"
	"github.com/clzhan/srlive3/rtmpconst"
	"github.com/clzhan/srlive3/stream"
)

type HttpHlsStream struct {
	notify     chan *int
	obj        *stream.StreamObject
	streamName string
	nsid       int
	err        error
	closed     chan bool

	first_video bool
	first_audio bool

	app   string
	live   string

	hlssource *SourceStream
}

func nsid() int {
	id, _ := conf.Snow.Next()
	return int(id)
}

func NewHttpHlsStream(streamName string) (s *HttpHlsStream) {
	stream := new(HttpHlsStream)
	stream.nsid = nsid()
	stream.notify = make(chan *int, 30)
	stream.closed = make(chan bool)
	stream.streamName = streamName

	stream.hlssource = NewSourceStream()

	stream.first_audio = false
	stream.first_video = false

	return stream
}

func (s *HttpHlsStream) SetObj(o *stream.StreamObject) {
	s.obj = o
}

func (s *HttpHlsStream) isClosed() bool {
	select {
	case _, ok := <-s.closed:
		if !ok {
			return true
		}
	default:
	}
	return false
}

func (s *HttpHlsStream) Close() error {
	log.Debug("HttpFlsStream Close Start ")

	log.Debug("HttpFlsStream Close streamName :", s.streamName)

	if s.isClosed() {
		return nil
	}

	close(s.closed)
	close(s.notify)
	log.Debug("HttpFlsStream Close End ")
	return nil
}

func (s *HttpHlsStream) NsID() int {
	return s.nsid
}

func (s *HttpHlsStream) Name() string {
	return s.streamName
}

func (s *HttpHlsStream) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v %v closed:%v", s.nsid, s.streamName, s.isClosed())
}

func (s *HttpHlsStream) StreamObject() *stream.StreamObject {
	return s.obj
}

func (s *HttpHlsStream) Notify(idx *int) error {

	if s.isClosed() {
		return errors.New("remote connection? " + " closed")
	}

	select {
	case s.notify <- idx:
		return nil
	default:
		log.Warn("romode addr?" + "httphls stream notify buffer full")
	}
	return nil
}
func (s *HttpHlsStream) OnPublish(app string, stream string) (err error) {
	s.app  = app
	s.live = stream
	return
}

func (s *HttpHlsStream) WriteLoop() {

	var (
		notify = s.notify
		opened bool
		idx    *int
		obj    *stream.StreamObject
		gop    *stream.MediaGop
		err    error
	)
	if s.hlssource != nil{
		err = s.hlssource.OnPublish(s.app, s.live)
	}

	obj = s.obj
	for {
		select {
		case idx, opened = <-notify:
			if !opened {
				return
			}
			gop = obj.ReadGop(idx)
			if gop != nil {
				frames := gop.Frames()[:]
				for _, frame := range frames {

					if frame.Type == rtmpconst.RTMP_MSG_VIDEO {

						if s.first_video == false{
							if obj.GetFirstVideoKeyFrame() != nil {
								if s.hlssource != nil{
									err = s.hlssource.OnVideo(obj.GetFirstVideoKeyFrame())
								}
							}
							s.first_video = true
						}

						if s.hlssource != nil{
							err = s.hlssource.OnVideo(frame)
						}
						//payload := frame.Payload.Bytes()

						//err = s.SendTag(w, r, payload, RTMP_MSG_VIDEO, frame.Timestamp)
						//log.Info("This is a Video......")
					} else if frame.Type == rtmpconst.RTMP_MSG_AUDIO {

						if s.first_audio == false{
							if obj.GetFirstAudioKeyFrame() != nil {
								if s.hlssource != nil{
									err = s.hlssource.OnAudio(obj.GetFirstAudioKeyFrame())
								}
							}
							s.first_audio = true
						}

						if s.hlssource != nil{
							err = s.hlssource.OnAudio(frame)
						}
						//payload := frame.Payload.Bytes()

						//err = s.SendTag(w, r, payload, RTMP_MSG_AUDIO, frame.Timestamp)
						//log.Info("This is a Audio......")
					}

					if err != nil {
						log.Error("http net write errr......")
						return
					}

				}
			}
		}
	}
}
