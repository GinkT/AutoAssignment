package db

import (
	"database/sql"
	"errors"
	"log"
)

var ErrorAlreadyInDB error = errors.New("Link is already in database!")

func AddLinkToDB(db *sql.DB, shortLink, longLink string) error {
	sqlStatement := `
		INSERT INTO public."links"
		VALUES
		($1, $2)
		ON CONFLICT DO NOTHING RETURNING "shortlink"
	`
	result, err := db.Exec(sqlStatement, shortLink, longLink)
	if err != nil {
		log.Println("[DB] Error happened while Exec", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error while getting Rows Affected")
	}
	if rowsAffected == 0 {
		log.Printf("[DB] Link %s is aldready in database\n", longLink)
		return ErrorAlreadyInDB
	}
	log.Printf("[DB] Succesfully added [%s] %s to DB\n", shortLink, longLink)
	return nil
}

func GetLinkFromDB(db *sql.DB, shortLink string) (string, error) {
	sqlStatement := `
		SELECT "longlink" FROM public."links"
		WHERE "shortlink" = $1
	`

	row := db.QueryRow(sqlStatement, shortLink)
	var longLink string
	if err := row.Scan(&longLink); err != nil {
		log.Println("[DB] Error happened while Scan", err)
		return "", err
	}
	log.Printf("[DB] Got long link %s for %s\n", longLink, shortLink)
	return longLink, nil
}
