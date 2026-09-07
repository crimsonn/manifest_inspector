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

	m.Display()
	return m, nil
}

func (m *MasterMedia) Display() {
	fmt.Printf("MediaType=%s\n", m.MediaType)
	fmt.Printf("Uri=%s\n", m.Uri)
	fmt.Printf("GroupID=%s\n", m.GroupID)
	fmt.Printf("Language=%s\n", m.Language)
	fmt.Printf("AssocLanguage=%s\n", m.AssocLanguage)
	fmt.Printf("Name=%s\n", m.Name)
	fmt.Printf("StableRenditionID=%s\n", m.StableRenditionID)
	fmt.Printf("Default=%s\n", m.Default)
	fmt.Printf("AutoSelect=%s\n", m.AutoSelect)
	fmt.Printf("Forced=%s\n", m.Forced)
	fmt.Printf("InStreamID=%s\n", m.InStreamID)
	fmt.Printf("BitDepth=%s\n", m.BitDepth)
	fmt.Printf("SampleRate=%f\n", m.SampleRate)
	fmt.Printf("Characteristics=%s\n", m.Characteristics)
	fmt.Printf("Channels=%s\n", m.Channels)
}
