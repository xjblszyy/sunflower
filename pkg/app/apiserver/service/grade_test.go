package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"sunflower/pkg/app/apiserver/dao"
	mock_dao "sunflower/pkg/app/apiserver/dao/mock"
	"sunflower/pkg/app/apiserver/model"
	dbClient "sunflower/pkg/simple/client/db"

	"sunflower/pkg/libs/testtools"
)

type GroupTestSuite struct {
	suite.Suite
	mock           sqlmock.Sqlmock
	mockedGradeDao *mock_dao.MockGradeDao
	svc            GradeSVC
}

func (s *GroupTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	dbClient.DB = testtools.MockedGORMDBForTest(s.T(), db)
	s.NoError(err)
	s.mock = mock
	ctrl := gomock.NewController(s.T())
	s.mockedGradeDao = mock_dao.NewMockGradeDao(ctrl)
	GradeDaoFunc := func(tx *gorm.DB, logger *zap.Logger) dao.GradeDao { return s.mockedGradeDao }
	logger := testtools.MockedZAPForTest(s.T())
	s.svc = NewGradeSVCImpl(dbClient.DB, logger, GradeDaoFunc)
}

// Run Test
func TestGroupTestSuite(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}

func (s *GroupTestSuite) TestGradeSVCImpl_Create() {
	s.T().Run("创建成绩-成功", func(t *testing.T) {
		var (
			class   = faker.Name()
			name    = faker.Name()
			score   = 90
			subject = faker.Name()
		)
		gs := []model.Grade{{
			Name:    name,
			Class:   class,
			Score:   score,
			Subject: subject,
		}}

		s.mockedGradeDao.EXPECT().CreateMany(gs).Return(nil)
		err := s.svc.Create(gs)
		s.NoError(err)
	})
}
