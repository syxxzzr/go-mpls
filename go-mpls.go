package go_mpls

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"strconv"
)

func parseUOMaskTable(rawData []byte) *UOMaskTable {
	return &UOMaskTable{
		MenuCall:                         (rawData[0] & (1 << 7)) != 0,
		TitleSearch:                      (rawData[0] & (1 << 6)) != 0,
		ChapterSearch:                    (rawData[0] & (1 << 5)) != 0,
		TimeSearch:                       (rawData[0] & (1 << 4)) != 0,
		SkipToNextPoint:                  (rawData[0] & (1 << 3)) != 0,
		SkipToPrevPoint:                  (rawData[0] & (1 << 2)) != 0,
		Stop:                             (rawData[0] & (1 << 0)) != 0,
		PauseOn:                          (rawData[1] & (1 << 7)) != 0,
		StillOff:                         (rawData[1] & (1 << 5)) != 0,
		ForwardPlay:                      (rawData[1] & (1 << 4)) != 0,
		BackwardPlay:                     (rawData[1] & (1 << 3)) != 0,
		Resume:                           (rawData[1] & (1 << 2)) != 0,
		MoveUpSelectedButton:             (rawData[1] & (1 << 1)) != 0,
		MoveDownSelectedButton:           (rawData[1] & (1 << 0)) != 0,
		MoveLeftSelectedButton:           (rawData[2] & (1 << 7)) != 0,
		MoveRightSelectedButton:          (rawData[2] & (1 << 6)) != 0,
		SelectButton:                     (rawData[2] & (1 << 5)) != 0,
		ActivateButton:                   (rawData[2] & (1 << 4)) != 0,
		SelectAndActivateButton:          (rawData[2] & (1 << 3)) != 0,
		PrimaryAudioStreamNumberChange:   (rawData[2] & (1 << 2)) != 0,
		AngleNumberChange:                (rawData[2] & (1 << 0)) != 0,
		PopupOn:                          (rawData[3] & (1 << 7)) != 0,
		PopupOff:                         (rawData[3] & (1 << 6)) != 0,
		PrimaryPGEnableDisable:           (rawData[3] & (1 << 5)) != 0,
		PrimaryPGStreamNumberChange:      (rawData[3] & (1 << 4)) != 0,
		SecondaryVideoEnableDisable:      (rawData[3] & (1 << 3)) != 0,
		SecondaryVideoStreamNumberChange: (rawData[3] & (1 << 2)) != 0,
		SecondaryAudioEnableDisable:      (rawData[3] & (1 << 1)) != 0,
		SecondaryAudioStreamNumberChange: (rawData[3] & (1 << 0)) != 0,
		SecondaryPGStreamNumberChange:    (rawData[4] & (1 << 6)) != 0,
	}
}

func parseStreamEntry(rawData []byte) *StreamEntry {
	length := int(rawData[0])
	streamType := int(rawData[1])

	refToSubPathID := 0
	refToSubClipID := 0
	refToStreamPID := 0
	if streamType == 0x01 {
		refToStreamPID = int(binary.BigEndian.Uint16(rawData[2:4]))
	} else if streamType == 0x02 {
		refToSubPathID = int(rawData[2])
		refToSubClipID = int(rawData[3])
		refToStreamPID = int(binary.BigEndian.Uint16(rawData[4:6]))
	} else if streamType == 0x03 || streamType == 0x04 {
		refToSubPathID = int(rawData[2])
		refToStreamPID = int(binary.BigEndian.Uint16(rawData[3:5]))
	}

	return &StreamEntry{
		Length:         length,
		StreamType:     streamType,
		RefToSubPathID: refToSubPathID,
		RefToSubClipID: refToSubClipID,
		RefToStreamPID: refToStreamPID,
	}
}

