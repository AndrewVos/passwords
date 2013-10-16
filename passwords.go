package main

import (
	"fmt"
	"github.com/AndrewVos/colour"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var filename string

func init() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: passwords <file>")
		os.Exit(1)
	} else {
		filename = os.Args[1]
	}
}

func main() {
	credentials, err := decrypt()
	if err != nil {
		fmt.Println("Error decrypting file...")
		return
	}
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	for {
		matchNextCredential(credentials)
	}
}

func matchNextCredential(credentials []Credential) {
	query := ""
	matches := []Credential{}

	ClearScreen()
	SetCursorPosition(1, 1)
	fmt.Printf("Start typing...")

	printQuery := func() {
		ClearScreen()
		SetCursorPosition(1, 1)
		if matches == nil {
			fmt.Printf("%v", query)
		} else {
			names := []string{}
			for _, credential := range matches {
				names = append(names, credential.SearchableContent)
			}
			fmt.Printf("%v", query)

			for i, name := range names {
				SetCursorPosition(i+2, 1)
				name = colouriseMatchInString(query, name)
				if i == 0 {
					fmt.Printf(colour.Green("=>")+" %v", name)
				} else {
					fmt.Printf("%v", name)
				}
			}
			SetCursorPosition(1, len(query)+1)
		}
	}

	for {
		b := WaitForNextByteFromStdin()
		if b == 127 {
			if query != "" {
				query = query[:len(query)-1]
			}
			matches = search(query, credentials)
			printQuery()
		} else if b == 10 {
			if len(matches) > 0 {
				displayCredential(matches[0])
				return
			}
		} else {
			matched, _ := regexp.MatchString("[0-9A-Za-z_\\-. ]", string(b))
			if matched {
				query = query + string(b)
				matches = search(query, credentials)
				printQuery()
			}
		}
	}
}

func colouriseMatchInString(query string, match string) string {
	query = strings.ToLower(query)
	match = strings.ToLower(match)
	parts := strings.Split(match, query)
	for i, part := range parts {
		parts[i] = colour.Blue(part)
	}
	return strings.Join(parts, colour.Red(query))
}

func displayCredential(credential Credential) {
	ClearScreen()
	SetCursorPosition(1, 1)
	fmt.Printf("Name:     %v\n", credential.Name)
	fmt.Printf("Site:     %v\n", credential.Site)
	fmt.Printf("Username: %v\n", credential.Username)
	fmt.Printf("Password: %v\n", credential.Password)
	fmt.Println()
	fmt.Println(colour.Blue("p = copy password to clipboard"))
	fmt.Println(colour.Blue("u = copy username to clipboard"))
	fmt.Println(colour.Blue("<enter> = to go back to search"))
	for {
		b := WaitForNextByteFromStdin()
		if b == 10 {
			break
		} else if string(b) == "u" {
			copyToClipboard(credential.Username)
			fmt.Println(colour.Yellow("Copied username to clipboard"))
		} else if string(b) == "p" {
			copyToClipboard(credential.Password)
			fmt.Println(colour.Yellow("Copied password to clipboard"))
		}
	}
}

func copyToClipboard(s string) {
	xclip := exec.Command("/usr/bin/xclip", "-selection", "clipboard")
	w, _ := xclip.StdinPipe()
	xclip.Start()
	w.Write([]byte(s))
	w.Close()
}

func search(query string, credentials []Credential) []Credential {
	if query == "" {
		return credentials
	}
	matches := []Credential{}
	for _, credential := range credentials {
		if strings.Contains(strings.ToLower(credential.SearchableContent), strings.ToLower(query)) {
			matches = append(matches, credential)
		}
	}
	return matches
}

func decrypt() ([]Credential, error) {
	output, err := exec.Command("/usr/bin/gpg", "--decrypt", filename).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n\n")
	credentials := []Credential{}
	for _, line := range lines {
		name := ""
		site := ""
		user := ""
		pass := ""
		parts := strings.Split(line, "\n")
		for _, part := range parts {
			if strings.HasPrefix(part, "a:") {
				name = strings.Replace(part, "a: ", "", 1)
			}
			if strings.HasPrefix(part, "s:") {
				site = strings.Replace(part, "s: ", "", 1)
			}
			if strings.HasPrefix(part, "u:") {
				user = strings.Replace(part, "u: ", "", 1)
			}
			if strings.HasPrefix(part, "p:") {
				pass = strings.Replace(part, "p: ", "", 1)
			}
		}
		credentials = append(credentials, NewCredential(name, site, user, pass))
	}
	return credentials, nil
}
