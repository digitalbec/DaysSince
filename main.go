package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {

	args := os.Args
	// check length of input date before we proceed
	if len(args) < 3 {
		log.Fatal("invalid input error format is: ./<build binary> dd/mm/yyyy")
	}

	now := time.Now().Format("02/01/2006")

	noOfdays, Err := NumberOfDays(args[1], now)
	if Err != nil {
		log.Fatal("error calculating difference in dates")
	}

	strDays := strconv.Itoa(noOfdays)

	fmt.Printf(strDays + " days since last reported accident")
}

const minDefinedYear = 1900
const maxDefinedYear = 2999
const ErrorMessage = "Please enter date range between 01/01/1900 and 31/12/2999"

type date struct {
	day   int
	month int
	year  int
}

var DaysInMonth = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

// checks if year is leap year or not
func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 && year%400 != 0 {
			return false
		}
		return true
	}

	return false
}

// validateDate - will verify date format and initialize date struct
func validateDate(inputDate string) (*date, error) {
	dateValidation := strings.Split(inputDate, "/")
	if len(dateValidation) != 3 {
		return nil, fmt.Errorf("date is not in the format of dd/dd/yyyy")
	}
	// parse string into dd/mm/yyyy format
	var dateSlice []int
	for _, value := range dateValidation {
		datebreak, strErr := strconv.Atoi(value)
		if strErr != nil {
			return nil, fmt.Errorf("string to int conversion failed: %v", strErr)
		}
		dateSlice = append(dateSlice, datebreak)
	}

	inputDates := DaysInMonth

	//leap year check
	if isLeapYear(dateSlice[2]) {
		inputDates[2] = 29
	}
	if dateSlice[2] <= minDefinedYear || dateSlice[2] >= maxDefinedYear {
		return nil, fmt.Errorf("year must be in between 1900 & 2999")
	}
	_, errorMonth := DaysInMonth[dateSlice[1]]
	if !errorMonth {
		return nil, fmt.Errorf("month must be in between 1-12")
	}
	if dateSlice[0] < 1 || dateSlice[0] > DaysInMonth[dateSlice[1]] {
		return nil, fmt.Errorf("day must be in between 1-31")
	}

	newDate := new(date)
	newDate.day = dateSlice[0]
	newDate.month = dateSlice[1]
	newDate.year = dateSlice[2]

	return newDate, nil
}

//func NumberOfDays find the number of days in a given range
func NumberOfDays(fDate, sDate string) (int, error) {

	startDate, startDateErr := validateDate(fDate)
	if startDateErr != nil {
		log.Fatal(fmt.Sprintf("First date input error: %s\n%s", startDateErr, ErrorMessage))
	}

	endDate, endDateErr := validateDate(sDate)
	if endDateErr != nil {
		log.Fatal(fmt.Sprintf("Second date input error: %s\n%s", endDateErr, ErrorMessage))
	}

	// Ensure date startdate is always smaller than enddate
	if (startDate.year > endDate.year) || (startDate.year == endDate.year && startDate.month > endDate.month) || (startDate.year == endDate.year && startDate.month == endDate.month && startDate.day > endDate.day) {
		startDate, endDate = endDate, startDate
	}

	NumberOfDays := 0
	validDates := DaysInMonth

	for {
		// leap year check
		if isLeapYear(startDate.year) {
			validDates[2] = 29
		} else {
			validDates[2] = 28
		}

		// date becomes equals then come out of loop
		if reflect.DeepEqual(startDate, endDate) {
			if NumberOfDays == 0 {
				break
			}
			NumberOfDays--
			break
		}
		startDate.day++

		//days complete within month, start day from 1 again and increment month
		if startDate.day > DaysInMonth[startDate.month] {
			startDate.month++
			startDate.day = 1
		}

		// Increment year
		if startDate.month > len(DaysInMonth) {
			startDate.year++
			startDate.month = 1
		}
		NumberOfDays++
	}
	return NumberOfDays, nil
}
