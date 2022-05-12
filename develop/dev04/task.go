package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// isAnagram - функция проверки если отсортированные слова равны то слова являются анаграммами
func isAnagram(firstWord string, secondWord string) bool {
	first := []rune(firstWord)
	second := []rune(secondWord)
	sort.Slice(first, func(i, j int) bool {
		return first[i] < first[j]
	})
	sort.Slice(second, func(i, j int) bool {
		return second[i] < second[j]
	})

	return string(first) == string(second)
}

// 1 решение квадратичная сложность(
func findAnagram(words []string) map[string][]string {
	result := make(map[string][]string)
	for idx, item := range words {
		words[idx] = strings.ToLower(item)
	}

	unique := true
	for _, word := range words {
		if _, ok := result[word]; !ok {

			for key := range result {
				if isAnagram(word, key) {
					unique = false
					result[key] = append(result[key], word)
				}
			}

			if unique {
				result[word] = []string{word}
			}
		}

		unique = true
	}

	return result
}

//2
type setOfAnagrams struct {
	anagrams map[string]struct{}
	arr      []string
}

func sortString(s string) string {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

func findAnagramFast(words []string) map[string][]string {
	result := make(map[string][]string)
	// key - отсортированное слово, value - множество анаграмм по этому слову
	set := make(map[string]*setOfAnagrams)
	sortedWord := ""

	for _, word := range words {
		word = strings.ToLower(word)
		sortedWord = sortString(word)
		if val, existSet := set[sortedWord]; existSet {
			if _, existWord := val.anagrams[word]; !existWord {
				val.anagrams[word] = struct{}{}
				val.arr = append(val.arr, word)
			}
		} else {
			newSet := new(setOfAnagrams)
			newSet.anagrams = make(map[string]struct{})
			newSet.anagrams[word] = struct{}{}
			newSet.arr = []string{word}
			set[sortedWord] = newSet
		}
	}

	for _, v := range set {
		if len(v.arr) == 1 {
			continue
		}
		key := v.arr[0]
		sort.Strings(v.arr)
		result[key] = v.arr
	}

	return result
}

func main() {
	dict := []string{"пятка", "листок", "пятак", "слиток", "столик", "тяпка", "пятка", "слово", "лсвоо"}
	res := findAnagramFast(dict)
	fmt.Println(res)
}