func parseStreamAttributes(rawData []byte) *StreamAttributes {
	length := int(rawData[0])
	streamCodingType := StreamCodingType(rawData[1])

	videoFormat := VideoFormat(0)
	frameRate := FrameRate(0)
	dynamicRangeType := DynamicRangeType(0)
	colorSpace := ColorSpace(0)
	crFlag := false
	hdrPlusFlag := false
	audioFormat := AudioFormat(0)
	sampleRate := SampleRate(0)
	languageCode := ""
	characterCode := CharacterCode(0)

	if streamCodingType == 0x24 {
		videoFormat = VideoFormat((rawData[2] & 0b11110000) >> 4)
		frameRate = FrameRate(rawData[2] & 0b00001111)
		dynamicRangeType = DynamicRangeType((rawData[3] & 0b11110000) >> 4)
		colorSpace = ColorSpace(rawData[3] & 0b00001111)
		crFlag = (rawData[4] & (1 << 7)) != 0
		hdrPlusFlag = (rawData[4] & (1 << 6)) != 0
	} else if streamCodingType == 0x92 {
		characterCode = CharacterCode(rawData[2])
		languageCode = string(rawData[3:6])
	} else if streamCodingType == 0x90 || streamCodingType == 0x91 {
		languageCode = string(rawData[2:5])
	} else if streamCodingType == 0x01 || streamCodingType == 0x02 || streamCodingType == 0x1b || streamCodingType == 0xea {
		videoFormat = VideoFormat((rawData[2] & 0b11110000) >> 4)
		frameRate = FrameRate(rawData[2] & 0b00001111)
	} else {
		audioFormat = AudioFormat((rawData[2] & 0b11110000) >> 4)
		sampleRate = SampleRate(rawData[2] & 0b00001111)
		languageCode = string(rawData[3:6])
	}

	return &StreamAttributes{
		Length:           length,
		StreamCodingType: streamCodingType,
		VideoFormat:      videoFormat,
		FrameRate:        frameRate,
		DynamicRangeType: dynamicRangeType,
		ColorSpace:       colorSpace,
		CRFlag:           crFlag,
		HDRPlusFlag:      hdrPlusFlag,
		AudioFormat:      audioFormat,
		SampleRate:       sampleRate,
		LanguageCode:     languageCode,
		CharacterCode:    characterCode,
	}
}

//func parseStream(rawData []byte) *Stream {
//	streamEntry := parseStreamEntry(rawData[0:])
//	streamAttributes := parseStreamAttributes(rawData[streamEntry.Length+1:])
//
//	return &Stream{
//		StreamEntry:      streamEntry,
//		StreamAttributes: streamAttributes,
//	}
//}

func parseStreamsList(rawData []byte, number int) ([]*Stream, int) {
	if number == 0 {
		return nil, 0
	}

	offset := 0
	var streamsList []*Stream
	for i := 0; i < number; i++ {
		streamEntry := parseStreamEntry(rawData[offset:])
		offset += streamEntry.Length + 1
		streamAttributes := parseStreamAttributes(rawData[offset:])
		offset += streamAttributes.Length + 1

		streamsList = append(streamsList, &Stream{
			StreamEntry:      streamEntry,
			StreamAttributes: streamAttributes,
		})
	}

	return streamsList, offset
}

func parseSTNTable(rawData []byte) *STNTable {
	length := int(binary.BigEndian.Uint16(rawData[:2]))
	numberOfPrimaryVideoStreams := int(rawData[4])
	numberOfPrimaryAudioStreams := int(rawData[5])
	numberOfPrimaryPGStreams := int(rawData[6])
	numberOfPrimaryIGStreams := int(rawData[7])
	numberOfSecondaryAudioStreams := int(rawData[8])
	numberOfSecondaryVideoStreams := int(rawData[9])
	numberOfSecondaryPGStreams := int(rawData[10])
	numberOfDVStreams := int(rawData[11])

	offset := 16
	primaryVideoStreamsList, o := parseStreamsList(rawData[offset:], numberOfPrimaryVideoStreams)
	offset += o
	primaryAudioStreamsList, o := parseStreamsList(rawData[offset:], numberOfPrimaryAudioStreams)
	offset += o
	primaryPGStreamsList, o := parseStreamsList(rawData[offset:], numberOfPrimaryPGStreams)
	offset += o
	secondaryPGStreamsList, o := parseStreamsList(rawData[offset:], numberOfSecondaryPGStreams)
	offset += o
	primaryIGStreamsList, o := parseStreamsList(rawData[offset:], numberOfPrimaryIGStreams)
	offset += o
	secondaryAudioStreamsList, o := parseStreamsList(rawData[offset:], numberOfSecondaryAudioStreams)
	offset += o
	secondaryVideoStreamsList, o := parseStreamsList(rawData[offset:], numberOfSecondaryVideoStreams)
	offset += o
	dvStreamsList, o := parseStreamsList(rawData[offset:], numberOfDVStreams)
	offset += o

	return &STNTable{
		Length:                        length,
		NumberOfPrimaryVideoStreams:   numberOfPrimaryVideoStreams,
		NumberOfPrimaryAudioStreams:   numberOfPrimaryAudioStreams,
		NumberOfPrimaryPGStreams:      numberOfPrimaryPGStreams,
		NumberOfPrimaryIGStreams:      numberOfPrimaryIGStreams,
		NumberOfSecondaryAudioStreams: numberOfSecondaryAudioStreams,
		NumberOfSecondaryVideoStreams: numberOfSecondaryVideoStreams,
		NumberOfSecondaryPGStreams:    numberOfSecondaryPGStreams,
		NumberOfDVStreams:             numberOfDVStreams,
		PrimaryVideoStreamsList:       primaryVideoStreamsList,
		PrimaryAudioStreamsList:       primaryAudioStreamsList,
		PrimaryPGStreamsList:          primaryPGStreamsList,
		SecondaryPGStreamsList:        secondaryPGStreamsList,
		PrimaryIGStreamsList:          primaryIGStreamsList,
		SecondaryAudioStreamsList:     secondaryAudioStreamsList,
		SecondaryVideoStreamsList:     secondaryVideoStreamsList,
		DVStreamsList:                 dvStreamsList,
	}
}

