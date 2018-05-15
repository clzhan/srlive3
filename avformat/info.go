package avformat

var (
	Samplerate = map[byte]string{
		0: "5.5kHz",
		1: "11kHz",
		2: "22kHz",
		3: "44kHz"}
	Samplelength = map[byte]string{
		0: "8Bit",
		1: "16Bit"}
	Audiotype = map[byte]string{
		0: "Mono",
		1: "Stereo"}


	Videoframetype = map[byte]string{
		1: "keyframe (for AVC, a seekable frame)",
		2: "inter frame (for AVC, a non-seekable frame)",
		3: "disposable inter frame (H.263 only)",
		4: "generated keyframe (reserved for server use only)",
		5: "video info/command frame"}
	Videocodec = map[byte]string{
		1: "JPEG (currently unused)",
		2: "Sorenson H.263",
		3: "Screen video",
		4: "On2 VP6",
		5: "On2 VP6 with alpha channel",
		6: "Screen video version 2",
		7: "AVC"}
	Audioformat = map[byte]string{
		0:  "Linear PCM, platform endian",
		1:  "ADPCM",
		2:  "MP3",
		3:  "Linear PCM, little endian",
		4:  "Nellymoser 16kHz mono",
		5:  "Nellymoser 8kHz mono",
		6:  "Nellymoser",
		7:  "G.711 A-law logarithmic PCM",
		8:  "G.711 mu-law logarithmic PCM",
		9:  "reserved",
		10: "AAC",
		11: "Speex",
		14: "MP3 8Khz",
		15: "Device-specific sound"}
)
