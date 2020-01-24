package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Program struct {
	Path string
}

type Test struct {
	Name   string
	Input  string
	Output string
}

func GetTestsNames(path string) ([]string, error) {
	tests, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, nil
	}
	var dirNames []string
	for _, t := range tests {

		if t.IsDir() {
			dirNames = append(dirNames, t.Name())
		}
	}

	return dirNames, nil
}

func GetTest(path string) (*Test, error) {
	testFolder, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var test Test
	subPaths := strings.Split(path, "/")
	test.Name = subPaths[len(subPaths)-1]
	if path[len(path)-1] != '/' {
		test.Name = subPaths[len(subPaths)-1]
		path += "/"
	}
	for _, t := range testFolder {
		if !t.IsDir() {
			filePath := path + t.Name()
			file, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			switch t.Name() {
			case "input.txt":
				test.Input = string(file)
			case "output.txt":
				test.Output = string(file)
			}
		}
	}
	return &test, nil
}

func (p *Program) Run(input string) (string, error) {
	cmd := exec.Command("bash", "-c", "docker run --rm -i frolvlad/alpine-gxx /bin/sh -c \"g++ -x c++ - && ./a.out\"")
	file, _ := os.Open(p.Path)
	cmd.Stdin = file
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", nil
	}
	return out.String(), nil
}

func main() {
	const testsPath string = "tests/"

	names, _ := GetTestsNames(testsPath)
	var tests []*Test
	for _, name := range names {
		test, _ := GetTest(testsPath + name)
		tests = append(tests, test)
	}

	for _, t := range tests {
		fmt.Printf("name: %s\n input: %s\n output: %s\n", t.Name, t.Input, t.Output)
	}

	fmt.Println("----------------------")

	for _, test := range tests {
		p := Program{
			Path: "task/main.cpp",
		}
		result, err := p.Run(test.Input)
		if err != nil {
			log.Printf("Error for %s:\n %s\n", test.Name, err.Error())
			continue
		}
		fmt.Printf("Result for %s: \n %s\n", test.Name, result)
	}

}
