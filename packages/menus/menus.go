package menus

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var foodmenuDataFile = os.Getenv("DDS_FOODMENUDATAFILE")

type FoodMenuItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type FoodCategory struct {
	Name          string         `json:"Name"`
	Note          string         `json:"Note"`
	FoodMenuItems []FoodMenuItem `json:"FoodMenuItems"`
}

type FoodMenu []FoodCategory

// Read all of the json data
func ReadFoodMenuDataFile() (FoodMenu, error) {
	jsonFoodMenu, err := ioutil.ReadFile(foodmenuDataFile)
	if err != nil {
		return nil, err
	}

	var data FoodMenu
	// Unmarshal json data into struct.
	err = json.Unmarshal(jsonFoodMenu, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Write all menu data out to file.
func WriteFoodMenuDataFile(menu FoodMenu) error {
	file, err := json.MarshalIndent(menu, "", " ")
	if err != nil {
		return err
	}
	// Write to file.
	err = ioutil.WriteFile(foodmenuDataFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
