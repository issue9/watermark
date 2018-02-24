watermark [![Build Status](https://travis-ci.org/issue9/watermark.svg?branch=master)](https://travis-ci.org/issue9/watermark)
======

处理上传文件，若是图片还可以设置水印。
```go
w, err := watermark.New("./path/to/watermark/file", 2, watermark.Center)
if err != nil{
    panic(err)
}

err = w.MarkFile("./path/to/file")
```


### 安装

```shell
go get github.com/issue9/watermark
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/watermark)
[![GoDoc](https://godoc.org/github.com/issue9/watermark?status.svg)](https://godoc.org/github.com/issue9/watermark)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
