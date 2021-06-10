package utils

import "time"

var LocalLocation *time.Location

func init() {
	LocalLocation = time.Now().Location()
}
