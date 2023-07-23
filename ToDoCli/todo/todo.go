package todo

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Item struct {
	Text     string
	Priority int
	Position int
}

func SaveItems(filename string, items []Item) error {
	b, err := json.Marshal(items)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadItems(filename string) ([]Item, error) {
	var Items []Item
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return []Item{}, err
	}
	if err := json.Unmarshal(b, &Items); err != nil {
		return []Item{}, err
	}
	for i, _ := range Items {
		Items[i].Position = i + 1
	}
	return Items, nil
}

func (i *Item) SetPriority(pri int) {
	switch pri {
	case 1:
		i.Priority = 1
	case 3:
		i.Priority = 3
	default:
		i.Priority = 2
	}
}

func (i *Item) PrettyP() string {
	if i.Priority == 1 {
		return "(1)"
	}
	if i.Priority == 3 {
		return "(3)"
	}
	return " "
}

func (i *Item) Label() string {
	return strconv.Itoa(i.Position) + "."

}
