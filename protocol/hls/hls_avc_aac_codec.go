package hls

import (
	"fmt"
	"log"

	"encoding/binary"

)

// the h264/avc and aac codec, for media stream.
//
// to demux the FLV/RTMP video/audio packet to sample,
// add each NALUs of h.264 as a sample unit to sample,
// while the entire aac raw data as a sample unit.
//
// for sequence header,
// demux it and save it in the avc_extra_data and aac_extra_data,
//
// for the codec info, such as audio sample rate,
// decode from FLV/RTMP header, then use codec info in sequence
// header to override it.
type avcAacCodec struct {
	// metadata specified
	duration  int
	width     int
	height    int
	frameRate int

	videoCodecID  int
	videoDataRate int // in bps
	audioCodecID  int
	audioDataRate int // in bps

	// video specified
	// profile_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcProfile uint8
	// level_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcLevel uint8
	// lengthSizeMinusOne, H.264-AVC-ISO_IEC_14496-15.pdf, page 16
	nalUnitLength               int8
	sequenceParameterSetLength  uint16
	sequenceParameterSetNALUnit []byte
	pictureParameterSetLength   uint16
	pictureParameterSetNALUnit  []byte

	// audio specified
	// 1.6.2.1 AudioSpecificConfig, in aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 33.
	// audioObjectType, value defines in 7.1 Profiles, aac-iso-13818-7.pdf, page 40.
	aacProfile uint8
	// samplingFrequencyIndex
	aacSampleRate uint8
	// channelConfiguration
	aacChannels uint8

	// the avc extra data, the AVC sequence header,
	// without the flv codec header,
	// @see: ffmpeg, AVCodecContext::extradata
	avcExtraSize int
	avcExtraData []byte
	// the aac extra data, the AAC sequence header,
	// without the flv codec header,
	// @see: ffmpeg, AVCodecContext::extradata
	aacExtraSize int
	aacExtraData []byte
}

func newAvcAacCodec() *avcAacCodec {
	return &avcAacCodec{
		aacSampleRate: hlsAacSampleRateUnset,
	}
}

