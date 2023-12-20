package customerimporter

import (
	"fmt"
	"strings"
	"testing"
)

func TestOpenCSV(t *testing.T) {
	type tmplTest struct {
		name        string
		nameFile    string
		expectedErr error
	}
	var tests = []tmplTest{
		{"correct", "customers_test.csv", nil},
		{"incorrect extenion", "valid.cv", fmt.Errorf("the file does not have the extension .csv")},
		{"incorrect file", "validOpen.csv", fmt.Errorf("failed to open file: validOpen.csv")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := openCSV(test.nameFile)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("want (%v) got (%v)", test.expectedErr, err)
			}

		})
	}
}

func TestReadCSV(t *testing.T) {
	type tmplTest struct {
		name            string
		nameFile        string
		expectedClients map[string]int
		expectedErr     error
	}
	var tests = []tmplTest{
		{"correct", "customers_test.csv", map[string]int{"github.io": 1, "cyberchimps.com": 1, "hubpages.com": 1, "360.cn": 1}, nil},
		{"blank file", "customers_test_blank.csv", map[string]int{}, nil},
	}

	for _, test := range tests {
		clients := make(map[string]int)
		t.Run(test.name, func(t *testing.T) {
			executeFile, err := openCSV(test.nameFile)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("want (%v) got (%v)", test.expectedErr, err.Error())
			}

			err = readCSV(executeFile, clients)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("want (%v) got (%v)", test.expectedErr, err.Error())
			}

			for key, expectedValue := range test.expectedClients {
				actualValue, ok := clients[key]
				if !ok {
					t.Errorf("Expected client %s not found", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("Incorrect number of clients for %s. Expected %d, got %d", key, expectedValue, actualValue)
				}
			}
		})

	}

}

func TestExtractDomainSuffix(t *testing.T) {
	type tmplTest struct {
		name           string
		email          string
		expectedDomain string
		expectedErr    error
	}
	var tests = []tmplTest{
		{"correct", "example@org.com", "org.com", nil},
		{"domain without @", "exampleorg.com", "", fmt.Errorf("email without @: exampleorg.com")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			domain, err := extractDomainSuffix(test.email)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("want (%v) got (%v)", test.expectedErr, err)
			}
			if strings.Compare(test.expectedDomain, domain) != 0 {
				t.Errorf("incorrect extractDomainSuffix()")
			}

		})
	}
}

func TestCompareDomain(t *testing.T) {
	type tmplTest struct {
		name              string
		email             []string
		expectedErrSuffix error
		expectedClients   map[string]int
	}
	var tests = []tmplTest{
		{"correct", []string{"example@org.com", "example@org.com", "example@abc.com", "example@xyz.com"}, nil, map[string]int{"org.com": 2, "abc.com": 1, "xyz.com": 1}},
		{"correct", []string{"example@org.com", "exampleorg.com", "example@abc.com", "example@xyz.com"}, fmt.Errorf("email without @: exampleorg.com"), nil},
	}

	for _, test := range tests {
		clients := make(map[string]int)
		t.Run(test.name, func(t *testing.T) {
			for _, email := range test.email {
				domain, err := extractDomainSuffix(email)
				if err != nil && err.Error() == test.expectedErrSuffix.Error() {
					t.Skip()
				}

				compareDomain(domain, clients)
			}

			for key, expectedValue := range test.expectedClients {
				actualValue, ok := clients[key]
				if !ok {
					t.Errorf("Expected client %s not found", key)
					continue
				}

				if actualValue != expectedValue {
					t.Errorf("Incorrect number of clients for %s. Expected %d, got %d", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestSortDomains(t *testing.T) {
	type tmplTest struct {
		name               string
		clients            map[string]int
		expectedSortedKeys []string
	}
	var tests = []tmplTest{
		{"sorted keys", map[string]int{"org.com": 2, "abc.com": 1, "xyz.com": 13, "qwe.com": 16}, []string{"abc.com", "org.com", "qwe.com", "xyz.com"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sortedKeys := sortDomains(test.clients)

			for i := range sortedKeys {
				if sortedKeys[i] != test.expectedSortedKeys[i] {
					t.Errorf("Expected sorted keys %v, got: %v", test.expectedSortedKeys, sortedKeys)
				}
			}
		})
	}
}
