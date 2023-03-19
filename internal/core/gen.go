package core

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/pkged"
	"github.com/sirupsen/logrus"
)

func Gen(genType string, specificPath string, overwrite bool) {

	// check if generate config file
	isGenConf := false
	if genType == "config" || genType == "conf" || genType == "artalk.yml" {
		isGenConf = true
		genType = "artalk.yml"
	}

	// get generation content
	var fileStr string
	if isGenConf {
		fileStr = GetConfTpl()
		// gen random `app_key`
		appKey := RandStringRunes(16)
		fileStr = strings.Replace(fileStr, `app_key: ""`, fmt.Sprintf(`app_key: "%s"`, appKey), 1)
	} else {
		fileStr = getEmbedFile(genType)
	}

	genFullPath := filepath.Base(genType) // generate file in work dir
	if specificPath != "" {
		genFullPath = specificPath
	}

	absPath, err := filepath.Abs(genFullPath)
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

	if _, err = dst.Write([]byte(fileStr)); err != nil {
		logrus.Fatal("Failed to write target file: ", err)
	}

	logrus.Info("File Generated: " + absPath)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*1234567890")

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

func getEmbedFile(filename string) string {
	file, err := pkged.FS().Open(strings.TrimPrefix(filename, "/"))
	if err != nil {
		logrus.Fatal("Invalid built-in resource `"+filename+"`: ", err)
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		logrus.Fatal("Read built-in resources `"+filename+"` error: ", err)
	}

	return string(buf)
}
