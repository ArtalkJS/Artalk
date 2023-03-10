package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/ArtalkJS/Artalk/internal/core"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create or edit an administrator account",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		core.LoadCore(cfgFile, workDir)

		fmt.Println("--------------------------------")
		fmt.Println(" " + i18n.T("Create admin account"))
		fmt.Println("--------------------------------")

		username, email, password, err := credentials()
		if err != nil {
			logrus.Fatal(err)
		}

		findUser := query.FindUser(username, email)
		if !findUser.IsEmpty() {
			findUser.SetPasswordEncrypt(password)
			if err := query.UpdateUser(&findUser); err != nil {
				logrus.Fatal(err)
			}

			logrus.Info(i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("Account")}) +
				", " + i18n.T("Password updated"))
			return
		}

		user := entity.User{
			Name:       username,
			Email:      email,
			IsAdmin:    true,
			BadgeName:  i18n.T("Admin"),
			BadgeColor: "#0083FF",
		}
		user.SetPasswordEncrypt(password)

		if err := query.CreateUser(&user); err != nil {
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

	fmt.Print(i18n.T("Enter {{name}}", map[string]interface{}{"name": i18n.T("Username")}) + ": ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	fmt.Print(i18n.T("Enter {{name}}", map[string]interface{}{"name": i18n.T("Email")}) + ": ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}
	if !utils.ValidateEmail(strings.TrimSpace(email)) {
		return "", "", "", errors.New("invalid email format")
	}

	fmt.Print(i18n.T("Enter {{name}}", map[string]interface{}{"name": i18n.T("Password")}) + ": ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()
	fmt.Print(i18n.T("Retype {{name}}", map[string]interface{}{"name": i18n.T("Password")}) + ": ")
	byteRePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))
	rePassword := strings.TrimSpace(string(byteRePassword))

	if rePassword != password {
		return "", "", "", errors.New("inconsistent password input")
	}

	return strings.TrimSpace(username), strings.TrimSpace(email), password, nil
}
