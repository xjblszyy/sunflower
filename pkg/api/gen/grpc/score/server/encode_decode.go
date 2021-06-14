// Code generated by goa v3.4.2, DO NOT EDIT.
//
// Score gRPC server encoders and decoders
//
// Command:
// $ goa gen sunflower/pkg/api/design -o pkg/api/

package server

import (
	"context"
	scorepb "sunflower/pkg/api/gen/grpc/score/pb"
	score "sunflower/pkg/api/gen/score"

	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeScoreListResponse encodes responses from the "Score" service
// "ScoreList" endpoint.
func EncodeScoreListResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*score.ScoreListResult)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Score", "ScoreList", "*score.ScoreListResult", v)
	}
	resp := NewScoreListResponse(result)
	return resp, nil
}

// DecodeScoreListRequest decodes requests sent to "Score" service "ScoreList"
// endpoint.
func DecodeScoreListRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *scorepb.ScoreListRequest
		ok      bool
	)
	{
		if message, ok = v.(*scorepb.ScoreListRequest); !ok {
			return nil, goagrpc.ErrInvalidType("Score", "ScoreList", "*scorepb.ScoreListRequest", v)
		}
		if err := ValidateScoreListRequest(message); err != nil {
			return nil, err
		}
	}
	var payload *score.ScoreListPayload
	{
		payload = NewScoreListPayload(message)
	}
	return payload, nil
}

// EncodeScoreDetailResponse encodes responses from the "Score" service
// "ScoreDetail" endpoint.
func EncodeScoreDetailResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*score.ScoreDetailResult)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Score", "ScoreDetail", "*score.ScoreDetailResult", v)
	}
	resp := NewScoreDetailResponse(result)
	return resp, nil
}

// DecodeScoreDetailRequest decodes requests sent to "Score" service
// "ScoreDetail" endpoint.
func DecodeScoreDetailRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *scorepb.ScoreDetailRequest
		ok      bool
	)
	{
		if message, ok = v.(*scorepb.ScoreDetailRequest); !ok {
			return nil, goagrpc.ErrInvalidType("Score", "ScoreDetail", "*scorepb.ScoreDetailRequest", v)
		}
	}
	var payload *score.ScoreDetailPayload
	{
		payload = NewScoreDetailPayload(message)
	}
	return payload, nil
}
