package service

import (
	"context"
	ai "github.com/trashwbin/dymall/rpc_gen/kitex_gen/ai"
	"testing"
)

func TestAnalyzeQuery_Run(t *testing.T) {
	/*	ctx := context.Background()
		s := NewAnalyzeQueryService(ctx)
		// init req and assert value

		req := &ai.AnalyzeQueryRequest{}
		resp, err := s.Run(req)
		t.Logf("err: %v", err)
		t.Logf("resp: %v", resp)
	*/
	// todo: edit your unit test
	// 创建上下文
	ctx := context.Background()

	// 创建 AnalyzeQueryService实例
	service := NewAnalyzeQueryService(ctx)

	// 定义测试用例
	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "Test Case 1: General Query",
			//input: "我的代码一直报错，感觉好沮丧，该怎么办？",
			input: "查询我的衣服订单",
		},
		{
			name:  "Test Case 2: Positive Feedback",
			input: "查询我前天买的内裤的订单",
		},
		{
			name:  "Test Case 3: Technical Question",
			input: "找一下前天我躲被窝里面偷偷下单的男士内裤",
		},
		{
			name:  "Test Case 3: Technical Question",
			input: "找一下前天我躲被窝里面偷偷用手机下单的男士内裤，袜子",
		},
	}

	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建请求对象
			req := &ai.AnalyzeQueryRequest{
				UserInput: tc.input,
			}

			// 调用 Run 方法
			resp, err := service.Run(req)
			if err != nil {
				t.Fatalf("Run failed: %v", err)
			}

			// 打印响应结果
			t.Logf("Response for %s: %+v", tc.name, resp)
		})
	}
}
