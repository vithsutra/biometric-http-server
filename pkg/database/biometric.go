package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
	"github.com/google/uuid"
)

func (q *Query) CreateBiometricDevice(biometric *models.Biometric) error {
	query1 := `INSERT INTO biometric (user_id,unit_id,online,label) VALUES ($1,$2,$3,$4)`
	query2 := fmt.Sprintf(`CREATE TABLE %s (student_id VARCHAR(100) NOT NULL , student_unit_id VARCHAR(100) NOT NULL , student_name VARCHAR(200) NOT NULL , student_usn VARCHAR(200) NOT NULL , department VARCHAR(100) NOT NULL , FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id) ON DELETE CASCADE)`, biometric.UnitId)

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, biometric.UserId, biometric.UnitId, biometric.Online, biometric.Label); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (q *Query) GetBiometricDevices(userId string) ([]*models.Biometric, error) {
	query := `SELECT user_id,unit_id,online,label FROM biometric WHERE user_id=$1`

	var biometrics []*models.Biometric

	rows, err := q.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var biometric models.Biometric
		if err := rows.Scan(&biometric.UserId, &biometric.UnitId, &biometric.Online, &biometric.Label); err != nil {
			return nil, err
		}
		biometrics = append(biometrics, &biometric)
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
	query2 := fmt.Sprintf(`DROP TABLE %s`, unitId)

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, unitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (q *Query) ClearBiometricDeviceData(userId string, unitId string) error {
	query1 := `DELETE FROM deletes WHERE unit_id=$1`
	query2 := `INSERT INTO deletes SELECT unit_id,student_unit_id FROM fingerprintdata WHERE unit_id=$1`
	uniqueId := uuid.NewString()
	query3 := fmt.Sprintf(`CREATE TEMP TABLE temp_%s AS SELECT * FROM biometric WHERE unit_id=$1`, uniqueId)
	query4 := `DELETE FROM biometric WHERE unit_id=$1`
	query5 := fmt.Sprintf(`INSERT INTO biometric SELECT * FROM temp_%s`, uniqueId)

	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.Exec(query1, unitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query2, unitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query3, unitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query4, unitId); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query5); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (q *Query) SwapBiometricData(fromMachineId string, toMachineId string) error {
	tx, err := q.db.Begin()

	if err != nil {
		return err
	}

	query1 := `DELETE FROM inserts WHERE unit_id = $1`

	_, err = tx.Exec(query1, fromMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query1, toMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query2 := `DELETE FROM deletes WHERE unit_id = $1`

	_, err = tx.Exec(query2, fromMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query2, toMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	uniqueId := uuid.NewString()

	query3 := `CREATE TEMP TABLE temp1_` + uniqueId + ` TABLE AS SELECT * FROM fingerprintdata WHERE unit_id=$1`

	_, err = tx.Exec(query3, fromMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query4 := `CREATE TEMP TABLE temp2_` + uniqueId + ` TABLE AS SELECT * FROM fingerprintdata WHERE unit_id=$1`

	_, err = tx.Exec(query4, toMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query5 := `CREATE TEMP TABLE temp3_` + uniqueId + ` TABLE AS SELECT * FROM biometric WHERE unit_id=$1`

	_, err = tx.Exec(query5, fromMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query6 := `CREATE TEMP TABLE temp4_` + uniqueId + ` TABLE AS SELECT * FROM biometric WHERE unit_id=$1`

	_, err = tx.Exec(query6, toMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query7 := `DELETE FROM biometric WHERE unit_id=$1`

	_, err = tx.Exec(query7, fromMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(query7, toMachineId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query8 := `INSERT INTO biometric SELECT * FROM temp3_` + uniqueId

	_, err = tx.Exec(query8)

	if err != nil {
		tx.Rollback()
		return err
	}

	query9 := `DROP TABLE temp3_` + uniqueId

	_, err = tx.Exec(query9)

	if err != nil {
		tx.Rollback()
		return err
	}

	query10 := `INSERT INTO biometric SELECT * FROM temp4_` + uniqueId

	_, err = tx.Exec(query10)

	if err != nil {
		tx.Rollback()
		return err
	}

	query11 := `DROP TABLE temp4_` + uniqueId

	_, err = tx.Exec(query11)

	if err != nil {
		tx.Rollback()
		return err
	}

	query12 := `INSERT INTO 
					fingerprintdata (student_id,student_unit_id,unit_id,fingerprint) 
					SELECT student_id,student_unit_id,` + toMachineId + `,fingerprint 	
				FROM temp1_` + uniqueId

	_, err = tx.Exec(query12)

	if err != nil {
		tx.Rollback()
		return err
	}

	query13 := `INSERT INTO 
					fingerprintdata (student_id,student_unit_id,unit_id,fingerprint) 
					SELECT student_id,student_unit_id,` + fromMachineId + `,fingerprint 	
				FROM temp2_` + uniqueId

	_, err = tx.Exec(query13)

	if err != nil {
		tx.Rollback()
		return err
	}

	query14 := `CREATE TEMP TABLE temp5_` + uniqueId + ` AS SELECT * FROM ` + fromMachineId

	_, err = tx.Exec(query14)

	if err != nil {
		tx.Rollback()
		return err
	}

	query15 := `DELETE FROM ` + fromMachineId

	_, err = tx.Exec(query15)

	if err != nil {
		tx.Rollback()
		return err
	}

	query16 := `INSERT INTO ` + fromMachineId + ` SELECT * FROM ` + toMachineId

	_, err = tx.Exec(query16)

	if err != nil {
		tx.Rollback()
		return err
	}

	query17 := `DELETE FROM ` + toMachineId

	_, err = tx.Exec(query17)

	if err != nil {
		tx.Rollback()
		return err
	}

	query18 := `INSERT INTO ` + toMachineId + ` SELECT * FROM temp5_` + uniqueId

	_, err = tx.Exec(query18)

	if err != nil {
		tx.Rollback()
		return err
	}

	query19 := `DROP TABLE temp5_` + uniqueId

	_, err = tx.Exec(query19)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
	// query20 := `INSERT INTO inserts `
}
