package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

func ReadCommand() []string {
	return strings.Split(ReadLine("-> "), " ")
}

func ReadLine(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func ReadPassword(msg string) string {
	fmt.Print(msg)
	text, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(text)
}

func Fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(-1)
}

func PrintEntries(pmap *PasswordMap, indices []int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCaption(true, fmt.Sprintf("total %d", len(indices)))
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Index", "Account Name", "E-Mail", "Username", "Password", "Note", "Last Edit (DD-MM-YYYY)"})
	for _, index := range indices {
		entry := pmap.Get(index)
		table.Append([]string{
			fmt.Sprintf("%d", index),
			entry.Name,
			entry.Mail,
			entry.User,
			entry.Pass,
			entry.Note,
			entry.Time.Format("02-01-2006 15:04:05"),
		})
	}
	table.Render()
}

func PrintPairs(arr []stringIntPair, first, second string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCaption(true, fmt.Sprintf("total %d", len(arr)))
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{first, second})
	for i := range arr {
		table.Append([]string{
			arr[i].s,
			fmt.Sprintf("%d", arr[i].n),
		})
	}
	table.Render()
}
