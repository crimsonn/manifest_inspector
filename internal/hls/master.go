package hls

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func (m *MasterPlaylist) Inspect() InspectionReport {
	report := InspectionReport{
		InspectTime: time.Now(),
		IsValid:     true,
		Issues:      make([]Issue, 0),
		Stats: MasterStats{
			VariantCount:    len(m.Variants),
			MediaCount:      len(m.Media),
			CodecsUsed:      make([]string, 0),
			ResolutionsUsed: make([]string, 0),
		},
	}

	if len(m.Variants) == 0 {
		report.addIssue(SeverityError, "ERR_NO_VARIANTS", "Master playlist does not contain any #EXT-X-STREAM-INF variant streams")
	}
	var minBW, maxBW int = -1, -1
	codecsMap := make(map[string]bool)
	for i, v := range m.Variants {
		if v.Codecs == "" {
			report.addIssue(SeverityWarning, "WARN_MISSING_CODECS", fmt.Sprintf("Variant #%d (URI: %s) is missing CODECS attribute", i+1, v.Playlist))
		} else {
			codecsMap[v.Codecs] = true
		}

		if v.Resolution != "" {
			report.Stats.ResolutionsUsed = append(report.Stats.ResolutionsUsed, v.Resolution)
		}
		if minBW == -1 || v.Bandwidth < minBW {
			minBW = v.Bandwidth
		}
		if maxBW == -1 || v.Bandwidth > maxBW {
			maxBW = v.Bandwidth
		}

		if i > 0 && v.Bandwidth < m.Variants[i-1].Bandwidth {
			report.addIssue(SeverityWarning, "WARN_UNSORTED_VARIANTS", "Variants are not sorted in ascending order by BANDWIDTH (may affect player ABR selection)")
		}
	}

	report.Stats.MinBandwidth = minBW
	report.Stats.MaxBandwidth = maxBW
	for codec := range codecsMap {
		report.Stats.CodecsUsed = append(report.Stats.CodecsUsed, codec)
	}

	for _, med := range m.Media {
		switch med.MediaType {
		case "AUDIO":
			report.Stats.HasAudioRendition = true
			report.Stats.AudioGroupCount++
			m.validateMedia(&report, med)
		case "SUBTITLES":
			report.Stats.HasSubtitles = true
			report.Stats.SubGroupCount++
			m.validateMedia(&report, med)

		}
	}

	return report
}

func (m *MasterPlaylist) validateMedia(report *InspectionReport, med MasterMedia) {
	stableRenditionIdRegex, err := regexp.Compile(`^[a-zA-Z0-9+/=._-]+$`)
	if err != nil {
		report.addIssue(SeverityError, "ERR_INVALID_STABLE_RENDITION_ID_REGEX", "Invalid stable rendition ID regex")
	}
	if med.StableRenditionID != "" && !stableRenditionIdRegex.MatchString(med.StableRenditionID) {
		report.addIssue(SeverityError, "ERR_INVALID_STABLE_RENDITION_ID", fmt.Sprintf("Invalid stable rendition ID for audio group %s, must be alphanumeric and contain only [a-zA-Z0-9+/=.-_]", med.GroupID))
	}
	if med.Default != "YES" && med.Default != "NO" {
		report.addIssue(SeverityError, "ERR_INVALID_DEFAULT", fmt.Sprintf("Invalid default attribute for audio group %s, must be YES or NO", med.GroupID))
	}
	if med.MediaType == "AUDIO" && med.GroupID != "" {
		// find that group id in variants
		found := false
		for _, v := range m.Variants {
			if v.Audio == med.GroupID {

				found = true
				break
			}
		}
		if !found {
			report.addIssue(SeverityError, "ERR_AUDIO_GROUP_NOT_FOUND", fmt.Sprintf("Audio group %s not found in variants", med.GroupID))
		}
	}

	if med.MediaType == "SUBTITLES" && med.GroupID != "" {
		// find that group id in variants
		found := false
		for _, v := range m.Variants {
			if v.Subtitles == med.GroupID {
				found = true
				break
			}
		}
		if !found {
			report.addIssue(SeverityError, "ERR_SUBTITLES_GROUP_NOT_FOUND", fmt.Sprintf("Subtitles group %s not found in variants", med.GroupID))
		}
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

	return nil
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
