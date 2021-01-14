package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tool-image/cache"
)

func main() {
	checkDirections()
	stReader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入操作类型（1=裁剪图片， 2=合并图片）:")
	kind, err := stReader.ReadString('\n')
	if err != nil {
		fmt.Println("There were errors reading, exiting program." + kind)
		return
	}
	kind = clearString(kind)
	k := parseToInt(kind)

	if k == 2 {
		fmt.Printf("请输入画布宽:")
		bgWidth, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + bgWidth)
			return
		}
		fmt.Printf("请输入画布高:")
		bgHeight, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + bgHeight)
			return
		}
		fmt.Printf("请输入图片宽:")
		imgWidth, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + imgWidth)
			return
		}
		fmt.Printf("请输入图片高:")
		imgHeight, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + imgHeight)
			return
		}
		fmt.Printf("请输入图片之间的间隔高度:")
		space, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + space)
			return
		}
		bgWidth = clearString(bgWidth)
		bgHeight = clearString(bgHeight)
		imgWidth = clearString(imgWidth)
		imgHeight = clearString(imgHeight)
		space = clearString(space)
		fmt.Println(fmt.Sprintf("即将合并图片，画布尺寸：{width: %s, height:%s}, 图片尺寸：{width: %s, height: %s}, 间隔：%s ", bgWidth, bgHeight, imgWidth,imgHeight, space))

		cache.MergeImages("images/input", "images/output/", "images/bg.jpg", cache.Vector2{X: parseToInt(imgWidth), Y: parseToInt(imgHeight)},
			cache.Vector2{X: parseToInt(bgWidth), Y: parseToInt(bgHeight)}, parseToInt(space))
		fmt.Println("处理完成！！！")
	}else if k == 1 {
		fmt.Printf("请输入图片宽:")
		imgWidth, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + imgWidth)
			return
		}
		fmt.Printf("请输入图片高:")
		imgHeight, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + imgHeight)
			return
		}
		fmt.Printf("请输入裁剪图片顶边距:")
		top, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + top)
			return
		}
		fmt.Printf("请输入裁剪图片左边距:")
		left, err := stReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program." + left)
			return
		}
		imgWidth = clearString(imgWidth)
		imgHeight = clearString(imgHeight)
		top = clearString(top)
		left = clearString(left)
		fmt.Println(fmt.Sprintf("即将裁剪图片，图片边距：{top: %s, left:%s}, 图片尺寸：{width: %s, height: %s}", top, left, imgWidth,imgHeight))
		cache.ClipImages("images/input","images/output/", cache.Vector2{X: parseToInt(imgWidth), Y: parseToInt(imgHeight)},
		cache.Vector2{X: parseToInt(top), Y: parseToInt(left)})
		fmt.Println("处理完成！！！")
	}else{
		fmt.Println("不合法的操作类型.")
	}

	for {
		time.Sleep(time.Second * time.Duration(100))
	}
}

func clearString(msg string) string {
	msg = strings.Replace(msg,"\n", "", -1)
	msg = strings.Replace(msg,"\r", "", -1)
	return msg
}

func parseToInt(msg string) uint {
	num,err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return uint(num)
}

func checkDirections()  {
	input := "images/input"
	output := "images/output"
	if !cache.PathIsExist(input) {
		fmt.Println("相对路径下找不到源文件夹：/images/origin")
		return
	}
	if !cache.PathIsExist(output) {
		err := os.MkdirAll(output, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
