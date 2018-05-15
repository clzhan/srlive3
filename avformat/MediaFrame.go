package avformat

import (
	"bytes"
	"fmt"
	"io"
)

func NewMediaFrame() *MediaFrame {
	x := &MediaFrame{
		count: 1,
		Payload: bytes.NewBuffer(nil)}
	return x
}

type MediaFrame struct {
	Idx            int
	Timestamp      uint32
	Type           byte //8 audio,9 video
	VideoFrameType byte //4bit
	VideoCodecID   byte //4bit
	AudioFormat    byte //4bit
	SamplingRate   byte //2bit
	SampleLength   byte //1bit
	AudioType      byte //1bit
	Payload        *bytes.Buffer
	StreamId       uint32
	count          int32
}

func (p *MediaFrame) IFrame() bool {
	return p.VideoFrameType == 1 || p.VideoFrameType == 4
}

//	RTMP_MSG_AUDIO = 8
//RTMP_MSG_VIDEO = 9

func (p *MediaFrame) String() string {
	if p == nil {
		return "<nil>"
	}
	if p.Type == 8/*RTMP_MSG_AUDIO*/{
		return fmt.Sprintf("%v Audio Frame Timestamp/%v Type/%v AudioFromat/%v SampleRate/%v SampleLength/%v AudioType/%v Payload/%v StreamId/%v", p.Idx, float64(p.Timestamp)/1000.0, p.Type, Audioformat[p.AudioFormat], Samplerate[p.SamplingRate], Samplelength[p.SampleLength], Audiotype[p.AudioType], p.Payload.Len(), p.StreamId)
	} else if p.Type == 9/*RTMP_MSG_VIDEO */{
		return fmt.Sprintf("%v Video Frame Timestamp/%v Type/%v VideoFrameType/%v VideoCodecID/%v Payload/%v StreamId/%v", p.Idx, float64(p.Timestamp)/1000.0, p.Type, Videoframetype[p.VideoFrameType], Videocodec[p.VideoCodecID], p.Payload.Len(), p.StreamId)
	}
	return fmt.Sprintf("%v Frame Timestamp/%v Type/%v Payload/%v StreamId/%v", p.Idx, p.Timestamp, p.Type, p.Payload.Len(), p.StreamId)
}

func (o *MediaFrame) Ref() *MediaFrame {
	//atomic.AddInt32(&o.count, 1)
	return o
}

func (o *MediaFrame) Release() {
	// if nc := atomic.AddInt32(&o.count, -1); nc <= 0 {
	// 	select {
	// 	case o.p.pool <- o:
	// 	default:
	// 	}
	// }
}

func (o *MediaFrame) Bytes() []byte {
	return o.Payload.Bytes()
}

func (o *MediaFrame) WriteTo(w io.Writer) (int, error) {
	return w.Write(o.Payload.Bytes())
}
