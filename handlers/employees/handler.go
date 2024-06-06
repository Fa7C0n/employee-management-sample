package employees

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation"
)

// type Service interface {
// 	ListEmployees()
// 	GetEmployeeById()
// 	CreateEmployee()
// 	UpdateEmployee()
// 	DeleteEmployee()
// }

// type service struct{}

// func (s *service) ListEmployees()   {}
// func (s *service) GetEmployeeById() {}
// func (s *service) CreateEmployee()  {}
// func (s *service) UpdateEmployee()  {}
// func (s *service) DeleteEmployee()  {}

// func NewService() Service {
// 	return &service{}
// }

type ErrorMessage struct {
	Message string `json:"message"`
}

type Handler struct{}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if errMsg, err := json.Marshal(ErrorMessage{Message: message}); err != nil {
		log.Println(err)
	} else {
		w.Write(errMsg)
	}
}

var (
	es     = make(EmployeeStore)
	esLock sync.Mutex
)

func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pageNumber, err := strconv.Atoi(params.Get("pageNumber"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "page number should be a number")
		return
	}

	pageSize, err := strconv.Atoi(params.Get("pageSize"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "page size should be a number")
		return
	}

	start := (pageNumber - 1) * pageSize
	end := start + pageSize

	var employeeList []Employee
	for _, employee := range es {
		employeeList = append(employeeList, employee)
	}

	total := len(employeeList)

	if end > total {
		end = total
	}

	if start >= total {
		data, err := json.Marshal(struct {
			Employees []Employee `json:"employees"`
			Total     int        `json:"total"`
		}{
			Employees: []Employee{},
			Total:     total,
		})

		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		w.Write(data)
		return
	}

	data, err := json.Marshal(struct {
		Employees []Employee `json:"employees"`
		Total     int        `json:"total"`
	}{
		Employees: employeeList[start:end],
		Total:     total,
	})

	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
	w.Write(data)
}

func (h *Handler) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := validation.Validate(id, validation.Required)
	if err != nil {
		writeError(w, http.StatusBadRequest, "employee id is required")
		return
	}

	idStr, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id should be a number")
		return
	}

	emp, err := es.GetEmployeeById(idStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg, err := json.Marshal(emp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := w.Write(msg); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Printf("%s %s %s 200", r.Method, r.RequestURI, r.RemoteAddr)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad data provided")
		return
	}

	var emp Employee
	if err := json.Unmarshal(body, &emp); err != nil {
		writeError(w, http.StatusBadRequest, "Bad data provided")
		return
	}

	esLock.Lock()
	defer esLock.Unlock()
	err = es.CreateEmployee(emp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("id")
	empNumId, err := strconv.Atoi(empId)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id should be a number")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad data provided")
		return
	}

	var emp Employee
	if err := json.Unmarshal(body, &emp); err != nil {
		writeError(w, http.StatusBadRequest, "Bad data provided")
		return
	}

	esLock.Lock()
	defer esLock.Unlock()
	emp, err = es.UpdateEmployee(empNumId, emp.Name, emp.Position, emp.Salary)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg, err := json.Marshal(emp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := w.Write(msg); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Printf("%s %s %s 200", r.Method, r.RequestURI, r.RemoteAddr)
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empId := r.PathValue("id")
	empNumId, err := strconv.Atoi(empId)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id should be a number")
		return
	}

	err = es.DeleteEmployee(empNumId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg, err := json.Marshal(struct {
		Message string `json:"message"`
	}{Message: fmt.Sprintf("deleted id %v", empId)})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(msg)
}

func NewHandler() *Handler {
	return &Handler{}
}
