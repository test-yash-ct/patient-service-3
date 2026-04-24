package models

type Patient struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	FullName     string `json:"full_name"`
	SSN          string `json:"ssn"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	MedicalNotes string `json:"medical_notes"`
}
