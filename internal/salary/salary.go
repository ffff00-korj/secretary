package salary

import "fmt"

func NewSalary(name string, sum, paymentDay int) *salary {
	return &salary{name: name, sum: sum, paymentDay: paymentDay}
}

func (s *salary) String() string {
	return fmt.Sprintf("%s, %d, %d", s.name, s.sum, s.paymentDay)
}

func (s *salary) GetName() string {
	return s.name
}

func (s *salary) GetSum() int {
	return s.sum
}

func (s *salary) GetPaymentDay() int {
	return s.paymentDay
}
