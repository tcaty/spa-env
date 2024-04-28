package passwd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// default path to passwd file
	path = "/etc/passwd"
	// default column separator in passwd file
	sep = ":"
	// default columns count in passwd file
	columns = 7
)

type Entry struct {
	// Username is used when user logs in
	Username string
	// An x character indicates that encrypted
	// and salted password is stored in /etc/shadow file.
	Password string
	// Each user must be assigned a user ID
	Uid string
	// The primary group ID (stored in /etc/group file)
	Gid string
	// User ID Info, the comment field
	// extra information about the users
	Gecos string
	// Home directory - the absolute path to the directory
	// the user will be in when they log in
	Home string
	// The absolute path of a command or shell (/bin/bash).
	Shell string
}

// Parse /etc/passwd file and return list of entries
func Parse() ([]*Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	var entries []*Entry

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry, err := parseLine(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("error occured while parsing %s: %v", path, err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// Parse single line of /etc/passwd file, return entry
// return error if the line has wrong format
func parseLine(line string) (*Entry, error) {
	c := strings.Split(line, sep)

	if len(c) != columns {
		return nil, fmt.Errorf("wrong %s format: columns excpected %d, but got %d", path, columns, len(c))
	}

	entry := &Entry{
		Username: c[0],
		Password: c[1],
		Uid:      c[2],
		Gid:      c[3],
		Gecos:    c[4],
		Home:     c[5],
		Shell:    c[6],
	}

	return entry, nil
}
