package cmd

import (
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfgFile string

var dataFile string

var (
	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A CLI To-Do-List",
		Long: `tri is simple CLI app that helps you keep track of your tasks. With tri app, you can create, edit, and delete tasks, as well as set due dates and priorities. You can also view your tasks by category or due date, and you can even search for specific tasks. To-Do List App is the perfect way to stay organized and on top of your to-do list.

Here are some of the features of tri CLI app:

Create, edit, and delete tasks
Set due dates and priorities
View tasks by category or due date
Search for specific tasks
Customize your to-do list
Sync your to-do list across devices
tri CLI app is the perfect way to stay organized and on top of your to-do list. Download it today and start getting things done!`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := homedir.Dir()

	if err != nil {
		log.Println("Unable to detect home directory. Please set data file using --datafile.")
	}

	rootCmd.PersistentFlags().StringVar(&dataFile, "datafile", home+string(os.PathSeparator)+"projects/Go-Cocepts/ToDoCli/.todo.json", "datafile to store todos")
}
