package stream

import (
	"errors"
	"sync"
	"time"

	"github.com/clzhan/srlive3/avformat"
	"github.com/clzhan/srlive3/log"
	"github.com/clzhan/srlive3/rtmpconst"
)

type StreamObject struct {
	name     string
	duration uint32
	list     []int
	cache    map[int]*MediaGop //cmap.ConcurrentMap
	gop      *MediaGop

	//gopcache map[int]*MediaGop
	subs    []NetStream
	subch   chan NetStream
	sublock sync.RWMutex

	notify chan *int

	lock               sync.RWMutex
	idx                int
	gidx               int
	csize              int
	MetaData           *avformat.MediaFrame
	FirstVideoKeyFrame *avformat.MediaFrame
	FirstAudioKeyFrame *avformat.MediaFrame
	//lastVideoKeyFrame  *MediaFrame

	streamid uint32

	httpsublock sync.RWMutex
	httpsubs    []NetStream

	httpHlssublock sync.RWMutex
	httpHlssubs    []NetStream
}

func New_streamObject(sid string, timeout time.Duration, record bool, csize int) (obj *StreamObject, err error) {
	obj = &StreamObject{
		name:   sid,
		list:   []int{},
		cache:  make(map[int]*MediaGop, csize),
		subs:   []NetStream{},
		notify: make(chan *int, csize*100),
		csize:  csize,
	}
	AddObject(obj)

	go obj.loop(timeout)

	return obj, nil
}

func (m *StreamObject) GetFirstVideoKeyFrame() *avformat.MediaFrame {
	return m.FirstVideoKeyFrame
}

func (m *StreamObject) GetFirstAudioKeyFrame() *avformat.MediaFrame {
	return m.FirstAudioKeyFrame
}

func (m *StreamObject) HttpAttach(c NetStream) {
	m.httpsublock.Lock()
	m.httpsubs = append(m.httpsubs, c)
	m.httpsublock.Unlock()
}

func (m *StreamObject) Attach(c NetStream) {
	m.sublock.Lock()
	m.subs = append(m.subs, c)
	m.sublock.Unlock()
}

func (m *StreamObject) HttpHlsAttach(c NetStream) {
	m.httpHlssublock.Lock()
	m.httpHlssubs = append(m.httpHlssubs, c)
	m.httpHlssublock.Unlock()
}

func (m *StreamObject) ReadGop(idx *int) *MediaGop {
	m.lock.RLock()
	if s, found := m.cache[*idx]; found {
		m.lock.RUnlock()
		return s
	}
	m.lock.RUnlock()
	log.Warn("Gop", m.name, *idx, "Not Found")
	return nil
}

func (m *StreamObject) WriteFrame(s *avformat.MediaFrame) (err error) {
	m.lock.Lock()
	if m.idx >= 0xffffffffffffff {
		m.idx = 0
	}
	s.Idx = m.idx
	m.idx += 1
	m.duration = s.Timestamp

	if s.VideoCodecID == 7 {
		payload := s.Payload.Bytes()
		if uint8(payload[1]) == 0 {
			log.Debug("--------- video first key :", m.name)
			m.FirstVideoKeyFrame = s
			m.streamid = s.StreamId
		}
	}

	if s.Type == rtmpconst.RTMP_MSG_AUDIO {

		payload := s.Payload.Bytes()
		if uint8(payload[1]) == 0 {
			//log.Debug("&&&&&&&&&& audio first key:", m.name)
			m.FirstAudioKeyFrame = s
		}
	}

	if s.Type == rtmpconst.RTMP_MSG_VIDEO && s.IFrame() && m.FirstVideoKeyFrame == nil {
		log.Debug(">>>>", s)
		m.FirstVideoKeyFrame = s
		m.streamid = s.StreamId
		m.lock.Unlock()
		return
	}
	if s.Type == rtmpconst.RTMP_MSG_AUDIO && m.FirstAudioKeyFrame == nil {
		log.Debug(">>>>", s)
		m.FirstAudioKeyFrame = s
		m.lock.Unlock()
		return
	}
	if s.Type == rtmpconst.RTMP_MSG_AMF_META && m.MetaData == nil {
		log.Debug(">>>>", s)
		m.MetaData = s
		m.lock.Unlock()
		return
	}

	if m.gop == nil {
		m.gop = &MediaGop{
			0,
			make([]*avformat.MediaFrame, 0),
			m.FirstVideoKeyFrame,
			m.FirstAudioKeyFrame,
			m.MetaData}
	}

	if len(m.list) >= m.csize {
		idx := m.list[0]
		if s, found := m.cache[idx]; found {
			s.Release()
			delete(m.cache, idx)
		}
		m.list = m.list[1:]
	}
	if s.IFrame() && m.gop.Len() > 0 {
		gop := m.gop
		m.list = append(m.list, gop.idx)
		m.cache[gop.idx] = gop
		//log.Debug("Gop", m.name, gop.idx, gop.Len(), len(m.list))
		m.gop = &MediaGop{
			gop.idx + 1,
			[]*avformat.MediaFrame{s},
			m.FirstVideoKeyFrame,
			m.FirstAudioKeyFrame,
			m.MetaData}

		// m.gop.chunk.wchunks = gop.chunk.wchunks
		// m.gop.freshChunk.writeMetadata(m.metaData)
		// m.gop.freshChunk.writeFullVideo(m.firstVideoKeyFrame)
		// m.gop.freshChunk.writeFullAudio(m.firstAudioKeyFrame)
		// m.gop.freshChunk.writeFullVideo(s)
		// m.gop.chunk.writeVideo(s)
		m.lock.Unlock()
		select {
		case m.notify <- &gop.idx:
		default:
			err = errors.New("buffer full")
		}
		return
	}
	m.gop.frames = append(m.gop.frames, s)
	// if s.Type == RTMP_MSG_VIDEO {
	// 	m.gop.freshChunk.writeVideo(s)
	// 	m.gop.chunk.writeVideo(s)
	// } else if s.Type == RTMP_MSG_AUDIO {
	// 	if !m.gop.audio {
	// 		m.gop.freshChunk.writeFullAudio(s)
	// 	} else {
	// 		m.gop.freshChunk.writeAudio(s)
	// 	}
	// 	m.gop.chunk.writeAudio(s)
	// }
	m.lock.Unlock()
	return
}

