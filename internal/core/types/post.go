package types

import "time"

type Post struct {
	ID          int
	Description string
	UserID      int
	Date        time.Time
	Media       []*Media
	Comments    []*Comment
}
