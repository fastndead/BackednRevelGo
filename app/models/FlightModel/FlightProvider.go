package FlightModel

import (
	"encoding/json"
	"errors"
	"database/sql"
	"strconv"
	"app/app/lib/dbManager"

)

type FlightProvider struct{
	Db *sql.DB
}
func (f *FlightProvider)GetAll() ([]Flight, error) {
	returnValue := []Flight{}

	sql := "SELECT flights.c_id, c_arrival_point, c_departure_point,c_fk_planes, c_name from airport.flights, airport.planes where planes.c_id = flights.c_fk_planes;"
	result, err := f.Db.Exec(sql)
	if err != nil{
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Нет ни одного рейса")
	}
	rows, err := f.Db.Query(sql)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var id, indexOfPlane int
		var arrivalPoint, departurePoint,  planeNameStr string

		if err := rows.Scan(&id, &arrivalPoint, &departurePoint,&indexOfPlane, &planeNameStr); err != nil{
			return nil, err
		}
		sql1 := "SELECT c_fk_pilot FROM airport.toc_flights_pilots WHERE " + strconv.Itoa(id) + " = toc_flights_pilots.c_fk_flight"
		PilotRows, err := f.Db.Query(sql1)
		if err != nil{
			return nil, err
		}
		defer PilotRows.Close()
		indexList := []int{}
		for PilotRows.Next(){
			var indexOfPilot int
			if err := PilotRows.Scan( &indexOfPilot); err != nil{
				return nil, err
			}
			indexList = append(indexList, indexOfPilot)
		}
		returnValue = append(returnValue, (Flight{Id: id, IdPilot: indexList, IdPlane: indexOfPlane, ArrivalPoint: arrivalPoint, DeparturePoint: departurePoint}))
	}
	return returnValue, nil
}

func (f *FlightProvider)GetById(index int) (Flight, error) {
	returnValue := Flight{}


	tx, err := f.Db.Begin()
	if err != nil{
		return Flight{}, err
	}
	defer transactionEnd(tx, err)
	sql := "SELECT flights.c_id, c_arrival_point, c_departure_point, c_fk_planes  from airport.flights WHERE flights.c_id = " + strconv.Itoa(index)
	result, err := f.Db.Exec(sql)
	if err != nil{
		return Flight{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return Flight{}, err
	}
	if  rowsAffected == 0 {
		return Flight{}, errors.New("Рейс не найден")
	}
	rows, err := tx.Query(sql)
	if err != nil{

		return Flight{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id, indexOfPlane int
		var arrivalPoint, departurePoint string

		if err := rows.Scan(&id, &arrivalPoint, &departurePoint, &indexOfPlane); err != nil{
			return Flight{}, err
		}
		sql1 := "SELECT c_fk_pilot FROM airport.toc_flights_pilots WHERE " + strconv.Itoa(id) + " = toc_flights_pilots.c_fk_flight"
		PilotRows, err := tx.Query(sql1)
		if err != nil{
			return Flight{}, err
		}
		defer PilotRows.Close()
		indexList := []int{}
		for PilotRows.Next(){
			var indexOfPilot int
			if err := PilotRows.Scan(&indexOfPilot); err != nil{
				return Flight{}, err
			}
			indexList = append(indexList, indexOfPilot)
		}

		returnValue = Flight{Id: id, IdPilot:indexList, IdPlane: indexOfPlane, ArrivalPoint: arrivalPoint, DeparturePoint: departurePoint}
	}
	return returnValue, nil
}

func (f *FlightProvider)Delete(index int) ([]Flight, error) {
	sql := "DELETE FROM airport.toc_flights_pilots CASCADE WHERE c_fk_flight = " + strconv.Itoa(index) + ";"
	sql += "DELETE FROM airport.flights CASCADE WHERE c_id = " + strconv.Itoa(index)+";"

	tx, err := f.Db.Begin()
	if err != nil{
		return nil, err
	}
	defer transactionEnd(tx, err)
	result, err := tx.Exec(sql)
	if err != nil{
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return nil, err
	}
	if  rowsAffected == 0 {
		return nil, errors.New("Рейс не найден")
	}
	return f.GetAll()
}

func (f *FlightProvider)Edit(index int, itemToAdd []byte) ([]Flight, error){
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправилные данные рейса")
	}

	tx, err := f.Db.Begin()
	if err != nil{

		return nil, err
	}
	defer transactionEnd(tx, err)
	sql1 := "UPDATE airport.flights SET c_arrival_point = '" + temp.ArrivalPoint + "', c_departure_point = '" + temp.DeparturePoint + "', c_fk_planes = " + strconv.Itoa(temp.IdPlane) + " WHERE c_id = " +strconv.Itoa(index)+ ";"
	_, err = tx.Exec(sql1)
	if err != nil{
		return nil, errors.New("Ошибка при вставке рейса")
	}
	curvalFlights, err := dbManager.GetCurVal("airport.flights_seq", tx)
	sql2 := "DELETE FROM airport.toc_flights_pilots WHERE c_fk_flight = "+ strconv.Itoa(curvalFlights) +";"
	_, err = tx.Exec(sql2)
	if err != nil{
		return nil, errors.New("Ошибка при удалении зависимостей")
	}
	for _,elem := range temp.IdPilot{
		err := addRelationFlightPilot(curvalFlights, elem, tx)
		if err != nil{
			return nil, err
		}
	}
	return f.GetAll()
}

func (f FlightProvider)Add(itemToAdd []byte) ([]Flight, error) {
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return nil, errors.New("Неправилные данные рейса")
	}

	tx, err := f.Db.Begin()
	if err != nil{
		return nil, err
	}
	defer transactionEnd(tx, err)
	sql1 := "INSERT INTO airport.flights(c_id, c_arrival_point, c_departure_point, c_fk_planes) VALUES (nextval('airport.flights_seq'),'" + temp.ArrivalPoint + "', '" + temp.DeparturePoint + "', " + strconv.Itoa(temp.IdPlane) + ");"
	_, err = tx.Exec(sql1)
	if err != nil{
		return nil, errors.New("Ошибка при вставке рейса")
	}
	curvalFlights, err := dbManager.GetCurVal("airport.flights_seq", tx)
	if err != nil{
		return nil, err
	}

	for _,elem := range temp.IdPilot{
		err := addRelationFlightPilot(curvalFlights, elem, tx)
		if err != nil{
			return nil, err
		}
	}
	return f.GetAll()
}

func addRelationFlightPilot(flight int, pilot int, tx *sql.Tx)error{
	sql3 := "INSERT INTO airport.toc_flights_pilots(c_id, c_fk_flight, c_fk_pilot) VALUES (nextval('airport.toc_flights_pilots_seq'), " + strconv.Itoa(flight) + ", " + strconv.Itoa(pilot) + ");"
	_, err := tx.Exec(sql3)
	if err != nil{
		return  err
	}
	return nil
}

func transactionEnd(tx *sql.Tx, err error){
	if err != nil{
		tx.Rollback()
	} else {
		tx.Commit()
	}
}