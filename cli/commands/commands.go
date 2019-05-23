package commands

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/winded/tyomaa/shared/api"

	"github.com/chzyer/readline"
	"github.com/urfave/cli"
	"github.com/winded/tyomaa/cli/settings"
	"github.com/winded/tyomaa/shared/api/client"
)

const (
	DateFormat = "02.01.2006 15:04:05"
)

func inputPrompt(reader *readline.Instance, prompt string, isPassword bool) string {
	if isPassword {
		value, _ := reader.ReadPassword(prompt)
		return string(value)
	} else {
		reader.SetPrompt(prompt)
		value, _ := reader.Readline()
		return value
	}
}

func validateAuth(settings settings.Settings) error {
	if settings.Api.Host == "" || settings.Api.Token == "" {
		return errors.New("No host or auth token found. Please login first")
	}
	return nil
}

func parseTime(value string) (time.Time, bool) {
	t, err := time.Parse(DateFormat, value)
	if err == nil {
		return t, true
	}

	t, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return t, true
	}

	t, err = time.Parse(time.UnixDate, value)
	if err == nil {
		return t, true
	}

	return time.Time{}, false
}

func displayEntries(w io.Writer, entries []api.TimeEntry) {
	var tw tabwriter.Writer
	tw.Init(w, 8, 8, 1, '\t', 0)

	fmt.Fprintf(&tw, "%v\t%v\t%v\t%v\t\n", "Entry ID", "Project name", "Start time", "End time")
	fmt.Fprintf(&tw, "%v\t%v\t%v\t%v\t\n", "----", "----", "----", "----")

	for _, entry := range entries {
		start := entry.Start.Local().Format(DateFormat)
		end := "None (ongoing)"
		if entry.End != nil {
			end = entry.End.Local().Format(DateFormat)
		}
		fmt.Fprintf(&tw, "%v\t%v\t%v\t%v\t\n", entry.ID, entry.Project, start, end)
	}

	tw.Flush()
}

