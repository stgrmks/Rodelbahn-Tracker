package cmd

import "github.com/spf13/cobra"

var initCrawl = &cobra.Command{
	Use:   "init",
	Short: "Updates data",
	Long:  `Insert's current state of the data into database without any notifications.'`,
	Run: func(cmd *cobra.Command, args []string) {
		RunInitCrawl()
	},
}

func RunInitCrawl() {
	notifyState := config.Notify
	config.Notify = false
	RunStartCrawler()
	config.Notify = notifyState
}
