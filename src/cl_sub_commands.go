package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	subcommand1 := flag.NewFlagSet("firstSub", flag.ExitOnError)
	subcommand2 := flag.NewFlagSet("secondSub", flag.ExitOnError)

	firstFlag := subcommand1.Bool("processing", false, "Command processing stauts")
	secondFlag := subcommand1.Int("bytes", 1024, "Byte length of result")

	flagSc2 := subcommand2.String("language", "en", "Language of result")

	if len(os.Args) < 2 {
		fmt.Println("Usage: app <subcommand> --flags")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "firstSub":
		subcommand1.Parse(os.Args[2:])
		fmt.Println("First subcommand")
		fmt.Println("Processing status: ", *firstFlag)
		fmt.Println("Byte length of result: ", *secondFlag)
	case "secondSub":
		subcommand2.Parse(os.Args[2:])
		fmt.Println("Second subcommand")
		fmt.Println("Language of result: ", *flagSc2)
	default:
		fmt.Println("Usage: app <subcommand> --flags")
		os.Exit(1)
	}

}
