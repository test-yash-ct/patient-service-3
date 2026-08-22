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
	return p, nil
}
