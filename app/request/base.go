package request

import (
	"encoding/json"
	"strings"
	"time"
)

/*
|========================================================================
|	This File For Generic Function or Method on another file or struct
|	and use it wisely !
|========================================================================
*/
type JsonDate time.Time

// implement Marshaller und Unmarshaller interface
func (j *JsonDate) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" || string(b) == `""` {
		return nil
	}
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Maybe a Format function for printing your date
func (j JsonDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
