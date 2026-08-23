package store

import (
	"context"
	"errors"

	"github.com/healthops/patient-service/internal/crypto"
	"github.com/healthops/patient-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientStore struct {
	pool *pgxpool.Pool
	enc  *crypto.FieldEncryptor
}

func New(pool *pgxpool.Pool, enc *crypto.FieldEncryptor) *PatientStore {
	return &PatientStore{pool: pool, enc: enc}
}

func (s *PatientStore) GetByID(ctx context.Context, tenantID, id string) (models.Patient, error) {
	const q = `SELECT id, tenant_id, full_name, ssn, email, phone, medical_notes FROM patients WHERE id = $1 AND tenant_id = $2`
	row := s.pool.QueryRow(ctx, q, id, tenantID)
	var p models.Patient
	err := row.Scan(&p.ID, &p.TenantID, &p.FullName, &p.SSN, &p.Email, &p.Phone, &p.MedicalNotes)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Patient{}, err
	}
	if err != nil {
		return models.Patient{}, err
	}
	if p.SSN, err = s.enc.Decrypt(p.SSN); err != nil {
		return models.Patient{}, err
	}
	if p.Email, err = s.enc.Decrypt(p.Email); err != nil {
		return models.Patient{}, err
	}
	if p.Phone, err = s.enc.Decrypt(p.Phone); err != nil {
		return models.Patient{}, err
	}
	if p.MedicalNotes, err = s.enc.Decrypt(p.MedicalNotes); err != nil {
		return models.Patient{}, err
	}
	if p.FullName, err = s.enc.Decrypt(p.FullName); err != nil {
		return models.Patient{}, err
	}
	return p, nil
}

func (s *PatientStore) IsCareTeamMember(ctx context.Context, tenantID, providerID, patientID string) (bool, error) {
	const q = `SELECT 1 FROM care_team_assignments WHERE tenant_id = $1 AND provider_id = $2 AND patient_id = $3`
	var one int
	err := s.pool.QueryRow(ctx, q, tenantID, providerID, patientID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
