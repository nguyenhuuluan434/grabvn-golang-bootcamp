package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type daoServiceImpl struct {
	repository DatabaseRepository
}

func NewRespositoryService(repository DatabaseRepository) RespositoryService {
	return &daoServiceImpl{repository: repository}
}

func (repo daoServiceImpl) CreatePost(req RequestPostData) (id string, err error) {
	id = uuid.New().String()
	post := Post{ID: id, Title: req.Title}
	tx := repo.repository.db().Begin()
	err = tx.Create(&post).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo daoServiceImpl) UpdatePost(postId string, req RequestPostData) (result bool, err error) {
	tx := repo.repository.db().Begin()
	post, err := repo.checkPostExist(postId, tx)
	if err != nil || post == (Post{}) {
		tx.Rollback()
		return false, err
	}
	if post.Title != req.Title {
		post.Title = req.Title
	}
	db := tx.Model(&post).Where("id = ?", postId).Update(post)

	if db.Error != nil {
		tx.Rollback()
		return false, db.Error
	}
	if db.RowsAffected > 0 {
		tx.Commit()
		return true, nil
	}
	tx.Rollback()
	return false, nil
}

func (repo daoServiceImpl) GetPost(postId string) (post Post, err error) {
	return repo.GetPostWithTransaction(postId, nil)
}

func (repo daoServiceImpl) GetPostWithoutTransaction(postId string) (post Post, err error) {
	err = repo.repository.db().Where(&Post{ID: postId}).Find(&post).Error
	return
}

func (repo daoServiceImpl) GetPostWithTransaction(postId string, tx *gorm.DB) (post Post, err error) {
	err = tx.Where(&Post{ID: postId}).Find(&post).Error
	return
}

func (repo daoServiceImpl) DeletePost(postId string) (result bool, err error) {
	tx := repo.repository.db().Begin()
	post, err := repo.checkPostExist(postId, tx)
	if err != nil || post == (Post{}) {
		tx.Rollback()
		return false, err
	}
	comments, err := repo.getCommentByPostId(postId, tx)
	if err != nil {
		tx.Rollback()
		return
	}
	if len(comments) > 0 {
		for _, comment := range comments {
			delCmtResult, err := repo.deleteCommentWithTransaction(comment.ID, tx)
			if err != nil || !delCmtResult {
				tx.Rollback()
				return
			}
		}
	}
	db := tx.Where("id = ?", postId).Delete(&postId)
	if db.Error != nil {
		tx.Rollback()
		return false, db.Error
	}
	if db.RowsAffected > 0 {
		tx.Commit()
		return true, nil
	}
	tx.Rollback()
	return false, errors.New("not found")
}

func (repo daoServiceImpl) CreateComment(postId string, req RequestCommentData) (result string, err error) {

	tx := repo.repository.db().Begin()
	_, err = repo.GetPostWithTransaction(postId, tx)
	if err != nil {
		tx.Rollback()
		return string(""), errors.New("not found")
	}
	id := uuid.New().String()
	comment := Comment{ID: id, Body: req.Body, PostID: postId}
	err = tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return string(""), errors.New(fmt.Sprintf("could not save commment"))
	}
	tx.Commit()
	return
}

func (repo daoServiceImpl) GetComment(commentId string) (comment Comment, err error) {
	db := repo.repository.db().Where("id = ?", commentId).Find(&comment)
	if db.Error != nil {
		return comment, db.Error
	}
	return
}

func (repo daoServiceImpl) getCommentWithoutTransaction(commentId string) (comment Comment, err error) {
	err = repo.repository.db().Where(&Post{ID: commentId}).Find(&comment).Error
	return
}

func (repo daoServiceImpl) getCommentWithTransaction(commentId string, tx *gorm.DB) (comment Comment, err error) {
	err = tx.Where(&Post{ID: commentId}).Find(&comment).Error
	return
}

func (repo daoServiceImpl) DeleteComment(commentId string) (result bool, err error) {
	tx := repo.repository.db().Begin()
	return repo.deleteCommentWithTransaction(commentId, tx)
}

func (repo daoServiceImpl) deleteCommentWithoutTransaction(commentId string) (result bool, err error) {
	tx := repo.repository.db()
	comment, err := repo.getCommentWithTransaction(commentId, tx)
	if err != nil || comment == (Comment{}) {
		tx.Rollback()
		return false, err
	}
	db := repo.repository.db().Where("id = ?", commentId).Delete(&comment)
	if db.Error != nil {
		tx.Rollback()
		return false, db.Error
	}
	if db.RowsAffected > 0 {
		tx.Commit()
		return true, nil
	}
	return false, errors.New("not found")
}
func (repo daoServiceImpl) deleteCommentWithTransaction(commentId string, tx *gorm.DB) (result bool, err error) {
	comment, err := repo.getCommentWithTransaction(commentId, tx)
	if err != nil || comment == (Comment{}) {
		return false, err
	}
	db := tx.Where("id = ?", commentId).Delete(&comment)
	if db.Error != nil {
		return false, db.Error
	}
	if db.RowsAffected > 0 {
		return true, nil
	}
	return false, errors.New("not found")
}

func (repo daoServiceImpl) checkPostExist(postId string, tx *gorm.DB) (post Post, err error) {
	post, err = repo.GetPostWithTransaction(postId, tx)
	if err != nil {
		return
	}
	if post == (Post{}) {
		return post, errors.New(fmt.Sprintf("not found post with Id %s", postId))
	}
	return post, nil
}

func (repo daoServiceImpl) getCommentByPostId(postId string, tx *gorm.DB) (comments []Comment, err error) {
	var db *gorm.DB
	if tx == nil {
		db = repo.repository.db().Where("id = ?", postId).Find(&comments)
	} else {
		db = tx.Where("id = ?", postId).Find(&comments)
	}
	err = db.Error
	return
}

func (repo daoServiceImpl) GetAllPostWithComments() (posts []PostWithComments, err error) {

	var rawQuery = "SELECT p.id,p.title,c.id,c.body,c. FROM posts p LEFT JOIN comments c ON p.id = c.post_id"

	rows, err := repo.repository.db().Raw(rawQuery).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	var postMap = make(map[string]*PostWithComments)
	for rows.Next() {
		var (
			postId      string
			title       string
			commentId   string
			commentBody string
		)
		rows.Scan(&postId, &title, &commentId, &commentBody)
		if postWithComment, ok := postMap[postId]; ok {
			postWithComment.Comments = append(postWithComment.Comments, Comment{ID: commentId, Body: commentBody})
			continue
		}
		postWithComments := PostWithComments{ID: postId, Title: title, Comments: []Comment{{ID: commentId, Body: commentBody}}}
		postMap[postId] = &postWithComments
	}
	if len(postMap) > 0 {
		posts = make([]PostWithComments, 1, 5)
		for _, v := range postMap {
			posts = append(posts, PostWithComments{ID: v.ID, Title: v.Title, Comments: v.Comments})
		}
	}

	return
}
