// Copyright © 2018 Anshul Sanghi <anshap1719@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/mgutz/ansi"
	"fmt"
	"time"
	"reflect"
	"github.com/go-cmd/cmd"
)

func runExternalCmd(name string, args []string) {
	c := cmd.NewCmd(name, args...)
	statusChan := c.Start()

	ticker := time.NewTicker(100 * time.Nanosecond)

	var previousLine string
	var previousError string
	var previousStderr = []string{""}
	var previousStdout = []string{""}

	var color func(string) string

	if name == "ng" || name == "npm" {
		color = ansi.ColorFunc("red+bh")
	} else if name == "go" || name == "gin" {
		color = ansi.ColorFunc("cyan+b")
	} else {
		color = ansi.ColorFunc("white")
	}

	var stderr bool

	go func() {
		for range ticker.C {
			status := c.Status()

			if status.Complete {
				c.Stop()
				break
			}

			if err := status.Error; err != nil {
				fmt.Errorf("error occurred: %s", err.Error())
				break
			}

			n := len(status.Stdout)
			n2 := len(status.Stderr)

			if n2 < 1 {
				stderr = false
			} else {
				stderr = true
			}

			var currentLine string
			var currentError string
			var currentStderr []string
			var currentStdout []string

			if n < 1 && n2 < 1 {
				continue
			}

			if n == 1 {
				currentLine = status.Stdout[n-1]
			}

			if n2 == 1 {
				currentError = status.Stderr[n2 - 1]
			}

			if n2 > 1 {
				currentStderr = status.Stderr
				if !reflect.DeepEqual(currentStderr, previousStderr) {
					for _, err := range status.Stderr {
						fmt.Println(color(err))
					}
				}
				previousStderr = currentStderr
			}

			if n > 1 {
				currentStdout = status.Stdout
				if !reflect.DeepEqual(currentStdout, previousStdout) {
					for _, err := range status.Stdout {
						fmt.Println(color(err))
					}
				}
				previousStdout = currentStdout
			}

			if n == 1 || n2 == 1 {
				if stderr && (previousError != currentError || previousError == "" && (currentError != "" && currentError != "\n")) {
					fmt.Println(color(currentError))
					previousError = currentError
				}

				if previousLine != currentLine || previousLine == "" && (currentLine != "" && currentLine != "\n") {
					fmt.Println(color(currentLine))
					previousLine = currentLine
				}

				continue
			} else {
				continue
			}
		}
	}()

	// Check if command is done
	select {
		case _ = <-statusChan:
			c.Stop()
		default:
			// no, still running
	}

	// Block waiting for command to exit, be stopped, or be killed
	_ = <-statusChan
}

//func runExternalCmd(name string, args []string) {
//	var stdoutBuf, stderrBuf bytes.Buffer
//	cmd := exec.Command(name, args...)
//
//	var color func(string) string
//
//	if name == "ng" {
//		color = ansi.ColorFunc("red+bh")
//	} else if name == "go" || name == "gin" {
//		color = ansi.ColorFunc("cyan+b")
//	}
//
//	stdoutIn, _ := cmd.StdoutPipe()
//	stderrIn, _ := cmd.StderrPipe()
//
//	var errStdout, errStderr error
//	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
//	stderr := io.MultiWriter(os.Stdin, &stderrBuf)
//	err := cmd.Start()
//	if err != nil {
//		log.Fatalf("cmd.Start() failed with '%s'\n", err)
//	}
//
//	go func() {
//		_, errStdout = io.Copy(stdout, stdoutIn)
//	}()
//
//	go func() {
//		_, errStderr = io.Copy(stderr, stderrIn)
//	}()
//
//	err = cmd.Wait()
//	if err != nil {
//		log.Fatalf("cmd.Run() failed with %s\n", err)
//	}
//	if errStdout != nil || errStderr != nil {
//		log.Fatal("failed to capture stdout or stderr\n")
//	}
//	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
//	if outStr != "" {
//		fmt.Printf(color(outStr))
//	}
//	if errStr != "" {
//		fmt.Printf(color(errStr))
//	}
//}

func RedFunc() func(string) string {
	return ansi.ColorFunc("red+bh")
}

func BlueFunc() func(string) string {
	return ansi.ColorFunc("cyan+b")
}
