package storage

import (
	"fmt"
	"jsonFileSearch/internal/log"
	"strconv"
	"strings"
)

func (db *Db) Search(str string) ([]CelestyMap, error) {
	searchStr := strings.ToLower(str)
	searchStr = strings.ReplaceAll(searchStr, ",", ".")
	if strings.HasSuffix(searchStr, "\"") {
		if strings.HasPrefix(searchStr, "\"") {
			return db.simpleSearchMap(strings.Trim(searchStr, "\" "), "")
		} else {
			search := strings.SplitN(searchStr, " ", 2)
			return db.simpleSearchMap(strings.Trim(search[1], "\" "), search[0])
		}
	}
	search := splitAndFilter(searchStr, " ")
	if len(search) == 0 {
		return nil, fmt.Errorf("Пустая строка поиска")
	}
	if len(search) == 1 {
		return db.simpleSearchMap(search[0], "")
	}
	if len(search) == 2 {
		if isNumeric(search[0]) {
			return db.rangeSearchMap(search[0], search[1], "")
		}
		return db.simpleSearchMap(search[1], search[0])
	}
	if len(search) == 3 {
		return db.rangeSearchMap(search[1], search[2], search[0])
	}

	db.PrintSearchInfo()

	return nil, fmt.Errorf("Слишком много параметров для поиска")
}

func splitAndFilter(str string, s string) []string {
	result := []string{}
	res := strings.Split(str, s)
	for _, v := range res {
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

func (db *Db) PrintSearchInfo() {
	fmt.Println("\nИщет записи в json следующими способами:")
	fmt.Println("- просто подстрока: нужно ввести искомую часть строки или числа, например: J95")
	fmt.Println("- подстрока параметра: нужно ввести через пробел название параметра и подстроку или число, например: Peri 130")
	fmt.Println("- диапазон чисел: нужно ввести через пробел два числа (от и до включительно), например: 0.8 0.9")
	fmt.Println("- диапазон чисел параметра: нужно ввести через пробел название параметра и два числа (от и до включительно), например: Year_of_perihelion 1998 2000")
	fmt.Println("- подстрока с пробелами: нужно ввести строку в кавычках, например: \"C/1995 O1 (Hale-Bopp)\"")
	fmt.Println("Внимание!\n" +
		"\tРежим поиска по подстроке: если ввести, например, 1, то будут находиться числа вроде 0.8481884 так как в них присутствует символ 1\n" +
		"\tРежим поиска по диапазону: тут уже идёт поиск по числам, например, если ввести 1 2, то будут находиться числа только между 1 и 2 по величине")
	fmt.Println()

}

func (db *Db) simpleSearchMap(s string, fieldName string) ([]CelestyMap, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("Пустая строка поиска")
	}
	if len(db.dataMap) == 0 {
		return nil, fmt.Errorf("Нет данных для поиска")
	}
	if _, ok := db.dataMap[0][formatFieldname(fieldName)]; fieldName != "" && !ok {
		return nil, fmt.Errorf("Указано неизвестное поле: %s", fieldName)
	}

	if fieldName == "" {
		log.Printf("Поиск по подстроке \"%s\" по всем полям", s)
	} else {
		log.Printf("Поиск по подстроке \"%s\" в поле \"%s\"", s, fieldName)
	}

	result := make([]CelestyMap, 0)
	for _, record := range db.dataMap {
		for k, v := range record {
			if (fieldName == "" || k == formatFieldname(fieldName)) && containsInAny(v, s) {
				result = append(result, record)
				break
			}
		}
	}

	return result, nil
}

func (db *Db) rangeSearchMap(from string, to string, fieldName string) ([]CelestyMap, error) {
	if len(db.dataMap) == 0 {
		return nil, fmt.Errorf("Нет данных для поиска")
	}
	if _, ok := db.dataMap[0][formatFieldname(fieldName)]; fieldName != "" && !ok {
		return nil, fmt.Errorf("Указано неизвестное поле: %s", fieldName)
	}

	d1, err := strconv.ParseFloat(from, 64)
	if err != nil {
		return nil, fmt.Errorf("Не удалось распознать число: %s, %v", from, err)
	}
	d2, err := strconv.ParseFloat(to, 64)
	if err != nil {
		return nil, fmt.Errorf("Не удалось распознать число: %s, %v", to, err)
	}
	if d1 > d2 {
		d1, d2 = d2, d1
	}

	if fieldName == "" {
		log.Printf("Поиск по диапазону от %f до %f по всем полям", d1, d2)
	} else {
		log.Printf("Поиск по диапазону от %f до %f в поле \"%s\"", d1, d2, fieldName)
	}

	result := make([]CelestyMap, 0)
	for _, record := range db.dataMap {
		for k, v := range record {
			if (fieldName == "" || k == formatFieldname(fieldName)) && inDiapasone(v, d1, d2) {
				result = append(result, record)
				break
			}
		}
	}

	return result, nil
}

func formatFieldname(name string) string {
	if len(name) == 0 {
		return name
	}
	name = strings.ToLower(name)
	if name != "e" && name != "i" {
		x := strings.ToTitle(string(name[0]))
		if len(name) > 1 {
			x += name[1:]
		}
		return x
	}
	return name
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		_, err = strconv.ParseFloat(s, 64)
	}
	return err == nil
}

func containsInAny(v interface{}, s string) bool {
	switch v := v.(type) {
	case string:
		return strings.Contains(strings.ToLower(v), s)
	case float64:
		return strings.Contains(fmt.Sprintf("%f", v), s)
	case int:
		return strings.Contains(fmt.Sprintf("%d", v), s)
	default:
		return false
	}
}

func inDiapasone(e interface{}, d1 float64, d2 float64) bool {
	switch v := e.(type) {
	case float64:
		return v >= d1 && v <= d2
	case int:
		return float64(v) >= d1 && float64(v) <= d2
	default:
		return false
	}
}
