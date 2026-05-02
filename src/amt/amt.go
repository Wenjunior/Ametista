package main

import (
	"github.com/spf13/cobra"
)

import (
	"amt/sub"
	"amt/scan"
)

func main() {
	commands := &cobra.Command {
		Short: "Reconnaissance tool",
	}

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

	scanOptions := scan.ScanOptions {}

	scan := &cobra.Command {
		Use: "scan",
		Short: "TCP port scanner",
		Run: func(command *cobra.Command, args []string) {
			scan.Run(scanOptions)
		},
	}

	scan.Flags().StringSliceVarP(&scanOptions.Targets, "targets", "t", []string{}, "Target hosts/networks")

	scan.Flags().StringVarP(&scanOptions.FileName, "list", "l", "", "File containing a list of target hosts/networks")

	scan.Flags().StringSliceVarP(&scanOptions.Patterns, "ports", "p", []string{"21-23", "25", "80", "110", "443", "445", "3036", "8080", "8443"}, "Only scan specified ports")

	scan.Flags().IntVarP(&scanOptions.BatchSize, "batch-size", "b", 1024, "Set a batch size")

	scan.Flags().IntVarP(&scanOptions.TimeOut, "timeout", "w", 3, "Set how many seconds it should wait")

	commands.AddCommand(scan)

	commands.Execute()
}