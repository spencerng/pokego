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
    return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func getIndicesOf(str string, target string) []int {
    var indices []int
    offset := 0
    for {
        i := strings.Index(str, target)

        if i == -1 {
            break
        }

        indices = append(indices, i+offset)
        str = str[i+len(target):]
        offset += i + len(target)
    }

    return indices
}

func convertString(str string, keywords StringDict) string {
    runeStr := []rune(str)
    for k, v := range keywords {
        indices := getIndicesOf(str, k)
        offset := 0
        for _, i := range indices {
            i += offset
            startValid := i == 0 || !isAlphaNum(runeStr[i-1])
            endValid := i+len(k) == len(str) || !isAlphaNum(runeStr[i+len(k)])
            if startValid && endValid {
                v := v.(string)
                str = str[:i] + v + str[i+len(k):]
                runeStr = []rune(str)
                offset += len(v) - len(k)
            }
        }
    }

    return str
}

func convertMultiline(lines []string, keywords StringDict) string {

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
        sourceContents := readFile(sourceNames[i] + ".pgo")
        newContents := convertMultiline(sourceContents, keywords)
        writeFile(sourceNames[i]+".go", newContents)
    }

    cmd := exec.Command("go", newArgs...)

    var out bytes.Buffer
    var errout bytes.Buffer
    cmd.Stderr = &errout

    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err := cmd.Run()

    for i := range sourceNames {
        deleteFile(sourceNames[i] + ".go")
    }

    if err != nil {
        keywords = StringDict{"go": "pokego", "Go": "Pokego"}
        erroutSplit := strings.Split(errout.String(), "\n")
        fmt.Printf("%s", convertMultiline(erroutSplit, keywords))
    } else {
        fmt.Printf("%s", out.String())
    }

}
