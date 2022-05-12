package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	fields    = flag.Int("f", 0, "выбрать поля (колонки)")
	delimiter = flag.String("d", "\t", "использовать другой разделитель")
	separated = flag.Bool("s", false, "только строки с разделителем")
)

func cut(line string) {
	columns := strings.Split(line, *delimiter)
	if *separated {
		if *fields > 0 && len(columns) >= *fields {
			fmt.Println(columns[*fields-1])
		} else {
			fmt.Println("")
		}
	} else {
		if *fields > 0 && len(columns) >= *fields {
			fmt.Println(columns[*fields-1])
		} else {
			fmt.Println(line)
		}
	}
}

func readInput() {
	signalChannel := make(chan os.Signal, 1)     // канал для обработки SIGINT
	signal.Notify(signalChannel, syscall.SIGINT) // запись в signalChannel если пришел SIGINT
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-signalChannel:
			fmt.Println("Завершение работы")
			close(signalChannel)
			return
		default:
			scanner.Scan()
			cut(scanner.Text())
		}
	}
}

func main() {
	flag.Parse()
	readInput()
}
