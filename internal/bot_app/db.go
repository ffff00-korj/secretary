package bot_app

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/ffff00-korj/secretary/internal/product"
)

func (app *bot_app) checkProductExists(p *product.Product) (bool, error) {
	query := `SELECT 1 AS EXISTS
    FROM products AS p
    WHERE p.name = $1
    LIMIT 1`
	var exists bool
	err := app.db.QueryRow(query, p.GetName()).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("Can't check products exists, %w", err)
	}
	return true, nil
}

func (app *bot_app) addProduct(p *product.Product) (int, error) {
	query := `INSERT INTO products(Name, Sum, PaymentDay)
    VALUES($1, $2, $3) RETURNING id`

	var id int
	err := app.db.QueryRow(query, p.GetName(), p.GetSum(), p.GetPaymentDay()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Can't add product %w", err)
	}
	return id, nil
}

func (app *bot_app) getExpensePeriod(day int) (int, int, bool, error) {
	query := `SELECT
        s.paymentDay AS paymentDay
    FROM salaries AS s
    ORDER BY s.paymentDay`

	rows, err := app.db.Query(query)
	if err != nil {
		return 0, 0, false, err
	}
	var (
		paymentDay  int
		paymentDays []int
	)
	for rows.Next() {
		rows.Scan(&paymentDay)
		paymentDays = append(paymentDays, paymentDay)
	}
	if paymentDays[len(paymentDays)-1] <= day {
		return paymentDays[len(paymentDays)-1], paymentDays[0], false, nil
	}
	return paymentDays[0], paymentDays[len(paymentDays)-1], true, nil
}

func (app *bot_app) currentPaymentPeriod() (time.Time, time.Time, error) {
	query := `SELECT
       i.paymentdate
    FROM
      (SELECT Cast(Make_date(2, $1, paymentday) AS DATE) AS paymentdate
       FROM salaries s
       UNION ALL SELECT cast(Make_date(2, $1, paymentday) - interval '1 month' AS date)
       FROM salaries s
       UNION ALL SELECT cast(make_date(2, $1, paymentday) + interval '1 month' AS date)
       FROM salaries s) AS i
    WHERE i.paymentdate >= make_date(2, $1, $2)
    LIMIT 2`
	rows, err := app.db.Query(query, time.Now().Month(), time.Now().Day())
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	var (
		prev time.Time
		next time.Time
	)
	rows.Next()
	rows.Scan(&prev)
	rows.Next()
	rows.Scan(&next)
	return prev, next, nil
}

func (app *bot_app) getExpenseReport() (string, error) {
	prev, next, err := app.currentPaymentPeriod()
	if err != nil {
		return "", fmt.Errorf("Can't get expense periods, %w", err)
	}
	query := `SELECT
       p.name,
       p.sum,
       p.paymentday
    FROM products p
    INNER JOIN
      (SELECT p.id,
              cast(make_date(2, $1, p.paymentday) AS date) AS paymentdate
       FROM products p
       UNION ALL SELECT p.id,
                        cast(make_date(2, $1, paymentday) - interval '1 month' AS date)
       FROM products p
       UNION ALL SELECT p.id,
                        cast(make_date(2, $1, paymentday) + interval '1 month' AS date)
       FROM products p) AS i ON p.id = i.id
    WHERE i.paymentdate > $2
      AND i.paymentdate <= $3
    ORDER BY i.paymentdate`
	rows, err := app.db.Query(query, time.Now().Month(), prev, next)
	if err != nil {
		return "", fmt.Errorf("Can't get expense report data, %w", err)
	}
	var (
		er         expenseReport
		total      int
		name       string
		sum        int
		paymentDay int
	)
	for rows.Next() {
		rows.Scan(&name, &sum, &paymentDay)
		er.rows = append(er.rows, expenseReportRow{Name: name, Sum: sum, PaymentDay: paymentDay})
		total += sum
	}
	er.total = expenseReportRow{Name: "total", Sum: total}

	return er.String(), nil
}

func (er *expenseReport) String() string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"name", "sum", "payment day"})
	for _, rr := range er.rows {
		t.AppendRow(table.Row{rr.Name, rr.Sum, rr.PaymentDay})
	}
	t.AppendFooter(table.Row{er.total.Name, er.total.Sum})

	return t.Render()
}
