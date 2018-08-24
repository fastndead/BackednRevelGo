package PilotModel

import (
	"encoding/json"
	"errors"
	"database/sql"
	"strconv"
	"fmt"
)

type PilotProvider struct{
	Db *sql.DB
}

func (p *PilotProvider)GetAll() ([]Pilot, error) {
	sql := "SELECT c_id, c_first_name, c_last_name FROM airport.pilots"
	rows, err := p.Db.Query(sql)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var PlaneList []Pilot
	for rows.Next(){
		var id int
		var firstName, lastName string

		if err := rows.Scan(&id, &firstName, &lastName); err != nil{
			return nil, err
		}

		PlaneList = append(PlaneList, Pilot{id, firstName, lastName})
	}
	return PlaneList, nil
}

func (p *PilotProvider)GetById(index int) (Pilot, error) {
	sql := "SELECT c_id, c_first_name, c_last_name FROM airport.pilots WHERE c_id = " + strconv.Itoa(index)
	rows, err := p.Db.Query(sql)
	if err != nil{
		return Pilot{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id int
		var firstName, lastName string

		if err := rows.Scan(&id, &firstName, &lastName); err != nil{
			return Pilot{}, err
		}
		return Pilot{id, firstName, lastName}, nil
	}
	return Pilot{}, errors.New("Пилот не найден")
}

func (p *PilotProvider)Delete(index int) ([]Pilot, error) {
	sql := "DELETE FROM airport.pilots CASCADE WHERE c_id = " + strconv.Itoa(index)
	result, err := p.Db.Exec(sql)
	if err != nil{
		return nil, errors.New("Данный пилот участвует в рейсе, его нельзя удалить")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Рейс не найден")
	}
	return p.GetAll()
}

func (p *PilotProvider)Edit(index int, itemToAdd []byte) ([]Pilot, error){
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправилные данные пилота")
	}
	sql := "UPDATE airport.pilots SET c_first_name = '" + temp.FirstName + "', c_last_name = '" + temp.LastName + "' WHERE pilots.c_id = "+ strconv.Itoa(index)
	result, err := p.Db.Exec(sql)
	if err != nil{
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	fmt.Println(rowsAffected)
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Пилот не найден")
	}
	return p.GetAll()
}

func (p *PilotProvider)Add(itemToAdd []byte) ([]Pilot, error) {
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправильные данные пилота")
	}
	sql := "INSERT INTO airport.pilots(c_id, c_first_name, c_last_name) VALUES (nextval('airport.planes_seq'),'" + temp.FirstName + "', '" + temp.LastName + "' )"
	result, err := p.Db.Exec(sql)
	if err != nil{
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	fmt.Println(rowsAffected)
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Пилот не найден")
	}
	return p.GetAll()
}
