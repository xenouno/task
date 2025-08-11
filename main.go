package main

import (
	"net/http"
	"task/auth"
	"github.com/gin-gonic/gin"
	"strconv"
)

type department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var departments = []department{
	{ID: "1", Name: "HR"},
	{ID: "2", Name: "Marketplace"},
}


type employee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Department string `json:"department"`
}

var employees=[]employee{
	{ID:1,Name:"Jatin",Age:34,Department:"HR"},
	{ID:2,Name:"Mohit",Age:24,Department:"Marketplace"},
	{ID:3,Name:"Sam",Age:25,Department:"HR"},
}




func getDepartments(c *gin.Context) {
	c.JSON(http.StatusOK, departments)
}
func getEmployees(c *gin.Context){
	c.IndentedJSON(http.StatusOK,employees)
}

func createDepartment(c *gin.Context) {
	var newItem department


	if err := c.BindJSON(&newItem); err != nil || newItem.ID == "" || newItem.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields (id, name)"})
		return
	}
	
	departments = append(departments, newItem)
	c.JSON(http.StatusCreated, newItem)
}
func createEmployee(c *gin.Context) {
	var newItem employee


	if err := c.BindJSON(&newItem); err != nil ||  newItem.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields (id, name)"})
		return
	}
	
	employees = append(employees, newItem)
	c.JSON(http.StatusCreated, newItem)
}
func returnEmps(c *gin.Context){
	var matchedEmps []employee

	dep:=c.Param("name")
	for _,emp:=range employees{
		if emp.Department==dep{
			matchedEmps = append(matchedEmps, emp)
		}
	}
	if len(matchedEmps) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No employees found for department"})
		return
	}

	c.IndentedJSON(http.StatusOK, matchedEmps)
}
func updateDepartments(c *gin.Context) {
	id := c.Param("id")
	var updated department

	if err := c.BindJSON(&updated); err != nil || updated.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'name'"})
		return
	}

	for i, dep := range departments {
		if dep.ID == id {
			departments[i].Name = updated.Name
			c.JSON(http.StatusOK, departments[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}
func updateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var updated employee
	if err := c.BindJSON(&updated); err != nil || updated.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields"})
		return
	}

	for i, emp := range employees {
		if emp.ID == id {
			updated.ID = id // keep the original ID
			employees[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
}

func deleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	for i, emp := range employees {
		if emp.ID == id {
			employees = append(employees[:i], employees[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
}


func deleteDepartment(c *gin.Context) {
	idParam := c.Param("id")

	for i, dep := range departments {
		if dep.ID == idParam {
			departments = append(departments[:i], departments[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Department deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
}


func main() {
	router := gin.Default()

	protected := router.Group("/", auth.AuthMiddleware())

	protected.GET("/departments", getDepartments)

	protected.GET("/employees", getEmployees)
	protected.POST("/departments", createDepartment)
	protected.POST("/employees", createEmployee)
	protected.PUT("/departments/:id", updateDepartments)
	protected.GET("/departments/:name",returnEmps)
    protected.PUT("/employees/:id", updateEmployee)
    protected.DELETE("/employees/:id", deleteEmployee)
    protected.DELETE("/departments/:id", deleteDepartment)

	router.Run("localhost:8080")
}
