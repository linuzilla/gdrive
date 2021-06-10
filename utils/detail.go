package utils

import (
	"bytes"
	"fmt"
	"reflect"
)

func Detail(something interface{}) string {
	var buffer bytes.Buffer

	s := reflect.ValueOf(something).Elem()
	typeOfT := s.Type()

	fmt.Fprintf(&buffer, ">> Type: %s\n", reflect.TypeOf(something).String())

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if f.CanInterface() {
			fmt.Fprintf(&buffer, ">>   %d: %s %s = %v `%s`", i,
				typeOfT.Field(i).Name,
				f.Type(),
				f.Interface(),
				typeOfT.Field(i).Tag)
		} else {
			fmt.Fprintf(&buffer, ">>   %d: %s %s `%s`", i,
				typeOfT.Field(i).Name,
				f.Type(),
				typeOfT.Field(i).Tag)
		}

		if typeOfT.Field(i).Anonymous {
			buffer.WriteString(" (embedded)")
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}
