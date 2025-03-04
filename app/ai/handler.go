package main

import (
	"context"
	"github.com/trashwbin/dymall/app/ai/biz/service"
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct{}

// AnalyzeQuery implements the AiServiceImpl interface.
func (s *AiServiceImpl) AnalyzeQuery(ctx context.Context, req *ai.AnalyzeQueryRequest) (resp *ai.AnalyzeQueryResponse, err error) {
	resp, err = service.NewAnalyzeQueryService(ctx).Run(req)

	return resp, err
}
