package database

import (
	"fmt"

	"github.com/VsenseTechnologies/biometric_http_server/internals/models"
)

func (q *Query) FetchAllBiometrics(userId string) ([]models.Biometric , error) {
	res , err := q.db.Query("SELECT user_id , unit_id , online FROM biometric WHERE user_id=$1" , userId)
	if err != nil {
		return nil , err
	}
	defer res.Close()
	var biometric models.Biometric
	var biometrics []models.Biometric
	for res.Next() {
		if err := res.Scan(&biometric.UserId , &biometric.UnitId , &biometric.Status); err != nil {
			return nil , err
		}
		biometrics = append(biometrics, biometric)
	}
	if res.Err() != nil {
		return nil , err
	}
	return biometrics , nil
}

func (q *Query) DeleteBiometricMachine(unitId string) error {
	tx ,err := q.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	} ()
	if err != nil {
		return err
	}
	_ , err = tx.Exec("DELETE FROM biometric WHERE unit_id=$1" , unitId)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DROP TABLE %s" , unitId)
	_ , err = tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) NewBiometricDevice(device models.Biometric) error {
	tx , err := q.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	} ()
	_ , err = tx.Exec("INSERT INTO biometric(unit_id , user_id , online) VALUES($1,$2,$3)" , device.UnitId , device.UserId , device.Status)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("CREATE TABLE %s (student_id VARCHAR(100) NOT NULL , student_unit_id VARCHAR(100) NOT NULL , student_name VARCHAR(200) NOT NULL , student_usn VARCHAR(200) NOT NULL , department VARCHAR(100) NOT NULL , FOREIGN KEY (student_id) REFERENCES fingerprintdata(student_id) ON DELETE CASCADE)" , device.UnitId)
	_ , err = tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}