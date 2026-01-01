package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

const flagGroupAnnotation = "group"

func cliUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", filepath.Base(os.Args[0]))

	// Group flags by annotation, default to "General Options"
	helpGroupLists := make(map[string][]*pflag.Flag)
	var helpGroupOrder []string
	var longestFlagName, longestHelpMessage, longestDefaultVal int

	pflag.CommandLine.VisitAll(func(f *pflag.Flag) {
		currentFlagAnnotations := f.Annotations[flagGroupAnnotation]
		flagGroup := "General Options"
		if len(currentFlagAnnotations) > 0 {
			flagGroup = currentFlagAnnotations[0]
		}

		if _, helpGroupExists := helpGroupLists[flagGroup]; !helpGroupExists {
			helpGroupLists[flagGroup] = []*pflag.Flag{}
			helpGroupOrder = append(helpGroupOrder, flagGroup)
		}
		helpGroupLists[flagGroup] = append(helpGroupLists[flagGroup], f)

		longestFlagName = max(longestFlagName, len(f.Name)+1)
		longestHelpMessage = max(longestHelpMessage, len(f.Usage)+1)
		longestDefaultVal = max(longestDefaultVal, len(getDefaultString(f))+1)
	})

	// Print each group
	for _, helpGroupName := range helpGroupOrder {
		flags := helpGroupLists[helpGroupName]
		if len(flags) == 0 {
			continue
		}

		fmt.Fprint(os.Stderr, colorText(hiYellow, helpGroupName+":\n"))
		for _, f := range flags {
			printFormattedFlag(
				f, longestFlagName, longestHelpMessage, longestDefaultVal)
		}
		fmt.Fprint(os.Stderr, "\n")
	}

	fmt.Fprintln(os.Stderr)
}

func printFormattedFlag(f *pflag.Flag, maxFlagName, maxHelpText, maxDef int) {
	defaultValue := getDefaultString(f)
	defaultValuePadding := strings.Repeat(" ", maxDef-len(defaultValue))

	helpPadding := strings.Repeat(" ", maxHelpText-len(f.Usage))
	defaultTxt := colorText(darkPurple, fmt.Sprintf(
		"%sDefault: %s%s", helpPadding, defaultValuePadding, defaultValue))

	flagPadding := strings.Repeat(" ", maxFlagName-len(f.Name))
	flagName := colorText(cyan, fmt.Sprintf("--%s%s", f.Name, flagPadding))

	usageText := colorText(green, f.Usage)

	fmt.Fprintf(os.Stderr, "\t%s %s   %s\n", flagName, usageText, defaultTxt)
}

// ANSI color codes

type color string

const (
	cyan       color = "\033[96m" // Bright cyan
	darkPurple color = "\033[38;5;55m"
	hiCyan     color = "\033[96m"
	hiYellow   color = "\033[93m" // Bright yellow
	green      color = "\033[92m" // Bright green
	white      color = "\033[97m" // Bright white
	hiBlack    color = "\033[90m" // Bright black (dim gray)
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	faint = "\033[2m"
)

func colorText(c color, text string) string { return string(c) + text + reset }

func getDefaultString(f *pflag.Flag) string {
	if f.DefValue == "" {
		return "\"\""
	}
	return f.DefValue
}

func addFlagToHelpGroup(flagName string, helpGroupName string) {
	lookupFlag := pflag.Lookup(flagName)
	if lookupFlag == nil {
		panic("unknown flag: " + flagName)
	}

	if lookupFlag.Annotations == nil {
		lookupFlag.Annotations = map[string][]string{}
	}
	lookupFlag.Annotations[flagGroupAnnotation] = []string{helpGroupName}
}
