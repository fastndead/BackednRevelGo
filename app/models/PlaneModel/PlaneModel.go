package PlaneModel

import (
	"encoding/json"
	"errors"
)

var planeList []Plane = []Plane{{1, "firstPlane"},
	{2, "secondPlane"},
	{3, "thirdPlane"},
	{4, "fourthPlane"},
	{5, "fifthPlane"}}

func GetAll() ([]Plane, error) {
	if len(planeList) > 0{
		return planeList, nil
	}
	return planeList, errors.New("Список самолётов пуст")
}

func GetById(index int) (Plane, error) {
	for id := range planeList {
		if planeList[id].Id == index {
			return planeList[id], nil
		}
	}
	return Plane{},errors.New("Самолёт не найден")
}

func Delete(index int) ([]Plane, error) {
	for id := range planeList {
		if planeList[id].Id == index {
			planeList = append(planeList[:id], planeList[id+1:]...)
			return planeList, nil
		}
	}
	return planeList, errors.New("Самолёт не найден")
}

func Edit(index int, itemToAdd []byte) ([]Plane, error){
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return planeList, errors.New("Неправилные данные самолёта")
	}
	for id := range planeList {
		if planeList[id].Id == index {
			planeList[id].Name = temp.Name
			return planeList, nil
		}
	}
	return planeList, errors.New("Не найден индекс")
}

func Add(itemToAdd []byte) ([]Plane, error) {
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return planeList, errors.New("Неправильные данные самолёта")
	}
	planeList = append(planeList, *temp)
	return planeList, nil
}
