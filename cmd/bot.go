package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var bot = &cobra.Command{
	Use:   "bot",
	Short: "Starts Bot",
	Long:  `Starts an interactive C&C Slackbot.'`,
	Run: func(cmd *cobra.Command, args []string) {
		RunInitCrawl()
	},
}

func RunBot() {
	fmt.Println("To be implemented!")
}
