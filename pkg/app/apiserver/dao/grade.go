//go:generate mockgen --source grade.go --destination mock/grade.mock.go
package dao

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	libsgorm "sunflower/pkg/libs/gormutil_v2"

	"sunflower/pkg/app/apiserver/model"
)

type GradeDaoImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (impl *GradeDaoImpl) FetchOne(scopes *GradeScope) (model.Grade, error) {
	res := model.Grade{}
	db := impl.db.Scopes(scopes.Scopes()...)
	if err := db.Model(res).Take(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (impl *GradeDaoImpl) CreateMany(grads []model.Grade) error {
	if err := impl.db.Create(&grads).Error; err != nil {
		impl.logger.Error("create grads failed", zap.Error(err))
		return err
	}
	return nil
}

func (impl *GradeDaoImpl) List(scopes *GradeScope, limit, cursor int) ([]model.Grade, *libsgorm.Page, error) {
	res := make([]model.Grade, 0)

	db := impl.db.Scopes(scopes.Scopes()...)
	query := db.Model(model.Grade{}).Find(&res)
	page, err := libsgorm.Pagination(query, limit, cursor, &res)
	if err != nil {
		impl.logger.Error("list failed", zap.Error(err))
		return nil, nil, err
	}
	return res, page, nil
}

type GradeDao interface {
	List(scopes *GradeScope, limit, cursor int) ([]model.Grade, *libsgorm.Page, error)
	CreateMany(grads []model.Grade) error
	FetchOne(scopes *GradeScope) (model.Grade, error)
}

type NewGradeDaoFunc = func(*gorm.DB, *zap.Logger) GradeDao

func NewGradeDaoImpl(db *gorm.DB, logger *zap.Logger) GradeDao {
	return &GradeDaoImpl{
		db:     db,
		logger: logger,
	}
}
