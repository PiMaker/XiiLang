package parser

import (
    "bufio"
    "os"
    "log"
    "errors"
    "strings"
)

type Token struct {
    Text string
    File string
    Line int
}

var (
    currentTokenLine int
    currentPath string
)

func TokenizeFile(path string) ([][]Token, error) {
    
    inFile, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer inFile.Close()

    scanner := bufio.NewScanner(inFile)
    scanner.Split(bufio.ScanLines)

    var tokens [][]Token

    log.Println("Tokenizing...")

    currentTokenLine = 0
    currentPath = path

    for scanner.Scan() {
        currentTokenLine++
        tokenize(scanner.Text(), &tokens)
    }

    log.Printf("%d lines processed\n", len(tokens))

    if len(tokens) == 0 {
        return nil, errors.New("No tokens found, is the file empty?")
    }

    return tokens, nil
}

func tokenize(line string, slice *[][]Token) {
    trimmed := strings.TrimSpace(line)
    if trimmed == "" {
        return
    }

    if strings.Index(trimmed, "#") == 0 {
        return
    }

    words := strings.Split(line, " ")
    *slice = append(*slice, stringSliceToTokenSlice(words))
}

func stringSliceToTokenSlice(vs []string) []Token {
    var vsm []Token
    for _, v := range vs {
        if strings.TrimSpace(v) != "" {
            vsm = append(vsm, Token{Text: v, File: currentPath, Line: currentTokenLine})
        }
    }
    return vsm
}