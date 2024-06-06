package employees

import "testing"

func TestCreateEmployee(t *testing.T) {
	es := make(EmployeeStore)
	emp := Employee{
		Id:       1,
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}
	err := es.CreateEmployee(emp)
	if err != nil {
		t.Errorf("Error creating employee: %v", err)
	}

	employee, exists := es[1]
	if !exists {
		t.Errorf("Employee failed to create")
	}

	if employee.Name != "John Doe" || employee.Position != "Developer" || employee.Salary != 60000 {
		t.Errorf("Employee data mismatch: %+v", employee)
	}
}

func TestGetEmployeeByID(t *testing.T) {
	es := make(EmployeeStore)
	emp := Employee{
		Id:       1,
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}
	err := es.CreateEmployee(emp)
	if err != nil {
		t.Errorf("Failed to create employee")
	}

	fetchedEmployee, err := es.GetEmployeeById(emp.Id)
	if err != nil {
		t.Errorf("Error fetching employee: %v", err)
	}
	if fetchedEmployee != emp {
		t.Errorf("Fetched employee mismatch: %+v", fetchedEmployee)
	}
}

func TestUpdateEmployee(t *testing.T) {
	es := make(EmployeeStore)
	emp := Employee{
		Id:       1,
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}
	err := es.CreateEmployee(emp)
	if err != nil {
		t.Errorf("Failed to create employee")
	}

	updatedEmployee, err := es.UpdateEmployee(emp.Id, "John Smith", "Senior Analyst", 75000)
	if err != nil {
		t.Errorf("Error updating employee: %v", err)
	}
	if updatedEmployee.Position != "Senior Analyst" || updatedEmployee.Salary != 75000 {
		t.Errorf("Updated employee data mismatch: %+v", updatedEmployee)
	}
}

func TestDeleteEmployee(t *testing.T) {
	es := make(EmployeeStore)
	emp := Employee{
		Id:       1,
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}
	err := es.CreateEmployee(emp)
	if err != nil {
		t.Errorf("Failed to create employee")
	}
	err = es.DeleteEmployee(emp.Id)
	if err != nil {
		t.Errorf("Error deleting employee: %v", err)
	}
	_, err = es.GetEmployeeById(emp.Id)
	if err == nil {
		t.Errorf("Employee should have been deleted")
	}
}
