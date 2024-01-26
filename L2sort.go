package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Row представляет структуру для хранения строки и ее ключей для сортировки
type Row struct {
	Original string
	Keys     []string
}

// RowSlice представляет срез строк для сортировки
type RowSlice []Row

func (s RowSlice) Len() int { return len(s) }

func (s RowSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s RowSlice) Less(i, j int) bool {
	for k := 0; k < len(s[i].Keys) && k < len(s[j].Keys); k++ {
		if s[i].Keys[k] == s[j].Keys[k] {
			continue
		}
		if numericSort && isNumeric(s[i].Keys[k]) && isNumeric(s[j].Keys[k]) {
			num1, _ := strconv.Atoi(s[i].Keys[k])
			num2, _ := strconv.Atoi(s[j].Keys[k])
			return num1 < num2
		}
		if monthSort {
			month1, err1 := time.Parse("January", s[i].Keys[k])
			month2, err2 := time.Parse("January", s[j].Keys[k])
			if err1 == nil && err2 == nil {
				return month1.Before(month2)
			}
		}
		return s[i].Keys[k] < s[j].Keys[k]
	}
	return false
}

var (
	keyColumn     int
	numericSort   bool
	reverseSort   bool
	uniqueLines   bool
	monthSort     bool
	ignoreBlanks  bool
	checkSorted   bool
	numericSuffix bool
)

func init() {
	flag.IntVar(&keyColumn, "k", 0, "Указание колонки для сортировки (по умолчанию 0)")
	flag.BoolVar(&numericSort, "n", false, "Сортировать по числовому значению")
	flag.BoolVar(&reverseSort, "r", false, "Сортировать в обратном порядке")
	flag.BoolVar(&uniqueLines, "u", false, "Не выводить повторяющиеся строки")
	flag.BoolVar(&monthSort, "M", false, "Сортировать по названию месяца")
	flag.BoolVar(&ignoreBlanks, "b", false, "Игнорировать хвостовые пробелы")
	flag.BoolVar(&checkSorted, "c", false, "Проверять отсортированы ли данные")
	flag.BoolVar(&numericSuffix, "h", false, "Сортировать по числовому значению с учетом суффиксов")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Использование: go run main.go [опции] файл")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filePath := args[0]
	lines, err := readLines(filePath)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		os.Exit(1)
	}

	rows := parseRows(lines)
	if checkSorted && isSorted(rows) {
		fmt.Println("Данные уже отсортированы.")
		os.Exit(0)
	}

	sort.Sort(RowSlice(rows))

	if reverseSort {
		reverse(rows)
	}

	if uniqueLines {
		rows = removeDuplicates(rows)
	}

	writeToFile(rows, filePath)
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseRows(lines []string) []Row {
	var rows []Row
	for _, line := range lines {
		keys := extractKeys(line)
		rows = append(rows, Row{Original: line, Keys: keys})
	}
	return rows
}

func extractKeys(line string) []string {
	if ignoreBlanks {
		line = strings.TrimSpace(line)
	}
	if keyColumn == 0 {
		return strings.Fields(line)
	}
	fields := strings.Fields(line)
	if keyColumn > len(fields) {
		return nil
	}
	return []string{fields[keyColumn-1]}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isSorted(rows []Row) bool {
	for i := 1; i < len(rows); i++ {
		if RowSlice(rows).Less(i, i-1) {
			return false
		}
	}
	return true
}

func reverse(rows []Row) {
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
}

func removeDuplicates(rows []Row) []Row {
	seen := make(map[string]bool)
	var result []Row
	for _, row := range rows {
		if !seen[row.Original] {
			seen[row.Original] = true
			result = append(result, row)
		}
	}
	return result
}

func writeToFile(rows []Row, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, row := range rows {
		file.WriteString(row.Original + "\n")
	}
}
