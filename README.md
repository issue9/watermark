watermark
[![Build Status](https://travis-ci.org/issue9/watermark.svg?branch=master)](https://travis-ci.org/issue9/watermark)
======

watermark 提供了简单的图片水印处理功能。支持处理 GIF、PNG 和 JPEG，水印也只支持这些类型的文件。

对于 GIF 水印，若被渲染图片为非 GIF 图片，则只取水印的第一帧作为水印内容；
若被渲染图片也是 GIF，则会将被渲染图片的第一帧与水印的第一帧合并，
水印的第二帧与被渲染图片的第二帧合并，依次类推。水印帧数不够的，则循环使用，
直到被渲染图片的帧数用完。

```go
w, err := watermark.New("./path/to/watermark/file", 2, watermark.Center)
if err != nil{
    panic(err)
}

err = w.MarkFile("./path/to/file")
```

安装
----

```shell
go get github.com/issue9/watermark
```

文档
----

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/watermark)
[![GoDoc](https://godoc.org/github.com/issue9/watermark?status.svg)](https://godoc.org/github.com/issue9/watermark)

版权
----

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
