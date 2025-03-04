package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/checkout/infra/rpc"
	ai "github.com/trashwbin/dymall/rpc_gen/kitex_gen/ai"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
	"os"
	"regexp"
)

type AnalyzeQueryService struct {
	ctx context.Context
}

// NewAnalyzeQueryService 创建一个新的 AnalyzeQueryService 实例
func NewAnalyzeQueryService(ctx context.Context) *AnalyzeQueryService {
	return &AnalyzeQueryService{ctx: ctx}
}

// Run 处理 AnalyzeQueryRequest 并返回 AnalyzeQueryResponse
func (s *AnalyzeQueryService) Run(req *ai.AnalyzeQueryRequest) (resp *ai.AnalyzeQueryResponse, err error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助分析文本中的内容，提取出用户要查询的物品的名称。 以 大模型所查询到的物品是“xx” 的格式输出  "),

		// 插入需要的对话历史（新对话的话这里不填）
		//schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(s.ctx, map[string]any{
		"role":     "语言分析师",
		"style":    "严谨且专业",
		"question": req.UserInput, // 使用请求中的用户输入
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("模板格式化失败: %v", err)
	}

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("未能加载.env文件: %v", err)
	}

	// 创建聊天模型
	chatModel, err := openai.NewChatModel(s.ctx, &openai.ChatModelConfig{
		Model:   "deepseek-chat",             // 使用的模型版本
		APIKey:  os.Getenv("OPENAI_API_KEY"), // 从环境变量获取 API 密钥
		BaseURL: "https://api.deepseek.com/v1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat model: %v", err)
	}

	// 生成聊天响应
	result, err := chatModel.Generate(s.ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("生成失败: %v", err)
	}
	fmt.Println("result的内容:", result)

	//读取result内容 并且将目标字段存入数组
	re := regexp.MustCompile(`“([^”]+)”`)
	matches := re.FindAllStringSubmatch(result.Content, -1)
	var items []string
	for _, match := range matches {
		if len(match) > 1 {
			items = append(items, match[1])
		}
	}
	if len(items) > 0 {
		fmt.Println("提取的内容:", items) // 输出: [男士内裤 袜子]
	} else {
		fmt.Println("未找到双引号内容")
	}

	for _, item := range items {
		// 创建 GetOrderReq 结构体实例
		getOrderReq := &order.GetOrderReq{
			OrderId: item,      // 使用 items 中的每个元素
			UserId:  uint32(1), // 确保 UserId 是 uint32 类型 //TODO
		}

		// 调用订单服务查询订单
		orderResp, err := rpc.OrderClient.GetOrder(s.ctx, getOrderReq)
		if err != nil {
			fmt.Printf("订单 %s 查询失败: %v\n", item, err)
			continue
		}

		fmt.Printf("订单数据: %+v\n", orderResp)
		// 返回第一个找到的订单
		return &ai.AnalyzeQueryResponse{
			OrderId:     orderResp.Order.OrderId,
			ProductName: "", //TODO 没有产品名字
			Intent:      "",
			Success:     true,
			Message:     "订单查询成功",
		}, nil
	}

	return &ai.AnalyzeQueryResponse{
		OrderId:     "",
		ProductName: "",
		Intent:      "",
		Success:     false,
		Message:     "未找到订单",
	}, nil
}
