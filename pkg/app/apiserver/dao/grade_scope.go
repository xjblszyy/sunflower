package dao

import (
	"gorm.io/gorm"

	"sunflower/pkg/simple/client/db"
)

type GradeScope struct {
	scopes []db.Scope
}

func NewGradeScope() *GradeScope {
	return &GradeScope{
		scopes: []db.Scope{},
	}
}

// 获取scopes
func (s *GradeScope) Scopes() []db.Scope {
	return s.scopes
}

func (s *GradeScope) AddSortBy(sortBy string) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Order(sortBy)
	}
	s.scopes = append(s.scopes, query)
	return s
}

func (s *GradeScope) AddName(name string) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
	s.scopes = append(s.scopes, query)
	return s
}

func (s *GradeScope) AddScore(score int) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Where("score = ?", score)
	}
	s.scopes = append(s.scopes, query)
	return s
}

func (s *GradeScope) AddClass(class string) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Where("class = ?", class)
	}
	s.scopes = append(s.scopes, query)
	return s
}

func (s *GradeScope) AddSubject(subject string) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Where("subject = ?", subject)
	}
	s.scopes = append(s.scopes, query)
	return s
}

func (s *GradeScope) AddPK(pk int) *GradeScope {
	query := func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", pk)
	}
	s.scopes = append(s.scopes, query)
	return s
}
