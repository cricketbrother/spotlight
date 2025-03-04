package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"spotlight/dir"
	"spotlight/downloader"
)

// main函数是程序的入口点。
func main() {
	var maxImages int
	var saveDir string
	var clean string

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "下载Windows聚焦壁纸")
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
	}

	// 定义命令行参数
	flag.IntVar(&maxImages, "n", 4, "每次下载的图片数量")
	flag.StringVar(&saveDir, "d", "images", "保存图片的目录")
	flag.StringVar(&clean, "c", "yes", "是否清空保存图片的目录")

	// 解析命令行参数
	flag.Parse()

	// 打印本次下载的图片数量和保存目录。
	log.Printf("本次将下载[%d]张图片到[%s]目录", maxImages, saveDir)

	// 如果需要清空保存图片的目录，则删除该目录里的所有文件和子目录，但不删除目录本身。
	if clean == "yes" {
		if err := dir.Clean(saveDir); err != nil {
			log.Printf("清空保存图片目录失败[%v]", err)
		}
	}

	// 创建保存图片的目录。
	if err := dir.Create(saveDir); err != nil {
		// 如果创建保存图片的目录失败，则记录错误并退出程序。
		log.Fatalf("创建保存图片的目录失败[%v]", err)
	}

	// 下载图片
	errorNum := downloader.DownloadImages(maxImages, saveDir)

	// 打印下载结果。
	log.Printf("下载成功[%d]张图片, 失败[%d]张图片", maxImages-errorNum, errorNum)
}
