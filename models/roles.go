package models

import (
	"fmt"
)

type Role string

const (
	Organizer     = `organizer`
	FileOrganizer = `fileOrganizer`
	Writer        = `write`
	Reader        = `reader`
	Commenter     = `commenter`
)

var AllRoles = []Role{
	Organizer,
	FileOrganizer,
	Writer,
	Reader,
	Commenter,
}

func (r Role) IsValid() error {
	switch r {
	case Organizer, FileOrganizer, Writer, Reader, Commenter:
		return nil
	}
	return fmt.Errorf("%s unknow role", r)
}