// demux the audio packet in aac codec.
// the packet mux in FLV/RTMP format defined in flv specification.
// demux the audio speicified data(sound_format, sound_size, ...) to sample.
// demux the aac specified data(aac_profile, ...) to codec from sequence header.
// demux the aac raw to sample units.
func (codec *avcAacCodec) audioAacDemux(data []byte, sample *codecSample) (err error) {

	sample.isVideo = false

	dataLen := len(data)
	if dataLen <= 0 {
		return
	}

	var offset int

	if dataLen-offset < 1 {
		return
	}

	// @see: E.4.2 Audio Tags, video_file_format_spec_v10_1.pdf, page 76
	soundFormat := data[offset]
	offset++

	soundType := soundFormat & 0x01
	soundSize := (soundFormat >> 1) & 0x01
	soundRate := (soundFormat >> 2) & 0x03
	soundFormat = (soundFormat >> 4) & 0x0f

	codec.audioCodecID = int(soundFormat)
	sample.soundType = int(soundType)
	sample.soundRate = int(soundRate)
	sample.soundSize = int(soundSize)

	// only support for aac
	if RtmpCodecAudioAAC != codec.audioCodecID {
		log.Println("hls only support audio aac, actual is ", codec.audioCodecID)
		return
	}

	if dataLen-offset < 1 {
		return
	}

	aacPacketType := data[offset]
	offset++
	sample.aacPacketType = int(aacPacketType)

	if RtmpCodecAudioTypeSequenceHeader == aacPacketType {
		// AudioSpecificConfig
		// 1.6.2.1 AudioSpecificConfig, in aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 33.
		codec.aacExtraSize = dataLen - offset
		if codec.aacExtraSize > 0 {
			codec.aacExtraData = make([]byte, codec.aacExtraSize)
			copy(codec.aacExtraData, data[offset:])
		}

		// only need to decode the first 2bytes:
		// audioObjectType, aac_profile, 5bits.
		// samplingFrequencyIndex, aac_sample_rate, 4bits.
		// channelConfiguration, aac_channels, 4bits
		if dataLen-offset < 2 {
			return
		}

		codec.aacProfile = data[offset]
		offset++
		codec.aacSampleRate = data[offset]
		offset++

		codec.aacChannels = (codec.aacSampleRate >> 3) & 0x0f
		codec.aacSampleRate = ((codec.aacProfile << 1) & 0x0e) | ((codec.aacSampleRate >> 7) & 0x01)
		codec.aacProfile = (codec.aacProfile >> 3) & 0x1f

		if 0 == codec.aacProfile || 0x1f == codec.aacProfile {
			err = fmt.Errorf("hls decdoe audio aac sequence header failed, aac profile=%d", codec.aacProfile)
			return
		}

		// the profile = object_id + 1
		// @see aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 78,
		//      Table 1. A.9 MPEG-2 Audio profiles and MPEG-4 Audio object types
		// so the aac_profile should plus 1, not minus 1, and nginx-rtmp used it to
		// downcast aac SSR to LC.
		codec.aacProfile--

	} else if RtmpCodecAudioTypeRawData == aacPacketType {
		// ensure the sequence header demuxed
		if 0 == len(codec.aacExtraData) {
			return
		}

		// Raw AAC frame data in UI8 []
		// 6.3 Raw Data, aac-iso-13818-7.pdf, page 28
		if err = sample.addSampleUnit(data[offset:]); err != nil {
			return
		}
	}

	// reset the sample rate by sequence header
	if codec.aacSampleRate != hlsAacSampleRateUnset {
		var aacSampleRates = []int{
			96000, 88200, 64000, 48000,
			44100, 32000, 24000, 22050,
			16000, 12000, 11025, 8000,
			7350, 0, 0, 0,
		}

		switch aacSampleRates[codec.aacSampleRate] {
		case 11025:
			sample.soundRate = RtmpCodecAudioSampleRate11025
		case 22050:
			sample.soundRate = RtmpCodecAudioSampleRate22050
		case 44100:
			sample.soundRate = RtmpCodecAudioSampleRate44100
		default:
		}
	}

	return
}

