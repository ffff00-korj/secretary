package product

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Product struct {
	name       string
	sum        int
	paymentDay int
}

func NewProduct(name string, sum, day int) *Product {
	return &Product{name: name, sum: sum, paymentDay: day}
}

func (p *Product) String() string {
	return fmt.Sprintf("Name: %s,\nSum: %d,\nPayment day: %d", p.name, p.sum, p.paymentDay)
}

func (p *Product) GetName() string {
	return p.name
}

func (p *Product) GetSum() int {
	return p.sum
}

func (p *Product) GetPaymentDay() int {
	return p.paymentDay
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
	return NewProduct(argl[0], sum, day), nil
}
