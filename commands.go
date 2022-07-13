package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
	"unicode"
)

type CommandProcessor struct {
	pmap *PasswordMap
}

func (cmd *CommandProcessor) Process(args []string) (error, bool) {
	switch args[0] {
	case "?", "h", "help":
		fmt.Println(HelpMessage)
	case "q", "quit", "exit":
		if cmd.pmap.edited {
			fmt.Print("There are unsaved changes. Are you sure? (yes/n): ")
			ch := ReadLine("")
			if ch == "yes" {
				return nil, true
			}
		} else {
			return nil, true
		}
	case "save":
		err := cmd.pmap.Save()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Saved successfully")
		}
	case "reset":
		err := cmd.Reset(args[1:])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "ls": // Shortcut
		err := cmd.List([]string{"all"})
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "list":
		err := cmd.List(args[1:])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "add":
		err := cmd.Add()
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "edit":
		err := cmd.Edit(args[1:])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "delete":
		err := cmd.Delete(args[1:])
		if err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Printf("Unknown command %s\n", args[0])
	}
	return nil, false
}

func (cmd *CommandProcessor) Reset(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("This command takes additional argument(s) (type help)")
	}
	switch args[0] {
	case "password":
		fmt.Println("Specify a new master password")
		pass1 := ReadPassword("Password: ")
		pass2 := ReadPassword("Password: ")
		if pass1 == pass2 {
			cmd.pmap.key = PasswordHash([]byte(pass1))
		} else {
			fmt.Println("Passwords are not same")
		}
	default:
		return fmt.Errorf("Unknown argument %s", args[0])
	}
	return nil
}

func mergeArgs(args []string) string {
	ret := ""
	for i, s := range args {
		ret += s
		if i < len(args)-1 {
			ret += " "
		}
	}
	return ret
}

type stringIntPair struct {
	s string
	n int
}

func sortUniqueMap(m map[string]int) []stringIntPair {
	arr := []stringIntPair{}
	for k, v := range m {
		arr = append(arr, stringIntPair{
			s: k,
			n: v,
		})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].s < arr[j].s
	})
	return arr
}

func (cmd *CommandProcessor) List(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("This command takes additional argument(s) (type help)")
	}
	switch args[0] {
	case "pattern":
		// Merge all args after pattern
		pattern := mergeArgs(args[1:])
		indices := cmd.pmap.Search(pattern)
		if len(indices) > 0 {
			PrintEntries(cmd.pmap, indices)
		} else {
			fmt.Println("Nothing found")
		}
	case "all":
		indices := cmd.pmap.All()
		if len(indices) > 0 {
			PrintEntries(cmd.pmap, indices)
		} else {
			fmt.Println("Nothing found")
		}
	case "accounts":
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueNames()), "Account Name", "Repeat")
	case "emails":
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueMails()), "Email", "Repeat")
	case "usernames":
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueUsers()), "Username", "Repeat")
	case "passwords":
		PrintPairs(sortUniqueMap(cmd.pmap.UniquePasswords()), "Password", "Repeat")
	default:
		return fmt.Errorf("Unknown argument %s", args[0])
	}
	return nil
}

func isWhitespace(s string) bool {
	for _, c := range []rune(s) {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}

func hasWhitespace(s string) bool {
	for _, c := range []rune(s) {
		if unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func (cmd *CommandProcessor) Add() error {
	name := ReadLine("Name: ")
	if isWhitespace(name) {
		return fmt.Errorf("Name required")
	}
	mail := ReadLine("Email: ")
	if isWhitespace(mail) {
		mail = ""
	}
	user := ReadLine("Username: ")
	if isWhitespace(user) {
		user = ""
	}
	pass := ReadLine("Password: ")
	if isWhitespace(pass) {
		return fmt.Errorf("Password required")
	}
	note := ReadLine("Note: ")
	cmd.pmap.Add(NewEntry(name, mail, user, pass, note))
	return nil
}

func (cmd *CommandProcessor) Edit(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("This command takes only one argument (type help)")
	}
	index, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Index has to be a valid number!")
	}
	if index >= 0 && index < int64(cmd.pmap.Length()) {
		entry := cmd.pmap.Get(int(index))
		edited := false
		name := ReadLine("Name: ")
		if !isWhitespace(name) {
			entry.Name = name
			edited = true
		}
		mail := ReadLine("Email: ")
		if !isWhitespace(mail) {
			entry.Mail = mail
			edited = true
		}
		user := ReadLine("Username: ")
		if !isWhitespace(user) {
			entry.User = user
			edited = true
		}
		pass := ReadLine("Password: ")
		if !isWhitespace(pass) {
			entry.Pass = pass
			edited = true
		}
		note := ReadLine("Note: ")
		if !isWhitespace(note) {
			entry.Note = note
			edited = true
		}
		if edited {
			entry.Time = time.Now()
			cmd.pmap.Set(int(index), entry)
		}
	} else {
		return fmt.Errorf("Index out of range")
	}
	return nil
}

func (cmd *CommandProcessor) Delete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("This command takes only one argument (type help)")
	}
	index, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Index has to be valid number!")
	}
	if index >= 0 && index < int64(cmd.pmap.Length()) {
		cmd.pmap.Delete(int(index))
	} else {
		return fmt.Errorf("Index out of range")
	}
	return nil
}
