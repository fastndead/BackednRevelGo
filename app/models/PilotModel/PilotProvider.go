package PilotModel

import (
	"encoding/json"
	"errors"
	"database/sql"
	"fmt"
	"app/app/lib/dbManager"
)

type PilotProvider struct {
	db *sql.DB
}


func (p *PilotProvider)Init()error{
	var err error
	p.db, err = dbManager.OpenConnection()
	if err != nil{
		return fmt.Errorf("Ошибка при подключении к базе: %err", err)
	}
	return nil
}

func (p *PilotProvider)Close()error {
	return dbManager.CloseConnection(p.db)
}


func (p *PilotProvider)GetAll() ([]Pilot, error) {
	request := "SELECT c_id, c_first_name, c_last_name FROM airport.pilots"
	rows, err := p.db.Query(request)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var PlaneList []Pilot
	for rows.Next(){
		var id sql.NullInt64
		var firstName, lastName sql.NullString

		if err := rows.Scan(&id, &firstName, &lastName); err != nil{
			return nil, err
		}

		PlaneList = append(PlaneList, Pilot{int(id.Int64), firstName.String, lastName.String})
	}
	return PlaneList, nil
}

func (p *PilotProvider)GetById(index int) (Pilot, error) {
	request := "SELECT c_id, c_first_name, c_last_name FROM airport.pilots WHERE c_id = $1"
	rows, err := p.db.Query(request, index)
	if err != nil{
		return Pilot{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id sql.NullInt64
		var firstName, lastName sql.NullString

		if err := rows.Scan(&id, &firstName, &lastName); err != nil{
			return Pilot{}, err
		}
		return Pilot{int(id.Int64), firstName.String, lastName.String}, nil
	}
	return Pilot{}, fmt.Errorf("Пилот не найден: $1", err)
}

func (p *PilotProvider)Delete(index int) ([]Pilot, error) {
	request := "DELETE FROM airport.pilots CASCADE WHERE c_id = $1"
	_, err := p.db.Exec(request, index)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, fmt.Errorf("Данный пилот участвует в рейсе, его нельзя удалить: $1", err)
	}
	return p.GetAll()
}

func (p *PilotProvider)Edit(index int, itemToAdd []byte) ([]Pilot, error){
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, fmt.Errorf("Неправилные данные пилота", err)
	}
	request := "UPDATE airport.pilots SET c_first_name = $1, c_last_name = $2 WHERE pilots.c_id = $3"
	_, err = p.db.Exec(request, temp.FirstName, temp.LastName, index)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Пилот не найден", err)
		}
		return nil, err
	}
	return p.GetAll()
}

func (p *PilotProvider)Add(itemToAdd []byte) ([]Pilot, error) {
	temp := &Pilot{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправильные данные пилота")
	}
	request := "INSERT INTO airport.pilots(c_id, c_first_name, c_last_name) VALUES (nextval('airport.planes_seq'),$1, $2 )"
	_, err = p.db.Exec(request, temp.FirstName, temp.LastName)
	if err != nil{
		return nil, err
	}
	return p.GetAll()
}
