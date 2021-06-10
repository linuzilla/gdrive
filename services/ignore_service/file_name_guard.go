package ignore_service

import (
	"path"
)

type fileNameGuardImpl struct {
	entries []*entry
}

func (svc *fileNameGuardImpl) ShouldIgnore(relativeDir string, fileName string) (*entry, bool) {
	fullName := path.Join(relativeDir, fileName)

	for _, entry := range svc.entries {
		if entry.Exactly {
			if entry.Pattern != nil {
				if entry.Pattern.MatchString(fullName) {
					//logger.Notice("%s full-match-pattern rule %d (%s)", fullName, i + 1, entry.filename)
					return entry, entry.Positive
				}
			} else {
				if entry.Filename == fullName {
					//logger.Notice("%s full-match-name rule %d (%s)", fullName, i + 1, entry.filename)
					return entry, entry.Positive
				}
			}
		} else {
			if entry.Pattern != nil {
				if entry.Pattern.MatchString(fileName) {
					//logger.Notice("%s match-pattern rule %d (%s)", fullName, i + 1, entry.filename)
					return entry, entry.Positive
				}
			} else {
				if entry.Filename == fileName {
					//logger.Notice("%s match-name rule %d (%s)", fullName, i + 1, entry.filename)
					return entry, entry.Positive
				}
			}
		}
	}
	return nil, false
}
