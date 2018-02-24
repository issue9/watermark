// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package watermark

import (
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

	w, err = New("./testdata/watermark.png", 10, -1)
	a.Equal(err, ErrInvalidPos).Nil(w)

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
	a.False(IsAllowExt("gif"))
}
