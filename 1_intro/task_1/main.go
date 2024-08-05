package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var data = make(map[string]float64)

/*
Takes a reader initialized with the standard input (os.Stdin) and returns
the input with all the whitespace trimmed.
*/
func GetInput(r *bufio.Reader) string {
	val, err := r.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	return strings.TrimSpace(val)
}

/*
Calculates and returns the average of the grades stored in `data`
*/
func CalculateAverage(data map[string]float64) float64 {
	var total float64
	for _, v := range data {
		total += v
	}

	return total / float64(len(data))
}

/*
Prints a summary of the collected data in a user-friendly format.
*/
func ReportData(name string, data map[string]float64) {
	fmt.Printf("\nSummary\nName: %v\nAdded subjects: %v\n\tSubject\t\tGrade\n", name, len(data))
	for sub, grade := range data {
		fmt.Printf("\t%v\t\t%v\n", sub, grade)
	}
	fmt.Printf("Average of all %v subjects: %.2f", len(data), CalculateAverage(data))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nName: ")
	name := GetInput(reader)
	fmt.Print("Number of subjects: ")
	subject_count, convErr := strconv.ParseInt(GetInput(reader), 10, 64)
	if convErr != nil {
		panic(convErr.Error())
	}

	for i := 0; i < int(subject_count); i++ {
		var subject string

		// loop until a new subject name is entered
		for {
			fmt.Printf("\n\tSubject #%v: ", i+1)
			subject = GetInput(reader)
			_, ok := data[subject]
			if ok {
				fmt.Println("\t\t**Subject has already been recorded.")
				continue
			}

			break
		}

		// loop until a valid grade is entered
		for {
			fmt.Printf("\t%v grade: ", subject)
			grade, convErr := strconv.ParseFloat(GetInput(reader), 64)
			if convErr != nil {
				fmt.Println("\t\t**Invalid input. Grade must be a number.")
				continue
			}

			if grade > 100 || grade < 0 {
				fmt.Println("\t\t**Invalid input. Grade must be between values 0 and 100.")
				continue
			}

			data[subject] = grade
			break
		}
	}

	ReportData(name, data)
}
