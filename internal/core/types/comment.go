package types

type Comment struct {
	ID       int
	Content  string
	ParentID int
	PostID   int
	UserID   int
}
