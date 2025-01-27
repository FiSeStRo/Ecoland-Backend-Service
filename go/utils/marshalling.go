package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/FiSeStRo/Ecoland-Backend-Service/database"
)

func UnmarshalFile(path string, data interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println("json file could not be read", err)
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println("error reading file data", err)
		return err
	}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Println("could not unmarshal file", err)
		return err
	}

	return nil
}

func SetupDefTable[T interface{}](
	path string,
	data []T,
	tableName string,
	tableCol string,
) error {
	db := database.GetDB()
	_, err := db.Exec("TRUNCATE TABLE " + tableName)
	if err != nil {
		return fmt.Errorf("failing to truncate %v, %w", tableName, err)
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("json file could not be read %w", err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("error reading file data %w", err)
	}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return fmt.Errorf("could not unmarshal file %w", err)
	}
	placeholders := placeholdersForStruct(data[0])
	for _, v := range data {

		stmt, err := db.Prepare(`INSERT INTO ` + tableName + tableCol + ` VALUE(` + placeholders + `)`)
		if err != nil {
			return fmt.Errorf("failing to prepare db %w", err)
		}
		args := structToInterfaceSlice(v)
		_, err = stmt.Exec(args...)
		if err != nil {
			return fmt.Errorf("failng to execute stmt %w", err)
		}
	}
	return nil
}

func placeholdersForStruct(v interface{}) string {
	val := reflect.ValueOf(v)
	numFields := val.NumField()
	placeholderSlice := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		placeholderSlice[i] = "?"
	}
	return strings.Join(placeholderSlice, ",")
}

func structToInterfaceSlice(obj interface{}) []interface{} {
	v := reflect.ValueOf(obj)
	fields := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		fields[i] = v.Field(i).Interface()
	}
	return fields
}
