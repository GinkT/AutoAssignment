package main

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)


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

func TestEnv_RedirectHandler(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer DB.Close()
	env := &Env{DB}

	row := sqlmock.NewRows([]string{"One"}).AddRow("https://github.com/DATA-DOG/go-sqlmock")
	sqlStatement := `
		SELECT "longlink" FROM public."links"
	`

	testShortLink := "hithere"
	mock.ExpectQuery(sqlStatement).WithArgs(testShortLink).WillReturnRows(row)
	req, err := http.NewRequest("GET", "/" + testShortLink, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.RedirectHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestEnv_ShortenerHandler(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer DB.Close()
	env := &Env{DB}
	mock.ExpectExec(`INSERT INTO public."links" VALUES`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("GET", "/?link=https://dou.ua/lenta/articles/golang-httptest", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.ShortenerHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// Middleware для валидации ссылки
func valiпdLink(shortener http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("link")
		_, err := url.ParseRequestURI(link)
		if err != nil {
			log.Printf("[Shortener] User link(%s) does not match URL\n", link)
			Json404Response(w, "Invalid link!")
			return
		}
		shortener.ServeHTTP(w, r)
	})
}
 */