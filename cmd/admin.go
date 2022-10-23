package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/lib/core"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "创建管理员账号",
	Long:  "根据提示创建管理员账号",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir)

		fmt.Println("--------------------------------")
		fmt.Println("         管理员账户创建         ")
		fmt.Println("--------------------------------")

		username, email, password, err := credentials()
		if err != nil {
			logrus.Fatal(err)
		}

		findUser := model.FindUser(username, email)
		if !findUser.IsEmpty() {
			findUser.SetPasswordEncrypt(password)
			if err := model.UpdateUser(&findUser); err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("该账户已存在，密码修改成功")
			return
		}

		user := model.User{
			Name:       username,
			Email:      email,
			IsAdmin:    true,
			BadgeName:  "管理员",
			BadgeColor: "#FF6C00",
		}
		user.SetPasswordEncrypt(password)

		if err := model.CreateUser(&user); err != nil {
			logrus.Fatal(err)
		}

		fmt.Println("--------------------------------")
		fmt.Println("  Name: " + username)
		fmt.Println("  Mail: " + email)
		fmt.Println("--------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
}

func credentials() (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Enter Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	if !lib.ValidateEmail(strings.TrimSpace(email)) {
		return "", "", "", errors.New("邮箱格式不合法")
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()
	fmt.Print("Retype Password: ")
	byteRePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))
	rePassword := strings.TrimSpace(string(byteRePassword))

	if rePassword != password {
		return "", "", "", errors.New("输入的密码不一致")
	}

	return strings.TrimSpace(username), strings.TrimSpace(email), password, nil
}
