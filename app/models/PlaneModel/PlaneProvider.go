package PlaneModel

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"app/app/lib/dbManager"
)

type PlaneProvider struct{
	db *sql.DB
}

func (p *PlaneProvider)Init()error{
	var err error
	p.db, err = dbManager.OpenConnection()
	if err != nil{
		return fmt.Errorf("Ошибка при подключении к базе: %err", err)
	}
	return nil
}

func (p *PlaneProvider)Close()error {
	return dbManager.CloseConnection(p.db)
}

func (p *PlaneProvider)GetAll() ([]Plane, error) {
	request := "SELECT c_id, c_name FROM airport.planes"
	rows, err := p.db.Query(request)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var PlaneList []Plane
	for rows.Next(){
		var id sql.NullInt64
		var planeName sql.NullString

		if err := rows.Scan(&id, &planeName); err != nil{
			return nil, err
		}

		PlaneList = append(PlaneList, Plane{int(id.Int64), planeName.String})
	}
	return PlaneList, nil
}

func (p *PlaneProvider)GetById(index int) (Plane, error) {
	request := "SELECT c_id, c_name FROM airport.planes WHERE c_id = $1"
	rows, err := p.db.Query(request,index)
	if err != nil{
		if err == sql.ErrNoRows{
			return Plane{}, nil
		}
		return Plane{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id sql.NullInt64
		var planeName sql.NullString

		if err := rows.Scan(&id, &planeName); err != nil{
			return Plane{}, err
		}
		return Plane{int(id.Int64), planeName.String}, nil
	}
	return Plane{}, fmt.Errorf("Самолёт не найден: %err", err)
}

func (p *PlaneProvider)Delete(index int) ([]Plane, error) {
	sql := "DELETE FROM airport.planes CASCADE WHERE c_id = $1"
	result, err := p.db.Exec(sql, index)
	if err != nil{
		return nil, fmt.Errorf("Данный самолёт учавствует в рейсе, его нельзя удалить: %err", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, fmt.Errorf("Самолёт не найден: %err", err)
	}
	return p.GetAll()
}

func (p *PlaneProvider)Edit(index int, itemToAdd []byte) (Plane, error){
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return Plane{}, fmt.Errorf("Неправильные данные самолёта: %err", err)
	}
	request := "UPDATE airport.planes SET c_name = $1 WHERE planes.c_id = $2"
	_, err = p.db.Exec(request,temp.Name, index)
	if err != nil{
		if err == sql.ErrNoRows{
			return Plane{}, fmt.Errorf("Самолёт не найден: $1", err)
		}
		return Plane{}, err
	}
	return p.GetById(index)
}

func (p *PlaneProvider)Add(itemToAdd []byte) (Plane, error) {
	temp := &Plane{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return Plane{}, fmt.Errorf("Неправильные данные самолёта: %err", err)
	}
	request := "INSERT INTO airport.planes(c_id, c_name) VALUES (nextval('airport.planes_seq'),$1 )"
	result, err := p.db.Exec(request, temp.Name)
	if err != nil{
		return Plane{}, err
	}
	rowsAffected, err := result.RowsAffected()
	fmt.Println(rowsAffected)
	if err != nil{
		return Plane{}, err
	}
	if  rowsAffected == 0 {
		return Plane{}, fmt.Errorf("Самолёт не найден: %err", err)
	}
	currval,err := dbManager.GetCurVal(sql.NullString{"airport.planes_seq", true}, p.db)
	return p.GetById(currval)
}
