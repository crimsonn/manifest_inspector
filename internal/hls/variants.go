package hls

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type MasterVariant struct {
	Bandwidth        int
	AverageBandwidth int
	Codecs           string
	Resolution       string
	FrameRate        int
	VideoRange       string
	Audio            string
	Subtitles        string
	ClosedCaptions   string
	Playlist         string
}

func parseVariant(text string, scanner *bufio.Scanner) (*MasterVariant, error) {

	variant := &MasterVariant{}
	for k, v := range parseAttributes(text) {
		switch strings.ToUpper(k) {
		case "BANDWIDTH":
			bandwidthNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, ParserErrorInvalidNumber
			}
			variant.Bandwidth = bandwidthNum
		case "AVERAGE_BANDWIDTH":
			bandwidthNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, ParserErrorInvalidNumber
			}
			variant.AverageBandwidth = bandwidthNum
		case "CODECS":
			variant.Codecs = v
		case "RESOLUTION":
			variant.Resolution = v
		case "FRAME-RATE":
			frameRate, err := strconv.Atoi(v)
			if err != nil {
				return nil, ParserErrorInvalidNumber
			}
			variant.FrameRate = frameRate
		case "VIDEO-RANGE":
			variant.VideoRange = v
		case "AUDIO":
			variant.Audio = v
		case "SUBTITLES":
			variant.Subtitles = v
		case "CLOSED-CATPIONS":
			variant.ClosedCaptions = v
		}
	}

	if scanner.Scan() {
		playlist := strings.TrimSpace(scanner.Text())
		variant.Playlist = playlist
	} else {
		return nil, io.ErrUnexpectedEOF
	}

	return variant, nil

}
