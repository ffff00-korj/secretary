package salary

import "fmt"

func NewSalary(name string, sum, paymentDay int) *Salary {
	return &Salary{name: name, sum: sum, paymentDay: paymentDay}
}

func (s *Salary) String() string {
	return fmt.Sprintf("%s, %d, %d", s.name, s.sum, s.paymentDay)
}

func (s *Salary) GetName() string {
	return s.name
}

func (s *Salary) GetSum() int {
	return s.sum
}

func (s *Salary) GetPaymentDay() int {
	return s.paymentDay
}
