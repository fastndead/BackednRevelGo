package PilotModel

import (
	"encoding/json"
	"errors"
)

var pilotList []Pilot = []Pilot{{1, "firstPilot", "firstPilot"},
	{2, "secondPilot", "secondPilot"},
	{3, "thirdPilot", "thirdPilot"}, 
	{4, "fourthPilot", "fourthPilot"}, 
	{5, "fifthPilot", "fifthPilot"}}

func GetAll() ([]Pilot, error) {
	if len(pilotList) > 0{
		return pilotList, nil
	}
	return pilotList, errors.New("Список пилотов пуст")
}

func GetById(index int) (Pilot, error) {
	for id := range pilotList {
		if pilotList[id].Id == index {
			return pilotList[id], nil
		}
	}
	return Pilot{},errors.New("Пилот не найден")
}

func Delete(index int) ([]Pilot, error) {
	for id := range pilotList {
		if pilotList[id].Id == index {
			pilotList = append(pilotList[:id], pilotList[id+1:]...)
			return pilotList, nil
		}
	}
	return pilotList, errors.New("Пилот не найден")
}

func Edit(index int, itemToAdd []byte) ([]Pilot, error){
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return pilotList, errors.New("Неправилные данные пилота")
	}
	for id := range pilotList {
		if pilotList[id].Id == index {
			pilotList[id].FirstName = temp.FirstName
			pilotList[id].LastName = temp.LastName
			return pilotList, nil
		}
	}
	return pilotList, errors.New("Не найден индекс")
}

func Add(itemToAdd []byte) ([]Pilot, error) {
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return pilotList, errors.New("Неправильные данные пилота")
	}
	pilotList = append(pilotList, *temp)
	return pilotList, nil
}
