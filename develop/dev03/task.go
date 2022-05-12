package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Parameters - структура для хранения флагов
type Parameters struct {
	columnNumber  int  // номер колонки для сортировки
	isNumericSort bool // сортировка по числовому занчению
	isReverse     bool // сортировка в обратном порядке
	isOnlyUnique  bool // использовать только уникальные строки
}

// line - структура для хранения строк из файла
// str - строка
// words - массив слов в строке заспличенных по пробелу
type line struct {
	str   string
	words []string
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
	for scanner.Scan() {
		lines = append(lines, line{str: scanner.Text(), words: strings.Fields(scanner.Text())})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// sortLines - функция реализующая утилиту sort
func sortLines(input []line, params Parameters) []line {
	// использование уникальных строк
	if params.isOnlyUnique {
		input = unique(input)
	}

	// если указан номер колонки которой сортируем
	if params.columnNumber > 0 {
		sort.Slice(input, func(i, j int) bool {
			// если в строке нет слов в указанной колонке
			if len(input[i].words) <= params.columnNumber-1 ||
				len(input[j].words) <= params.columnNumber-1 {
				return true
			}
			// сортировка по числовому значению (ищем в строке float в нужном столбце и сортируем по нему)
			if params.isNumericSort {
				first, err := strconv.ParseFloat(input[i].words[params.columnNumber-1], 64)
				if err != nil {
					return true
				}
				second, err := strconv.ParseFloat(input[j].words[params.columnNumber-1], 64)
				if err != nil {
					return true
				}
				return first < second
			}
			// сортировка без числового признака
			return input[i].words[params.columnNumber-1] < input[j].words[params.columnNumber-1]
		})

		if params.isReverse {
			reverse(input)
		}

		return input
	}

	// если не указан столбец ищем первое числовое значение в строке
	if params.isNumericSort {
		re := regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")
		sort.Slice(input, func(i, j int) bool {
			firstStr := re.FindString(input[i].str)
			secondStr := re.FindString(input[j].str)
			// если числа нет
			if firstStr == "" || secondStr == "" {
				return true
			}
			first, _ := strconv.Atoi(firstStr)
			second, _ := strconv.Atoi(secondStr)
			return first < second
		})
		if params.isReverse {
			reverse(input)
		}
		return input
	}

	// в случае если ни один из флагов не был использован
	sort.Slice(input, func(i, j int) bool {
		return input[i].str < input[j].str
	})

	if params.isReverse {
		reverse(input)
	}

	return input
}

// unique - функция которая оставляет только уникальные строки
func unique(lines []line) []line {
	var res []line
	keys := make(map[string]struct{})
	for _, l := range lines {
		if _, ok := keys[l.str]; !ok {
			keys[l.str] = struct{}{}
			res = append(res, l)
		}
	}
	return res
}

// reverse - функция для реверса вывода
func reverse(lines []line) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

func main() {
	columnNumber := flag.Int("k", 0, "указание колонки для сортировки")
	isNumericSort := flag.Bool("n", false, "сортировать по числовому значению")
	isReverse := flag.Bool("r", false, "сортировать в обратном порядке")
	isOnlyUnique := flag.Bool("u", false, "не выводить повторяющиеся строки")
	flag.Parse()
	params := Parameters{
		columnNumber:  *columnNumber,
		isNumericSort: *isNumericSort,
		isReverse:     *isReverse,
		isOnlyUnique:  *isOnlyUnique,
	}

	lines, err := readFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	result := sortLines(lines, params)
	for _, line := range result {
		fmt.Println(line.str)
	}
}
