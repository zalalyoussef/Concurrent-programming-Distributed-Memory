package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode"
)

type Employee struct {
	ID     int
	Salary float64
	Name   string
	Hash   string
}

func NewEmployee(ID int, Salary float64, Name string) Employee {
	return Employee{
		ID:     ID,
		Salary: Salary,
		Name:   Name,
		Hash:   randomString(30),
	}
}

func randomString(length int) string {
	const chars = "896"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func ReadData() []Employee {
	filePath := "Employees.txt"
	var employees []Employee

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return employees
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) == 3 {
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Printf("Error parsing ID: %v", err)
				continue
			}

			salary, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Printf("Error parsing Salary: %v", err)
				continue
			}

			name := parts[2]

			employee := NewEmployee(id, salary, name)
			employees = append(employees, employee)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v", err)
	}

	return employees
}

func PrintEmployeesToFile(employees []Employee, filePath string) {
	// Create or open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	// Create a tabwriter with padding and minwidth settings
	w := tabwriter.NewWriter(file, 0, 0, 2, ' ', tabwriter.TabIndent)

	// Write headers with tab-separated columns
	fmt.Fprintln(w, "ID\tSalary\tName\tHash")

	// Write the employee details to the tabwriter
	for _, employee := range employees {
		fmt.Fprintf(w, "%d\t%.2f\t%s\t%s\n", employee.ID, employee.Salary, employee.Name, employee.Hash)
	}

	// Flush the tabwriter
	w.Flush()

	fmt.Println("Employee details have been written to", filePath)
}

const NUMBER_OF_THREADS = 5

var completed = Employee{-1, 0, "END", ""}

func main() {
	// Create channels
	employees := ReadData()
	var insertCh = make(chan Employee)
	var removeCh = make(chan Employee)
	var filteredCh = make(chan Employee)
	var resultToMainChannel = make(chan []Employee)
	var writeFlag = make(chan int)

	for i := 0; i < NUMBER_OF_THREADS; i++ {
		go workerThread(removeCh, filteredCh, writeFlag)
	}

	go dataThread(insertCh, removeCh, writeFlag)
	go resultThread(filteredCh, resultToMainChannel, writeFlag)

	// Sends data
	for _, employee := range employees {

		insertCh <- employee
	}
	// Inform DataThread that there is no more employees
	insertCh <- completed

	// Wait for receiving the results
	var result []Employee = <-resultToMainChannel

	PrintEmployeesToFile(result, "Results.txt")
}

func workerThread(removeCh <-chan Employee, filteredCh chan<- Employee, writeflag chan<- int) {
	var finished = false
	for !finished {
		writeflag <- 1
		var received = <-removeCh
		if received == completed {

			finished = true
		} else {

			if !unicode.IsDigit(rune(received.Hash[0])) {
				filteredCh <- received

			}
		}
	}
	filteredCh <- completed
}

func dataThread(insertCh <-chan Employee, removeCh chan<- Employee, writeFlag <-chan int) {
	const DATA_SIZE_ARRAY = 10
	var data [DATA_SIZE_ARRAY]Employee
	var dataLast = -1
	var noMoreInMain = false
	var valueAsked = 0
	var dataFromMain = Employee{0, 0, "DEFAULT", ""}
	for valueAsked != 2 {

		if dataLast >= 0 && dataLast < DATA_SIZE_ARRAY-1 && !noMoreInMain {
			select {
			case valueAsked = <-writeFlag:
				if valueAsked == 1 {
					if dataLast == -1 && noMoreInMain {
						removeCh <- completed

					} else {
						removeCh <- data[dataLast]
						dataLast--
					}
				}
			case dataFromMain = <-insertCh:
				if dataFromMain != completed {
					dataLast++
					data[dataLast] = dataFromMain
				} else {
					noMoreInMain = true
				}
			}
		} else if dataLast == -1 && !noMoreInMain {
			// Array is empty but we main didn't finished to add items
			dataFromMain = <-insertCh
			if dataFromMain != completed {
				dataLast++
				data[dataLast] = dataFromMain
			} else {
				noMoreInMain = true
			}
		} else {
			// Main finished to add items or array is full, sending end signal if array is empty
			valueAsked = <-writeFlag
			if valueAsked == 1 {
				if dataLast == -1 && noMoreInMain {
					removeCh <- completed
				} else {
					removeCh <- data[dataLast]
					dataLast--
				}
			}
		}

	}
}

func resultThread(filteredCh <-chan Employee, resultToMainChannel chan<- []Employee, writeFlag chan int) {
	result := make([]Employee, 25)
	nbResults := 0

	var receivedFinishSignals = 0
	for receivedFinishSignals < NUMBER_OF_THREADS {
		var received = <-filteredCh
		if received == completed {
			receivedFinishSignals++
		} else {
			var iElement = nbResults
			for iElement > 0 && result[iElement-1].ID > received.ID {
				result[iElement] = result[iElement-1]
				result[iElement].ID = iElement
				iElement--
			}
			result[iElement] = received
			result[iElement].ID = iElement
			nbResults++
		}
	}

	//writeFlag <- 2
	if nbResults > 0 {
		resultToMainChannel <- result[0:nbResults]
	} else {
		resultToMainChannel <- []Employee{}
	}
}
