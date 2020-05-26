package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

var defaultLevel = gzip.DefaultCompression

// SetLevel 设置压缩级别
func SetLevel(level int) {
	defaultLevel = level
}

// GzipCompress gzip 压缩
func GzipCompress(in []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer, err := gzip.NewWriterLevel(&buffer, defaultLevel)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = writer.Write(in)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GzipUncompress gzip 解压
func GzipUncompress(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}
