package employees

import "errors"

var (
	errEmployeeExists   = errors.New("employee already exists")
	errEmployeeNotFound = errors.New("employee not found")
)

type Employee struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

type EmployeeStore map[int]Employee

func (es EmployeeStore) CreateEmployee(employee Employee) error {
	if _, exists := es[employee.Id]; exists {
		return errEmployeeExists
	}
	es[employee.Id] = employee

	return nil
}

func (es EmployeeStore) GetEmployeeById(id int) (Employee, error) {
	employee, exists := es[id]
	if !exists {
		return Employee{}, errEmployeeNotFound
	}

	return employee, nil
}

func (es EmployeeStore) UpdateEmployee(id int, name, position string, salary float64) (Employee, error) {
	emp, exists := es[id]

	if !exists {
		return Employee{}, errEmployeeNotFound
	}

	if name != "" {
		emp.Name = name
	}

	if position != "" {
		emp.Position = position
	}

	if salary > 0 {
		emp.Salary = salary
	}

	es[id] = emp

	return emp, nil
}

func (es EmployeeStore) DeleteEmployee(id int) error {
	_, exists := es[id]

	if !exists {
		return errEmployeeNotFound
	}

	delete(es, id)
	return nil
}
