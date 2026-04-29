package main

import (
	"github.com/spf13/cobra"
)

import (
	"amt/sub"
)

func main() {
	commands := &cobra.Command {}

	commands.PersistentFlags().BoolP("help", "h", false, "")

	commands.PersistentFlags().Lookup("help").Hidden = true

	commands.SetHelpCommand(&cobra.Command {
		Use: "no-help",
		Hidden: true,
	})

	commands.CompletionOptions.DisableDefaultCmd = true

	subOptions := sub.SubOptions {}

	sub := &cobra.Command {
		Use: "sub",
		Short: "Passive subdomain enumeration",
		Run: func(command *cobra.Command, args []string) {
			sub.Run(subOptions)
		},
	}

	sub.Flags().StringSliceVarP(&subOptions.Domains, "domains", "d", []string{}, "Target root domain names")

	sub.Flags().StringVarP(&subOptions.FileName, "list", "l", "", "File containing a list of target root domain names")

	sub.Flags().IntVarP(&subOptions.TimeOut, "timeout", "t", 10, "Set how many seconds it should wait for a response")

	sub.Flags().StringVarP(&subOptions.Output, "output", "o", "", "File to write results to")

	commands.AddCommand(sub)

	commands.Execute()
}