func (m *StreamObject) Close() {
	log.Info(m.name, "StreamObject Close")
	RemoveObject(m.name)
	close(m.notify)
}
func (m *StreamObject) NotifyHlsClose() {
	var (
		hls      NetStream
		hlsnsubs = []NetStream{}
		hlssubs  = []NetStream{}
	)

	////hls add
	m.httpHlssublock.Lock()
	hlsnsubs = hlsnsubs[0:0]
	hlssubs = m.httpHlssubs[:]
	m.httpHlssublock.Unlock()
	log.Debug("htthls players", m.name, len(hlssubs))
	for _, hls = range hlssubs {
		hls.Close()
	}
	m.httpHlssublock.Lock()
	m.httpHlssubs = hlsnsubs[:]
	m.httpHlssublock.Unlock()

}

func (m *StreamObject) loop(timeout time.Duration) {
	log.Info(m.name, "stream object is runing")
	defer log.Info(m.name, "stream object is stopped")

	var (
		opened bool
		idx    *int
		w      NetStream
		err    error
		nsubs  = []NetStream{}
		subs   = []NetStream{}
	)

	var (
		h      NetStream
		herr   error
		hnsubs = []NetStream{}
		hsubs  = []NetStream{}
	)

	var (
		hls      NetStream
		hlserr   error
		hlsnsubs = []NetStream{}
		hlssubs  = []NetStream{}
	)

	defer m.clear()
	for {
		select {
		case idx, opened = <-m.notify:
			if !opened {
				//TODO
				m.NotifyHlsClose()
				return
			}
			m.sublock.Lock()
			nsubs = nsubs[0:0]
			subs = m.subs[:]
			m.sublock.Unlock()
			//log.Debug("players", m.name, len(subs))
			for _, w = range subs {
				if err = w.Notify(idx); err != nil {
					log.Error(w, err)
					w.Close()
				} else {
					nsubs = append(nsubs, w)
				}
			}
			m.sublock.Lock()
			m.subs = nsubs[:]
			m.sublock.Unlock()

			////hsh add
			m.httpsublock.Lock()
			hnsubs = hnsubs[0:0]
			hsubs = m.httpsubs[:]
			m.httpsublock.Unlock()
			//log.Debug("httflv players", m.name, len(hsubs))
			for _, h = range hsubs {
				if herr = h.Notify(idx); herr != nil {
					log.Error(h, herr)
					h.Close()
				} else {
					hnsubs = append(hnsubs, h)
				}
			}
			m.httpsublock.Lock()
			m.httpsubs = hnsubs[:]
			m.httpsublock.Unlock()

			////hls add
			m.httpHlssublock.Lock()
			hlsnsubs = hlsnsubs[0:0]
			hlssubs = m.httpHlssubs[:]
			m.httpHlssublock.Unlock()
			//log.Debug("htthls players", m.name, len(hlssubs))
			for _, hls = range hlssubs {
				if hlserr = hls.Notify(idx); hlserr != nil {
					log.Error(hls, herr)
					hls.Close()
				} else {
					hlsnsubs = append(hlsnsubs, hls)
				}
			}
			m.httpHlssublock.Lock()
			m.httpHlssubs = hlsnsubs[:]
			m.httpHlssublock.Unlock()

		case <-time.After(timeout):
			m.Close()
		}
	}
}

func (m *StreamObject) clear() {
	m.sublock.Lock()
	for _, w := range m.subs {
		w.Close()
	}
	m.subs = m.subs[0:0]
	m.sublock.Unlock()

	//hsh add
	m.httpsublock.Lock()
	for _, h := range m.httpsubs {
		h.Close()
	}
	m.httpsubs = m.httpsubs[0:0]
	m.httpsublock.Unlock()
}


type NetStream interface {
	NsID() int
	Name() string
	String() string
	Notify(idx *int) error
	Close() error
	StreamObject() *StreamObject
}