func parsePlayItem(rawData []byte) *PlayItem {
	length := int(binary.BigEndian.Uint16(rawData[:2]))
	clipInfoFileName := string(rawData[2:7])
	clipCodecIdentifier := string(rawData[7:11])
	isMultiAngle := (rawData[12] & (1 << 4)) != 0
	connectionCondition := int(rawData[12] & 0b00001111)
	refToSTCID := int(rawData[13])
	inTime := float32(binary.BigEndian.Uint32(rawData[14:18])) / 45000
	outTime := float32(binary.BigEndian.Uint32(rawData[18:22])) / 45000
	userOperationMaskTable := parseUOMaskTable(rawData[22:])
	playItemRandomAccessFlag := (rawData[30] & (1 << 7)) != 0
	stillMode := int(rawData[31])

	stillTime := float32(0)
	if stillMode == 0x01 {
		stillTime = float32(binary.BigEndian.Uint16(rawData[32:34])) / 45000
	}

	stnTableStart := 34
	numberOfAngles := 0
	isDifferentAudios := false
	isSeamlessAngleChange := false
	var angleList []*Angle = nil
	if isMultiAngle {
		numberOfAngles = int(rawData[34])
		isDifferentAudios = (rawData[35] & (1 << 1)) != 0
		isSeamlessAngleChange = (rawData[35] & (1 << 0)) != 0
		for i := 0; i < numberOfAngles; i++ {
			angleList = append(angleList, &Angle{
				ClipInformationFileName: string(rawData[36+10*i : 41+10*i]),
				ClipCodecIdentifier:     string(rawData[41+10*i : 45+10*i]),
				RefToSTCID:              int(rawData[45+10*i]),
			})
		}
		stnTableStart = 34 + 10*numberOfAngles
	}
	stnTable := parseSTNTable(rawData[stnTableStart:])

	return &PlayItem{
		Length:                   length,
		ClipInformationFileName:  clipInfoFileName,
		ClipCodecIdentifier:      clipCodecIdentifier,
		IsMultiAngle:             isMultiAngle,
		ConnectionCondition:      connectionCondition,
		RefToSTCID:               refToSTCID,
		INTime:                   inTime,
		OUTTime:                  outTime,
		UserOperationMaskTable:   userOperationMaskTable,
		PlayItemRandomAccessFlag: playItemRandomAccessFlag,
		StillMode:                stillMode,
		StillTime:                stillTime,
		NumberOfAngles:           numberOfAngles,
		IsDifferentAudios:        isDifferentAudios,
		IsSeamlessAngleChange:    isSeamlessAngleChange,
		AnglesList:               angleList,
		STNTable:                 stnTable,
	}
}