// demux the video packet in h.264 codec.
// the packet mux in FLV/RTMP format defined in flv specification.
// demux the video specified data(frame_type, codec_id, ...) to sample.
// demux the h.264 sepcified data(avc_profile, ...) to codec from sequence header.
// demux the h.264 NALUs to sampe units.
func (codec *avcAacCodec) videoAvcDemux(data []byte, sample *codecSample) (err error) {

	sample.isVideo = true

	maxLen := len(data)

	if maxLen <= 0 {
		return
	}

	var offset int

	// video decode
	if maxLen-offset < 1 {
		return
	}

	// @see: E.4.3 Video Tags, video_file_format_spec_v10_1.pdf, page 78
	frameType := data[offset]
	offset++

	codecID := frameType & 0x0f
	frameType = (frameType >> 4) & 0x0f

	sample.frameType = int(frameType)

	// ignore info frame without error
	if RtmpCodecVideoAVCFrameVideoInfoFrame == sample.frameType {
		return
	}

	// only support h.264/avc
	if RtmpCodecVideoAVC != codecID {
		return
	}
	codec.videoCodecID = int(codecID)

	if maxLen-offset < 4 {
		return
	}

	avcPacketType := data[offset]
	offset++

	// 3bytes,
	compositionTime := int(data[offset])<<16 + int(data[offset+1])<<8 + int(data[offset+2])
	offset += 3

	// pts = dts + cts
	sample.cts = compositionTime
	sample.avcPacketType = int(avcPacketType)

	if RtmpCodecVideoAVCTypeSequenceHeader == avcPacketType {
		// AVCDecoderConfigurationRecord
		// 5.2.4.1.1 Syntax, H.264-AVC-ISO_IEC_14496-15.pdf, page 16
		codec.avcExtraSize = maxLen - offset
		if codec.avcExtraSize > 0 {
			codec.avcExtraData = make([]byte, codec.avcExtraSize)
			copy(codec.avcExtraData, data[offset:offset+codec.avcExtraSize])
		}

		if maxLen-offset < 6 {
			return
		}

		// configurationVersion
		offset++

		// AVCProfileIndication
		codec.avcProfile = data[offset]
		offset++

		// profile_compatibility
		offset++

		// AVCLevelIndication
		codec.avcLevel = data[offset]
		offset++

		// parse the NALU size.
		lengthSizeMinusOne := data[offset]
		offset++

		lengthSizeMinusOne &= 0x03
		codec.nalUnitLength = int8(lengthSizeMinusOne)

		// 1 sps
		if maxLen-offset < 1 {
			return
		}

		numOfSequenceParameterSets := data[offset]
		offset++
		numOfSequenceParameterSets &= 0x1f
		if numOfSequenceParameterSets != 1 {
			err = fmt.Errorf("hsl decode avc sequence header size failed")
			return
		}

		if maxLen-offset < 2 {
			return
		}

		codec.sequenceParameterSetLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		if maxLen-offset < int(codec.sequenceParameterSetLength) {
			err = fmt.Errorf("hsl decode avc sequence header data failed")
			return
		}

		if codec.sequenceParameterSetLength > 0 {
			codec.sequenceParameterSetNALUnit = make([]byte, codec.sequenceParameterSetLength)
			copy(codec.sequenceParameterSetNALUnit, data[offset:offset+int(codec.sequenceParameterSetLength)])
			offset += int(codec.sequenceParameterSetLength)
		}

		// 1 pps
		if maxLen-offset < 1 {
			return
		}

		numOfPictureParameterSets := data[offset]
		offset++

		if numOfPictureParameterSets != 1 {
			err = fmt.Errorf("hls decode video avc sequence header pps failed")
			return
		}

		if maxLen-offset < 2 {
			return
		}

		codec.pictureParameterSetLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		if maxLen-offset < int(codec.pictureParameterSetLength) {
			return
		}

		if codec.pictureParameterSetLength > 0 {
			codec.pictureParameterSetNALUnit = make([]byte, codec.pictureParameterSetLength)
			copy(codec.pictureParameterSetNALUnit, data[offset:offset+int(codec.pictureParameterSetLength)])
			offset += int(codec.pictureParameterSetLength)
		}

	} else if RtmpCodecVideoAVCTypeNALU == avcPacketType {
		// ensure the sequence header demuxed
		if len(codec.pictureParameterSetNALUnit) <= 0 {
			return
		}

		// One or more NALUs (Full frames are required)
		// 5.3.4.2.1 Syntax, H.264-AVC-ISO_IEC_14496-15.pdf, page 20

		pictureLength := maxLen - offset
		for i := 0; i < pictureLength; {
			if maxLen-offset < int(codec.nalUnitLength)+1 {
				return
			}

			nalUnitLength := 0
			switch codec.nalUnitLength {
			case 3:
				nalUnitLength = int(binary.BigEndian.Uint32(data[offset : offset+4]))
				offset += 4
			case 2:
				nalUnitLength = int(data[offset])<<16 + int(data[offset+1])<<8 + int(data[offset+2])
				offset += 3
			case 1:
				nalUnitLength = int(binary.BigEndian.Uint16(data[offset : offset+2]))
				offset += 2
			default:
				nalUnitLength = int(data[offset])
				offset++
			}

			// maybe stream is AnnexB format.
			if nalUnitLength < 0 {
				return
			}

			// NALUnit
			if maxLen-offset < nalUnitLength {
				return
			}

			// 7.3.1 NAL unit syntax, H.264-AVC-ISO_IEC_14496-10.pdf, page 44.
			if err = sample.addSampleUnit(data[offset : offset+nalUnitLength]); err != nil {
				err = fmt.Errorf("hls add video sample failed")
				return
			}
			offset += nalUnitLength

			i += int(codec.nalUnitLength) + 1 + nalUnitLength
		}

	}

	return
}
