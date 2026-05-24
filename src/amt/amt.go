package main

import (
	"os"
	"fmt"
)

import (
	"amt/sub"
	"amt/flag"
	"amt/scan"
	"amt/probe"
	"amt/spider"
)

func main() {
	defer func() {
		err := recover()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	} ()

	subCommand := flag.NewFlagSet("sub", flag.ExitOnError)

	subOptions := sub.Options {}

	subCommand.StringVar(&subOptions.Domain, "d", "", "Target root domain name")

	subCommand.StringVar(&subOptions.FileName, "l", "", "File containing a list of target root domain names")

	subCommand.IntVar(&subOptions.Seconds, "t", 10, "Set a timeout in seconds")

	subCommand.StringVar(&subOptions.Output, "o", "", "File to write results to")

	scanCommand := flag.NewFlagSet("scan", flag.ExitOnError)

	scanOptions := scan.Options {}

	scanCommand.StringVar(&scanOptions.Target, "t", "", "Target host")

	scanCommand.StringVar(&scanOptions.FileName, "l", "", "File containing a list of target hosts")

	scanCommand.StringVar(&scanOptions.Patterns, "p", "", "Only scan specified ports")

	scanCommand.IntVar(&scanOptions.BatchSize, "b", 3000, "Set a batch size")

	scanCommand.IntVar(&scanOptions.Seconds, "w", 3, "Set a timeout in seconds")

	scanCommand.StringVar(&scanOptions.Output, "o", "", "File to write results to")

	probeCommand := flag.NewFlagSet("probe", flag.ExitOnError)

	probeOptions := probe.Options {}

	show := probe.Show {}

	probeCommand.StringVar(&probeOptions.FileName, "l", "", "File containing a list of target URLs")

	probeCommand.IntVar(&probeOptions.BatchSize, "b", 3000, "Set a batch size")

	probeCommand.BoolVar(&show.IPAddress, "i", false, "Show IP address")

	probeCommand.IntVar(&probeOptions.Seconds, "w", 10, "Set a timeout in seconds")

	probeCommand.BoolVar(&show.StatusCode, "s", false, "Show status code")

	probeCommand.BoolVar(&show.Server, "server", false, "Show Server header")

	probeCommand.BoolVar(&show.XPoweredBy, "x", false, "Show X-Powered-By header")

	probeCommand.BoolVar(&show.Location, "location", false, "Show Location header")

	probeCommand.BoolVar(&show.ContentLength, "cl", false, "Show content length")

	probeCommand.BoolVar(&show.ContentType, "ct", false, "Show Content-Type")

	probeCommand.BoolVar(&show.Title, "t", false, "Show page title")

	probeCommand.StringVar(&probeOptions.Output, "o", "", "File to write results to")

	spiderCommand := flag.NewFlagSet("spider", flag.ExitOnError)

	spiderOptions := spider.Options {}

	spiderCommand.StringVar(&spiderOptions.URL, "u", "", "Target URL")

	spiderCommand.StringVar(&spiderOptions.FileName, "l", "", "File containing a list of target URLs")

	spiderCommand.IntVar(&spiderOptions.Seconds, "t", 10, "Set a timeout in seconds")

	spiderCommand.IntVar(&spiderOptions.BatchSize, "b", 1500, "Set a batch size")

	spiderCommand.StringVar(&spiderOptions.Output, "o", "", "File to write results to")

	if len(os.Args) == 1 {
		panic("Too few arguments")
	}

	switch os.Args[1] {
	case "sub":
		subCommand.Parse(os.Args[2:])

		sub.Run(subOptions)
	case "scan":
		scanCommand.Parse(os.Args[2:])

		scan.Run(scanOptions)
	case "probe":
		probeCommand.Parse(os.Args[2:])

		probe.Run(probeOptions, show)
	case "spider":
		spiderCommand.Parse(os.Args[2:])

		spider.Run(spiderOptions)
	default:
		if os.Args[1] == "-h" {
			fmt.Println("Usage: amt [subcommand] [options]\n\nCommands:\n  sub\t\tPassive subdomain enumeration\n  scan\t\tTCP port scanner\n  probe\t\tHTTP(S) probing\n  spider\tWeb crawler")

			return
		}

		panic("Unknown option(s)/subcommand")
	}
}