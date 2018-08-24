package dbManager

import ("database/sql"
	"github.com/revel/revel"
	"go_modules/src/github.com/pkg/errors"
)

func OpenConnection() (*sql.DB, error){
	constr := revel.Config.StringDefault("connectionString", "")
	db, err := sql.Open("postgres", constr )
	if err != nil{
		return nil, errors.New("Ошибка при подключении к базе")
	}
	err = db.Ping()
	if err != nil{
		return nil, errors.New("Ошибка при проверке подключения к базе")
	}
	return db, nil
}

func CloseConnection(db *sql.DB){
	db.Close()
}

func GetCurVal(seq string, db *sql.Tx)(int, error){
	sql := "SELECT last_value FROM "+seq+";"
	rows, err := db.Query(sql)
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


