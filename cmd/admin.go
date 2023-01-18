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
		fmt.Println(" " + i18n.T("Create an admin account"))
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

			logrus.Info(i18n.T("The account already exists and the password was changed successfully"))
			return
		}

		user := entity.User{
			Name:       username,
			Email:      email,
			IsAdmin:    true,
			BadgeName:  i18n.T("Admin"),
			BadgeColor: "#FF6C00",
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
	if !utils.ValidateEmail(strings.TrimSpace(email)) {
		return "", "", "", errors.New(i18n.T("Invalid email format"))
	}

	fmt.Print(i18n.T("Enter Password") + ": ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()
	fmt.Print(i18n.T("Retype Password") + ": ")
	byteRePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))
	rePassword := strings.TrimSpace(string(byteRePassword))

	if rePassword != password {
		return "", "", "", errors.New(i18n.T("Inconsistent password input"))
	}

	return strings.TrimSpace(username), strings.TrimSpace(email), password, nil
}
