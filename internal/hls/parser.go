package hls

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
		// Master Playlist
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
	masterPlaylist := NewMasterPlaylist()
	masterPlaylist.scan(scanner)

	return nil
}
