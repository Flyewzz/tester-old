package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type Program struct {
	Path        string
	MemoryLimit string
	DiskLimit   string
	CpuLimit    string
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

	// filepath.Dir(p.Path))) - a directory contains the executing program
	cmd := exec.Command("bash", "-c", fmt.Sprintf("docker run --rm -i --memory=%s --memory-swap %s --cpus=%s "+
		"-v %s:/program frolvlad/alpine-gxx "+
		"/bin/sh -c \"g++ program/main.cpp && ./a.out\"",
		p.MemoryLimit, p.DiskLimit, p.CpuLimit, filepath.Dir(p.Path)))
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	cmd.ProcessState.ExitCode()
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
			Path:        "/Users/alpha/Desktop/programs/tester/task/main.cpp",
			MemoryLimit: "30m",
			DiskLimit:   "31m",
			CpuLimit:    ".50",
		}
		result, err := p.Run(test.Input)
		if err != nil {
			log.Printf("Error for %s:\n %s\n", test.Name, err.Error())
			continue
		}
		fmt.Printf("Result for %s: \n %s\n", test.Name, result)
	}

}
