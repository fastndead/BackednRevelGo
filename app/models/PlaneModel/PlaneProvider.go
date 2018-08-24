package PlaneModel

import (
	"encoding/json"
	"errors"
	"database/sql"
	"strconv"
	"fmt"
)

type PlaneProvider struct{
	Db *sql.DB
}

func (p *PlaneProvider)GetAll() ([]Plane, error) {
	sql := "SELECT c_id, c_name FROM airport.planes"
	rows, err := p.Db.Query(sql)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var PlaneList []Plane
	for rows.Next(){
		var id int
		var planeName string

		if err := rows.Scan(&id, &planeName); err != nil{
			return nil, err
		}

		PlaneList = append(PlaneList, Plane{id, planeName})
	}
	return PlaneList, nil
}

func (p *PlaneProvider)GetById(index int) (Plane, error) {
	sql := "SELECT c_id, c_name FROM airport.planes WHERE c_id = " + strconv.Itoa(index)
	rows, err := p.Db.Query(sql)
	if err != nil{
		return Plane{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id int
		var planeName string

		if err := rows.Scan(&id, &planeName); err != nil{
			return Plane{}, err
		}
		return Plane{id, planeName}, nil
	}
	return Plane{}, errors.New("Самолёт не найден")
}

func (p *PlaneProvider)Delete(index int) ([]Plane, error) {
	sql := "DELETE FROM airport.planes CASCADE WHERE c_id = " + strconv.Itoa(index)
	result, err := p.Db.Exec(sql)
	if err != nil{
		return nil, errors.New("Данный самолёт участвует в рейсе, его нельзя удалить")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Самолёт не найден")
	}
	return p.GetAll()
}

func (p *PlaneProvider)Edit(index int, itemToAdd []byte) ([]Plane, error){
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправилные данные самолёта")
	}
	sql := "UPDATE airport.planes SET c_name = '" + temp.Name + "' WHERE planes.c_id = "+ strconv.Itoa(index)
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
		return nil, errors.New("Самолёт не найден")
	}
	return p.GetAll()
}

func (p *PlaneProvider)Add(itemToAdd []byte) ([]Plane, error) {
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправильные данные самолёта")
	}
	sql := "INSERT INTO airport.planes(c_id, c_name) VALUES (nextval('airport.planes_seq'),'" + temp.Name + "' )"
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
		return nil, errors.New("Самолёт не найден")
	}
	return p.GetAll()
}
