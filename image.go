package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

func resizeImage(imgType string, imgData []byte, width, height int) []byte {

	var srcImg image.Image
	var dstImg image.Image
	var out bytes.Buffer
	var err error

	switch imgType {
	case "image/png":
		if srcImg, err = png.Decode(bytes.NewBuffer(imgData)); nil == err {
			dstImg = imaging.Resize(srcImg, width, height, imaging.Lanczos)
			err = imaging.Encode(&out, dstImg, imaging.PNG)
		}
	case "image/jpeg", "image/jpg":
		if srcImg, err = jpeg.Decode(bytes.NewBuffer(imgData)); nil == err {
			dstImg = imaging.Resize(srcImg, width, height, imaging.Lanczos)
			err = imaging.Encode(&out, dstImg, imaging.JPEG)
		}
	case "image/gif":
		if srcImg, err = gif.Decode(bytes.NewBuffer(imgData)); nil == err {
			dstImg = imaging.Resize(srcImg, width, height, imaging.Lanczos)
			err = imaging.Encode(&out, dstImg, imaging.GIF)
		}
	default:
		err = fmt.Errorf("unsupported image type: %s", imgType)
	}

	// 不认识或者解码失败的直接返回原始数据
	if nil != err {
		return imgData
	}

	// 正常缩放后的
	return out.Bytes()
}
