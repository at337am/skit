package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GoogleTranslateAPIURL 是 Google 翻译的非官方 API 地址
const GoogleTranslateAPIURL = "https://translate.googleapis.com/translate_a/single"

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type GoogleTranslate struct {
	client HTTPClient
}

func NewGoogleTranslate(c HTTPClient) *GoogleTranslate {
	return &GoogleTranslate{client: c}
}

// Translate 使用结构体持有的 client 进行翻译
func (g *GoogleTranslate) Translate(text, sourceLang, targetLang string) (string, error) {
	requestURL := buildRequestURL(text, sourceLang, targetLang)

	// 2. 发送请求并获取数据
	body, err := fetchTranslationData(g.client, requestURL)
	if err != nil {
		return "", err
	}

	// 3. 解析响应并返回结果
	return parseTranslationResponse(body)
}

// buildRequestURL 根据翻译参数构造请求 URL
func buildRequestURL(text, sourceLang, targetLang string) string {
	return fmt.Sprintf(
		"%s?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		GoogleTranslateAPIURL,
		sourceLang,
		targetLang,
		url.QueryEscape(text),
	)
}

// fetchTranslationData 使用一个实现了 HTTPClient 接口的客户端发送请求
func fetchTranslationData(client HTTPClient, requestURL string) ([]byte, error) {
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("请求翻译 API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("翻译 API 返回错误状态: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

// parseTranslationResponse 解析 API 响应并提取翻译文本
func parseTranslationResponse(body []byte) (string, error) {
	var result []any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 响应失败: %w", err)
	}

	// 提取翻译文本
	if len(result) > 0 {
		if innerSlice, ok := result[0].([]any); ok && len(innerSlice) > 0 {
			if textSlice, ok := innerSlice[0].([]any); ok && len(textSlice) > 0 {
				if translated, ok := textSlice[0].(string); ok {
					return translated, nil
				}
			}
		}
	}

	return "", fmt.Errorf("无法从 API 响应中提取翻译文本")
}
