package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Agency struct {
	Id             int
	Name           string
	Address        string
	phoneNumber    string
	registeryDate  string
	numberEmployee int
}

const Path = "./data.txt"

var regionAgencies map[string][]Agency

func checkFlagRegionAndExist(inputRegion string) (string, error) {
	fRegion, err := checkFlagRegion(inputRegion)
	checkErr(err)
	// check region is Exist:
	if _, ok := regionAgencies[fRegion]; !ok {
		return "", fmt.Errorf("region %s is not Exist!\n", fRegion)
	}
	return fRegion, nil
}
func checkFlagRegion(inputRegion string) (string, error) {
	var flagRegion string
	if inputRegion == "" {
		fmt.Printf("You Don't enter region name, Please Enter Specific Region\nRegion Name: ")
		fmt.Scanln(&flagRegion)

		if flagRegion == "" {
			return "", fmt.Errorf("region Name is Empty!\n")
		}
	} else {
		flagRegion = inputRegion
	}

	return flagRegion, nil
}
func list(region string) []Agency {
	return regionAgencies[region]
}
func get(region string, id int) (Agency, error) {
	for _, v := range regionAgencies[region] {
		if v.Id == id {
			return v, nil
		}
	}
	return Agency{}, fmt.Errorf("Agency with Id: %d is not Found!", id) //what is the better code ???
}
func create(region string, agency Agency) error {
	file, err := os.OpenFile(Path, os.O_APPEND|os.O_WRONLY, 0600)
	checkErr(err)
	defer file.Close()

	newLine := fmt.Sprintf("%s,%d,%s,%s,%s,%s,%d\n",
		region, agency.Id, agency.Name, agency.Address, agency.phoneNumber, agency.registeryDate, agency.numberEmployee)
	if _, err = file.WriteString(newLine); err != nil {
		return err
	}

	return nil
}
func edit(region string, id int, agency Agency) error {
	file, err := os.OpenFile(Path, os.O_RDONLY, 0600)
	checkErr(err)
	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	var existId string
	var tmpLine string
	var editedFile []string

	for fileScanner.Scan() {
		tmpLine = fileScanner.Text() + "\n"
		existId = strings.Split(tmpLine, ",")[1]
		if existId == strconv.Itoa(id) {
			tmpLine = fmt.Sprintf("%s,%d,%s,%s,%s,%s,%d\n",
				region, agency.Id, agency.Name, agency.Address, agency.phoneNumber, agency.registeryDate, agency.numberEmployee)
		}
		editedFile = append(editedFile, tmpLine)
	}

	file.Close()
	os.Remove(Path)
	// what is the best solution ??????
	os.Create(Path)
	file, err = os.OpenFile(Path, os.O_WRONLY, 0600)
	checkErr(err)
	for _, v := range editedFile {
		if _, err = file.WriteString(v); err != nil {
			return err
		}
	}
	file.Close()
	return nil
}
func status(region string) (int, int) {
	numberAllEmploye := 0
	for _, v := range regionAgencies[region] {
		numberAllEmploye += v.numberEmployee
	}
	return len(regionAgencies[region]), numberAllEmploye
}
func textFileToMapStructure() error {
	file, err := os.OpenFile(Path, os.O_RDONLY, 0600)
	checkErr(err)
	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	var tmpSlice []string
	var newListAgency []Agency
	var tmpAgency Agency

	for fileScanner.Scan() {
		tmpSlice = strings.Split(fileScanner.Text(), ",")
		tmpAgency.Id, err = strconv.Atoi(tmpSlice[1])
		checkErr(err)
		tmpAgency.Name = tmpSlice[2]
		tmpAgency.Address = tmpSlice[3]
		tmpAgency.phoneNumber = tmpSlice[4]
		tmpAgency.registeryDate = tmpSlice[5]
		tmpAgency.numberEmployee, err = strconv.Atoi(tmpSlice[6])
		checkErr(err)
		if _, ok := regionAgencies[tmpSlice[0]]; ok {
			regionAgencies[tmpSlice[0]] = append(regionAgencies[tmpSlice[0]], tmpAgency)
		} else {
			newListAgency = append(newListAgency, tmpAgency)
			regionAgencies[tmpSlice[0]] = newListAgency
			newListAgency = nil
		}
	}
	file.Close()
	return nil
}
func calcId() int {
	lastId := 0
	for _, v := range regionAgencies {
		lastId += len(v)
	}
	return lastId
}
func checkErr(err error){
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	regionAgencies = make(map[string][]Agency)

	command := flag.String("command", "", "choose the best command")
	region := flag.String("region", "", "choose the region you want")
	flag.Parse()

	err := textFileToMapStructure()
	checkErr(err)

	lastId := calcId()

	switch *command {
	case "list":
		{
			fRegion, err := checkFlagRegionAndExist(*region)
			checkErr(err)

			var listAgency []Agency
			listAgency = list(fRegion)
			for agen := range listAgency {
				fmt.Printf("Name of Agency is: %s\n", listAgency[agen].Name)
			}
		}
	case "get":
		{
			fRegion, err := checkFlagRegionAndExist(*region)
			checkErr(err)
			var inputId int
			fmt.Printf("Please enter ID of Agency: ")
			fmt.Scanln(&inputId)
			chooseAgency, err := get(fRegion, inputId)
			checkErr(err)
			fmt.Println("==============================================")
			fmt.Println("region is:", fRegion,
				"\nID of Agency is:", chooseAgency.Id,
				"\nName of Agency is:", chooseAgency.Name,
				"\nAddress of Agency is:", chooseAgency.Address,
				"\nPhone Number of Agency is:", chooseAgency.phoneNumber,
				"\nRegistery Date is:", chooseAgency.registeryDate,
				"\nNumber of employee is:", chooseAgency.numberEmployee)
		}
	case "create":
		{
			fRegion, err := checkFlagRegion(*region)
			checkErr(err)
			var newAgency Agency
			newAgency.Id = lastId + 1
			fmt.Print("Please enter this values:\nName of Agency: ")
			fmt.Scanln(&newAgency.Name)
			fmt.Print("Address of Agency: ")
			fmt.Scanln(&newAgency.Address)
			fmt.Print("Phone Number of Agency: ")
			fmt.Scanln(&newAgency.phoneNumber)
			fmt.Print("Registery Date of Agency: ")
			fmt.Scanln(&newAgency.registeryDate)
			fmt.Print("Number of employee: ")
			fmt.Scanln(&newAgency.numberEmployee)

			err = create(fRegion, newAgency)
			checkErr(err)
		}
	case "edit":
		{
			fRegion, err := checkFlagRegionAndExist(*region)
			checkErr(err)
			var inputId int
			fmt.Printf("Please enter ID of Agency: ")
			fmt.Scanln(&inputId)
			chooseAgency, err := get(fRegion, inputId)
			checkErr(err)

			var tmpAgency Agency
			checkValue := func(new string, old string) string {
				if strings.ToLower(new) == "no" || strings.ToLower(new) == "n" {
					return old
				} else {
					return new
				}
			}

			tmpAgency.Id = chooseAgency.Id
			var tmp string
			fmt.Printf("Name of Agency is: %s, Are you change it?(No or new_value): ", chooseAgency.Name)
			fmt.Scanln(&tmp)
			tmpAgency.Name = checkValue(tmp, chooseAgency.Name)

			fmt.Printf("Address of Agency is: %s, Are you change it?(No or new_value): ", chooseAgency.Address)
			fmt.Scanln(&tmp)
			tmpAgency.Address = checkValue(tmp, chooseAgency.Address)

			fmt.Printf("Phone Number of Agency is: %s, Are you change it?(No or new_value): ", chooseAgency.phoneNumber)
			fmt.Scanln(&tmp)
			tmpAgency.phoneNumber = checkValue(tmp, chooseAgency.phoneNumber)

			fmt.Printf("Registery Date is: %s, Are you change it?(No or new_value): ", chooseAgency.registeryDate)
			fmt.Scanln(&tmp)
			tmpAgency.registeryDate = checkValue(tmp, chooseAgency.registeryDate)

			fmt.Printf("Number of employee is: %d, Are you change it?(No or new_value): ", chooseAgency.numberEmployee)
			fmt.Scanln(&tmp)
			tmpAgency.numberEmployee, err = strconv.Atoi(checkValue(tmp, strconv.Itoa(chooseAgency.numberEmployee)))
			checkErr(err)

			err = edit(fRegion, inputId, tmpAgency)
			checkErr(err)

		}
	case "status":
		{
			fRegion, err := checkFlagRegionAndExist(*region)
			checkErr(err)
			numberAllAgency, numberAllEmploye := status(fRegion)
			fmt.Printf("in region %s, number of total Agency is: %d, number of total Employee is: %d\n",
				fRegion, numberAllAgency, numberAllEmploye)
		}
	default:
		log.Fatalf("Error occurd!! Command %s not Found", *command)
	}
}