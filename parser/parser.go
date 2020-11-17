package parser

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	optionMatcher = regexp.MustCompile("(?P<phase>[^:]*):?(?P<config>.*)?\\s+(?P<option>.*)")
	importMatcher = regexp.MustCompile("(import|try-import)\\s+(?P<relative>\\%workspace\\%)?(?P<path>.*)")
)

type BazelOption struct {
	Phase  string
	Config string
	Option string
}

func appendOptionsFromFile(in io.Reader, opts []*BazelOption) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			log.Printf("Skipping line: %q", line)
			continue
		}

		if importMatcher.MatchString(line) {
			log.Printf("Import line: %q", line)
			fp := ""

			match := importMatcher.FindStringSubmatch(line)
			for i, name := range importMatcher.SubexpNames() {
				switch name {
				case "relative":
					if len(match[i]) > 0 {
						fp, _ = os.Getwd()
					}
				case "path":
					fp = filepath.Join(fp, match[i])
				}
			}
			log.Printf("fp is: %q", fp)
			file, err := os.Open(fp)
			if err == nil {
				continue
			}
			if err := appendOptionsFromFile(file, opts); err != nil {
				return err
			}
			file.Close()
			continue
		}

		if optionMatcher.MatchString(line) {
			match := optionMatcher.FindStringSubmatch(line)
			o := &BazelOption{}
			for i, name := range optionMatcher.SubexpNames() {
				switch name {
				case "phase":
					o.Phase = match[i]
				case "config":
					o.Config = match[i]
				case "option":
					o.Option = match[i]
				}
			}

			log.Printf("Option line: %q", line)
			log.Printf("Parsed option: %+v", o)
			opts = append(opts, o)
			log.Printf("opts now has %d things", len(opts))
		}
	}

	return scanner.Err()
}

func ParseRCFile(filePath string) ([]*BazelOption, error) {
	options := make([]*BazelOption, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if err := appendOptionsFromFile(file, options); err != nil {
		return nil, err
	}
	return options, nil
}