func parseSubPlayItem(rawData []byte) *SubPlayItem {
	length := int(binary.BigEndian.Uint16(rawData[:2]))
	clipInformationFileName := string(rawData[2:7])
	clipCodecIdentifier := string(rawData[7:11])
	connectionCondition := int((rawData[14] & 0b00011110) >> 1)
	isMultiClipEntries := (rawData[14] & (1 << 0)) != 0
	refToSTCID := int(rawData[15])
	inTime := float32(binary.BigEndian.Uint32(rawData[16:20])) / 45000
	outTime := float32(binary.BigEndian.Uint32(rawData[20:24])) / 45000
	syncPlayItemID := int(binary.BigEndian.Uint16(rawData[24:26]))
	syncStartPTS := int(binary.BigEndian.Uint32(rawData[26:30]))

	numberOfMultiClipEntries := 0
	var multiClipEntriesList []*MultiClipEntry = nil
	if isMultiClipEntries {
		numberOfMultiClipEntries = int(rawData[31])
		for i := 0; i < numberOfMultiClipEntries; i++ {
			multiClipEntriesList = append(multiClipEntriesList, &MultiClipEntry{
				ClipInformationFileName: string(rawData[31+10*i : 36+10*i]),
				ClipCodecIdentifier:     string(rawData[36+10*i : 40+10*i]),
				RefToSTCID:              int(rawData[40+10*i]),
			})
		}
	}

	return &SubPlayItem{
		Length:                   length,
		ClipInformationFileName:  clipInformationFileName,
		ClipCodecIdentifier:      clipCodecIdentifier,
		ConnectionCondition:      connectionCondition,
		IsMultiClipEntries:       isMultiClipEntries,
		RefToSTCID:               refToSTCID,
		INTime:                   inTime,
		OUTTime:                  outTime,
		SyncPlayItemID:           syncPlayItemID,
		SyncStartPTS:             syncStartPTS,
		NumberOfMultiClipEntries: numberOfMultiClipEntries,
		MultiClipEntriesList:     multiClipEntriesList,
	}
}

func parseSubPath(rawData []byte) *SubPath {
	length := int(binary.BigEndian.Uint32(rawData[:4]))
	subPathType := SubPathType(int(rawData[5]))
	isRepeatSubPath := (rawData[7] & (1 << 0)) != 0
	numberOfSubPathItems := int(rawData[9])

	offset := 10
	var subPlayItemsList []*SubPlayItem = nil
	for i := 0; i < numberOfSubPathItems; i++ {
		subPlayItem := parseSubPlayItem(rawData[offset:])
		subPlayItemsList = append(subPlayItemsList, subPlayItem)
		offset += subPlayItem.Length + 2
	}

	return &SubPath{
		Length:               length,
		SubPathType:          subPathType,
		IsRepeatSubPath:      isRepeatSubPath,
		NumberOfSubPlayItems: numberOfSubPathItems,
		SubPlayItemsList:     subPlayItemsList,
	}
}

func parseAppInfoPlayList(rawData []byte, channel chan *AppInfoPlayList) {
	length := int(binary.BigEndian.Uint32(rawData[:4]))
	playbackType := PlaybackType(rawData[5])

	var playbackCount int
	if playbackType == 1 {
		playbackCount = 0
	} else {
		playbackCount = int(binary.BigEndian.Uint16(rawData[6:8]))
	}
	userOperationMaskTable := parseUOMaskTable(rawData[8:])

	channel <- &AppInfoPlayList{
		Length:                        length,
		PlaybackType:                  playbackType,
		PlaybackCount:                 playbackCount,
		UOMaskTable:                   userOperationMaskTable,
		RandomAccessFlag:              (rawData[16] & (1 << 7)) != 0,
		AudioMixFlag:                  (rawData[16] & (1 << 6)) != 0,
		LosslessBypassFlag:            (rawData[16] & (1 << 5)) != 0,
		MVCBaseViewRFlag:              (rawData[16] & (1 << 4)) != 0,
		SDRConversionNotificationFlag: (rawData[16] & (1 << 3)) != 0,
	}
}

func parsePlayList(rawData []byte, channel chan *PlayList) {
	length := int(binary.BigEndian.Uint32(rawData[:4]))
	numberOfPlayItems := int(binary.BigEndian.Uint16(rawData[6:8]))
	numberOfSubPaths := int(binary.BigEndian.Uint16(rawData[8:10]))

	var playItemList []*PlayItem
	playItemStart := 10
	for i := 0; i < numberOfPlayItems; i++ {
		playItem := parsePlayItem(rawData[playItemStart:])
		playItemList = append(playItemList, playItem)
		playItemStart += playItem.Length + 2
	}

	var subPathsList []*SubPath
	subPathStart := playItemStart
	for i := 0; i < numberOfSubPaths; i++ {
		subPath := parseSubPath(rawData[subPathStart:])
		subPathsList = append(subPathsList, subPath)
		subPathStart += subPath.Length
	}

	channel <- &PlayList{
		Length:            length,
		NumberOfPlayItems: numberOfPlayItems,
		NumberOfSubPaths:  numberOfSubPaths,
		PlayItemList:      playItemList,
		SubPathsList:      subPathsList,
	}
}

