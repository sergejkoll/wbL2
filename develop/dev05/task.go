package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	after      = flag.Int("A", 0, "печатать +N строк после совпадения")
	before     = flag.Int("B", 0, "печатать +N строк до совпадения")
	context    = flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	count      = flag.Bool("c", false, "вывод количества найденных строк")
	ignoreCase = flag.Bool("i", false, "игнорирование регистра")
	invert     = flag.Bool("v", false, "инвертировать результат")
	fixed      = flag.Bool("F", false, "искать точное совпадение со строкой, не паттерн")
	lineNum    = flag.Bool("n", false, "печатать номер строки")
)

type line struct {
	str string
	num int
}

// readFile - функция чтения данных из файла
func readFile(path string) ([]line, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []line
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		lines = append(lines, line{str: scanner.Text(), num: idx})
		idx += 1
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// addAfter - функция добваления строк после найденной
// input - изначальный файл
// set - множество уникальных строк которые в итоге запишутся в ответ
// idx - текующий индекс найденного слова
// count - число строк после
func addAfter(input []line, set map[string]struct{}, idx int, count int) []string {
	res := make([]string, 0)
	var s string
	if idx+count > len(input) {
		for _, val := range input[idx:] {
			// обработка с номером строки
			if *lineNum {
				s = strconv.Itoa(val.num) + " " + val.str
			} else {
				s = val.str
			}
			if _, ok := set[s]; ok {
				continue
			} else {
				res = append(res, s)
				set[s] = struct{}{}
			}
		}
		return res
	}

	for _, val := range input[idx : idx+count] {
		// обработка с номером строки
		if *lineNum {
			s = strconv.Itoa(val.num) + " " + val.str
		} else {
			s = val.str
		}
		if _, ok := set[s]; ok {
			continue
		} else {
			res = append(res, s)
			set[s] = struct{}{}
		}
	}
	return res
}

// addBefore - функция добваления строк перед найденной
// input - изначальный файл
// set - множество уникальных строк которые в итоге запишутся в ответ
// idx - текующий индекс найденного слова
// count - число строк до
func addBefore(input []line, set map[string]struct{}, idx int, count int) []string {
	res := make([]string, 0)
	var s string

	if idx-count < 0 {
		for _, val := range input[:idx] {
			// обработка с номером строки
			if *lineNum {
				s = strconv.Itoa(val.num) + " " + val.str
			} else {
				s = val.str
			}
			if _, ok := set[s]; ok {
				continue
			} else {
				res = append(res, s)
				set[s] = struct{}{}
			}
		}
		return res
	}

	for _, val := range input[idx-count : idx] {
		// обработка с номером строки
		if *lineNum {
			s = strconv.Itoa(val.num) + " " + val.str
		} else {
			s = val.str
		}
		if _, ok := set[s]; ok {
			continue
		} else {
			res = append(res, s)
			set[s] = struct{}{}
		}
	}
	return res
}

// addContext - печатать ±N строк вокруг совпадения
func addContext(input []line, set map[string]struct{}, idx int) []string {
	res := make([]string, 0)
	var s string
	res = append(res, addBefore(input, set, idx, *context)...)
	if *lineNum {
		s = strconv.Itoa(input[idx].num) + " " + input[idx].str
	} else {
		s = input[idx].str
	}
	if _, ok := set[s]; !ok {
		res = append(res, s)
		set[s] = struct{}{}
	}
	res = append(res, addAfter(input, set, idx, *context)...)

	return res
}

// flagHandler - обработчик флагов After, Before и Context
func flagHandler(input []line, set map[string]struct{}, idx int) (result []string) {
	if *context > 0 {
		return addContext(input, set, idx)
	}
	if *before > 0 {
		result = append(result, addBefore(input, set, idx, *before)...)
	}

	var s string
	if *lineNum {
		s = strconv.Itoa(input[idx].num) + " " + input[idx].str
	} else {
		s = input[idx].str
	}
	if _, ok := set[s]; !ok {
		result = append(result, s)
		set[s] = struct{}{}
	}

	if *after > 0 {
		result = append(result, addAfter(input, set, idx, *after)...)
	}
	return result
}

// grep - функция поиска строки в файле
func grep(input []line, pattern string) error {
	counter := 0
	result := make([]string, 0)
	set := make(map[string]struct{})

	if !*fixed {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}

		for idx, value := range input {
			line := value.str
			if *ignoreCase {
				line = strings.ToLower(line)
			}
			if re.MatchString(line) && !*invert {
				result = append(result, flagHandler(input, set, idx)...)
				counter += 1
			}
			if !re.MatchString(line) && *invert {
				result = append(result, flagHandler(input, set, idx)...)
				counter += 1
			}
		}

	} else {
		for idx, value := range input {
			line := value.str
			if *ignoreCase {
				line = strings.ToLower(line)
			}
			if strings.Contains(line, pattern) && !*invert {
				result = append(result, flagHandler(input, set, idx)...)
				counter += 1
			}
			if !strings.Contains(line, pattern) && *invert {
				result = append(result, flagHandler(input, set, idx)...)
				counter += 1
			}
		}
	}

	if *count {
		fmt.Println(counter)
		return nil
	}

	for _, val := range result {
		fmt.Println(val)
	}

	return nil
}

func main() {
	flag.Parse()
	filePath := flag.Arg(0)
	pattern := flag.Arg(1)

	input, err := readFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = grep(input, pattern)
	if err != nil {
		log.Fatal(err)
	}
}
