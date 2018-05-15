package hls

// flv format
// AACPacketType IF SoundFormat == 10 UI8
// The following values are defined:
//     0 = AAC sequence header
//     1 = AAC raw
const (
	// RtmpCodecAudioTypeReserved set to the max value to reserved, for array map.
	RtmpCodecAudioTypeReserved = 2

	// RtmpCodecAudioTypeSequenceHeader audio type sequence header
	RtmpCodecAudioTypeSequenceHeader = 0
	// RtmpCodecAudioTypeRawData audio raw data
	RtmpCodecAudioTypeRawData = 1
)

// E.4.3.1 VIDEODATA
// Frame Type UB [4]
// Type of video frame. The following values are defined:
//     1 = key frame (for AVC, a seekable frame)
//     2 = inter frame (for AVC, a non-seekable frame)
//     3 = disposable inter frame (H.263 only)
//     4 = generated key frame (reserved for server use only)
//     5 = video info/command frame
const (
	// RtmpCodecVideoAVCFrameReserved set to the max value to reserved, for array map.
	RtmpCodecVideoAVCFrameReserved = 0
	// RtmpCodecVideoAVCFrameReserved1 .
	RtmpCodecVideoAVCFrameReserved1 = 6

	// RtmpCodecVideoAVCFrameKeyFrame video h264 key frame
	RtmpCodecVideoAVCFrameKeyFrame = 1
	// RtmpCodecVideoAVCFrameInterFrame video h264 inter frame
	RtmpCodecVideoAVCFrameInterFrame = 2
	// RtmpCodecVideoAVCFrameDisposableInterFrame .
	RtmpCodecVideoAVCFrameDisposableInterFrame = 3
	// RtmpCodecVideoAVCFrameGeneratedKeyFrame .
	RtmpCodecVideoAVCFrameGeneratedKeyFrame = 4
	// RtmpCodecVideoAVCFrameVideoInfoFrame .
	RtmpCodecVideoAVCFrameVideoInfoFrame = 5
)

// AVCPacketType IF CodecID == 7 UI8
// The following values are defined:
//     0 = AVC sequence header
//     1 = AVC NALU
//     2 = AVC end of sequence (lower level NALU sequence ender is
//         not required or supported)
const (
	// set to the max value to reserved, for array map.
	RtmpCodecVideoAVCTypeReserved = 3

	// RtmpCodecVideoAVCTypeSequenceHeader .
	RtmpCodecVideoAVCTypeSequenceHeader = 0
	// RtmpCodecVideoAVCTypeNALU .
	RtmpCodecVideoAVCTypeNALU = 1
	// RtmpCodecVideoAVCTypeSequenceHeaderEOF .
	RtmpCodecVideoAVCTypeSequenceHeaderEOF = 2
)

// E.4.3.1 VIDEODATA
// CodecID UB [4]
// Codec Identifier. The following values are defined:
//     2 = Sorenson H.263
//     3 = Screen video
//     4 = On2 VP6
//     5 = On2 VP6 with alpha channel
//     6 = Screen video version 2
//     7 = AVC
//     13 = HEVC
const (
	// RtmpCodecVideoReserved set to the max value to reserved, for array map.
	RtmpCodecVideoReserved = 0
	// RtmpCodecVideoReserved1 .
	RtmpCodecVideoReserved1 = 1
	// RtmpCodecVideoReserved2 .
	RtmpCodecVideoReserved2 = 8

	// RtmpCodecVideoSorensonH263 .
	RtmpCodecVideoSorensonH263 = 2
	// RtmpCodecVideoScreenVideo .
	RtmpCodecVideoScreenVideo = 3
	// RtmpCodecVideoOn2VP6 .
	RtmpCodecVideoOn2VP6 = 4
	// RtmpCodecVideoOn2VP6WithAlphaChannel .
	RtmpCodecVideoOn2VP6WithAlphaChannel = 5
	// RtmpCodecVideoScreenVideoVersion2 .
	RtmpCodecVideoScreenVideoVersion2 = 6
	// RtmpCodecVideoAVC .
	RtmpCodecVideoAVC = 7
	// RtmpCodecVideoHEVC h265
	RtmpCodecVideoHEVC = 13
)

