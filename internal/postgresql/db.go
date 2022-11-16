package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id          uint `json:"user_id"`
	Cash, Count uint
	OrderID     uint `json:"order_id"`
	ServiceID   uint `json:"service_id"`
}

func (u *User) checkUser() (ok bool, err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("select id from users where id=$1")
	if err != nil {
		return
	}
	var temp int
	err = stmt.QueryRow(u.Id).Scan(&temp)
	if err != nil {

		return
	}
	return true, nil
}

func DbConnection() (*sql.DB, error) {
	connStr := "user=test password=123 dbname=postgres sslmode=disable host=postgres port=5432"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetBalance(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok != true) {
		return
	}

	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("select cash from users where id=$1")
	if err != nil {
		return
	}

	err = stmt.QueryRow(user.Id).Scan(&user.Cash)
	if err != nil {
		return
	}
	return
}

func Replenish(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok != true) {
		return err
	}

	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into transactions (id, user_id,cost,type) values (default,$1,$2,'replenishment')returning id")
	if err != nil {
		return
	}
	var val string
	err = stmt.QueryRow(user.Id, user.Count).Scan(&val)
	if err != nil {
		return
	}
	return
}

func WriteTransaction(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok != true) {
		return
	}

	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into transactions values (default,$1,$2,$3,$4, 'buy')returning id")
	if err != nil {
		return
	}
	var val string
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&val)
	if err != nil {
		return
	}

	return
}

func RecognizeRevenue(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok != true) {
		return
	}

	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(" select id from transactions where user_id=$1 and service_id=$2 and order_id=$3 and cost=$4 and type='buy'")
	if err != nil {
		return
	}

	var result int
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&result)
	if err != nil {
		return
	}
	fmt.Println(result)

	stmt, err = db.Prepare("insert into transactions values (default,$1,$2,$3,$4, 'revenue')returning id")
	if err != nil {
		return
	}
	var val string
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&val)
	if err != nil {
		return
	}

	return
}
