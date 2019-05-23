package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"github.com/winded/tyomaa/backend/db"
)

func commands() []cli.Command {
	return []cli.Command{
		{
			Name: "users",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "List users",
					Action: func(c *cli.Context) error {
						var users []db.User
						if err := db.Instance.Find(&users).Error; err != nil {
							return err
						}

						fmt.Println("Listing " + strconv.Itoa(len(users)) + " users")
						fmt.Println("--------------------------------")

						for _, user := range users {
							fmt.Println(user.Name)
						}

						return nil
					},
				},
				{
					Name:      "create",
					Usage:     "Create new user",
					ArgsUsage: "<name> <password>",
					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							return errors.New("Missing name or password argument")
						}

						var user db.User
						user.Name = c.Args().Get(0)
						err := user.SetPassword(c.Args().Get(1))
						if err != nil {
							return err
						}

						if err = db.Instance.Save(&user).Error; err != nil {
							return err
						}

						fmt.Printf("Created user with ID: %v", user.ID)
						return nil
					},
				},
			},
		},
	}
}

func main() {
	dbConn, err := db.Init()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	db.AutoMigrate(dbConn)

	app := cli.NewApp()
	app.Name = "tyomaa-admin-cli"
	app.Usage = "Admin CLI interface for tyomaa"
	app.Commands = commands()

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
