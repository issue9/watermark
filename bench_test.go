// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package watermark

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

// go1.10 BenchmarkWater_MakeImage_500xJPEG-4   	  200000	      5689 ns/op
func BenchmarkWater_MakeImage_500xJPEG(b *testing.B) {
	a := assert.New(b)

	copyBackgroundFile(a, "./testdata/output/bench.jpg", "./testdata/background.jpg")

	w, err := New("./testdata/watermark.jpg", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.jpg", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".jpg")
	}
}

// go1.10 BenchmarkWater_MakeImage_500xPNG-4    	  300000	      3540 ns/op
func BenchmarkWater_MakeImage_500xPNG(b *testing.B) {
	a := assert.New(b)

	copyBackgroundFile(a, "./testdata/output/bench.png", "./testdata/background.png")

	w, err := New("./testdata/watermark.png", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.png", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".png")
	}
}

// go1.10 BenchmarkWater_MakeImage_500xGIF-4    	  200000	      7020 ns/op
func BenchmarkWater_MakeImage_500xGIF(b *testing.B) {
	a := assert.New(b)

	copyBackgroundFile(a, "./testdata/output/bench.gif", "./testdata/background.gif")

	w, err := New("./testdata/watermark.gif", 10, TopLeft)
	a.NotError(err).NotNil(w)

	file, err := os.OpenFile("./testdata/output/bench.gif", os.O_RDWR, os.ModePerm)
	a.NotError(err).NotNil(file)
	defer file.Close()

	for i := 0; i < b.N; i++ {
		w.Mark(file, ".gif")
	}
}
