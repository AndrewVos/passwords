package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func matchNextCredential(credentials []Credential) {
	query := ""
	lastLinePrinted := 0
	matches := []Credential{}

	printQuery := func() {
		fmt.Printf("\r")
		for i := 0; i <= lastLinePrinted; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\r")
		line := query
		if matches != nil {
			names := []string{}
			for _, credential := range matches {
				names = append(names, credential.Name)
			}
			// line = fmt.Sprintf("%v => [%v]", query, strings.Join(names, ","))
			line = fmt.Sprintf("%v => [%v]", query, matches[0].Name)
		}
		lastLinePrinted = len(line)
		fmt.Printf("%v", line)
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
			fmt.Println("p = copy password to clipboard")
			fmt.Println("u = copy username to clipboard")
			for {
				b = waitForNextByteFromStdin()
				if b == 10 {
					return
				} else if string(b) == "u" {
					copyToClipboard(matches[0].Username)
					fmt.Println("Copied username to clipboard")
				} else if string(b) == "p" {
					copyToClipboard(matches[0].Password)
					fmt.Println("Copied password to clipboard")
				}
			}
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

func copyToClipboard(s string) {
	xclip := exec.Command("/usr/bin/xclip", "-selection", "clipboard")
	w, _ := xclip.StdinPipe()
	fmt.Println(xclip.Start())
	w.Write([]byte(s))
	w.Close()
}

func search(query string, credentials []Credential) []Credential {
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
