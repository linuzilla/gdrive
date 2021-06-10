package ignore_service

import (
	"bufio"
	"github.com/linuzilla/gdrive/constants"
	"os"
	"regexp"
	"strings"
)

type ignoreServiceImpl struct {
}

func New() IgnoreService {
	return &ignoreServiceImpl{}
}

func (svc *ignoreServiceImpl) LoadRules(baseDir string) FileNameGuard {
	file, err := os.Open(baseDir + `/` + constants.GDriveIgnoreFile)
	if err != nil {
		return nil
	}

	defer file.Close()

	var entries []*entry

	entries = append(entries, &entry{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		positive := true
		exactly := false
		if len(text) == 0 {
			continue
		}
		if strings.HasPrefix(text, "#") {
			continue
		} else if strings.HasPrefix(text, "!") {
			positive = false
			text = strings.TrimSpace(text[1:])
		}

		if strings.HasPrefix(text, `/`) {
			exactly = true
		}

		if strings.Contains(text, `*`) {
			if pattern, err := regexp.Compile(`^` +
				strings.ReplaceAll(strings.ReplaceAll(
					text, `.`, `\.`),
					`*`, `.*`) +
				`$`); err != nil {
				entries = append(entries, &entry{
					Index:    len(entries) + 1,
					Filename: text,
					Exactly:  exactly,
					Positive: positive,
				})
			} else {
				entries = append(entries, &entry{
					Index:    len(entries) + 1,
					Filename: text,
					Pattern:  pattern,
					Exactly:  exactly,
					Positive: positive,
				})
			}
		} else {
			entries = append(entries, &entry{
				Index:    len(entries) + 1,
				Filename: text,
				Exactly:  exactly,
				Positive: positive,
			})
		}
	}

	if len(entries) == 0 {
		return nil
	} else {
		return &fileNameGuardImpl{
			entries: entries,
		}
	}
}
