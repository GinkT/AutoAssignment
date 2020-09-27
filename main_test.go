package main

import "testing"


// Проверяю чтобы функция возвращала то же значение
func TestShrinkLink(t *testing.T) {
	testLinks := []string {
		"https://dev.to/ale_ukr/how-to-test-database-interactions-in-golang-applications-3041",
		"https://www.google.com/search?q=how+to+test+database+connect+function+golang&oq=how+to+test+database+connect+function+golang&aqs=chrome..69i57.7331j0j7&sourceid=chrome&ie=UTF-8",
		"https://www.youtube.com/",
		"https://www.msn.com/ru-ru/news/article/%d0%bf%d0%b5%d1%87%d0%b0%d0%bb%d1%8c%d0%bd%d1%8b%d0%b5-%d1%80%d0%b5%d0%ba%d0%be%d1%80%d0%b4%d1%8b-%d0%b2-%d0%bc%d0%be%d1%81%d0%ba%d0%b2%d0%b5-%d0%b2%d1%8b%d1%80%d0%be%d1%81%d0%bb%d0%be-%d1%87%d0%b8%d1%81%d0%bb%d0%be-%d0%b7%d0%b0%d0%b1%d0%be%d0%bb%d0%b5%d0%b2%d1%88%d0%b8%d1%85-covid-19/ar-BB19t3kJ?ocid=msedgntp",
		"https://lms.mtuci.ru/",
	}
	testData := make([]string, 0)
	for _, value := range testLinks {
		testData = append(testData, ShrinkLink(value))
	}

	for idx, value := range testData {
		if ShrinkLink(testLinks[idx]) != value {
			t.Fatalf("Got %s, expected %s\n", testLinks[idx], value)
		}
	}
}