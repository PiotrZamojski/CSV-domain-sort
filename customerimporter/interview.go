// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func openCSV(csvFile string) (*os.File, error) {
	if !strings.HasSuffix(csvFile, ".csv") {
		return nil, fmt.Errorf("the file does not have the extension .csv")
	}

	executeFile, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", csvFile)
	}

	return executeFile, nil
}

func readCSV(executeFile *os.File, clients map[string]int) error {
	defer executeFile.Close()
	r := csv.NewReader(executeFile)
	isFirstLine := true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file: %v", executeFile)
		}
		if isFirstLine {
			isFirstLine = false
			continue
		}
		email, err := extractDomainSuffix(record[2])
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		compareDomain(email, clients)
	}
	return nil
}

func extractDomainSuffix(email string) (string, error) {
	index := strings.Index(email, "@")
	if index == -1 {
		return "", fmt.Errorf("email without @: %s", email)
	}

	return email[index+1:], nil
}

func compareDomain(email string, clients map[string]int) {
	_, ok := clients[email]
	if !ok {
		clients[email] = 1
	} else {
		clients[email]++
	}
}

func sortDomains(clients map[string]int) []string {
	var sortedKeys []string
	for key := range clients {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	return sortedKeys
}

func printNumberOfClients(sortedKeys []string, clients map[string]int) {
	for _, key := range sortedKeys {
		fmt.Printf("%s, %+v\n", key, clients[key])
	}
}
