package repository

import (
	"github.com/hanifmaliki/chat-app/internal/entity"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"

	"gorm.io/gorm"
)

type RoomUserRepository interface {
	Create(data *entity.RoomUser, by string) error
	Delete(conds *entity.RoomUser, by string) error
	Find(conds *entity.RoomUser, query *pkg_model.Query) ([]*entity.RoomUser, error)
	FindOne(conds *entity.RoomUser, query *pkg_model.Query) (*entity.RoomUser, error)
}

type roomUserRepository struct {
	db *gorm.DB
}

func NewRoomUserRepository(db *gorm.DB) RoomUserRepository {
	return &roomUserRepository{db: db}
}

func (r *roomUserRepository) Create(data *entity.RoomUser, by string) error {
	data.CreatedBy = by
	data.UpdatedBy = by
	return r.db.Create(data).Error
}

func (r *roomUserRepository) Delete(conds *entity.RoomUser, by string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.RoomUser{}).Omit("updated_at").Where(conds).Update("deleted_by", by).Error
		if err != nil {
			return err
		}
		return tx.Where(conds).Delete(&entity.RoomUser{}).Error
	})
}

func (r *roomUserRepository) Find(conds *entity.RoomUser, query *pkg_model.Query) ([]*entity.RoomUser, error) {
	var datas []*entity.RoomUser
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

func (r *roomUserRepository) FindOne(conds *entity.RoomUser, query *pkg_model.Query) (*entity.RoomUser, error) {
	data := &entity.RoomUser{}
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
