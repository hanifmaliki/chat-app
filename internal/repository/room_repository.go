package repository

import (
	"github.com/hanifmaliki/chat-app/internal/entity"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(data *entity.Room, by string) error
	Delete(conds *entity.Room, by string) error
	Find(conds *entity.Room, query *pkg_model.Query) ([]*entity.Room, error)
	FindOne(conds *entity.Room, query *pkg_model.Query) (*entity.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(data *entity.Room, by string) error {
	data.CreatedBy = by
	data.UpdatedBy = by
	return r.db.Create(data).Error
}

func (r *roomRepository) Delete(conds *entity.Room, by string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.Room{}).Omit("updated_at").Where(conds).Update("deleted_by", by).Error
		if err != nil {
			return err
		}
		return tx.Where(conds).Delete(&entity.Room{}).Error
	})
}

func (r *roomRepository) Find(conds *entity.Room, query *pkg_model.Query) ([]*entity.Room, error) {
	var datas []*entity.Room
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

func (r *roomRepository) FindOne(conds *entity.Room, query *pkg_model.Query) (*entity.Room, error) {
	data := &entity.Room{}
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
