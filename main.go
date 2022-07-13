package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed help.txt
var HelpMessage string

func main() {
	file := flag.String("file", "pasman.bin", "The file name contains encrypted passwords and other information")
	flag.Parse()
	var pmap *PasswordMap
	_, err := os.Stat(*file)
	if os.IsNotExist(err) {
		// First use
		fmt.Println("Welcome to the pasman. pasman is a password manager lives in your terminal.")
		fmt.Println("First you must create a password for encrypting your pasman file.")
		fmt.Println("This password will be your master password and if you forget or lose it")
		fmt.Println("you will not be able to recover your passwords.")
		fmt.Println("We will ask you same password twice, for security :)")
		pass1 := ReadPassword("Password: ")
		pass2 := ReadPassword("Password: ")
		if pass1 == pass2 {
			pmap = EmptyMap(PasswordHash([]byte(pass1)), *file)
		}
	} else if err != nil {
		// Error
		Fatalf("Could not open file %s because of error: %v\n", *file, err)
	} else {
		// Not first use and no errors
		abs, err := filepath.Abs(*file)
		if err != nil {
			abs = *file
		}
		fmt.Println("Using file", abs)
		pass := ReadPassword("Password: ")
		pmap = EmptyMap(PasswordHash([]byte(pass)), *file)
		err = pmap.Load()
		if err != nil {
			Fatalf("%v\n", err)
		}
	}
	// Processor processes user commands
	cmd := CommandProcessor{
		pmap: pmap,
	}
	// Main loop
	for {
		err, quit := cmd.Process(ReadCommand())
		if err != nil {
			Fatalf("%v\n", err)
		}
		if quit {
			break
		}
	}
}
