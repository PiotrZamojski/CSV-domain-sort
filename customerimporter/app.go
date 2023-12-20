package customerimporter

import "fmt"

func Application() {
	clients := make(map[string]int)
	csvFile := "customerimporter/customers.csv"

	executeFile, err := openCSV(csvFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = readCSV(executeFile, clients)
	if err != nil {
		fmt.Println(err)
		return
	}

	sortedKeys := sortDomains(clients)
	printNumberOfClients(sortedKeys, clients)
}
