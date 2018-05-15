package hls

import (

	"github.com/clzhan/srlive2/log"
	"github.com/clzhan/srlive3/avformat"
)

// SourceStream delivery RTMP stream to HLS(m3u8 and ts),
type SourceStream struct {
	muxer *hlsMuxer
	cache *hlsCache

	codec  *avcAacCodec
	sample *codecSample
	jitter *TimeJitter

	// we store the stream dts,
	// for when we notice the hls cache to publish,
	// it need to know the segment start dts.
	//
	// for example. when republish, the stream dts will
	// monotonically increase, and the ts dts should start
	// from current dts.
	//
	// or, simply because the HlsCache never free when unpublish,
	// so when publish or republish it must start at stream dts,
	// not zero dts.
	streamDts int64
}

// NewSourceStream new a hls source stream
func NewSourceStream() *SourceStream {
	return &SourceStream{
		muxer: newHlsMuxer(),
		cache: newHlsCache(),

		codec:  newAvcAacCodec(),
		sample: newCodecSample(),
		jitter: NewTimeJitter(),
	}
}

// OnAudio process on audio data, mux to ts
func (hls *SourceStream) OnAudio(msg *avformat.MediaFrame) (err error) {

	hls.sample.clear()
	payload := msg.Payload.Bytes()
	if err = hls.codec.audioAacDemux(payload, hls.sample); err != nil {
		log.Error("hls codec demux audio failed, err=", err)
		return
	}

	if hls.codec.audioCodecID != RtmpCodecAudioAAC {
		log.Error("codec audio codec id is not aac, codeID=", hls.codec.audioCodecID)
		return
	}

	// ignore sequence header
	if RtmpCodecAudioTypeSequenceHeader == hls.sample.aacPacketType {
		if err = hls.cache.onSequenceHeader(hls.muxer); err != nil {
			log.Error("hls cache on sequence header failed, err=", err)
			return
		}

		return
	}

	hls.jitter.Correct(msg, 0, 0, RtmpTimeJitterFull)

	// the pts calc from rtmp/flv header
	pts := int64(msg.Timestamp * 90)

	// for pure audio, update the stream dts also
	hls.streamDts = pts

	if err = hls.cache.writeAudio(hls.codec, hls.muxer, pts, hls.sample); err != nil {
		log.Error("hls cache write audio failed, err=", err)
		return
	}

	return
}

// OnVideo process on video data, mux to ts
func (hls *SourceStream) OnVideo(msg *avformat.MediaFrame) (err error) {

	hls.sample.clear()
	payload := msg.Payload.Bytes()
	if err = hls.codec.videoAvcDemux(payload, hls.sample); err != nil {
		log.Error("hls codec demuxer video failed, err=", err)
		return
	}

	// ignore info frame,
	if RtmpCodecVideoAVCFrameVideoInfoFrame == hls.sample.frameType {
		return
	}

	if hls.codec.videoCodecID != RtmpCodecVideoAVC {
		return
	}

	// ignore sequence header
	if RtmpCodecVideoAVCFrameKeyFrame == hls.sample.frameType &&
		RtmpCodecVideoAVCTypeSequenceHeader == hls.sample.frameType {
		return hls.cache.onSequenceHeader(hls.muxer)
	}

	hls.jitter.Correct(msg, 0, 0, RtmpTimeJitterFull)

	//PTS,和DTS的计算。如果是根据ES流计算怎么算，我还没想到容易的方法。
	// 但是如果根据FLV格式转换，PTS就是(flvTagHeader.timestamp +videoTagHeader.CompositionTime) * 90 , 为啥是90呢？
	// flv里面的时间戳的单位都是毫秒的，1/1000秒。mpegts的系统时钟为27MHZ，这里要除以300(规定的除以300，参考ISO-13818-1)。
	// 也就是90000hz，一秒钟90000个周期，所以，PTS代表的周期flv的时间戳*90000/1000 ，90也就是这么来的
	dts := msg.Timestamp * 90
	hls.streamDts = int64(dts)

	if err = hls.cache.writeVideo(hls.codec, hls.muxer, int64(dts), hls.sample); err != nil {
		log.Error("hls cache write video failed")
		return
	}

	return
}

// OnPublish publish stream event, continue to write the m3u8,
// for the muxer object not destroyed.
func (hls *SourceStream) OnPublish(app string, stream string) (err error) {

	if err = hls.cache.onPublish(hls.muxer, app, stream, hls.streamDts); err != nil {
		log.Error("err.........")
		return
	}
	return
}

// OnUnPublish the unpublish event, only close the muxer, donot destroy the
// muxer, for when we continue to publish, the m3u8 will continue.
func (hls *SourceStream) OnUnPublish() (err error) {

	if err = hls.cache.onUnPublish(hls.muxer); err != nil {
		return
	}

	return
}

func (hls *SourceStream) hlsMux() {

	return
}