func parsePlayListMark(rawData []byte, channel chan *PlayListMark) {
	length := int(binary.BigEndian.Uint32(rawData[:4]))
	numberOfPlayListMarks := int(binary.BigEndian.Uint16(rawData[4:6]))

	var playListMarksList []*PlayListMarkItem = nil
	for i := 0; i < numberOfPlayListMarks; i++ {
		playListMarksList = append(playListMarksList, &PlayListMarkItem{
			MarkType:        int(rawData[7+14*i]),
			RefToPlayItemID: int(binary.BigEndian.Uint16(rawData[8+14*i : 10+14*i])),
			MarkTimeStamp:   float32(binary.BigEndian.Uint32(rawData[10+14*i:14+14*i])) / 45000,
			EntryESPID:      int(binary.BigEndian.Uint16(rawData[14+14*i : 16+14*i])),
			Duration:        int(binary.BigEndian.Uint32(rawData[16+14*i : 20+14*i])),
		})
	}

	channel <- &PlayListMark{
		Length:                length,
		NumberOfPlayListMarks: numberOfPlayListMarks,
		PlayListMarksList:     playListMarksList,
	}
}

func parseExtensionData(rawData []byte, channel chan *ExtensionData) {
	length := int(binary.BigEndian.Uint32(rawData[:4]))
	if length == 0 {
		channel <- nil
	}
	dataBlockStartAddress := int(binary.BigEndian.Uint32(rawData[4:8]))

	numberOfExtDataEntries := int(rawData[11])
	var extDataEntryItemsList []*ExtDataEntryItem = nil
	for i := 0; i < numberOfExtDataEntries; i++ {
		extDataType := int(binary.BigEndian.Uint16(rawData[12+12*i : 14+12*i]))
		extDataVersion := int(binary.BigEndian.Uint16(rawData[14+12*i : 16+12*i]))
		extDataStartAddress := int(binary.BigEndian.Uint32(rawData[16+12*i : 20+12*i]))
		extDataLength := int(binary.BigEndian.Uint32(rawData[20+12*i : 24+12*i]))
		extDataEntry := rawData[extDataStartAddress : extDataStartAddress+extDataLength]

		extDataEntryItem := ExtDataEntryItem{
			ExtDataType:         extDataType,
			ExtDataVersion:      extDataVersion,
			ExtDataStartAddress: extDataStartAddress,
			ExtDataLength:       extDataLength,
			ExtDataEntry:        extDataEntry,
		}
		extDataEntryItemsList = append(extDataEntryItemsList, &extDataEntryItem)
	}

	channel <- &ExtensionData{
		Length:                 length,
		DataBlockStartAddress:  dataBlockStartAddress,
		NumberOfExtDataEntries: numberOfExtDataEntries,
		ExtDataEntryItemsList:  extDataEntryItemsList,
	}
}

func Parse(path string) (*MPLS, error) {
	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(rawData[:4], []byte("MPLS")) {
		return nil, errors.New("invalid file")
	}

	versionNumber, err := strconv.Atoi(string(rawData[0x04:0x08]))

	if err != nil {
		return nil, err
	}

	playlistStartAddress := int(binary.BigEndian.Uint32(rawData[0x08:0x0c]))
	playlistMarkStartAddress := int(binary.BigEndian.Uint32(rawData[0x0c:0x10]))
	extensionDataStartAddress := int(binary.BigEndian.Uint32(rawData[0x10:0x14]))

	applicationInfoPlayList := make(chan *AppInfoPlayList)
	go parseAppInfoPlayList(rawData[0x28:0x39], applicationInfoPlayList)

	playList := make(chan *PlayList)
	go parsePlayList(rawData[playlistStartAddress:], playList)

	playListMark := make(chan *PlayListMark)
	go parsePlayListMark(rawData[playlistMarkStartAddress:], playListMark)

	extensionData := make(chan *ExtensionData, 1)
	if extensionDataStartAddress != 0 {
		go parseExtensionData(rawData[extensionDataStartAddress:], extensionData)
	} else {
		extensionData <- nil
	}

	return &MPLS{
		FilePath:                  path,
		RawData:                   rawData,
		VersionNumber:             versionNumber,
		PlaylistStartAddress:      playlistStartAddress,
		PlaylistMarkStartAddress:  playlistMarkStartAddress,
		ExtensionDataStartAddress: extensionDataStartAddress,
		ApplicationInfoPlaylist:   <-applicationInfoPlayList,
		PlayList:                  <-playList,
		PlayListMark:              <-playListMark,
		ExtensionData:             <-extensionData,
	}, nil
}
