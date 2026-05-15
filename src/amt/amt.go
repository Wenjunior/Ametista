package main

import (
	"github.com/spf13/cobra"
)

import (
	"amt/sub"
	"amt/scan"
	"amt/probe"
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

	sub.Flags().IntVarP(&subOptions.Seconds, "timeout", "t", 10, "Set a timeout in seconds")

	sub.Flags().StringVarP(&subOptions.Output, "output", "o", "", "File to write results to")

	commands.AddCommand(sub)

	scanOptions := scan.ScanOptions {}

	scan := &cobra.Command {
		Use: "scan",
		Short: "Simple TCP port scanner",
		Run: func(command *cobra.Command, args []string) {
			scan.Run(scanOptions)
		},
	}

	scan.Flags().StringSliceVarP(&scanOptions.Targets, "targets", "t", []string{}, "Target hosts/networks")

	scan.Flags().StringVarP(&scanOptions.FileName, "list", "l", "", "File containing a list of target hosts/networks")

	scan.Flags().StringSliceVarP(&scanOptions.Patterns, "ports", "p", []string{}, "Only scan specified ports")

	scan.Flags().IntVarP(&scanOptions.BatchSize, "batch-size", "b", 3000, "Set a batch size")

	scan.Flags().IntVarP(&scanOptions.Seconds, "timeout", "w", 3, "Set a timeout in seconds")

	scan.Flags().StringVarP(&scanOptions.Output, "output", "o", "", "File to write results to")

	commands.AddCommand(scan)

	probeOptions := probe.ProbeOptions {}

	show := probe.Show {}

	probe := &cobra.Command {
		Use: "probe",
		Short: "HTTP/HTTPS probing",
		Run: func(command *cobra.Command, args []string) {
			probe.Run(probeOptions, show)
		},
	}

	probe.Flags().StringSliceVarP(&probeOptions.URLs, "urls", "u", []string{}, "Target URLs")

	probe.Flags().StringVarP(&probeOptions.FileName, "list", "l", "", "File containing a list of target URLs")

	probe.Flags().IntVarP(&probeOptions.BatchSize, "batch-size", "b", 3000, "Set a batch size")

	probe.Flags().IntVarP(&probeOptions.Seconds, "timeout", "w", 3, "Set a timeout in seconds")

	probe.Flags().BoolVarP(&show.StatusCode, "status-code", "s", false, "Show status code")

	probe.Flags().BoolVarP(&show.Server, "server", "S", false, "Show Server header")

	probe.Flags().BoolVarP(&show.XPoweredBy, "x-powered-by", "x", false, "Show X-Powered-By header")

	probe.Flags().BoolVarP(&show.Location, "location", "L", false, "Show Location header")

	probe.Flags().BoolVarP(&show.ContentLength, "content-length", "c", false, "Show content length")

	probe.Flags().BoolVarP(&show.ContentType, "content-type", "C", false, "Show Content-Type")

	probe.Flags().BoolVarP(&show.Title, "title", "t", false, "Show page title")

	probe.Flags().StringVarP(&probeOptions.Output, "output", "o", "", "File to write results to")

	commands.AddCommand(probe)

	commands.Execute()
}