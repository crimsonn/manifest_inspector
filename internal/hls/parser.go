package hls

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	ErrUnableToFetchURL = errors.New("unable to fetch HLS/MPD manifest from url")
	ErrUnableToParse    = errors.New("unable to parse HLS/MPD manifest")

	ParserErrorInvalidNumber = errors.New("invalid number")
	ParserErrorInvalidTag    = errors.New("invalid tag")
)

type Parser interface {
	Parse(line string) error
}

type parser struct {
}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnableToFetchURL, err)
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnableToParse, err)
	}
	if strings.Contains(string(buf), "#EXT-X-STREAM-INF") {
		err := p.parseMasterPlaylist(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseMasterPlaylist(buf []byte) error {
	scanner := bufio.NewScanner(strings.NewReader(string(buf)))
	if scanner.Err() != nil {
		return fmt.Errorf("%w: %v", ErrUnableToParse, scanner.Err())
	}
	master := NewMasterPlaylist()
	err := master.scan(scanner)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	report := master.Inspect()
	if !report.IsValid {
		fmt.Println("Manifest failed validation!")
		for _, issue := range report.Issues {
			fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Code, issue.Message)
		}
	} else {
		fmt.Printf("Valid Master Playlist with %d variants (Bitrate range: %d bps - %d bps)\n",
			report.Stats.VariantCount, report.Stats.MinBandwidth, report.Stats.MaxBandwidth)
	}

	return nil
}
