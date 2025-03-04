package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spotlight/types"
	"strconv"
	"time"
)

// getImageBatchInfo 根据提供的数字获取图像批次信息。返回ImageBatchInfo结构体指针，如果发生错误则返回错误信息。
//
//	num：要获取的图像批次数量
func GetImageBatchInfo(num int) (*types.ImageBatchInfo, error) {
	// 创建一个新的HTTP客户端。
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建一个新的HTTP GET请求。
	req, err := http.NewRequest("GET", "https://fd.api.iris.microsoft.com/v4/api/selection", nil)
	if err != nil {
		return nil, err
	}

	// 设置请求的头部。
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 验证 num 是否为有效数字并在 1-4 范围内
	switch {
	case num < 1:
		num = 1
	case num > 4:
		num = 4
	}

	// 将参数添加到请求的查询字符串中。
	q := req.URL.Query()
	q.Add("placement", "88000820")
	q.Add("country", "CN")
	q.Add("locale", "zh-CN")
	q.Add("fmt", "json")
	q.Add("bcnt", strconv.Itoa(num))
	req.URL.RawQuery = q.Encode()

	// 发送HTTP请求。
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保在函数返回前关闭响应体。
	defer resp.Body.Close()

	// 检查HTTP响应状态码是否为200（OK）。
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败HTTP状态码[%d]", resp.StatusCode)
	}

	// 读取HTTP响应体。
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 将响应体解析为 ImageBatchInfo 结构。
	var res *types.ImageBatchInfo
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	// 返回解析后的图像批次信息。
	return res, nil
}
