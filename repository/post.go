package repository

import (
	"ELKExample/models"
	"fmt"
	"time"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	insertOp    = "insert"
	deleteOp    = "delete"
	updateOp    = "update"
)

func (db Databases) SavePost(post *models.Posts) error {
	if err := db.Conn.Table("posts").Save(post).Error; err != nil {
		return err
	}

	db.SaveLogs(post.ID, "Insert Data")

	return nil
}

func (db Databases) UpdatePost(post *models.Posts) error {
	if err := db.Conn.Exec("update posts set title = ?, body = ? where id = ?", post.Title, post.Body, post.ID).Error; err != nil {
		return err
	}

	db.SaveLogs(post.ID, "Update Data")

	return nil
}

func (db Databases) DeletePost(id int) error {
	if err := db.Conn.Delete(&models.Posts{}, id).Error; err != nil {
		return err
	}

	db.SaveLogs(id, "Delete Data")

	return nil
}

func (db Databases) GetPost() ([]*models.Posts, error) {
	var listPosts []*models.Posts
	if err := db.Conn.Table("posts").Find(&listPosts).Error; err != nil {
		return nil,err
	}

	return listPosts, nil
}

func (db Databases) SaveLogs(id int, operation string) {
	postLogs := &models.PostLog{
		PostID:    id,
		Operation: operation,
		Create_at: time.Now(),
	}
	if err := db.Conn.Table("post_logs").Save(postLogs).Error; err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
}
