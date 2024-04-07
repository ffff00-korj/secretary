package product

import (
	"errors"
	"fmt"
	"strconv"
)

func NewProduct(args []string) (*Product, error) {
	if len(args) != 3 {
		return nil, errors.New("Not enough arguments!")
	}
	sum, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Second argument should be a number!")
	}
	day, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument should be a number!")
	}
	return &Product{name: args[0], sum: sum, paymentDay: day}, nil
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
