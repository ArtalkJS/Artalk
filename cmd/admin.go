package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewAdminCommand(app *ArtalkCmd) *cobra.Command {
	adminCmd := &cobra.Command{
		Use:   "admin",
		Short: "Create or edit an administrator account",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("--------------------------------")
			fmt.Println(" " + i18n.T("Create admin account"))
			fmt.Println("--------------------------------")

			// get from flags
			username, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")

			if username == "" || email == "" || password == "" {
				var err error
				username, email, password, err = credentials()
				if err != nil {
					log.Fatal(err)
				}
			}

			findUser := app.Dao().FindUser(username, email)
			if !findUser.IsEmpty() {
				findUser.SetPasswordEncrypt(password)
				findUser.IsAdmin = true
				if err := app.Dao().UpdateUser(&findUser); err != nil {
					log.Fatal(err)
				}

				log.Info(i18n.T("{{name}} already exists", map[string]interface{}{"name": i18n.T("Account")}) +
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

			if err := app.Dao().CreateUser(&user); err != nil {
				log.Fatal(err)
			}

			fmt.Println("--------------------------------")
			fmt.Println("  Name: " + username)
			fmt.Println("  Mail: " + email)
			fmt.Println("--------------------------------")
		},
	}

	flag(adminCmd, "name", "", i18n.T("Username"))
	flag(adminCmd, "email", "", i18n.T("Email"))
	flag(adminCmd, "password", "", i18n.T("Password"))

	return adminCmd
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
		return "", "", "", fmt.Errorf("invalid email format")
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
		return "", "", "", fmt.Errorf("inconsistent password input")
	}

	return strings.TrimSpace(username), strings.TrimSpace(email), password, nil
}
