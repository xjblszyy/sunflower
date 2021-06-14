// Code generated by goa v3.4.2, DO NOT EDIT.
//
// Score gRPC server types
//
// Command:
// $ goa gen sunflower/pkg/api/design -o pkg/api/

package server

import (
	scorepb "sunflower/pkg/api/gen/grpc/score/pb"
	score "sunflower/pkg/api/gen/score"

	goa "goa.design/goa/v3/pkg"
)

// NewScoreListPayload builds the payload of the "ScoreList" endpoint of the
// "Score" service from the gRPC request type.
func NewScoreListPayload(message *scorepb.ScoreListRequest) *score.ScoreListPayload {
	v := &score.ScoreListPayload{
		Cursor: int(message.Cursor),
		Limit:  int(message.Limit),
	}
	if message.SortField != "" {
		v.SortField = &message.SortField
	}
	if message.SortOrder != "" {
		v.SortOrder = &message.SortOrder
	}
	if message.Name != "" {
		v.Name = &message.Name
	}
	if message.Class != "" {
		v.Class = &message.Class
	}
	if message.Scores != 0 {
		scoresptr := int(message.Scores)
		v.Scores = &scoresptr
	}
	if message.Subject != "" {
		v.Subject = &message.Subject
	}
	return v
}

// NewScoreListResponse builds the gRPC response type from the result of the
// "ScoreList" endpoint of the "Score" service.
func NewScoreListResponse(result *score.ScoreListResult) *scorepb.ScoreListResponse {
	message := &scorepb.ScoreListResponse{
		Errcode: int32(result.Errcode),
		Errmsg:  result.Errmsg,
	}
	if result.NextCursor != nil {
		message.NextCursor = int32(*result.NextCursor)
	}
	if result.Total != nil {
		message.Total = int32(*result.Total)
	}
	if result.Data != nil {
		message.Data = make([]*scorepb.GradeResult, len(result.Data))
		for i, val := range result.Data {
			message.Data[i] = &scorepb.GradeResult{
				Id:        val.ID,
				Class:     val.Class,
				Name:      val.Name,
				Score:     int32(val.Score),
				Subject:   val.Subject,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
			}
		}
	}
	return message
}

// NewScoreDetailPayload builds the payload of the "ScoreDetail" endpoint of
// the "Score" service from the gRPC request type.
func NewScoreDetailPayload(message *scorepb.ScoreDetailRequest) *score.ScoreDetailPayload {
	v := &score.ScoreDetailPayload{
		ID: int(message.Id),
	}
	return v
}

// NewScoreDetailResponse builds the gRPC response type from the result of the
// "ScoreDetail" endpoint of the "Score" service.
func NewScoreDetailResponse(result *score.ScoreDetailResult) *scorepb.ScoreDetailResponse {
	message := &scorepb.ScoreDetailResponse{
		Errcode: int32(result.Errcode),
		Errmsg:  result.Errmsg,
	}
	if result.Data != nil {
		message.Data = svcScoreGradeResultToScorepbGradeResult(result.Data)
	}
	return message
}

// ValidateScoreListRequest runs the validations defined on ScoreListRequest.
func ValidateScoreListRequest(message *scorepb.ScoreListRequest) (err error) {
	if message.Limit > 100 {
		err = goa.MergeErrors(err, goa.InvalidRangeError("message.limit", message.Limit, 100, false))
	}
	if message.SortField != "" {
		if !(message.SortField == "class" || message.SortField == "name" || message.SortField == "score" || message.SortField == "subject" || message.SortField == "createdAt") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.sortField", message.SortField, []interface{}{"class", "name", "score", "subject", "createdAt"}))
		}
	}
	if message.SortOrder != "" {
		if !(message.SortOrder == "asc" || message.SortOrder == "desc") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.sortOrder", message.SortOrder, []interface{}{"asc", "desc"}))
		}
	}
	return
}

// svcScoreGradeResultToScorepbGradeResult builds a value of type
// *scorepb.GradeResult from a value of type *score.GradeResult.
func svcScoreGradeResultToScorepbGradeResult(v *score.GradeResult) *scorepb.GradeResult {
	if v == nil {
		return nil
	}
	res := &scorepb.GradeResult{
		Id:        v.ID,
		Class:     v.Class,
		Name:      v.Name,
		Score:     int32(v.Score),
		Subject:   v.Subject,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}

	return res
}

// protobufScorepbGradeResultToScoreGradeResult builds a value of type
// *score.GradeResult from a value of type *scorepb.GradeResult.
func protobufScorepbGradeResultToScoreGradeResult(v *scorepb.GradeResult) *score.GradeResult {
	if v == nil {
		return nil
	}
	res := &score.GradeResult{
		ID:        v.Id,
		Class:     v.Class,
		Name:      v.Name,
		Score:     int(v.Score),
		Subject:   v.Subject,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}

	return res
}
