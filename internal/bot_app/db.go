package bot_app

import (
	"database/sql"
	"time"

	"github.com/ffff00-korj/secretary/internal/product"
)

func (app *bot_app) checkProductExists(p *product.Product) (bool, error) {
	query := `SELECT
        1 AS exists
    FROM
        products AS p
    WHERE
        p.name = $1
    LIMIT
        1`

	var exists bool
	err := app.db.QueryRow(query, p.GetName()).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (app *bot_app) addProduct(p *product.Product) (int, error) {
	query := `INSERT INTO products(Name, Sum, PaymentDay)
        VALUES($1, $2, $3) RETURNING id`

	var id int
	err := app.db.QueryRow(query, p.GetName(), p.GetSum(), p.GetPaymentDay()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *bot_app) getTotal() (int, error) {
	dayNow := time.Now().Day()
	var query string
	if dayNow < 5 || dayNow >= 20 {
		query = `SELECT
            SUM(p.sum) AS total
        FROM
            products AS p
        WHERE
            p.paymentDay >= 5 AND
            p.paymentDay < 20`
	} else {
		query = `SELECT
            SUM(p.sum) AS total
        FROM
            products AS p
        WHERE
            p.paymentDay < 5 OR
            p.paymentDay >= 20`
	}
	var total int
	err := app.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
