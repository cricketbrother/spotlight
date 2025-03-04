package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"spotlight/api"
	"spotlight/types"
	"strings"
	"sync"
	"time"
)

// downloadImage 下载给定URL的图片并保存到指定路径。如果下载或保存图片过程中发生错误，则返回错误信息；否则返回nil。
//
//	url: 图片的URL地址
//	imagePath: 图片保存的路径
func downloadImage(url, imagePath string) error {
	log.Printf("开始下载图片[%s]", url)

	// 创建一个新的HTTP客户端。
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建一个新的HTTP GET请求。
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败[%v]", err)
	}

	// 设置请求的头部。
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 发送HTTP请求并获取响应。
	res, err := client.Do(req)
	if err != nil {
		// 如果发送请求失败，返回错误信息。
		return fmt.Errorf("发送请求失败[%v]", err)
	}
	// 确保在函数返回前关闭响应体。
	defer res.Body.Close()

	// 检查HTTP响应状态码是否为200（成功）。
	if res.StatusCode != http.StatusOK {
		// 如果响应状态码不是200，返回错误信息。
		return fmt.Errorf("请求失败HTTP状态码[%d]", res.StatusCode)
	}

	// 读取响应体中的所有数据。
	data, err := io.ReadAll(res.Body)
	if err != nil {
		// 如果读取响应失败，返回错误信息。
		return fmt.Errorf("读取响应失败[%v]", err)
	}

	// 将读取的数据保存到指定的图片路径。
	if err := os.WriteFile(imagePath, data, 0644); err != nil {
		// 如果保存图片失败，返回错误信息。
		return fmt.Errorf("保存图片失败[%v]", err)
	}

	// 记录图片下载完成的日志。
	log.Printf("图片下载完成[%s]\n", imagePath)

	// 图片下载并保存成功，返回nil。
	return nil
}

// downloadImageBatch 下载一批图片。如果下载过程中发生错误，则返回出错的数量。
//
//	imageBatchInfo: 包含待下载图片信息的结构体指针。
//	saveDir: 保存图片的目录路径。
func downloadImageBatch(imageBatchInfo *types.ImageBatchInfo, saveDir string) int {
	// 使用WaitGroup等待所有图片下载完成。
	var wg sync.WaitGroup
	// 使用Mutex保护errors切片的并发访问。
	var mu sync.Mutex
	// errors切片用于存储下载过程中发生的错误。
	var errors []error

	// 遍历图片信息，为每个图片启动一个goroutine进行下载。
	for _, item := range imageBatchInfo.Batchrsp.Items {
		wg.Add(1)
		go func(item types.Item) {
			defer wg.Done()

			// 解析图片信息。
			var imageInfo *types.ImageInfo
			err := json.Unmarshal([]byte(item.Item), &imageInfo)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("解析JSON失败[%v]", err))
				mu.Unlock()
				return
			}

			// 如果图片URL不为空，则下载图片。
			if imageInfo.AD.LandscapeImage.Asset != "" {
				imagePath := filepath.Join(saveDir, strings.ReplaceAll(filepath.Clean(fmt.Sprintf("%s_%s%s", imageInfo.AD.EntityID, imageInfo.AD.Title, filepath.Ext(imageInfo.AD.LandscapeImage.Asset))), "..", ""))
				if err := downloadImage(imageInfo.AD.LandscapeImage.Asset, imagePath); err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}
		}(item)
	}

	// 等待所有goroutine完成。
	wg.Wait()

	return len(errors)
}

// DownloadImages 下载指定数量的图片。
//
//	maxImages: 要下载的图片数量。
//	saveDir: 保存图片的目录路径。
func DownloadImages(maxImages int, saveDir string) int {
	errorNum := 0
	for i := 0; i < maxImages/4; i++ {
		// 获取图片批次信息,
		imageBatchInfo, err := api.GetImageBatchInfo(4)
		if err != nil {
			log.Printf("获取图片批次信息失败[%v]", err)
			errorNum += 4
			continue
		}

		// 下载图片批次
		errorNum += downloadImageBatch(imageBatchInfo, saveDir)
	}

	if maxImages%4 != 0 {
		imageBatchInfo, err := api.GetImageBatchInfo(maxImages % 4)
		if err != nil {
			log.Printf("获取图片批次信息失败[%v]", err)
			errorNum += maxImages % 4
			return errorNum
		}

		errorNum += downloadImageBatch(imageBatchInfo, saveDir)
	}

	return errorNum
}
