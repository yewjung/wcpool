package partyRepoSql

import (
	"database/sql"
	"log"
)

// partyId -> userId (email)
type PartyRepo struct {
	DB *sql.DB
}

// create new user
func (partyRepo *PartyRepo) AddMemberToParty(partyid string, email string) error {
	_, err := partyRepo.DB.Exec("INSERT INTO party(partyid, email) VALUES($1, $2)", partyid, email)
	if err != nil {
		log.Default().Panic(err)
	}
	return err
}

// update user
func (partyRepo *PartyRepo) UpdateEmail(partyid string, oldEmail string, newEmail string) error {
	_, err := partyRepo.DB.Exec("UPDATE party SET email=$1 WHERE partyid=$2 AND email=$3", newEmail, partyid, oldEmail)
	if err != nil {
		log.Default().Panic(err)
	}
	return err
}

func (partyRepo *PartyRepo) GetPartyMemberIDs(partyId string) ([]string, error) {
	rows, err := partyRepo.DB.Query("SELECT email FROM party WHERE partyid = $1", partyId)
	if err != nil {
		log.Default().Panic(err)
		return nil, err
	}
	emails := make([]string, 0)
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}
	return emails, nil
}
