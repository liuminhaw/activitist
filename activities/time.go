package activities

import (
	// "fmt"
	"database/sql"
	"strings"
	"time"
)

type activityTime time.Time

// Implement Marshaler and Unmarshaler interface
func (at *activityTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	// fmt.Printf("trimmed time: %s\n", s)
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*at = activityTime(t)

	return nil
}

func (at *activityTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(*at).Format("2006-01-02 15:04:05") + "\""), nil
}

// Format function for printing your date with input format
func (at *activityTime) String() string {
	t := time.Time(*at)
	return t.Format("2006-01-02 15:04:05")
}

func (at *activityTime) IsZero() bool {
	return time.Time(*at).IsZero()
}

func (at *activityTime) NewNullTime() sql.NullTime {
	if at.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  time.Time(*at),
		Valid: true,
	}
}
