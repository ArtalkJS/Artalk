package core

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/pkged"
	"github.com/sirupsen/logrus"
)

func Gen(genType string, specificPath string, overwrite bool) {
	// 参数
	if genType == "config" || genType == "conf" {
		genType = "artalk-go.example.yml"
	}

	genPath := filepath.Base(genType)
	if specificPath != "" {
		genPath = specificPath
	}

	file, err := pkged.Open("/" + strings.TrimPrefix(genType, "/"))
	if err != nil {
		logrus.Fatal("无效的内置资源: "+genType+" ", err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal("读取内置资源: "+genType+" 失败 ", err)
	}

	// 自动生成 app_key
	if strings.Contains(filepath.Base(genType), "artalk-go.example.yml") {
		str := string(buf)
		appKey := RandStringRunes(16)
		str = strings.Replace(str, `app_key: ""`, fmt.Sprintf(`app_key: "%s"`, appKey), 1)
		buf = []byte(str)
	}

	absPath, err := filepath.Abs(genPath)
	if err != nil {
		logrus.Fatal(err)
	}
	if s, err := os.Stat(absPath); err == nil && s.IsDir() {
		absPath = filepath.Join(absPath, filepath.Base(genType))
	}

	if CheckFileExist(absPath) && !overwrite {
		logrus.Fatal("已存在文件: " + absPath)
	}

	dst, err := os.Create(absPath)
	if err != nil {
		logrus.Fatal("创建目标文件失败 ", err)
	}
	defer dst.Close()

	if _, err = dst.Write(buf); err != nil {
		logrus.Fatal("写入目标文件失败 ", err)
	}

	logrus.Info("生成文件: " + absPath)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
