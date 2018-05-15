package avformat

const (
	//当前的音频编码
	SOUND_MP3                   = 2
	SOUND_NELLYMOSER_16KHZ_MONO = 4
	SOUND_NELLYMOSER_8KHZ_MONO  = 5
	SOUND_NELLYMOSER            = 6
	SOUND_ALAW                  = 7
	SOUND_MULAW                 = 8
	SOUND_AAC                   = 10
	SOUND_SPEEX                 = 11

	AAC_SEQHDR = 0 //理解为第一帧的序列包
	AAC_RAW    = 1 //裸数据
)

const (
	AVC_SEQHDR = 0 //理解为第一帧序列包，需要解析关键参数
	AVC_NALU   = 1
	AVC_EOS    = 2

	FRAME_KEY   = 1
	FRAME_INTER = 2

	VIDEO_H264 = 7
)
