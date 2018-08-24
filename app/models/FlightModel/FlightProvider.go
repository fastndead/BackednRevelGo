package FlightModel

import (
	"encoding/json"
	"database/sql"
	"app/app/lib/dbManager"
	"fmt"
)

type FlightProvider struct{
	db *sql.DB
}

func (f *FlightProvider)Init()error{
	var err error
	f.db, err = dbManager.OpenConnection()
	if err != nil{
		return fmt.Errorf("Ошибка при подключении к базе: %err", err)
	}
	return nil	
}

func (f *FlightProvider)Close()error{
	return  dbManager.CloseConnection(f.db)
}

func (f *FlightProvider)GetAll() ([]Flight, error) {
	returnValue := []Flight{}

	request := "SELECT flights.c_id, c_arrival_point, c_departure_point,c_fk_planes, c_name from airport.flights, airport.planes where planes.c_id = flights.c_fk_planes;"
	rows, err := f.db.Query(request)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var id, indexOfPlane sql.NullInt64
		var arrivalPoint, departurePoint,  planeNameStr sql.NullString

		if err := rows.Scan(&id, &arrivalPoint, &departurePoint,&indexOfPlane, &planeNameStr); err != nil{
			return nil, err
		}
		indexList, err := getPilots(f.db, id)
		if err != nil{
			return nil, err
		}
		returnValue = append(returnValue, (Flight{Id: int(id.Int64), IdPilot: indexList, IdPlane: int(indexOfPlane.Int64), ArrivalPoint: arrivalPoint.String, DeparturePoint: departurePoint.String}))
	}
	return returnValue, nil
}

func (f *FlightProvider)GetById(index int) (Flight, error) {
	returnValue := Flight{}

	request := "SELECT flights.c_id, c_arrival_point, c_departure_point, c_fk_planes  from airport.flights WHERE flights.c_id = $1"
	rows, err := f.db.Query(request, index)
	if err != nil{
		if err == sql.ErrNoRows{
			return Flight{}, nil
		}
		return Flight{}, err
	}
	defer rows.Close()
	for rows.Next(){
		var id, indexOfPlane sql.NullInt64
		var arrivalPoint, departurePoint sql.NullString

		if err := rows.Scan(&id, &arrivalPoint, &departurePoint, &indexOfPlane); err != nil{
			return Flight{}, err
		}
		sql1 := "SELECT c_fk_pilot FROM airport.toc_flights_pilots WHERE $1 = toc_flights_pilots.c_fk_flight"
		PilotRows, err := f.db.Query(sql1, id)
		if err != nil{
			return Flight{}, err
		}
		defer PilotRows.Close()
		indexList, err := getPilots(f.db, id)
		if err != nil{
			return Flight{}, err
		}

		returnValue = Flight{Id: int(id.Int64), IdPilot: indexList, IdPlane: int(indexOfPlane.Int64), ArrivalPoint: arrivalPoint.String, DeparturePoint: departurePoint.String}
	}
	return returnValue, nil
}

func delete(db *sql.DB, index int)error{
	tx, err := db.Begin()
	if err != nil{
		return err
	}
	defer transactionEnd(tx, err)

	request := "DELETE FROM airport.toc_flights_pilots CASCADE WHERE c_fk_flight = $1"
	_, err = tx.Exec(request, index)
	if err != nil{
		return fmt.Errorf("Ошибка при удалении зависимостей: $1",err)
	}
	request = "DELETE FROM airport.flights CASCADE WHERE c_id = $1"
	_, err = tx.Exec(request, index)
	if err != nil{
		if err == sql.ErrNoRows{
			return fmt.Errorf("Рейс не найден: $1",err)
		}
		return fmt.Errorf("Ошибка при удалении рейса: $1",err)
	}
	return nil
}

func (f *FlightProvider)Delete(index int) ([]Flight, error) {
	err := delete(f.db, index)
	if err != nil{
		return nil, err
	}
	return f.GetAll()
}

