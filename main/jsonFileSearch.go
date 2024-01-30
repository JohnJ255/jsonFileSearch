package main

import (
	"encoding/json"
	"fmt"
	"jsonFileSearch/internal/log"
	"jsonFileSearch/internal/storage"
	"jsonFileSearch/internal/terminal"
	"os"
	"strings"
)

func main() {
	log.Printf("Программа поиска по базе данных JSON")
	if len(os.Args) < 2 {
		log.Printf("Использование: jsonDbSearch.exe <путь_и_имя_файла_базы_данных>")
		return
	}
	log.Printf("Читается указанный в параметрах запуска файл...")

	dbFileName := os.Args[1]

	db, err := storage.Connect()
	if err != nil {
		log.Printf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	db.LoadFromJsonFile(dbFileName)

	log.Printf("Нажмите Ctrl+C для выхода или напишите слово \"выход\" вместо строки для поиска")

	db.PrintSearchInfo()

	searchString := ""

	for searchString != "выход" {
		log.Printf("Введите строку для поиска (или \"выход\"):")
		searchString, err = terminal.InputString()
		if err != nil {
			log.Printf("Ошибка при вводе: %v", err)
			continue
		}
		if strings.ToLower(strings.Trim(searchString, " ")) == "выход" {
			break
		}

		records, err := db.Search(searchString)
		if err != nil {
			log.Printf("Ошибка при поиске: %v", err)
			continue
		}
		resultSaved := false
		for !resultSaved {
			log.Printf("Найдено записей: %d", len(records))
			if len(records) == 0 {
				break
			}
			resultSaved, err = showResultDialog(db, records)
			if err != nil {
				log.Printf("Ошибка при сохранении результата: %v", err)
			}
		}
	}
}

func showResultDialog(db *storage.Db, records []storage.CelestyMap) (bool, error) {
	n := 3
	if len(records) < n {
		n = len(records)
	}
	yes := true
	if n >= 3 {
		log.Printf("Показать первые %d записи на экране? (да/enter=нет). Напишите \"отмена\" если результат не нужно показывать или сохранять", n)
		var err error
		var answer string
		answer, yes, err = terminal.InputConfirm()
		if err != nil {
			return false, err
		}
		if answer == "отмена" {
			return true, nil
		}
	}
	if yes {
		log.Printf("Результат:\n%s", prepareResult(records[0:n]))
	}
	resultFileName := "output.json"
	log.Printf("Укажите название файла (или просто нажмите enter чтобы сохранить в %s) или напишите \"отмена\":", resultFileName)
	answer, err := terminal.InputString()
	if err != nil {
		return false, err
	}
	if answer == "отмена" {
		return true, nil
	}

	if answer != "" {
		resultFileName = answer
	}

	if fileExists(resultFileName) {
		log.Printf("Файл уже существует, перезаписать? (да/enter=нет)")
		answer, yes, err = terminal.InputConfirm()
		if err != nil {
			return false, err
		}
		if answer == "отмена" {
			return true, nil
		}
		if !yes {
			return false, nil
		}
	}

	err = db.SaveToFileGracefully(records, resultFileName)
	if err != nil {
		log.Printf("Ошибка при сохранении результата в файл")
		return false, err
	}
	return true, nil
}

func prepareResult(records []storage.CelestyMap) string {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при выводе данных:", err)
		return ""
	}
	return string(jsonData)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
