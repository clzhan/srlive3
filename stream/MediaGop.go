package stream

import "github.com/clzhan/srlive3/avformat"

type MediaGop struct {
	idx    int
	frames []*avformat.MediaFrame
	//freshChunk *RtmpChunker
	//chunk      *RtmpChunker
	//audio bool
	videoConfig *avformat.MediaFrame
	audioConfig *avformat.MediaFrame
	metaConfig  *avformat.MediaFrame
}

func (o *MediaGop) Frames() []*avformat.MediaFrame {
	return o.frames
}

func (o *MediaGop) Release() {
	for _, f := range o.frames {
		f.Release()
	}
	o.frames = o.frames[0:0]
	//o.freshChunk.Reset()
	//o.chunk.Reset()
	//o.freshChunk = nil
	//o.chunk = nil
}

func (o *MediaGop) Len() int {
	return len(o.frames)
}

