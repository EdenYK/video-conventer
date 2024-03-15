package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {

		file, _ := c.FormFile("video")
		log.Println(file.Filename)

		dst := "./" + file.Filename
		// 上传文件至指定的完整文件路径
		err := c.SaveUploadedFile(file, dst)
		fmt.Println("保存文件成功")
		if err != nil {
			return
		}

		cmd := exec.Command("./ffmpeg", "-i", dst, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "dist/output.m3u8")

		// 创建管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("无法创建标准输出管道:", err)
			return
		}

		// 将标准错误输出重定向到标准输出
		cmd.Stderr = cmd.Stdout

		// 启动命令
		if err := cmd.Start(); err != nil {
			fmt.Println("命令启动失败:", err)
			return
		}

		// 将命令的标准输出传递给标准输出和标准错误
		go io.Copy(os.Stdout, stdout)

		// 等待命令执行完成
		if err := cmd.Wait(); err != nil {
			fmt.Println("命令执行失败:", err)
			return
		}

		fmt.Println("音频格式转换和压缩完成！")

		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	})
	router.Run()
}
