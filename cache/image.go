package cache

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Vector2 struct {
	X uint
	Y uint
}

type Range struct {
	Left int
	Top int
	Width int
	Height int
}

type ImageInfo struct {
	Name string
	Path string
	Format string
	Data image.Image
}

func (mine *ImageInfo)ID() int {
	arr := strings.Split(mine.Name, "(")
	if arr == nil || len(arr) < 1 {
		return -1
	}
	msg := arr[1]
	inx := strings.Index(msg, ")")
	num,err := strconv.ParseInt(msg[:inx], 10, 32)
	if err != nil {
		return -1
	}
	return int(num)
}

func ClipImages(input, output string, wh, start Vector2)  {
	if !PathIsExist(input) {
		fmt.Println("找不到该文件夹："+ input)
		return
	}
	files, _ := ioutil.ReadDir(input)
	list := make([]*ImageInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			img, fm := readImage(filepath.Join(input, file.Name()))
			if img != nil {
				list = append(list, &ImageInfo{Name: file.Name(), Path: input, Format: fm, Data: img})
			}
		}
	}
	for i := 0;i < len(list);i += 1 {
		path := output + fmt.Sprintf("image-%d.jpg", i)
		clipImage(list[i].Data, Range{Left: int(start.X), Top: int(start.Y), Width: int(wh.X), Height: int(wh.Y)}, path, list[i].Format)
	}
}

func MergeImages(input, output, bg string,imgWH, bgWH Vector2, space uint) {
	if !PathIsExist(input) {
		fmt.Println("找不到该文件夹："+ input)
		return
	}
	bgImg,_ := readImage(bg)
	files, _ := ioutil.ReadDir(input)
	list := make([]*ImageInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			img,fm := readImage(filepath.Join(input, file.Name()))
			if img != nil {
				list = append(list, &ImageInfo{Name: file.Name(), Format: fm, Path: input, Data: img})
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].ID() > list[j].ID() {
			return false
		}else {
			return true
		}
	})
	max := len(list)
	length := max / 2
	if max % 2 != 0 {
		length = length + 1
	}
	for i := 0;i < length;i += 1{
		index := i * 2
		output := output + fmt.Sprintf("image_%d.jpg", i)
		if index == max - 1 {
			one := list[index]
			_ = mergeBothImage(output, one.Data, nil, bgImg, imgWH, bgWH, space)
		}else{
			one := list[index]
			two := list[index + 1]
			_ = mergeBothImage(output, one.Data, two.Data, bgImg, imgWH, bgWH, space)
		}
	}
}

func readImage(fullPath string) (image.Image, string) {
	reader, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("Impossible to open the file:", err)
		return nil,""
	}
	defer reader.Close()
	img, fm, err := image.Decode(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", fullPath, err)
		return nil,""
	}
	return img,fm
}

func mergeBothImage(output string, img1, img2, bg image.Image, imgWH, bgWH Vector2, space uint) error {
	if img1 == nil || bg == nil {
		return errors.New("the image is nil")
	}
	m0 := resize.Resize(bgWH.X, bgWH.Y, bg, resize.Lanczos3)
	m1 := resize.Resize(imgWH.X, imgWH.Y, img1, resize.Lanczos3)
	var m2 image.Image
	if img2 != nil {
		m2 = resize.Resize(imgWH.X, imgWH.Y, img2, resize.Lanczos3)
	}

	newImg := image.NewNRGBA(image.Rect(0, 0, int(bgWH.X), int(bgWH.Y)))
	draw.Draw(newImg, newImg.Bounds(), m0, m0.Bounds().Min, draw.Src)
	xx := int((bgWH.X - imgWH.X) / 2)
	yy := int((bgWH.Y - imgWH.Y * 2 - space) / 2)
	draw.Draw(newImg, newImg.Bounds(), m1, m1.Bounds().Min.Sub(image.Pt(xx, yy)), draw.Over)
	if m2 != nil {
		draw.Draw(newImg, newImg.Bounds(), m2, m2.Bounds().Min.Sub(image.Pt(xx, m1.Bounds().Max.Y + yy + int(space))), draw.Over)
	}

	return saveJPG(output, newImg)
}

func clipImage(data image.Image, rect Range, output, fm string) error {
	// canvas := resize.Resize(uint(rect.Width), uint(rect.Height), img, resize.Lanczos3)
	canvas := data
	switch fm {
	case "jpeg":
		img := canvas.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(rect.Left, rect.Top, rect.Width, rect.Height)).(*image.YCbCr)
		//buf := bytes.NewBuffer(nil)
		//_ = png.Encode(buf, subImg)
		//dist := base64.StdEncoding.EncodeToString(buf.Bytes())
		//fmt.Print(dist)
		return saveJPG(output, subImg)
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(rect.Left, rect.Top, rect.Width, rect.Height)).(*image.NRGBA)
			return savePNG(output, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(rect.Left, rect.Top, rect.Width, rect.Height)).(*image.RGBA)
			return savePNG(output, subImg)
		}
	case "gif":
		img := canvas.(*image.Paletted)
		subImg := img.SubImage(image.Rect(rect.Left, rect.Top, rect.Width, rect.Height)).(*image.Paletted)
		return saveGIF(output, subImg)
	case "bmp":
		img := canvas.(*image.RGBA)
		subImg := img.SubImage(image.Rect(rect.Left, rect.Top, rect.Width, rect.Height)).(*image.RGBA)
		return saveBMP(output, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

func fixImageSize(width, height float64) error {
	return nil
}

func saveJPG(path string, data image.Image) error {
	file,err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, data, &jpeg.Options{Quality: 100})
}

func savePNG(path string, data image.Image) error {
	file,err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, data)
}

func saveBMP(path string, data image.Image) error {
	file,err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return bmp.Encode(file, data)
}

func saveGIF(path string, data image.Image) error {
	file,err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.Encode(file, data, &gif.Options{})
}