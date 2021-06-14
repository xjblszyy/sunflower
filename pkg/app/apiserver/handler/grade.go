package handler

import (
	"context"
	"io"
	"mime/multipart"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"sunflower/pkg/api/gen/log"
	"sunflower/pkg/api/gen/score"
	"sunflower/pkg/app/apiserver/dao"
	"sunflower/pkg/app/apiserver/model"
	"sunflower/pkg/app/apiserver/serializer"
	"sunflower/pkg/app/apiserver/service"
	dbClient "sunflower/pkg/simple/client/db"
)

type scoresrvc struct {
	logger *log.Logger
}

func (s *scoresrvc) ScoreList(ctx context.Context, p *score.ScoreListPayload) (res *score.ScoreListResult, err error) {
	res = &score.ScoreListResult{}
	logger := L(ctx, s.logger)
	logger.Info("score.ScoreList")

	svc := ScoreSVC(dbClient.DB, logger)
	data, page, err := svc.List(p.SortField, p.SortOrder, p.Name, p.Class, p.Subject, p.Scores, p.Limit, p.Cursor)
	if err != nil {
		return res, err
	}

	res.Data = serializer.ModelGradesToGoa(data)
	res.Total = &page.TotalRecord
	res.NextCursor = &page.NextCursor
	return res, nil
}

func (s *scoresrvc) ScoreDetail(ctx context.Context, p *score.ScoreDetailPayload) (res *score.ScoreDetailResult, err error) {
	res = &score.ScoreDetailResult{}
	logger := L(ctx, s.logger)
	logger.Info("score.ScoreDetail")

	svc := ScoreSVC(dbClient.DB, logger)
	data, err := svc.FetchOne(p.ID)
	if err != nil {
		return res, err
	}

	res.Data = serializer.ModelGradeTOGoa(data)
	return res, nil
}

func (s *scoresrvc) Upload(ctx context.Context, p *score.UploadPayload, req io.ReadCloser) (res *score.UploadResult, err error) {
	res = &score.UploadResult{}
	logger := L(ctx, s.logger)
	logger.Info("score.Upload")

	files, err := multipartFile(ctx, req, p.ContentType)
	if err != nil {
		logger.Error("multipart file failed", zap.Error(err))
		return res, err
	}

	for _, file := range files {
		tx := dbClient.DB.Begin()
		svc := ScoreSVC(tx, logger)
		grades, err := s.buildGrade(file)
		if err != nil {
			tx.Rollback()
			return res, err
		}

		if _err := svc.Create(grades); _err != nil {
			tx.Rollback()
			return res, _err
		}

		tx.Commit()
	}

	data := score.SuccessResult{
		OK: true,
	}
	res.Data = &data
	return res, nil
}

func (s *scoresrvc) buildGrade(file multipart.File) ([]model.Grade, error) {
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		s.logger.Error("open from reader failed", zap.Error(err))
		return nil, service.ErrGradeInternalError
	}

	// 获取工作簿表格名称
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		s.logger.Error("get rows failed", zap.Error(err))
		return nil, service.ErrGradeInternalError
	}
	res := make([]model.Grade, 0, len(rows))

	for i := 1; i < len(rows); i++ {
		sc, err := strconv.Atoi(rows[i][2])
		if err != nil {
			s.logger.Error("string to int failed", zap.Error(err))
			return nil, service.ErrGradeInternalError
		}

		grade := model.Grade{
			Class:   rows[i][0],
			Name:    rows[i][1],
			Score:   sc,
			Subject: rows[i][3],
		}

		res = append(res, grade)
	}
	return res, nil
}

// NewUser returns the User service implementation.
func NewScore(logger *log.Logger) score.Service {
	return &scoresrvc{
		logger: logger,
	}
}

var ScoreSVC = func(db *gorm.DB, logger *zap.Logger) service.GradeSVC {
	return service.NewGradeSVCImpl(db, logger, dao.NewGradeDaoImpl)
}
