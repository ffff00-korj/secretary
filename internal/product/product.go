package product

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Product struct {
	Name       string `db:"Name"`
	Sum        int    `db:"Sum"`
	PaymentDay int    `db:"PaymentDay"`
}

func (p *Product) String() string {
	return fmt.Sprintf("Name: %s,\nSum: %d,\nPayment day: %d", p.Name, p.Sum, p.PaymentDay)
}

func NewProductFromArgs(args string) (*Product, error) {
	argl := strings.Split(args, " ")
	if len(argl) != 3 {
		return nil, errors.New("Not enough arguments!")
	}
	sum, err := strconv.Atoi(argl[1])
	if err != nil {
		return nil, errors.New("Second argument should be a number!")
	}
	day, err := strconv.Atoi(argl[2])
	if err != nil {
		return nil, errors.New("Third argument should be a number!")
	}
	return &Product{Name: argl[0], Sum: sum, PaymentDay: day}, nil
}
