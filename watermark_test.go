// SPDX-License-Identifier: MIT

package watermark

import (
	"image"
	"io"
	"os"
	"testing"

	"github.com/issue9/assert"
)

// 复制文件到 output 目录下，并重命名。
func copyBackgroundFile(a *assert.Assertion, dest, src string) {
	destFile, err := os.Create(dest)
	a.NotError(err).NotNil(destFile)

	srcFile, err := os.Open(src)
	a.NotError(err).NotNil(srcFile)

	n, err := io.Copy(destFile, srcFile)
	a.NotError(err).True(n >= 0)

	destFile.Close()
	srcFile.Close()
}

// 输出各种组合的水印图片。
// bgExt 表示背景图片的扩展名。
// water 表示水印图片的扩展名。
func output(a *assert.Assertion, pos Pos, bgExt, waterExt string) {
	water := "./testdata/watermark" + waterExt
	src := "./testdata/background" + bgExt
	dest := "./testdata/output/" + waterExt[1:] + bgExt

	copyBackgroundFile(a, dest, src)

	// 添加水印
	w, err := New(water, 10, pos)
	a.NotError(err).NotNil(w)
	a.NotError(w.MarkFile(dest))
}

func TestNew(t *testing.T) {
	a := assert.New(t)

	w, err := New("./testdata/watermark.unsupported", 10, TopLeft)
	a.Equal(err, ErrUnsupportedWatermarkType).Nil(w)

	a.Panic(func() {
		w, err = New("./testdata/watermark.png", 10, -1)
	})

	src := "./testdata/background.unsupported"
	dest := "./testdata/output/unsupported.unsupported"
	copyBackgroundFile(a, dest, src)

	w, err = New("./testdata/watermark.png", 10, TopLeft)
	a.NotError(err).NotNil(w)
	err = w.MarkFile(dest)
	a.Equal(err, ErrUnsupportedWatermarkType)
}

func TestWatermark_MarkFile(t *testing.T) {
	a := assert.New(t)

	output(a, TopLeft, ".jpg", ".jpg")
	output(a, TopRight, ".jpg", ".png")
	output(a, Center, ".jpg", ".gif")

	output(a, BottomLeft, ".png", ".jpg")
	output(a, BottomRight, ".png", ".png")
	output(a, Center, ".png", ".gif")

	output(a, BottomLeft, ".gif", ".jpg")
	output(a, BottomRight, ".gif", ".png")
	output(a, Center, ".gif", ".gif")
}

func TestIsAllowExt(t *testing.T) {
	a := assert.New(t)

	a.True(IsAllowExt(".jpg"))
	a.True(IsAllowExt(".JPeG"))
	a.True(IsAllowExt(".png"))
	a.True(IsAllowExt(".Gif"))

	a.Panic(func() { IsAllowExt("") })
	a.Panic(func() { IsAllowExt("gif") })
}

func TestWater_checkTooLarge(t *testing.T) {
	a := assert.New(t)

	w, err := New("./testdata/watermark.png", 10, BottomRight)
	a.NotError(err).NotNil(w)
	dst := image.Rect(0, 0, w.image.Bounds().Dx(), w.image.Bounds().Dy())
	a.Equal(w.checkTooLarge(image.Point{X: 0, Y: 0}, dst), ErrWatermarkTooLarge)

	// padding 为 0 正好 1：1 覆盖
	w.padding = 0
	dst = image.Rect(0, 0, w.image.Bounds().Dx(), w.image.Bounds().Dy())
	a.NotError(w.checkTooLarge(image.Point{X: 0, Y: 0}, dst))
}
