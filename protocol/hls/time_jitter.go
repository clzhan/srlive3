package hls

import (
	"github.com/clzhan/srlive3/avformat"
	"github.com/clzhan/srlive3/rtmpconst"
)

// TimeJitter time jitter detect and correct to make sure the rtmp stream is monotonically
type TimeJitter struct {
	lastPktTime        int64
	lastPktCorrectTime int64
}

// NewTimeJitter create a new time jitter
func NewTimeJitter() *TimeJitter {
	return &TimeJitter{}
}

// Correct  detect the time jitter and correct it.
// tba, the audio timebase, used to calc the "right" delta if jitter detected.
// tbv, the video timebase, used to calc the "right" delta if jitter detected.
// start_at_zero whether ensure stream start at zero.
// mono_increasing whether ensure stream is monotonically inscreasing.
func (tj *TimeJitter) Correct(msg *avformat.MediaFrame, tba float64, tbv float64, timeJitter uint32) {

	if nil == msg {
		return
	}

	if RtmpTimeJitterFull != timeJitter {
		// all jitter correct features is disabled, ignore.
		if RtmpTimeJitterOff == timeJitter {
			return
		}

		// start at zero, but donot ensure monotonically increasing.
		if RtmpTimeJitterZero == timeJitter {
			// for the first time, last_pkt_correct_time is zero.
			// while when timestamp overflow, the timestamp become smaller,
			// reset the last_pkt_correct_time.
			if tj.lastPktCorrectTime <= 0 || tj.lastPktCorrectTime > int64(msg.Timestamp) {
				tj.lastPktCorrectTime = int64(msg.Timestamp)
			}

			msg.Timestamp -= uint32(tj.lastPktCorrectTime)

			return
		}
	}

	// full jitter algorithm, do jitter correct.
	// set to 0 for metadata.
	//if !msg.Header.IsAudio() && !msg.Header.IsVideo() {
	//	msg.Header.Timestamp = 0
	//	return
	//}

	sampleRate := tba
	frameRate := tbv

	/**
	 * we use a very simple time jitter detect/correct algorithm:
	 * 1. delta: ensure the delta is positive and valid,
	 *     we set the delta to DefaultFrameTimeMs,
	 *     if the delta of time is nagative or greater than MaxJitterMs.
	 * 2. last_pkt_time: specifies the original packet time,
	 *     is used to detect next jitter.
	 * 3. last_pkt_correct_time: simply add the positive delta,
	 *     and enforce the time monotonically.
	 */
	timeLocal := msg.Timestamp
	delta := int64(timeLocal) - tj.lastPktTime

	// if jitter detected, reset the delta.
	if delta < 0 || delta > MaxJitterMs {
		// calc the right diff by audio sample rate
		if msg.Type == rtmpconst.RTMP_MSG_AUDIO && sampleRate > 0 {
			delta = (int64)(float64(delta) * 1000.0 / sampleRate)
		} else if msg.Type ==  rtmpconst.RTMP_MSG_VIDEO && frameRate > 0 {
			delta = (int64)(float64(delta) * 1.0 / frameRate)
		} else {
			delta = DefaultFrameTimeMs
		}
	}

	// sometimes, the time is absolute time, so correct it again.
	if delta < 0 || delta > MaxJitterMs {
		delta = DefaultFrameTimeMs
	}

	if tj.lastPktCorrectTime+delta > 0 {
		tj.lastPktCorrectTime = tj.lastPktCorrectTime + delta
	} else {
		tj.lastPktCorrectTime = 0
	}

	msg.Timestamp = uint32(tj.lastPktCorrectTime)
	tj.lastPktTime = int64(timeLocal)

}
