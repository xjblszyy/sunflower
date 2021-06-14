package serializer

import (
	"strconv"
	"time"

	"sunflower/pkg/api/gen/score"
	"sunflower/pkg/app/apiserver/model"
)

func ModelGradeTOGoa(data model.Grade) *score.GradeResult {
	res := score.GradeResult{
		// ID
		ID: strconv.Itoa(data.ID),
		// 班级
		Class: data.Class,
		// 姓名
		Name: data.Name,
		// 得分
		Score: data.Score,
		// 科目
		Subject: data.Subject,
	}
	if !data.CreatedAt.IsZero() {
		res.CreatedAt = data.CreatedAt.Format(time.RFC3339)
	}
	if !data.UpdatedAt.IsZero() {
		res.UpdatedAt = data.UpdatedAt.Format(time.RFC3339)
	}
	return &res
}

func ModelGradesToGoa(data []model.Grade) []*score.GradeResult {
	res := make([]*score.GradeResult, 0, len(data))
	for i := 0; i < len(data); i++ {
		content := ModelGradeTOGoa(data[i])
		res = append(res, content)
	}
	return res

}
