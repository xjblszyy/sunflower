package service

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	goa "goa.design/goa/v3/pkg"
	"gorm.io/gorm"

	libsgorm "sunflower/pkg/libs/gormutil_v2"

	"sunflower/pkg/api/gen/score"
	"sunflower/pkg/app/apiserver/dao"
	"sunflower/pkg/app/apiserver/model"
)

var (
	ErrGradeInternalError = score.MakeInternalServerError(errors.New("服务器内部错误"))
	ErrGradeNotFoundError = score.MakeBadRequest(errors.New("记录不存在"))
	ErrGradeError         = func(err string) *goa.ServiceError {
		return score.MakeBadRequest(fmt.Errorf("%s", err))
	}
)

type GradeSVC interface {
	List(sortField, sortOrder, name, class, subject *string, score *int, limit, cursor int) ([]model.Grade, *libsgorm.Page, error)
	Create(data []model.Grade) error
	FetchOne(id int) (model.Grade, error)
}

func NewGradeSVCImpl(db *gorm.DB, logger *zap.Logger,
	gradeDao dao.NewGradeDaoFunc,
) GradeSVC {
	return &GradeSVCImpl{
		db:              db,
		logger:          logger,
		newGradeDaoFunc: gradeDao,
	}
}

type GradeSVCImpl struct {
	db              *gorm.DB
	logger          *zap.Logger
	newGradeDaoFunc dao.NewGradeDaoFunc
}

func (impl *GradeSVCImpl) FetchOne(id int) (model.Grade, error) {
	gradeDao := impl.newGradeDaoFunc(impl.db, impl.logger)
	res, err := gradeDao.FetchOne(dao.NewGradeScope().AddPK(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, ErrGradeNotFoundError
		}
		impl.logger.Error("fetch one failed", zap.Error(err))
		return res, ErrGradeInternalError
	}
	return res, nil
}

func (impl *GradeSVCImpl) Create(data []model.Grade) error {
	gradeDao := impl.newGradeDaoFunc(impl.db, impl.logger)
	if err := gradeDao.CreateMany(data); err != nil {
		if model.IsDuplicateError(err) {
			// TODO 错误提示需要更精确
			return ErrGradeError("数据已存在")
		}
	}
	return nil
}

func (impl *GradeSVCImpl) List(sortField, sortOrder, name, class, subject *string, score *int,
	limit, cursor int) ([]model.Grade, *libsgorm.Page, error) {
	gradeDao := impl.newGradeDaoFunc(impl.db, impl.logger)

	scopes := impl.buildListScopes(sortField, sortOrder, name, class, subject, score)
	res, page, err := gradeDao.List(scopes, limit, cursor)
	if err != nil {
		impl.logger.Error("list failed", zap.Error(err))
		return nil, nil, ErrGradeInternalError
	}
	return res, page, nil
}

func (impl *GradeSVCImpl) buildListScopes(sortField, sortOrder, name, class, subject *string, score *int) *dao.GradeScope {
	scopes := dao.NewGradeScope()
	_sortOrder := "desc"
	if sortOrder != nil {
		_sortOrder = *sortOrder
	}

	_sortField := "created_at"
	if sortField != nil {
		_sortField = *sortField
	}
	scopes.AddSortBy(fmt.Sprintf("%s %s", _sortField, _sortOrder))

	if name != nil {
		scopes = scopes.AddName(*name)
	}

	if class != nil {
		scopes = scopes.AddClass(*class)
	}
	if score != nil {
		scopes = scopes.AddScore(*score)
	}
	if subject != nil {
		scopes = scopes.AddSubject(*subject)
	}

	return scopes
}
