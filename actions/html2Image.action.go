package actions

import (
	"fmt"
	"github.com/itskingori/go-wkhtml/wkhtmltox"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
)

var baseDir, _ = os.Getwd()

const (
	// 模版名称
	modName = "image_mod"
	// Html文件
	htmlFile = "views/word.html"
)

func ConvertWordsToImage(requestId string, words []string) ([]byte, error) {

	tpl, err := template.New(modName).
		Funcs(template.FuncMap{"unescaped": Unescaped}).
		ParseFiles(fmt.Sprintf("%s/%s", baseDir, htmlFile))
	if err != nil {
		fmt.Printf("template.New err: %v", err)
		return nil, err
	}

	dstModHtml := GetFile(requestId, "html")
	//defer RemoveFile(dstModHtml)

	dstFile, err := os.Create(dstModHtml)
	if err != nil {
		fmt.Printf("os.Create err: %v", err)
		return nil, err
	}

	err = tpl.ExecuteTemplate(
		dstFile,
		modName,
		map[string]interface{}{
			"List":      words,
			"LastIndex": len(words) - 1,
		},
	)
	if err != nil {
		fmt.Printf("tpl.ExecuteTemplate err: %v", err)
		return nil, err
	}

	format := "png"
	quality := 100
	width := 610
	opts := wkhtmltox.ImageOptions{
		// 图片格式
		Format: &format,
		// 图片质量
		Quality: &quality,
		// 图片宽度
		Width: &width,
	}
	ifs := wkhtmltox.NewImageFlagSetFromOptions(&opts)

	dstImage := GetFile(requestId, format)
	//defer RemoveFile(dstImage)
	output, _ := ifs.Generate(dstModHtml, dstImage)
	fmt.Println(output)

	return output, nil
}

// Unescaped html标签不转义
func Unescaped(x string) interface{} {
	return template.HTML(x)
}

// GetFile 获取文件地址
func GetFile(fileName, ext string) string {
	path := fmt.Sprintf("%s/storage/", baseDir)
	fileInfo, _ := os.Stat(path)
	if fileInfo == nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logrus.WithFields(logrus.Fields{"fileName": fileName, "ext": ext}).Error(err)
		}
	}
	return path + fileName + "." + ext
}

// RemoveFile 删除文件
func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		logrus.WithField("filePath", filePath).Error(err)
	}
}
