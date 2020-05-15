package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

type StringDict map[string]interface{}

func getKeywordsMap() StringDict {
	jsonFile, _ := os.Open("dict.json")
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	var keywords map[string]interface{}
	json.Unmarshal(jsonBytes, &keywords)

	return keywords
}

func parseArgs() ([]string, []string) {
	var newArgs []string
	var sourceNames []string

	args := os.Args[1:]
	for i := range args {
		arg := args[i]
		if strings.Contains(arg, ".pgo") {
			baseName := strings.Split(arg, ".pgo")[0]
			sourceNames = append(sourceNames, baseName)
			arg = baseName + ".go"
		}
		newArgs = append(newArgs, arg)
	}

	return newArgs, sourceNames
}

func readFile(filename string) []string {
	file, _ := os.Open(filename)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()

	return lines
}

func isAlphaNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsLetter(r)
}

func convertString(str string, keywords StringDict) string {
	runeStr := []rune(str)
	for k, v := range keywords {
		i := strings.Index(str, k)
		if i != -1 {
			startValid := i == 0 || !isAlphaNum(runeStr[i-1])
			endValid := i+len(k) == len(str)-1 || !isAlphaNum(runeStr[i+len(k)])
			if startValid && endValid {
				v := v.(string)
				str = strings.Replace(str, k, v, 1)
				return convertString(str, keywords)
			}
		}
	}

	return str
}

func convertFile(filename string, keywords StringDict) string {

	lines := readFile(filename)
	var newLines []string

	for _, line := range lines {
		tokens := strings.Split(line, "\"")

		for i := range tokens {
			if i%2 == 1 {
				continue
			}

			tokens[i] = convertString(tokens[i], keywords)
		}
		newLine := strings.Join(tokens, "\"")
		newLines = append(newLines, newLine)
	}

	newText := strings.Join(newLines, "\n")
	return newText
}

func writeFile(filename string, contents string) {
	file, _ := os.Create(filename)
	file.WriteString(contents)
	file.Sync()
}

func deleteFile(filename string) {
	os.Remove(filename)
}

func main() {

	newArgs, sourceNames := parseArgs()
	keywords := getKeywordsMap()
	for i := range sourceNames {
		newContents := convertFile(sourceNames[i]+".pgo", keywords)
		writeFile(sourceNames[i]+".go", newContents)
	}

	cmd := exec.Command("go", newArgs...)
	
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	
	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s", errout.String())
	} else {
		fmt.Printf("%s", out.String())
	}

	for i := range sourceNames {
		deleteFile(sourceNames[i] + ".go")
	}
}
