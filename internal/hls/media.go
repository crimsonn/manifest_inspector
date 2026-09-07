package hls

import (
	"fmt"
	"strconv"
	"strings"
)

type MasterMedia struct {
	MediaType         string  `json:"type"`
	Uri               string  `json:"uri"`
	GroupID           string  `json:"group_id"`
	Language          string  `json:"language"`
	AssocLanguage     string  `json:"assoc_language"`
	Name              string  `json:"name"`
	StableRenditionID string  `json:"stable_rendition_id"`
	Default           string  `json:"defualt"`
	AutoSelect        string  `json:"auto_select"`
	Forced            string  `json:"forced"`
	InStreamID        string  `json:"instream_id"`
	BitDepth          string  `json:"bit_depth"`
	SampleRate        float32 `json:"sample_rate"`
	Characteristics   string  `json:"characteristics"`
	Channels          string  `json:"channels"`
}

func parseMedia(line string) (*MasterMedia, error) {
	m := &MasterMedia{}
	for k, v := range parseAttributes(line) {
		switch strings.ToUpper(k) {
		case "TYPE":
			m.MediaType = v
		case "URI":
			m.Uri = v
		case "GROUP-ID":
			m.GroupID = v
		case "LANGUAGE":
			m.Language = v
		case "ASSOC-LANGUAGE":
			m.AssocLanguage = v
		case "NAME":
			m.Name = v
		case "STABLE-RENDITION-ID":
			m.StableRenditionID = v
		case "DEFAULT":
			m.Default = v
		case "AUTO-SELECT":
			m.AutoSelect = v
		case "FORCED":
			m.Forced = v
		case "INSTREAM-ID":
			m.InStreamID = v
		case "BIT-DEPTH":
			m.BitDepth = v
		case "SAMPLE-RATE":
			val, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid sample rate %q: %w", v, err)
			}
			m.SampleRate = float32(val)
		case "CHARACTERISTICS":
			m.Characteristics = v
		case "CHANNELS":
			m.Channels = v
		}
	}

	return m, nil
}
