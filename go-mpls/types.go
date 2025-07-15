package go_mpls

type PlaybackType int
type StreamCodingType int
type VideoFormat int
type FrameRate int
type DynamicRangeType int
type ColorSpace int
type AudioFormat int
type SampleRate int
type CharacterCode int
type SubPathType int

const (
	StandardPlay PlaybackType = iota + 1
	RandomPlay
	ShufflePlay
)

const (
	MPEG1Video               StreamCodingType = 0x01
	MPEG2Video               StreamCodingType = 0x02
	MPEG4AVCVideo            StreamCodingType = 0x1b
	MPEG4MVCVideo            StreamCodingType = 0x20
	SMTPEVC1Video            StreamCodingType = 0xea
	HEVCVideo                StreamCodingType = 0x24
	MPEG1Audio               StreamCodingType = 0x03
	MPEG2Audio               StreamCodingType = 0x04
	LPCMAudio                StreamCodingType = 0x80
	DolbyDigitalAudio        StreamCodingType = 0x81
	DTSAudio                 StreamCodingType = 0x82
	DolbyDigitalTureHDAudio  StreamCodingType = 0x83
	DolbyDigitalPlusAudioPri StreamCodingType = 0x84
	DTSHDHighResolutionAudio StreamCodingType = 0x85
	DTSHDMasterAudio         StreamCodingType = 0x86
	DolbyDigitalPlusAudioSec StreamCodingType = 0xa1
	DTSHDAudio               StreamCodingType = 0xa2
	PresentationGraphics     StreamCodingType = 0x90
	InteractiveGraphics      StreamCodingType = 0x91
	TextSubtitle             StreamCodingType = 0x92
)

const (
	VF480I VideoFormat = iota + 1
	VF576I
	VF480P
	VF1080I
	VF720P
	VF1080P
	VF576P
	VF2160P
)

const (
	FR23D98FPS FrameRate = iota + 1
	FR24FPS
	FR25FPS
	FR29D97FPS
	FR50FPS
	FR59D94FPS
)

const (
	SDR DynamicRangeType = iota
	HDR10
	DolbyVision
)

const (
	Reserved ColorSpace = iota
	BT709
	BT2020
)

const (
	Mono                  AudioFormat = 0x01
	Stereo                            = 0x03
	MultiChannel                      = 0x06
	StereoAndMultiChannel             = 0x0c
)

const (
	SR48KHz       SampleRate = 0x01
	SR96KHz                  = 0x04
	SR192KHz                 = 0x05
	SR48And192KHz            = 0x0c
	SR48And96KHz             = 0x0e
)

const (
	UTF8 CharacterCode = iota + 0x01
	UTF16BE
	ShiftJIS
	KSC5601
	GB18030
	GB2312
	BIG5
)

const (
	PrimaryAudio SubPathType = iota + 0x02
	InteractiveGraphicsMenu
	TextSubtitlePath
	OutMuxAndSyncTypeOfStreams
	OutMuxAndAsyncTypeOfPIP
	InMuxAndSyncTypeOfPIP
	StereoscopicVideo
	StereoscopicIGMenu
	DolbyVisionEnhancement
)

type MPLS struct {
	FilePath                  string
	RawData                   []byte
	VersionNumber             int
	PlaylistStartAddress      int
	PlaylistMarkStartAddress  int
	ExtensionDataStartAddress int
	ApplicationInfoPlaylist   *AppInfoPlayList
	PlayList                  *PlayList
	PlayListMark              *PlayListMark
	ExtensionData             *ExtensionData
}

type ExtDataEntryItem struct {
	ExtDataType         int
	ExtDataVersion      int
	ExtDataStartAddress int
	ExtDataLength       int
	ExtDataEntry        []byte
}

type ExtensionData struct {
	Length                 int
	DataBlockStartAddress  int
	NumberOfExtDataEntries int
	ExtDataEntryItemsList  []*ExtDataEntryItem
}

type PlayListMark struct {
	Length                int
	NumberOfPlayListMarks int
	PlayListMarksList     []*PlayListMarkItem
}

type PlayListMarkItem struct {
	MarkType        int
	RefToPlayItemID int
	MarkTimeStamp   float32
	EntryESPID      int
	Duration        int
}

type AppInfoPlayList struct {
	Length                        int
	PlaybackType                  PlaybackType
	PlaybackCount                 int
	UOMaskTable                   *UOMaskTable
	RandomAccessFlag              bool
	AudioMixFlag                  bool
	LosslessBypassFlag            bool
	MVCBaseViewRFlag              bool
	SDRConversionNotificationFlag bool
}

