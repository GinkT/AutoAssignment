package db

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
)

func TestAddLinkToDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
		INSERT INTO public."links"
		`

	testShortLink, testLongLink := "dhasd123h", "https://github.com/DATA-DOG/go-sqlmock"
	mock.ExpectExec(sqlStatement).WithArgs(testShortLink, testLongLink).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(sqlStatement).WithArgs(testShortLink, testLongLink).WillReturnResult(sqlmock.NewResult(0, 0))

	AddLinkToDB(db, testShortLink, testLongLink)
	if err := AddLinkToDB(db, testShortLink, testLongLink); err != ErrorAlreadyInDB {
		t.Fatalf("Expected to see %s", ErrorAlreadyInDB)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLinkFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
		SELECT "longlink" FROM public."links"
		`

	testShortLink := "381238123"
	row := sqlmock.NewRows([]string{"One"}).AddRow("https://github.com/DATA-DOG/go-sqlmock")
	mock.ExpectQuery(sqlStatement).WithArgs(testShortLink).WillReturnRows(row)

	GetLinkFromDB(db, testShortLink)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
