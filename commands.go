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
	case "info":
		fmt.Println("Last save time:", cmd.pmap.saved.LastSaveTime.Format("02-01-2006 15:04:05"))
	case "save":
		err := cmd.pmap.Save()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Saved successfully")
		}
	case "reset":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.Reset(args[1])
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
		} else {
			fmt.Println("Added successfully")
		}
	case "edit":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.Edit(args[1])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "delete":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.Delete(args[1])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Deleted successfully")
		}
	case "sort":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.Sort(args[1])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Sorted successfully")
		}
	case "import":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.pmap.Import(args[1])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Imported successfully")
		}
	case "export":
		if len(args) < 2 {
			fmt.Println("This command requires an argument (type help)")
			break
		}
		err := cmd.pmap.Export(args[1])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Exported successfully")
		}
	default:
		fmt.Printf("Unknown command %s\n", args[0])
	}
	return nil, false
}

func (cmd *CommandProcessor) Reset(arg string) error {
	switch arg {
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
		return fmt.Errorf("Unknown argument %s", arg)
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
	case "search":
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
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueAccounts()), "Account Name", "Repeat")
	case "emails":
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueMails()), "Email", "Repeat")
	case "usernames":
		PrintPairs(sortUniqueMap(cmd.pmap.UniqueUsernames()), "Username", "Repeat")
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

func (cmd *CommandProcessor) Edit(arg string) error {
	index, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return fmt.Errorf("Index has to be a valid number!")
	}
	if index >= 0 && index < int64(cmd.pmap.Length()) {
		entry := cmd.pmap.Get(int(index))
		fmt.Println("Leave empty to not change")
		edited := false
		fmt.Println("Current account name is", entry.Name)
		name := ReadLine("Name: ")
		if !isWhitespace(name) {
			entry.Name = name
			edited = true
		}
		fmt.Println("Current email is", entry.Mail)
		mail := ReadLine("Email: ")
		if !isWhitespace(mail) {
			entry.Mail = mail
			edited = true
		}
		fmt.Println("Current username is", entry.User)
		user := ReadLine("Username: ")
		if !isWhitespace(user) {
			entry.User = user
			edited = true
		}
		fmt.Println("Current password is", entry.Pass)
		pass := ReadLine("Password: ")
		if !isWhitespace(pass) {
			entry.Pass = pass
			edited = true
		}
		fmt.Println("Current note is", entry.Note)
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

func (cmd *CommandProcessor) Delete(arg string) error {
	index, err := strconv.ParseInt(arg, 10, 64)
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

func (cmd *CommandProcessor) Sort(arg string) error {
	switch arg {
	case "accounts":
		cmd.pmap.SortByAccounts()
	case "emails":
		cmd.pmap.SortByMails()
	case "usernames":
		cmd.pmap.SortByUsernames()
	case "passwords":
		cmd.pmap.SortByPasswords()
	case "time":
		cmd.pmap.SortByTime()
	default:
		return fmt.Errorf("Unknown argument %s", arg)
	}
	return nil
}
