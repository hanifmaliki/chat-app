package repository

import (
	"github.com/hanifmaliki/chat-app/internal/entity"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(data *entity.User, by string) error
	Delete(conds *entity.User, by string) error
	Find(conds *entity.User, query *pkg_model.Query) ([]*entity.User, error)
	FindOne(conds *entity.User, query *pkg_model.Query) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(data *entity.User, by string) error {
	data.CreatedBy = by
	data.UpdatedBy = by
	return r.db.Create(data).Error
}

func (r *userRepository) Delete(conds *entity.User, by string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.User{}).Omit("updated_at").Where(conds).Update("deleted_by", by).Error
		if err != nil {
			return err
		}
		return tx.Where(conds).Delete(&entity.User{}).Error
	})
}

func (r *userRepository) Find(conds *entity.User, query *pkg_model.Query) ([]*entity.User, error) {
	var datas []*entity.User
	db := r.db

	// Handle preload/expansion of related entities
	for _, expand := range query.Expand {
		db = db.Preload(expand)
	}

	// Handle sorting
	if query.SortBy != "" {
		db = db.Order(query.SortBy)
	}

	// Execute the query
	if err := db.Where(conds).Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}

func (r *userRepository) FindOne(conds *entity.User, query *pkg_model.Query) (*entity.User, error) {
	data := &entity.User{}
	db := r.db
	for _, expand := range query.Expand {
		db = db.Preload(expand)
	}
	if query.SortBy != "" {
		db = db.Order(query.SortBy)
	}
	err := db.Where(conds).First(data).Error
	return data, err
}
