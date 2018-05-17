package rtmp


import (
	"time"

	"github.com/clzhan/srlive3/log"
	"github.com/clzhan/srlive3/rtmpconst"
	"github.com/clzhan/srlive3/stream"
	"errors"
)

type ServerHandler interface {
	OnPublishing(s *RtmpNetStream) error
	OnPlaying(s *RtmpNetStream) error
	OnClosed(s *RtmpNetStream)
	OnError(s *RtmpNetStream, err error)
}

type ClientHandler interface {
	OnPublishStart(s *RtmpNetStream) error
	OnPlayStart(s *RtmpNetStream) error
	OnClosed(s *RtmpNetStream)
	OnError(s *RtmpNetStream, err error)
}

type DefaultClientHandler struct {
}

func (this *DefaultClientHandler) OnPublishStart(s *RtmpNetStream) error {
	//先去找是否有推送的
	if obj, found := stream.FindObject(s.streamName); !found {
		return errors.New("OnPublishStart not found")
	} else {
		s.obj = obj
		s.obj.Attach(s)
		return nil
	}

	return nil
}
func (this *DefaultClientHandler) OnPlayStart(s *RtmpNetStream) error {
	if obj, found := stream.FindObject(s.streamName); !found {
		log.Info("not found %s, New_streamObject ",s.streamName)
		obj, err := stream.New_streamObject(s.streamName, 90*time.Second, true, 10)
		if err != nil {
			return err
		}
		s.obj = obj
	} else {
		s.obj = obj
	}
	return nil
}
func (this *DefaultClientHandler) OnClosed(s *RtmpNetStream) {
	if s.mode == rtmpconst.MODE_PRODUCER {
		log.Infof("RtmpNetStream client Publish %s %s closed", s.conn.remoteAddr, s.path)
		if obj, found := stream.FindObject(s.streamName); found {
			obj.Close()
		}
	}
}
func (this *DefaultClientHandler) OnError(s *RtmpNetStream, err error) {
	log.Errorf("RtmpNetStream %v %s %s %+v", s.mode, s.conn.remoteAddr, s.path, err)
	s.Close()
}

type DefaultServerHandler struct {
}

// 发布者成功发布流后,就启动广播
func (p *DefaultServerHandler) OnPublishing(s *RtmpNetStream) error {
	// 在广播中发现这个广播已经存在,那么就认为这个广播是无效的.(例如已经发布ip/myapp/mystream这个广播,再次发布ip/app/mystream,就认为这个广播是无效的)
	if obj, found := stream.FindObject(s.streamName); !found {
		obj, err := stream.New_streamObject(s.streamName, 90*time.Second, true, 10)
		if err != nil {
			return err
		}
		s.obj = obj
	} else {
		s.obj = obj
	}

	return nil
}

// 订阅者成功订阅流后,就将订阅者添加进广播中
func (p *DefaultServerHandler) OnPlaying(s *RtmpNetStream) error {
	// 根据订阅者(s)提供的信息,来查找订阅者需要订阅的广播,如果找到了,那么就让这个广播添加这个订阅者
	if obj, found := stream.FindObject(s.streamName); !found {

		log.Info("OnPlaying FindObject not found...",s.streamName)

		obj, err := stream.New_streamObject(s.streamName, 90*time.Second, true, 10)
		if err != nil {
			return err
		}
		s.obj = obj
	} else {

		log.Info("OnPlaying FindObject found...",s.streamName)

		s.obj = obj
	}
	s.obj.Attach(s)

	go s.writeLoop()
	return nil
}

func (p *DefaultServerHandler) OnClosed(s *RtmpNetStream) {
	// mode := "UNKNOWN"
	// if s.mode == MODE_CONSUMER {
	// 	mode = "CONSUMER"
	// } else if s.mode == MODE_PROXY {
	// 	mode = "PROXY"
	// } else if s.mode == MODE_CONSUMER|MODE_PRODUCER {
	// 	mode = "PRODUCER|CONSUMER"
	// }
	// log.Infof("RtmpNetStream %v %s %s closed", mode, s.conn.remoteAddr, s.path)
	// if d, ok := find_broadcast(s.path); ok {
	// 	if s.mode == MODE_PRODUCER {
	// 		d.stop()
	// 	} else if s.mode == MODE_CONSUMER {
	// 		d.removeConsumer(s)
	// 	} else if s.mode == MODE_CONSUMER|MODE_PRODUCER {
	// 		d.removeConsumer(s)
	// 		d.stop()
	// 	}
	// }
	if s.mode ==rtmpconst.MODE_PRODUCER {
		log.Infof("RtmpNetStream server Publish %s %s closed", s.conn.remoteAddr, s.path)
		if obj, found := stream.FindObject(s.streamName); found {
			obj.Close()
		}
	}
}

func (p *DefaultServerHandler) OnError(s *RtmpNetStream, err error) {
	log.Errorf("RtmpNetStream %s %s %+v", s.conn.remoteAddr, s.path, err)
	s.Close()
}
