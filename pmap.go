package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type Entry struct {
	Name string    // Where this account belongs to (google, microsoft, war thunder etc.)
	Mail string    // Email used for creation of this account
	User string    // Username
	Pass string    // Password
	Note string    // Additional note (recovery ket etc.)
	Time time.Time // Last edit time of this entry
}

func NewEntry(name, mail, user, pass, note string) Entry {
	return Entry{
		Name: name,
		Mail: mail,
		User: user,
		Pass: pass,
		Note: note,
		Time: time.Now(),
	}
}

type PasswordMap struct {
	edited bool
	key    []byte
	file   string
	saved  struct {
		LastSaveTime time.Time
		Entries      []Entry
	}
}

func EmptyMap(key []byte, file string) *PasswordMap {
	pmap := new(PasswordMap)
	pmap.key = key
	pmap.file = file
	pmap.saved.Entries = make([]Entry, 0)
	return pmap
}

func (pmap *PasswordMap) Add(entry Entry) {
	pmap.saved.Entries = append(pmap.saved.Entries, entry)
	pmap.edited = true
}

func (pmap *PasswordMap) Delete(index int) {
	pmap.saved.Entries = append(pmap.saved.Entries[:index], pmap.saved.Entries[index+1:]...)
	pmap.edited = true
}

func (pmap *PasswordMap) Length() int {
	return len(pmap.saved.Entries)
}

func (pmap *PasswordMap) Get(index int) Entry {
	return pmap.saved.Entries[index]
}

func (pmap *PasswordMap) Set(index int, entry Entry) {
	pmap.saved.Entries[index] = entry
	pmap.edited = true
}

func (pmap *PasswordMap) Search(pattern string) []int {
	indices := []int{}
	for i, entry := range pmap.saved.Entries {
		if strings.Contains(entry.Name, pattern) ||
			strings.Contains(entry.Mail, pattern) ||
			strings.Contains(entry.User, pattern) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (pmap *PasswordMap) All() []int {
	indices := []int{}
	for i := range pmap.saved.Entries {
		indices = append(indices, i)
	}
	return indices
}

func (pmap *PasswordMap) UniqueNames() map[string]int {
	m := map[string]int{}
	for _, entry := range pmap.saved.Entries {
		i, ok := m[entry.Name]
		if ok {
			m[entry.Name] = i + 1
		} else {
			m[entry.Name] = 1
		}
	}
	return m
}

func (pmap *PasswordMap) UniqueMails() map[string]int {
	m := map[string]int{}
	for _, entry := range pmap.saved.Entries {
		i, ok := m[entry.Mail]
		if ok {
			m[entry.Mail] = i + 1
		} else {
			m[entry.Mail] = 1
		}
	}
	return m
}

func (pmap *PasswordMap) UniqueUsers() map[string]int {
	m := map[string]int{}
	for _, entry := range pmap.saved.Entries {
		i, ok := m[entry.User]
		if ok {
			m[entry.User] = i + 1
		} else {
			m[entry.User] = 1
		}
	}
	return m
}

func (pmap *PasswordMap) UniquePasswords() map[string]int {
	m := map[string]int{}
	for _, entry := range pmap.saved.Entries {
		i, ok := m[entry.Pass]
		if ok {
			m[entry.Pass] = i + 1
		} else {
			m[entry.Pass] = 1
		}
	}
	return m
}

func (pmap *PasswordMap) Save() error {
	jsonData, err := json.Marshal(pmap.saved)
	if err != nil {
		return err
	}
	data, err := Encrypt(jsonData, pmap.key)
	if err != nil {
		return err
	}
	err = os.WriteFile(pmap.file, data, 0666)
	if err != nil {
		return err
	}
	pmap.edited = false
	return nil
}

func (pmap *PasswordMap) Load() error {
	data, err := os.ReadFile(pmap.file)
	if err != nil {
		return err
	}
	jsonData, err := Decrypt(data, pmap.key)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &pmap.saved)
	if err != nil {
		return err
	}
	return nil
}
