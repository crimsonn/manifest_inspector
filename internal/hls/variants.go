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
	if !scanner.Scan() {
		return nil, io.ErrUnexpectedEOF
	}

	variant := &MasterVariant{}
	for k, v := range parseAttributes(text) {
		switch strings.ToUpper(k) {
		case "BANDWIDTH":
			bandwidthNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, ParserErrorInvalidNumber
			}
			variant.Bandwidth = bandwidthNum
		}
	}

	return variant, nil

}
