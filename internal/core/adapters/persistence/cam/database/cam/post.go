package cam

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	"github.com/jmoiron/sqlx"
	"time"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type Post struct {
	camXO.Post
}

func PostFromModel(p *types.Post) *Post {
	description := sql.NullString{
		String: p.Description,
	}
	if p.Description != "" {
		description.Valid = true
	}
	return &Post{
		Post: camXO.Post{
			PostID:      p.ID,
			Description: description,
			UserID:      p.UserID,
		},
	}
}

func (p *Post) ToModel() *types.Post {
	return &types.Post{
		ID:          p.PostID,
		Description: p.Description.String,
		UserID:      p.UserID,
		Date:        p.Date,
	}
}

func (d *Driver) CreatePost(post *Post) (*Post, error) {
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
	post.Date = time.Now()
	err = post.Insert(context.Background(), d.db)
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
		return post, nil
	}
}

func (d *Driver) RetrievePostByID(id int) (*Post, error) {
	p, err := camXO.PostByPostID(context.Background(), d.db, id)
	if err != nil {
		return nil, err
	}
	return &Post{
		Post: *p,
	}, nil
}

func (d *Driver) RetrieveSubscribedPostsPaginated(userID int, pageNumber int, batchSize int) ([]*Post, error) {
	var whereClause = ``
	if userID != 0 {
		whereClause = fmt.Sprintf("WHERE UserID IN (SELECT FollowedUserID FROM captureamoment.subscription WHERE subscription.UserID = %d)", userID)
	}
	var sqlStr = `SELECT ` +
		`PostID, Description, Date, UserID ` +
		`FROM captureamoment.post ` +
		whereClause +
		`ORDER BY PostID DESC ` +
		`LIMIT ?,?;`
	offset := (pageNumber - 1) * batchSize
	rows, err := d.db.Queryx(sqlStr, offset, batchSize)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	posts, err := postsFromRows(rows)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (d *Driver) RetrieveUserPostsPaginated(userID int, pageNumber int, batchSize int) ([]*Post, error) {
	var sqlstr = `SELECT ` +
		`PostID, Description, Date, UserID ` +
		`FROM captureamoment.post ` +
		`WHERE UserID = ? ` +
		`ORDER BY PostID DESC ` +
		`LIMIT ?,?;`
	offset := (pageNumber - 1) * batchSize
	rows, err := d.db.Queryx(sqlstr, userID, offset, batchSize)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	posts, err := postsFromRows(rows)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func postsFromRows(rows *sqlx.Rows) ([]*Post, error) {
	var posts []*Post
	for rows.Next() {
		var pid, uid int
		var date time.Time
		var desc sql.NullString
		err := rows.Scan(&pid, &desc, &date, &uid)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &Post{
			Post: camXO.Post{
				PostID:      pid,
				Description: desc,
				Date:        date,
				UserID:      uid,
			},
		},
		)
	}
	return posts, nil
}
