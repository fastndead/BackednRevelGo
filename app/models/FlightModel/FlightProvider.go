package FlightModel

import (
	"encoding/json"
	"errors"
)

var flightList []Flight = []Flight{{1, "first", "first", "first", "first"}, {2, "second", "second", "second", "second"},
	{3, "third", "third", "third", "third"}, {4, "fourth", "fourth", "fourth", "fourth"}, {5, "fifth", "fifth", "fifth", "fifth"}}

func GetAll() ([]Flight, error) {
	if len(flightList) > 0{
		return flightList, nil
	}
	return flightList, errors.New("Список рейсов пуст")
}

func GetById(index int) (Flight, error) {
	for id := range flightList {
		if flightList[id].Id == index {
			return flightList[id], nil
		}
	}
	return Flight{},errors.New("Рейс не найден")
}

func Delete(index int) ([]Flight, error) {
	for id := range flightList {
		if flightList[id].Id == index {
			flightList = append(flightList[:id], flightList[id+1:]...)
			return flightList, nil
		}
	}
	return flightList, errors.New("Рейс не найден")
}

func Edit(index int, itemToAdd []byte) ([]Flight, error){
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return flightList, errors.New("Неправилные данные рейса")
	}
	for id := range flightList {
		if flightList[id].Id == index {
			flightList[id].ArrivalPoint = temp.ArrivalPoint
			flightList[id].DeparturePoint = temp.DeparturePoint
			flightList[id].Pilot = temp.Pilot
			flightList[id].Plane = temp.Plane
			return flightList, nil
		}
	}
	return flightList, errors.New("Не найден индекс")
}

func Add(itemToAdd []byte) ([]Flight, error) {
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return flightList, errors.New("Неправильные данные рейса")
	}
	flightList = append(flightList, *temp)
	return flightList, nil
}
