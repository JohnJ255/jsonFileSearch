package storage

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"jsonFileSearch/internal/log"
	"os"
)

type CelestyMap map[string]interface{}

type Db struct {
	dataMap []CelestyMap
}

func Connect() (*Db, error) {
	return &Db{
		dataMap: make([]CelestyMap, 0),
	}, nil
}

func (db *Db) Close() {
}

func (db *Db) LoadFromJsonFile(dbFileName string) {
	dbData, err := os.ReadFile(dbFileName)
	if err != nil {
		log.Printf("Ошибка чтения файла базы данных: %v", err)
		return
	}
	log.Printf("Файл базы данных успешно загружен в память.")
	log.Printf("Идёт распознавание данных...")

	err = json.Unmarshal(dbData, &db.dataMap)
	dbData = []byte{}
	if err != nil {
		log.Printf("Ошибка при распаковке JSON данных: %v", err)
		return
	}

	log.Printf("Данные успешно распознаны, всего %d записей.", len(db.dataMap))
}

func (db *Db) SaveToFile(records []CelestyMap, resultFileName string) error {
	resultData, err := json.Marshal(&records)
	if err != nil {
		return fmt.Errorf("Ошибка при преобразовании данных в JSON: %v", err)
	}

	// Запись результата в файл
	err = os.WriteFile(resultFileName, resultData, 0644)
	if err != nil {
		return fmt.Errorf("Ошибка при записи результата в файл: %v", err)
	}

	log.Printf("Результат сохранен в файл: %s", resultFileName)
	return nil
}

func (db *Db) SaveToFileGracefully(records []CelestyMap, resultFileName string) error {
	file, err := os.Create(resultFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(records)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) SaveToBinaryFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(db.dataMap)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) LoadFromBinaryFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&db.dataMap)
	if err != nil {
		return err
	}

	return nil
}
