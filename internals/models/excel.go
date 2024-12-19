package models

type Excel struct {
	Name       string `json:"name"`
	USN        string `json:"usn"`
	Attendance map[string]string 
}