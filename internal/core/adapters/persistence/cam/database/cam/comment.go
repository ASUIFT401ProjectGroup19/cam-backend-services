package cam

import (
	"context"
	"database/sql"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type Comment struct {
	camXO.Comment
}

func CommentFromModel(comment *types.Comment) *Comment {
	c := &Comment{
		Comment: camXO.Comment{
			CommentID: comment.ID,
			CommentText: sql.NullString{
				String: comment.Content,
				Valid:  true,
			},
			ParentID: sql.NullInt64{},
			PostID:   comment.PostID,
			UserID:   comment.UserID,
		},
	}
	if comment.ParentID != 0 {
		c.ParentID = sql.NullInt64{
			Int64: int64(comment.ParentID),
			Valid: true,
		}
	}
	return c
}

func (c *Comment) ToModel() *types.Comment {
	comment := &types.Comment{
		ID:       c.CommentID,
		Content:  c.CommentText.String,
		ParentID: int(c.ParentID.Int64),
		PostID:   c.PostID,
		UserID:   c.UserID,
	}
	return comment
}

func (d *Driver) CreateComment(comment *Comment) (*Comment, error) {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return nil, &BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = comment.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return nil, &Unknown{Message: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, &Exists{Message: err.Error()}
		} else {
			return nil, &InsertRecord{Message: err.Error()}
		}
	case nil:
		return comment, nil
	}
}

func (d *Driver) ReadComment(commentID int) (*Comment, error) {
	c, err := camXO.CommentByCommentID(context.Background(), d.db, commentID)
	if err != nil {
		return nil, err
	}
	return &Comment{
		Comment: *c,
	}, nil
}

func (d *Driver) ReadCommentsByPostID(postID int) ([]*Comment, error) {
	c, err := camXO.CommentByPostID(context.Background(), d.db, postID)
	if err != nil {
		return nil, err
	}
	comments := make([]*Comment, len(c))
	for i, v := range c {
		comments[i] = &Comment{
			Comment: *v,
		}
	}
	return comments, err
}
