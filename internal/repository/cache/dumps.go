package cache

import (
	"encoding/json"
	"log"
	"os"
)

func makeDump(filename string, pool any) {
	data, err := json.Marshal(pool)
	if err != nil {
		log.Panicln("ERROR can not marshall pull to "+filename+": ", err.Error())
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Println("ERROR can not write pull to  "+filename+": ", err.Error())
		return
	}
	log.Println("INFO dump saved successfully to " + filename)
}

func loadFromDump(filename string, pool any) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &pool)
	return err
}
