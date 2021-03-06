package provider

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type matcher func(string) bool

func Prompt(prompt string, sensitive bool) (string, error) {
	fmt.Fprintf(ProviderOut, "%s: ", prompt)
	if sensitive {
		var input []byte
		input, err := terminal.ReadPassword(1)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(input)), nil
	}
	reader := bufio.NewReader(ProviderIn)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func PromptMulti(choices []string) (string, int) {
	match := func(_ string) bool {
		return false
	}
	return PromptMultiMatch(choices, match)
}

func PromptMultiMatch(choices []string, match matcher) (string, int) {
	if len(choices) == 0 {
		return "error", -1
	} else if len(choices) == 1 {
		return choices[0], 0
	}
	str := ""
	for i, c := range choices {
		if match(c) {
			return c, i
		}
		// not the optimal way to biuld a string, I know
		str = fmt.Sprintf("%s[%3d] %s\n", str, i, c)
	}
	fmt.Fprintf(ProviderOut, "%sChoice: ", str)
	for {
		var selection int = -1
		fmt.Fscanf(ProviderIn, "%d", &selection)
		if selection < 0 || selection >= len(choices) {
			fmt.Fprintf(ProviderOut, "Invalid choice, try again\n")
		} else {
			return choices[selection], selection
		}
	}
}

func PromptMultiMatchRole(choices []string, opt string) (string, int) {
	re, err := regexp.Compile("role/(keycloak-)?(" + opt + ")$")
	if err != nil {
		fmt.Fprintf(ProviderErr, "Error interpreting requested role: %s", opt)
		return PromptMulti(choices)
	}
	match := func(c string) bool {
		return re.Match([]byte(c))
	}
	choice, n := PromptMultiMatch(choices, match)
	m := regexp.MustCompile("role/(keycloak-)?(.*)").FindStringSubmatch(choice)
	return m[len(m)-1], n
}