func GetCommands(settings settings.Settings, apiClient *client.ApiClient, inputReader *readline.Instance, saveSettings func(settings.Settings)) []cli.Command {
	return []cli.Command{
		cli.Command{
			Name: "auth",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "status",
					Usage: "Show authentication status",
					Action: func(c *cli.Context) error {
						if settings.Api.Host == "" || settings.Api.Token == "" {
							fmt.Println("Not authenticated")
							return nil
						}

						response, err := apiClient.AuthTokenGet()
						if err != nil {
							return err
						}

						if response.User == nil {
							fmt.Println("Not authenticated")
							return nil
						}

						fmt.Println("Authenticated as " + response.User.Name)
						return nil
					},
				},
				cli.Command{
					Name:  "login",
					Usage: "Authenticate to a server",
					Action: func(c *cli.Context) error {
						host := inputPrompt(inputReader, "Host URL: ", false)
						username := inputPrompt(inputReader, "Username: ", false)
						password := inputPrompt(inputReader, "Password: ", true)

						newSettings := settings
						newSettings.Api.Host = host
						apiClient.Settings = newSettings.Api

						request := api.TokenPostRequest{
							Username: username,
							Password: password,
						}
						response, err := apiClient.AuthTokenPost(request)
						if err != nil {
							return err
						}

						newSettings.Api.Token = response.Token
						apiClient.Settings = newSettings.Api
						saveSettings(newSettings)

						fmt.Println("Login successful")
						return nil
					},
				},
			},
		},

		cli.Command{
			Name: "clock",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "status",
					Usage: "Show current clock status",
					Action: func(c *cli.Context) error {
						response, err := apiClient.ClockGet()
						if err != nil {
							return err
						}

						if response.Entry == nil {
							fmt.Println("Clock not running")
							return nil
						}

						duration := time.Since(response.Entry.Start)

						fmt.Printf("Clock running for %v (%v)\n", response.Entry.Project, duration.Truncate(time.Second).String())
						return nil
					},
				},
				cli.Command{
					Name:      "start",
					Usage:     "Start clock on specified project",
					ArgsUsage: "<project>",
					Action: func(c *cli.Context) error {
						project := c.Args().Get(0)
						if project == "" {
							return errors.New("Project argument missing")
						}

						request := api.ClockStartPostRequest{
							Project: project,
						}
						response, err := apiClient.ClockStartPost(request)
						if err != nil {
							return err
						}

						fmt.Printf("Clock started for project %v\n", response.Entry.Project)
						return nil
					},
				},
				cli.Command{
					Name:  "stop",
					Usage: "Stop active clock",
					Action: func(c *cli.Context) error {
						response, err := apiClient.ClockStopPost()
						if err != nil {
							return err
						}

						fmt.Printf("Clock stopped for project %v\n", response.Entry.Project)
						return nil
					},
				},
			},
		},

		cli.Command{
			Name: "entries",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "list",
					Usage: "Show entries",
					Action: func(c *cli.Context) error {
						response, err := apiClient.EntriesGet(api.EntriesGetRequest{})
						if err != nil {
							return err
						}

						if len(response.Entries) == 0 {
							fmt.Println("No entries found")
							return nil
						}

						displayEntries(c.App.Writer, response.Entries)
						return nil
					},
				},
				cli.Command{
					Name:      "create",
					Usage:     "Create a new entry. Start and end date formats can be RFC 3339, Unix or Finnish format",
					ArgsUsage: "<project> <start> <end>",
					Action: func(c *cli.Context) error {
						project := c.Args().Get(0)
						sStart := c.Args().Get(1)
						sEnd := c.Args().Get(2)
						if project == "" {
							return errors.New("Project arugment missing")
						} else if sStart == "" {
							return errors.New("Start arugment missing")
						} else if sEnd == "" {
							return errors.New("End arugment missing")
						}

						start, ok := parseTime(sStart)
						if !ok {
							return errors.New("Failed to parse start date")
						}
						end, ok := parseTime(sEnd)
						if !ok {
							return errors.New("Failed to parse end date")
						}

						if start == end || start.After(end) {
							return errors.New("Start time needs to be before end time")
						}

						request := api.EntriesPostRequest{
							Entry: api.TimeEntry{
								Project: project,
								Start:   start,
								End:     &end,
							},
						}
						response, err := apiClient.EntriesPost(request)
						if err != nil {
							return err
						}

						fmt.Printf("Entry created with ID %v\n", response.Entry.ID)
						return nil
					},
				},
				cli.Command{
					Name:      "view",
					Usage:     "Show a specific entry",
					ArgsUsage: "<entryid>",
					Action: func(c *cli.Context) error {
						sEntryID := c.Args().Get(0)
						if sEntryID == "" {
							return errors.New("Entry ID arugment missing")
						}

						entryID, err := strconv.Atoi(sEntryID)
						if err != nil {
							return err
						}

						response, err := apiClient.EntriesSingleGet(uint(entryID))
						if err != nil {
							return err
						}

						displayEntries(c.App.Writer, []api.TimeEntry{response.Entry})
						return nil
					},
				},
				cli.Command{
					Name:      "edit",
					Usage:     "Modify an existing entry",
					ArgsUsage: "<entryid>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "project, p",
							Usage: "Change the project of the entry",
						},
						cli.StringFlag{
							Name:  "start, s",
							Usage: "Change the start time of the entry",
						},
						cli.StringFlag{
							Name:  "end, e",
							Usage: "Change the end time of the entry",
						},
					},
					Action: func(c *cli.Context) error {
						sEntryID := c.Args().Get(0)
						if sEntryID == "" {
							return errors.New("Entry ID arugment missing")
						}

						entryID, err := strconv.Atoi(sEntryID)
						if err != nil {
							return err
						}

						response, err := apiClient.EntriesSingleGet(uint(entryID))
						if err != nil {
							return err
						}

						entry := response.Entry

						var ok bool
						if c.IsSet("project") {
							entry.Project = c.String("project")
						}
						if c.IsSet("start") {
							if entry.Start, ok = parseTime(c.String("start")); !ok {
								return err
							}
						}
						if c.IsSet("end") {
							entry.End = &time.Time{}
							if *entry.End, ok = parseTime(c.String("end")); !ok {
								return err
							}
						}

						_, err = apiClient.EntriesSinglePost(uint(entryID), api.EntriesSinglePostRequest{
							Entry: entry,
						})
						if err != nil {
							return err
						}

						fmt.Println("Entry modified successfully")
						return nil
					},
				},
			},
		},

		cli.Command{
			Name: "projects",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "list",
					Usage: "Show projects",
					Action: func(c *cli.Context) error {
						response, err := apiClient.ProjectsGet()
						if err != nil {
							return err
						}

						if len(response.Projects) == 0 {
							fmt.Println("No projects found")
							return nil
						}

						var tw tabwriter.Writer
						tw.Init(c.App.Writer, 8, 8, 1, '\t', 0)

						fmt.Fprintf(&tw, "%v\t%v\t\n", "Project name", "Total time")
						fmt.Fprintf(&tw, "%v\t%v\t\n", "----", "----")

						for _, project := range response.Projects {
							fmt.Fprintf(&tw, "%v\t%v\t\n", project.Name, project.TotalTime.Truncate(time.Second).String())
						}

						tw.Flush()
						return nil
					},
				},
			},
		},
	}
}
