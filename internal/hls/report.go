package hls

import "time"

type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

type Issue struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Line     int      `json:"line"`
}

type MasterStats struct {
	VariantCount      int      `json:"variant_count"`
	MediaCount        int      `json:"media_count"`
	AudioGroupCount   int      `json:"audio_group_count"`
	SubGroupCount     int      `json:"sub_group_count"`
	MinBandwidth      int      `json:"min_bandwidth"`
	MaxBandwidth      int      `json:"max_bandwidth"`
	CodecsUsed        []string `json:"codecs_used"`
	ResolutionsUsed   []string `json:"resolutions_used"`
	HasAudioRendition bool     `json:"has_audio_rendition"`
	HasSubtitles      bool     `json:"has_subtitles"`
}

type InspectionReport struct {
	InspectTime time.Time   `json:"inspect_time"`
	IsValid     bool        `json:"is_valid"` // True if no SeverityError exists
	Issues      []Issue     `json:"issues"`
	Stats       MasterStats `json:"stats"`
}

func (r *InspectionReport) addIssue(sev Severity, code string, msg string) {
	if sev == SeverityError {
		r.IsValid = false
	}
	r.Issues = append(r.Issues, Issue{
		Code:     code,
		Severity: sev,
		Message:  msg,
	})
}
