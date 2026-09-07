package hls

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type MasterPlaylist struct {
	Version  int
	Media    []MasterMedia
	Variants []MasterVariant
}

func NewMasterPlaylist() *MasterPlaylist {
	return &MasterPlaylist{
		Version: 0,
		Media:   make([]MasterMedia, 0),
	}
}

func (m *MasterPlaylist) scan(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#EXTM3U") {
			continue
		}

		tag, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch tag {
		case "#EXT-X-VERSION":
			versionNum, err := strconv.Atoi(value)
			if err != nil {
				return ParserErrorInvalidNumber
			}
			m.setVersion(versionNum)
		case "#EXT-X-MEDIA":
			media, err := parseMedia(value)
			if err != nil {
				return err
			}
			m.setMedia(*media)
		case "#EXT-X-STREAM-INF":
			variant, err := parseVariant(value, scanner)
			if err != nil {
				return err
			}
			m.setVariants(*variant)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("Variants: %+v\n", m.Variants)

	return nil
}

func (m *MasterPlaylist) setMedias(media []MasterMedia) {
	m.Media = media
}

func (m *MasterPlaylist) setMedia(media MasterMedia) {
	m.Media = append(m.Media, media)
}

func (m *MasterPlaylist) setVersion(version int) {
	m.Version = version
}

func (m *MasterPlaylist) setVariants(variant MasterVariant) {
	m.Variants = append(m.Variants, variant)
}
