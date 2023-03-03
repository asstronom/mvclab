package consoleview

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/pkg/domain"
	"github.com/manifoldco/promptui"
)

func (app *CMDApp) auth() error {
	selectItems := []Handler{app.handlers["signup"], app.handlers["login"]}
	prompt := promptui.Select{
		Label:     "Select action",
		Items:     selectItems,
		Templates: handlerSelect,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error running prompt: %w", err)
	}
	app.Serve(selectItems[i])
	return nil
}

func (app *CMDApp) signup() error {
	usernamePrompt := promptui.Prompt{
		Label: "Enter username",
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting username: %w", err)
	}
	passwordPrompt := promptui.Prompt{
		Label: "Enter password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting for password: %w", err)
	}
	confirmPassword, err := passwordPrompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting for password confirm: %w", err)
	}
	if password != confirmPassword {
		fmt.Println("passwords don't match")
		app.Serve(app.handlers["signup"])
		return nil
	}

	id, err := app.service.CreateUser(context.Background(), &domain.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	log.Println("new user id:", id)
	token, err := app.service.GenerateToken(context.Background(), username, password)
	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}
	os.Setenv("authToken", token)
	app.Serve(app.handlers["profile"])
	return nil
}

func (app *CMDApp) login() error {
	usernamePrompt := promptui.Prompt{
		Label: "Enter username",
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting username: %w", err)
	}
	passwordPrompt := promptui.Prompt{
		Label: "Enter password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting for password: %w", err)
	}

	token, err := app.service.GenerateToken(context.Background(), username, password)
	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}
	os.Setenv("authToken", token)
	app.Serve(app.handlers["profile"])
	return nil
}

func (app *CMDApp) profile() error {
	token, ok := os.LookupEnv("authToken")
	if !ok {
		return fmt.Errorf("error: no authToken in env vars")
	}
	userid, err := app.service.ParseToken(context.Background(), token)
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}
	user, err := app.service.GetUser(context.Background(), userid)
	if err != nil {
		return fmt.Errorf("error getting user by id: %w", err)
	}

	fmt.Printf(`Username: %s
	First Name: %s
	LastName: %s
	Email: %s
	Status: %s
	Description: %s
	Do you want to edit profile?`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Status,
		user.Description)

	prompt := promptui.Prompt{
		Label:     "Edit",
		IsConfirm: true,
	}
	res, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error prompting for edit: %w", err)
	}
	if strings.ToLower(res) == "n" {
		return nil
	}

	app.Serve(app.handlers["editProfile"])
	return nil
}

func (app *CMDApp) editProfile() error {
	token, ok := os.LookupEnv("authToken")
	if !ok {
		return fmt.Errorf("error: no authToken in env vars")
	}
	userid, err := app.service.ParseToken(context.Background(), token)
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}
	var user domain.User
	prompt := promptui.Prompt{}
	prompt.Label = "First Name"
	user.FirstName, _ = prompt.Run()
	prompt.Label = "Last Name"
	user.LastName, _ = prompt.Run()
	prompt.Label = "Email"
	user.Email, _ = prompt.Run()
	prompt.Label = "Status"
	user.Status, _ = prompt.Run()
	prompt.Label = "Description"
	user.Description, _ = prompt.Run()
	user.ID = userid
	err = app.service.UpdateUserDetails(context.Background(), &user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	app.Serve(app.handlers["profile"])
	return nil
}
