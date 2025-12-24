package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// DeepSeek API配置
const (
	DeepSeekAPIURL = "https://api.deepseek.com/chat/completions"
)

// AIRequest DeepSeek API请求结构
type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	Stream   bool        `json:"stream"`
}

// AIMessage 消息结构
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse DeepSeek API响应结构
type AIResponse struct {
	Choices []struct {
		Message AIMessage `json:"message"`
	} `json:"choices"`
}

// GenerateEvaluationScore 调用AI生成评价分数（员工自评）
// 参数：
//   - selfComment: 员工的自评内容/学习心得
//   - understanding: 理解程度（1-5）
//   - difficulty: 难度感受（1-5）
//   - satisfaction: 满意度（1-5）
//   - courseName: 课程名称（用于上下文）
//
// 返回：0-100的分数
func GenerateEvaluationScore(selfComment string, understanding, difficulty, satisfaction int, courseName string) (float64, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`你是一个专业的培训评估专家。请根据学员对课程《%s》的自评内容，评估其学习掌握程度，给出0-100的分数。

评分标准：
- 90-100分：深刻理解课程内容，能够举一反三，有独到见解
- 80-89分：很好地掌握了课程要点，能够应用所学知识
- 70-79分：基本掌握课程内容，理解程度良好
- 60-69分：了解课程大致内容，但理解不够深入
- 60分以下：未能掌握课程核心内容

学员自评内容：
%s

学员自评指标：
- 理解程度：%d/5
- 难度感受：%d/5
- 满意度：%d/5

请只返回一个0-100之间的整数分数，不要包含任何其他文字、符号或解释。`, 
		courseName, selfComment, understanding, difficulty, satisfaction)

	// 调用AI接口
	score, err := callDeepSeekAPI(prompt)
	if err != nil {
		// AI调用失败时，使用简单算法估算分数
		return calculateFallbackScore(selfComment, understanding, difficulty, satisfaction), nil
	}

	return score, nil
}

// GenerateTeacherScore 调用AI生成教师评分（讲师对学员的评价）
// 参数：
//   - teacherComment: 讲师的评价内容
//   - courseName: 课程名称
//
// 返回：0-100的分数
func GenerateTeacherScore(teacherComment string, courseName string) (float64, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`你是一个专业的培训评估专家。请根据讲师对学员在课程《%s》中表现的评价，给出0-100的分数。

评分标准：
- 90-100分：表现优异，积极参与，掌握全面
- 80-89分：表现良好，认真学习，掌握较好
- 70-79分：表现合格，基本掌握课程内容
- 60-69分：表现一般，需要加强学习
- 60分以下：表现不佳，未能掌握核心内容

讲师评价内容：
%s

请只返回一个0-100之间的整数分数，不要包含任何其他文字、符号或解释。`, 
		courseName, teacherComment)

	// 调用AI接口
	score, err := callDeepSeekAPI(prompt)
	if err != nil {
		// AI调用失败时，使用默认分数
		return 75.0, nil
	}

	return score, nil
}

// callDeepSeekAPI 调用DeepSeek API
func callDeepSeekAPI(prompt string) (float64, error) {
	// 获取API Key
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("DEEPSEEK_API_KEY未设置")
	}

	// 构建请求体
	reqBody := AIRequest{
		Model: "deepseek-chat",
		Messages: []AIMessage{
			{
				Role:    "system",
				Content: "你是一个专业的教育评估专家，擅长根据学习内容评估学员的掌握程度。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API返回错误: %s", string(body))
	}

	// 解析响应
	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(aiResp.Choices) == 0 {
		return 0, fmt.Errorf("API未返回结果")
	}

	// 提取分数（去除可能的非数字字符）
	scoreStr := strings.TrimSpace(aiResp.Choices[0].Message.Content)
	scoreStr = strings.Trim(scoreStr, "分.")
	
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析分数失败: %v, 原始内容: %s", err, scoreStr)
	}

	// 确保分数在0-100范围内
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return score, nil
}

// calculateFallbackScore AI调用失败时的备用评分算法
func calculateFallbackScore(selfComment string, understanding, difficulty, satisfaction int) float64 {
	// 基础分：根据内容长度（50-200字为最佳）
	contentLength := len([]rune(selfComment))
	var baseScore float64
	if contentLength < 50 {
		baseScore = 50
	} else if contentLength <= 200 {
		baseScore = 70
	} else if contentLength <= 500 {
		baseScore = 75
	} else {
		baseScore = 70
	}

	// 理解程度权重：30%
	understandingScore := float64(understanding) * 20 * 0.3

	// 难度感受权重：10%（难度越高，说明学习越认真）
	difficultyScore := float64(difficulty) * 20 * 0.1

	// 满意度权重：10%
	satisfactionScore := float64(satisfaction) * 20 * 0.1

	// 内容质量权重：50%
	contentScore := baseScore * 0.5

	totalScore := understandingScore + difficultyScore + satisfactionScore + contentScore

	// 确保分数在60-95范围内（备用算法不给极端分数）
	if totalScore < 60 {
		totalScore = 60
	} else if totalScore > 95 {
		totalScore = 95
	}

	return totalScore
}
