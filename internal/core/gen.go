package core

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/sirupsen/logrus"
)

func Gen(genType string, specificPath string, overwrite bool) {
	// 参数
	if genType == "config" || genType == "conf" || genType == "artalk.yml" {
		genType = "artalk.example.yml"
	}

	genPath := filepath.Base(genType)
	if specificPath != "" {
		genPath = specificPath
	}

	file, err := pkged.FS().Open(strings.TrimPrefix(genType, "/"))
	if err != nil {
		logrus.Fatal("Invalid built-in resource `"+genType+"`: ", err)
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal("Read built-in resources `"+genType+"` error: ", err)
	}

	// 自动生成 app_key
	if strings.Contains(filepath.Base(genType), "artalk.example.yml") {
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
		logrus.Fatal(i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("File")}) + ": " + absPath)
	}

	dst, err := os.Create(absPath)
	if err != nil {
		logrus.Fatal("Failed to create target file: ", err)
	}
	defer dst.Close()

	if _, err = dst.Write(buf); err != nil {
		logrus.Fatal("Failed to write target file: ", err)
	}

	logrus.Info("File Generated: " + absPath)
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
