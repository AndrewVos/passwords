package main

import (
	"fmt"
	"github.com/AndrewVos/colour"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func clearScreen() {
	fmt.Printf("\x1b[2J")
}
func setCursorPosition(line int, column int) {
	fmt.Printf("\x1b[" + strconv.Itoa(line) + ";" + strconv.Itoa(column) + "H")
}
func matchNextCredential(credentials []Credential) {
	query := ""
	matches := []Credential{}

	clearScreen()
	setCursorPosition(1, 1)
	fmt.Printf("Start typing...")

	printQuery := func() {
		clearScreen()
		setCursorPosition(1, 1)
		if matches == nil {
			fmt.Printf("%v", query)
		} else {
			names := []string{}
			for _, credential := range matches {
				names = append(names, credential.Name)
			}
			fmt.Printf("%v", query)

			for i, name := range names {
				setCursorPosition(i+2, 1)
				name = colouriseMatchInString(query, name)
				if i == 0 {
					fmt.Printf(colour.Green("=>")+" %v", name)
				} else {
					fmt.Printf(name)
				}
			}
			setCursorPosition(1, len(query)+1)
		}
	}

	for {
		b := waitForNextByteFromStdin()
		if b == 127 {
			if query != "" {
				query = query[:len(query)-1]
			}
			matches = search(query, credentials)
			printQuery()
		} else if b == 10 {
			displayCredential(matches[0])
			return
		} else {
			matched, _ := regexp.MatchString("[a-zA-Z _\\-]", string(b))
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
	clearScreen()
	setCursorPosition(1, 1)
	fmt.Printf("Name:     %v\n", credential.Name)
	fmt.Printf("Site:     %v\n", credential.Site)
	fmt.Printf("Username: %v\n", credential.Username)
	fmt.Printf("Password: %v\n", credential.Password)
	fmt.Println()
	fmt.Println(colour.Blue("p = copy password to clipboard"))
	fmt.Println(colour.Blue("u = copy username to clipboard"))
	fmt.Println(colour.Blue("<enter> = to go back to search"))
	for {
		b := waitForNextByteFromStdin()
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
		if strings.Contains(strings.ToLower(credential.Name), strings.ToLower(query)) {
			matches = append(matches, credential)
		}
	}
	return matches
}

func waitForNextByteFromStdin() byte {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return b[0]
}

func main() {
	credentials := decrypt()
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	for {
		matchNextCredential(credentials)
	}
}

type Credential struct {
	Name     string
	Site     string
	Username string
	Password string
}

func decrypt() []Credential {
	output, err := exec.Command("/usr/bin/gpg", "--decrypt", "").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(output), "\n\n")
	credentials := []Credential{}
	for _, line := range lines {
		credential := Credential{}
		parts := strings.Split(line, "\n")
		for _, part := range parts {
			if strings.HasPrefix(part, "a:") {
				credential.Name = strings.Replace(part, "a: ", "", 1)
			}
			if strings.HasPrefix(part, "s:") {
				credential.Site = strings.Replace(part, "s: ", "", 1)
			}
			if strings.HasPrefix(part, "u:") {
				credential.Username = strings.Replace(part, "u: ", "", 1)
			}
			if strings.HasPrefix(part, "p:") {
				credential.Password = strings.Replace(part, "p: ", "", 1)
			}
		}
		credentials = append(credentials, credential)
	}
	return credentials
}