// SoundFormat UB [4]
// Format of SoundData. The following values are defined:
//     0 = Linear PCM, platform endian
//     1 = ADPCM
//     2 = MP3
//     3 = Linear PCM, little endian
//     4 = Nellymoser 16 kHz mono
//     5 = Nellymoser 8 kHz mono
//     6 = Nellymoser
//     7 = G.711 A-law logarithmic PCM
//     8 = G.711 mu-law logarithmic PCM
//     9 = reserved
//     10 = AAC
//     11 = Speex
//     14 = MP3 8 kHz
//     15 = Device-specific sound
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher.
// Speex is supported in Flash Player 10 and higher.
const (
	// RtmpCodecAudioReserved1 set to the max value to reserved, for array map.
	RtmpCodecAudioReserved1 = 16

	// RtmpCodecAudioLinearPCMPlatformEndian .
	RtmpCodecAudioLinearPCMPlatformEndian = 0
	// RtmpCodecAudioADPCM .
	RtmpCodecAudioADPCM = 1
	// RtmpCodecAudioMP3 .
	RtmpCodecAudioMP3 = 2
	// RtmpCodecAudioLinearPCMLittleEndian .
	RtmpCodecAudioLinearPCMLittleEndian = 3
	// RtmpCodecAudioNellymoser16kHzMono .
	RtmpCodecAudioNellymoser16kHzMono = 4
	// RtmpCodecAudioNellymoser8kHzMono .
	RtmpCodecAudioNellymoser8kHzMono = 5
	// RtmpCodecAudioNellymoser .
	RtmpCodecAudioNellymoser = 6
	// RtmpCodecAudioReservedG711AlawLogarithmicPCM .
	RtmpCodecAudioReservedG711AlawLogarithmicPCM = 7
	// RtmpCodecAudioReservedG711MuLawLogarithmicPCM .
	RtmpCodecAudioReservedG711MuLawLogarithmicPCM = 8
	// RtmpCodecAudioReserved .
	RtmpCodecAudioReserved = 9
	// RtmpCodecAudioAAC .
	RtmpCodecAudioAAC = 10
	// RtmpCodecAudioSpeex .
	RtmpCodecAudioSpeex = 11
	// RtmpCodecAudioReservedMP3Of8kHz .
	RtmpCodecAudioReservedMP3Of8kHz = 14
	// RtmpCodecAudioReservedDeviceSpecificSound .
	RtmpCodecAudioReservedDeviceSpecificSound = 15
)

// the FLV/RTMP supported audio sample rate.
// Sampling rate. The following values are defined:
// 0 = 5.5 kHz = 5512 Hz
// 1 = 11 kHz = 11025 Hz
// 2 = 22 kHz = 22050 Hz
// 3 = 44 kHz = 44100 Hz
const (
	// RtmpCodecAudioSampleRateReserved set to the max value to reserved, for array map.
	RtmpCodecAudioSampleRateReserved = 4

	// RtmpCodecAudioSampleRate5512 .
	RtmpCodecAudioSampleRate5512 = 0
	// RtmpCodecAudioSampleRate11025 .
	RtmpCodecAudioSampleRate11025 = 1
	// RtmpCodecAudioSampleRate22050 .
	RtmpCodecAudioSampleRate22050 = 2
	// RtmpCodecAudioSampleRate44100 .
	RtmpCodecAudioSampleRate44100 = 3
)

// the FLV/RTMP supported audio sample size.
// Size of each audio sample. This parameter only pertains to
// uncompressed formats. Compressed formats always decode
// to 16 bits internally.
// 0 = 8-bit samples
// 1 = 16-bit samples
const (
	// RtmpCodecAudioSampleSizeReserved set to the max value to reserved, for array map.
	RtmpCodecAudioSampleSizeReserved = 2

	// RtmpCodecAudioSampleSize8bit .
	RtmpCodecAudioSampleSize8bit = 0
	// RtmpCodecAudioSampleSize16bit .
	RtmpCodecAudioSampleSize16bit = 1
)

// the FLV/RTMP supported audio sound type/channel.
// Mono or stereo sound
// 0 = Mono sound
// 1 = Stereo sound
const (
	// RtmpCodecAudioSoundTypeReserved set to the max value to reserved, for array map.
	RtmpCodecAudioSoundTypeReserved = 2

	// RtmpCodecAudioSoundTypeMono .
	RtmpCodecAudioSoundTypeMono = 0
	// RtmpCodecAudioSoundTypeStereo .
	RtmpCodecAudioSoundTypeStereo = 1
)

// TokenStr token for auth
const TokenStr = "?token="

const (
	// RtmpMaxFmt0HeaderSize is the max rtmp header size:
	//   1bytes basic header,
	//   11bytes message header,
	//   4bytes timestamp header,
	//   that is, 1+11+4=16bytes.
	RtmpMaxFmt0HeaderSize = 16
)

//RtmpRole define
const (
	// RtmpRoleUnknown role undefined
	RtmpRoleUnknown = 0
	// RtmpRoleFMLEPublisher role fmle publisher
	RtmpRoleFMLEPublisher = 1
	// RtmpRoleFlashPublisher role flash publisher
	RtmpRoleFlashPublisher = 2
	// RtmpRolePlayer role player
	RtmpRolePlayer = 3
)

const (
	// RtmpTimeJitterFull time jitter full mode, to ensure stream start at zero, and ensure stream monotonically increasing.
	RtmpTimeJitterFull = 0x01
	// RtmpTimeJitterZero zero mode, only ensure sttream start at zero, ignore timestamp jitter.
	RtmpTimeJitterZero = 0x02
	// RtmpTimeJitterOff off mode, disable the time jitter algorithm, like atc.
	RtmpTimeJitterOff = 0x03
)

const (
	// PureAudioGuessCount for 26ms per audio packet,
	// 115 packets is 3s.
	PureAudioGuessCount = 115
)

const (
	// MaxJitterMs max time delta, which is the between localtime and last packet time
	MaxJitterMs = 500
	// DefaultFrameTimeMs default time delta, which is the between localtime and last packet time
	DefaultFrameTimeMs = 40
)
