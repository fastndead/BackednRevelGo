package dbManager

import ("database/sql"
	"github.com/revel/revel"
	"fmt"
)

func OpenConnection() (*sql.DB, error){
	constr := revel.Config.StringDefault("connectionString", "")
	db, err := sql.Open("postgres", constr )
	if err != nil{
		return nil, fmt.Errorf("Ошибка при подключении к базе: %err", err)
	}
	err = db.Ping()
	if err != nil{
		return nil, fmt.Errorf("Ошибка при проверке подключения к базе: %err", err)
	}
	return db, nil
}

func CloseConnection(db *sql.DB)error{
	err := db.Close()
	if err != nil{
		return err
	}
	return nil
}

func GetCurVal(seq sql.NullString, db *sql.DB)(int, error){
	request := "SELECT last_value FROM " + seq.String
	rows, err := db.Query(request)
	if err != nil{
		return -1, err
	}
	defer rows.Close()
	var index int
	for rows.Next() {
		if err := rows.Scan(&index); err != nil{
			return -1, err
		}
	}
	return index, nil
}


