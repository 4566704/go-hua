package common

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipTest() {
	//zip测试
	err := Zip("D:\\zip\\aaa", "d:/zip/aaa.zip")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = Unzip("d:/zip/aaa.zip", "D:\\zip\\ccc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Zip 打包成zip文件 可以是目录或单个文件
func Zip(srcFile string, destZip string) error {
	// 预防：旧文件无法覆盖
	os.RemoveAll(destZip)

	//创建zip文件
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	//打开zip文件
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	//遍历目录
	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		//println(path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+string(os.PathSeparator))
		println(header.Name)
		// header.Name = path
		if info.IsDir() {
			header.Name += "\\"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return err
	})

	return err
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		var decodeName string
		if f.Flags == 0 {
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := io.ReadAll(decoder)
			decodeName = string(content)
		} else {
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			decodeName = f.Name
		}
		fpath := filepath.Join(destDir, decodeName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
