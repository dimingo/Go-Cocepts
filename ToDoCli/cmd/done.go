/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"example/ToDOCli/todo"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"do"},
	Short:   "Mark Item as Done",

	Run: doneRun,
}

func doneRun(cmd *cobra.Command, args []string) {
	items, err := todo.ReadItems(dataFile)
	i, err := strconv.Atoi(args[0])

	if i > 0 && i < len(items) {
		items[i-1].Done = true
		fmt.Printf("%q %v\n", items[i-1].Text, "marked done")

		sort.Sort(todo.ByPri(items))
		todo.SaveItems(dataFile, items)
	} else {
		log.Println(i, "doesn't match any items")
	}

	if err != nil {
		log.Fatal(args[0], "is not a valid label\n", err)
	}
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