func edit(db *sql.DB,index int, itemToAdd []byte)(error){
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return fmt.Errorf("Неправильные данные рейса: %err", err)
	}

	tx, err := db.Begin()
	if err != nil{

		return err
	}
	defer transactionEnd(tx, err)
	sql1 := "UPDATE airport.flights SET c_arrival_point = $1, c_departure_point = $2, c_fk_planes = $3 WHERE c_id = $4;"
	_, err = tx.Exec(sql1, temp.ArrivalPoint, temp.DeparturePoint, temp.IdPlane, index)
	if err != nil{
		return fmt.Errorf("Ошибка при вставке рейса: %err", err)
	}
	sql2 := "DELETE FROM airport.toc_flights_pilots WHERE c_fk_flight = $1;"
	_, err = tx.Exec(sql2, index)
	if err != nil{
		return fmt.Errorf("Ошибка при удалении зависимостей: %err", err)
	}
	for _,elem := range temp.IdPilot{
		err := addRelationFlightPilot(sql.NullInt64{int64(index), true}, sql.NullInt64{int64(elem), true}, tx)
		if err != nil{
			return err
		}
	}
	return nil
}

func (f *FlightProvider)Edit(index int, itemToAdd []byte) (Flight, error){
	err := edit(f.db, index, itemToAdd)
	if err != nil{
		return Flight{}, err
	}
	return f.GetById(index)
}

func add(db *sql.DB, itemToAdd []byte)error{
	temp := &Flight{}
	err := json.Unmarshal(itemToAdd, temp)
	if err != nil {
		return fmt.Errorf("Неправильные данные рейса: %err", err)
	}

	tx, err := db.Begin()
	if err != nil{
		return err
	}
	defer transactionEnd(tx, err)
	sql1 := "INSERT INTO airport.flights(c_id, c_arrival_point, c_departure_point, c_fk_planes) VALUES (nextval('airport.flights_seq'),'$1', '$2', $3);"
	_, err = tx.Exec(sql1, temp.ArrivalPoint, temp.DeparturePoint, temp.IdPlane)
	if err != nil{
		return fmt.Errorf("Ошибка при вставке рейса: %err", err)
	}
	curvalFlights, err := dbManager.GetCurVal(sql.NullString{"airport.flights_seq", true}, db)
	if err != nil{
		return err
	}

	for _,elem := range temp.IdPilot{
		err := addRelationFlightPilot(sql.NullInt64{int64(curvalFlights), true}, sql.NullInt64{int64(elem	), true}, tx)
		if err != nil{
			return err
		}
	}
	return nil
}

func (f *FlightProvider)Add(itemToAdd []byte) (Flight, error) {
	err := add(f.db, itemToAdd)
	if err != nil{
		return Flight{}, err
	}
	curvalFlights, err := dbManager.GetCurVal(sql.NullString{"airport.flights_seq", true}, f.db)
	if err != nil{
		return Flight{},err
	}
	return f.GetById(curvalFlights)
}

func addRelationFlightPilot(flight sql.NullInt64, pilot sql.NullInt64, tx *sql.Tx)error{
	sql3 := "INSERT INTO airport.toc_flights_pilots(c_id, c_fk_flight, c_fk_pilot) VALUES (nextval('airport.toc_flights_pilots_seq'), $1, $2);"
	_, err := tx.Exec(sql3, flight, pilot)
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

func getPilots(db *sql.DB, id sql.NullInt64)([]int, error){
	sql1 := "SELECT c_fk_pilot FROM airport.toc_flights_pilots WHERE $1 = toc_flights_pilots.c_fk_flight"
	PilotRows, err := db.Query(sql1, id)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Нет ни одного пилота: $1", err)
		}
		return nil, err
	}
	defer PilotRows.Close()
	indexList := []int{}
	for PilotRows.Next(){
		var indexOfPilot sql.NullInt64
		if err := PilotRows.Scan( &indexOfPilot); err != nil{
			return nil, err
		}
		indexList = append(indexList, int(indexOfPilot.Int64))
	}
	return indexList, nil
}