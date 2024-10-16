// SPDX-FileCopyrightText: 2018-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package watermark 提供一个简单的水印功能
package watermark

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 水印的位置
const (
	TopLeft Pos = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

var (
	// ErrUnsupportedWatermarkType 不支持的水印类型
	ErrUnsupportedWatermarkType = errors.New("不支持的水印类型")

	// ErrWatermarkTooLarge 当水印位置距离右下角的范围小于水印图片时，返回错误。
	ErrWatermarkTooLarge = errors.New("水印太大")
)

// 允许做水印的图片类型
var allowExts = []string{
	".gif", ".jpg", ".jpeg", ".png",
}

// Pos 表示水印的位置
type Pos int

// Watermark 用于给图片添加水印功能
//
// 目前支持 gif、jpeg 和 png 三种图片格式。
// 若是 gif 图片，则只取图片的第一帧；png 支持透明背景。
type Watermark struct {
	image   image.Image // 水印图片
	gifImg  *gif.GIF    // 如果是 GIF 图片，image 保存第一帧的图片， gifImg 保存全部内容
	padding int         // 水印留的边白
	pos     Pos         // 水印的位置
}

// NewFromFile 从文件声明一个 [Watermark] 对象
//
// path 为水印文件的路径；
// padding 为水印在目标图像上的留白大小；
// pos 水印的位置。
func NewFromFile(path string, padding int, pos Pos) (*Watermark, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return New(f, filepath.Ext(path), padding, pos)
}

// NewFromFS 从文件系统初始化 [Watermark] 对象
func NewFromFS(fsys fs.FS, path string, padding int, pos Pos) (*Watermark, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return New(f, filepath.Ext(path), padding, pos)
}

// New 声明 [Watermark] 对象
//
// r 为水印图片内容；
// ext 为水印图片的扩展名，会根据扩展名判断图片类型；
// padding 为水印在目标图像上的留白大小；
// pos 图片位置；
func New(r io.Reader, ext string, padding int, pos Pos) (w *Watermark, err error) {
	if pos < TopLeft || pos > Center {
		panic("无效的 pos 值")
	}

	var img image.Image
	var gifImg *gif.GIF
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(r)
	case ".png":
		img, err = png.Decode(r)
	case ".gif":
		gifImg, err = gif.DecodeAll(r)
		img = gifImg.Image[0]
	default:
		return nil, ErrUnsupportedWatermarkType
	}
	if err != nil {
		return nil, err
	}

	return &Watermark{
		image:   img,
		gifImg:  gifImg,
		padding: padding,
		pos:     pos,
	}, nil
}

// IsAllowExt 该扩展名的图片是否允许使用水印
//
// ext 必须带上 . 符号
func IsAllowExt(ext string) bool {
	if ext == "" {
		panic("参数 ext 不能为空")
	}

	if ext[0] != '.' {
		panic("参数 ext 必须以 . 开头")
	}

	ext = strings.ToLower(ext)

	for _, e := range allowExts {
		if e == ext {
			return true
		}
	}
	return false
}

// MarkFile 给指定的文件打上水印
func (w *Watermark) MarkFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return w.Mark(file, strings.ToLower(filepath.Ext(path)))
}

// Mark 将水印写入 src 中，由 ext 确定当前图片的类型。
func (w *Watermark) Mark(src io.ReadWriteSeeker, ext string) (err error) {
	var srcImg image.Image

	ext = strings.ToLower(ext)
	switch ext {
	case ".gif":
		return w.markGIF(src) // GIF 另外单独处理
	case ".jpg", ".jpeg":
		srcImg, err = jpeg.Decode(src)
	case ".png":
		srcImg, err = png.Decode(src)
	default:
		return ErrUnsupportedWatermarkType
	}
	if err != nil {
		return err
	}

	bound := srcImg.Bounds()
	point := w.getPoint(bound.Dx(), bound.Dy())

	if err = w.checkTooLarge(point, bound); err != nil {
		return err
	}

	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)

	if _, err = src.Seek(0, 0); err != nil {
		return err
	}

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(src, dstImg, nil)
	case ".png":
		return png.Encode(src, dstImg)
	default:
		return ErrUnsupportedWatermarkType
	}
}

func (w *Watermark) markGIF(src io.ReadWriteSeeker) error {
	srcGIF, err := gif.DecodeAll(src)
	if err != nil {
		return err
	}
	bound := srcGIF.Image[0].Bounds()
	point := w.getPoint(bound.Dx(), bound.Dy())

	if err = w.checkTooLarge(point, bound); err != nil {
		return err
	}

	if w.gifImg == nil {
		for index, img := range srcGIF.Image {
			dstImg := image.NewPaletted(img.Bounds(), img.Palette)
			draw.Draw(dstImg, dstImg.Bounds(), img, image.Point{}, draw.Src)
			draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)
			srcGIF.Image[index] = dstImg
		}
	} else { // 水印也是 GIF
		windex := 0
		wmax := len(w.gifImg.Image)
		for index, img := range srcGIF.Image {
			dstImg := image.NewPaletted(img.Bounds(), img.Palette)
			draw.Draw(dstImg, dstImg.Bounds(), img, image.Point{}, draw.Src)

			// 获取对应帧数的水印图片
			if windex >= wmax {
				windex = 0
			}
			draw.Draw(dstImg, dstImg.Bounds(), w.gifImg.Image[windex], point, draw.Over)
			windex++

			srcGIF.Image[index] = dstImg
		}
	}

	if _, err = src.Seek(0, 0); err != nil {
		return err
	}
	return gif.EncodeAll(src, srcGIF)
}

func (w *Watermark) checkTooLarge(start image.Point, dst image.Rectangle) error {
	// 允许的最大高宽
	width := dst.Dx() - start.X - w.padding
	height := dst.Dy() - start.Y - w.padding

	if width < w.image.Bounds().Dx() || height < w.image.Bounds().Dy() {
		return ErrWatermarkTooLarge
	}
	return nil
}

func (w *Watermark) getPoint(width, height int) image.Point {
	var point image.Point

	switch w.pos {
	case TopLeft:
		point = image.Point{X: -w.padding, Y: -w.padding}
	case TopRight:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -w.padding,
		}
	case BottomLeft:
		point = image.Point{
			X: -w.padding,
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case BottomRight:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case Center:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()) / 2,
			Y: -(height - w.padding - w.image.Bounds().Dy()) / 2,
		}
	default:
		panic("无效的 pos 值")
	}

	return point
}
