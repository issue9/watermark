// SPDX-FileCopyrightText: 2018-2024 caixw
//
// SPDX-License-Identifier: MIT

package watermark

import (
	"os"
	"testing"

	"github.com/issue9/assert/v4"
)

func BenchmarkWater_MakeImage_500xJPEG(b *testing.B) {
	a := assert.New(b, false)

	copyBackgroundFile(a, "./testdata/output/bench.jpg", "./testdata/background.jpg")

	w, err := NewFromFile("./testdata/watermark.jpg", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.jpg", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".jpg")
	}
}

func BenchmarkWater_MakeImage_500xPNG(b *testing.B) {
	a := assert.New(b, false)

	copyBackgroundFile(a, "./testdata/output/bench.png", "./testdata/background.png")

	w, err := NewFromFile("./testdata/watermark.png", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.png", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".png")
	}
}

func BenchmarkWater_MakeImage_500xGIF(b *testing.B) {
	a := assert.New(b, false)

	copyBackgroundFile(a, "./testdata/output/bench.gif", "./testdata/background.gif")

	w, err := NewFromFile("./testdata/watermark.gif", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.gif", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".gif")
	}
}
