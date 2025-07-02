package database

import (
	"log"
	"strconv"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) CreateBiometricDevice(biometric *models.Biometric) error {
	query1 := `INSERT INTO biometric (user_id,unit_id,online,label) VALUES ($1,$2,$3,$4)`
	query2 := `INSERT INTO student_unit_numbers (unit_id, student_unit_id) VALUES ($1, $2)`

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, biometric.UserId, biometric.UnitId, biometric.Online, biometric.Label); err != nil {
		tx.Rollback()
		return err
	}

	for i := 1; i <= 1000; i++ {
		id := strconv.Itoa(i)
		if _, err := tx.Exec(query2, biometric.UnitId, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return err
}

func (q *Query) GetBiometricDevices(userId string) ([]*models.Biometric, error) {
	query := `SELECT user_id, unit_id, online, label FROM biometric WHERE user_id=$1`

	var biometrics []*models.Biometric

	rows, err := q.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		biometric := new(models.Biometric)
		if err := rows.Scan(&biometric.UserId, &biometric.UnitId, &biometric.Online, &biometric.Label); err != nil {
			return nil, err
		}
		biometrics = append(biometrics, biometric)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return biometrics, nil
}


func (q *Query) UpdateBiometricLabel(unit_id string, label string) error {
	query := `UPDATE biometric SET label=$2 WHERE unit_id=$1`
	_, err := q.db.Exec(query, unit_id, label)
	return err
}

func (q *Query) DeleteBiometricDevice(unitId string) error {
	query1 := `DELETE FROM biometric WHERE unit_id=$1`
	query2 := `DELETE FROM student_unit_numbers WHERE unit_id = $1`

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(query1, unitId); err != nil {
		log.Printf("Error while creating Biometric Device : %v", err)
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(query2, unitId); err != nil {
		log.Printf("Error while creating Student_Unit_Numbers : %v", err)
		tx.Rollback()
	}

	tx.Commit()
	return err
}

func (q *Query) GetAvailableStudentUnitIDs(unitId string) ([]string, bool, error) {
	query := `SELECT student_unit_id FROM student_unit_numbers WHERE unit_id = $1 AND availability = TRUE`

	rows, err := q.db.Query(query, unitId)
	if err != nil {
		return nil, false, nil
	}

	var student_unit_ids []string
	var student_unit_id string

	for rows.Next() {
		if err := rows.Scan(&student_unit_id); err != nil {
			return nil, true, err
		}
		student_unit_ids = append(student_unit_ids, student_unit_id)
	}

	return student_unit_ids, true, nil
}

func (q *Query) UpdateAvailableStudentUnitIDs(unitID string, student_unit_ids []string, updateTo bool) error {
	query := `UPDATE student_unit_numbers SET availability = $3 WHERE unit_id = $1 AND student_unit_id = $2`

	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	for i := range student_unit_ids {
		if _, err := tx.Exec(query, unitID, student_unit_ids[i], updateTo); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