type UOMaskTable struct {
	MenuCall                         bool
	TitleSearch                      bool
	ChapterSearch                    bool
	TimeSearch                       bool
	SkipToNextPoint                  bool
	SkipToPrevPoint                  bool
	Stop                             bool
	PauseOn                          bool
	StillOff                         bool
	ForwardPlay                      bool
	BackwardPlay                     bool
	Resume                           bool
	MoveUpSelectedButton             bool
	MoveDownSelectedButton           bool
	MoveLeftSelectedButton           bool
	MoveRightSelectedButton          bool
	SelectButton                     bool
	ActivateButton                   bool
	SelectAndActivateButton          bool
	PrimaryAudioStreamNumberChange   bool
	AngleNumberChange                bool
	PopupOn                          bool
	PopupOff                         bool
	PrimaryPGEnableDisable           bool
	PrimaryPGStreamNumberChange      bool
	SecondaryVideoEnableDisable      bool
	SecondaryVideoStreamNumberChange bool
	SecondaryAudioEnableDisable      bool
	SecondaryAudioStreamNumberChange bool
	SecondaryPGStreamNumberChange    bool
}

type Angle struct {
	ClipInformationFileName string
	ClipCodecIdentifier     string
	RefToSTCID              int
}

type StreamEntry struct {
	Length         int
	StreamType     int
	RefToSubPathID int
	RefToSubClipID int
	RefToStreamPID int
}

type StreamAttributes struct {
	Length           int
	StreamCodingType StreamCodingType
	VideoFormat      VideoFormat
	FrameRate        FrameRate
	DynamicRangeType DynamicRangeType
	ColorSpace       ColorSpace
	CRFlag           bool
	HDRPlusFlag      bool
	AudioFormat      AudioFormat
	SampleRate       SampleRate
	LanguageCode     string
	CharacterCode    CharacterCode
}

type Stream struct {
	StreamEntry      *StreamEntry
	StreamAttributes *StreamAttributes
}

type STNTable struct {
	Length                        int
	NumberOfPrimaryVideoStreams   int
	NumberOfPrimaryAudioStreams   int
	NumberOfPrimaryPGStreams      int
	NumberOfPrimaryIGStreams      int
	NumberOfSecondaryAudioStreams int
	NumberOfSecondaryVideoStreams int
	NumberOfSecondaryPGStreams    int
	NumberOfDVStreams             int
	PrimaryVideoStreamsList       []*Stream
	PrimaryAudioStreamsList       []*Stream
	PrimaryPGStreamsList          []*Stream
	SecondaryPGStreamsList        []*Stream
	PrimaryIGStreamsList          []*Stream
	SecondaryAudioStreamsList     []*Stream
	SecondaryVideoStreamsList     []*Stream
	DVStreamsList                 []*Stream
}

type PlayList struct {
	Length            int
	NumberOfPlayItems int
	NumberOfSubPaths  int
	PlayItemList      []*PlayItem
	SubPathsList      []*SubPath
	RefToSubPathID    int
	RefToSubClipID    int
	RefToStreamPID    int
}

type PlayItem struct {
	Length                   int
	ClipInformationFileName  string
	ClipCodecIdentifier      string
	IsMultiAngle             bool
	ConnectionCondition      int
	RefToSTCID               int
	INTime                   float32
	OUTTime                  float32
	UserOperationMaskTable   *UOMaskTable
	PlayItemRandomAccessFlag bool
	StillMode                int
	StillTime                float32
	NumberOfAngles           int
	IsDifferentAudios        bool
	IsSeamlessAngleChange    bool
	AnglesList               []*Angle
	STNTable                 *STNTable
}

type MultiClipEntry struct {
	ClipInformationFileName string
	ClipCodecIdentifier     string
	RefToSTCID              int
}

type SubPlayItem struct {
	Length                   int
	ClipInformationFileName  string
	ClipCodecIdentifier      string
	ConnectionCondition      int
	IsMultiClipEntries       bool
	RefToSTCID               int
	INTime                   float32
	OUTTime                  float32
	SyncPlayItemID           int
	SyncStartPTS             int
	NumberOfMultiClipEntries int
	MultiClipEntriesList     []*MultiClipEntry
}

type SubPath struct {
	Length               int
	SubPathType          SubPathType
	IsRepeatSubPath      bool
	NumberOfSubPlayItems int
	SubPlayItemsList     []*SubPlayItem
}
