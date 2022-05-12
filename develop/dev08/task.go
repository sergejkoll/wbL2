package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// readInput - функция чтения из stdin и опредления команды: pipeline, fork или обычная команда
func readInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		input = strings.TrimSuffix(input, "\n")

		if strings.Contains(input, "|") {
			err = Pipeline(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}

		if strings.Contains(input, "&") {
			err = Fork(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}

		if err = execInput(input, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// execInput - вызов команды переданной в инпуте
func execInput(input string, out io.Writer) error {
	cmdArgs := strings.Split(input, " ")
	switch cmdArgs[0] {
	case "cd":
		if len(cmdArgs) < 2 {
			return os.Chdir(os.Getenv("HOME"))
		}
		return os.Chdir(cmdArgs[1])
	case "q":
		os.Exit(0)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = out

	return cmd.Run()
}

// getCommands - вспомогательная функция для Pipeline которая формирует массив *exec.Cmd
func getCommands(input string) []*exec.Cmd {
	commands := make([]*exec.Cmd, 0)

	strCommands := strings.Split(input, "|")
	for _, cmd := range strCommands {
		cmdArgs := strings.Split(cmd, " ")
		commands = append(commands, exec.Command(cmdArgs[0], cmdArgs[1:]...))
	}

	return commands
}

// Pipeline - функция обработки пайпа для каждой следующей команды указываем stdin как out от прошлой и затем в цикле
// запускаем каждую команду по очереди и затем ждем результат который записывается в буфер
func Pipeline(input string) (err error) {
	commands := getCommands(input)
	if len(commands) < 1 {
		return nil
	}

	var output, stderr bytes.Buffer

	for i, cmd := range commands[:len(commands)-1] {
		if commands[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return err
		}
		cmd.Stderr = &stderr
	}

	commands[len(commands)-1].Stdout, commands[len(commands)-1].Stderr = &output, &stderr

	for _, cmd := range commands {
		if err = cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range commands {
		if err = cmd.Wait(); err != nil {
			return err
		}
	}

	if len(output.Bytes()) > 0 {
		fmt.Fprintln(os.Stdout, string(output.Bytes()))
	}

	if len(stderr.Bytes()) > 0 {
		fmt.Fprintln(os.Stderr, string(output.Bytes()))
	}

	return nil
}

// Fork - функция обработки форк команд в дочернем процессе запускаем переданную команду
func Fork(input string) error {
	input = strings.TrimSuffix(input, " &")
	fmt.Println(input)
	id, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		os.Exit(1)
	}
	if id == 0 {
		err := execInput(input, nil)
		if err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}

func main() {
	readInput()
}
