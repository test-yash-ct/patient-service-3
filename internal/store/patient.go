package store

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/healthops/patient-service/internal/models"
	"github.com/healthops/patient-service/internal/obs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientStore struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PatientStore {
	return &PatientStore{pool: pool}
}

func (s *PatientStore) GetByID(ctx context.Context, tenantID, id string) (models.Patient, error) {
	requestID := obs.RequestIDFromContext(ctx)
	const q = `SELECT id, tenant_id, full_name, ssn, email, phone, medical_notes
		FROM patients WHERE id = $1 AND tenant_id = $2`
	row := s.pool.QueryRow(ctx, q, id, tenantID)
	var p models.Patient
	err := row.Scan(&p.ID, &p.TenantID, &p.FullName, &p.SSN, &p.Email, &p.Phone, &p.MedicalNotes)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf(`{"event":"patient_not_found","request_id":"%s","tenant":"%s","patient_id":"%s"}`, requestID, tenantID, id)
		return models.Patient{}, err
	}
	if err != nil {
		log.New(os.Stdout, "", 0).Printf(`{"event":"patient_query_error","request_id":"%s","tenant":"%s","patient_id":"%s"}`, requestID, tenantID, id)
	}
	return p, err
}

func (s *PatientStore) ListByTenant(ctx context.Context, tenantID string) ([]models.Patient, error) {
	requestID := obs.RequestIDFromContext(ctx)
	const q = `SELECT id, tenant_id, full_name, ssn, email, phone, medical_notes
		FROM patients WHERE tenant_id = $1`
	rows, err := s.pool.Query(ctx, q, tenantID)
	if err != nil {
		log.New(os.Stdout, "", 0).Printf(`{"event":"patient_list_error","request_id":"%s","tenant":"%s"}`, requestID, tenantID)
		return nil, err
	}
	defer rows.Close()
	var out []models.Patient
	for rows.Next() {
		var p models.Patient
		if err := rows.Scan(&p.ID, &p.TenantID, &p.FullName, &p.SSN, &p.Email, &p.Phone, &p.MedicalNotes); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *PatientStore) Update(ctx context.Context, tenantID, id, fullName, email, phone, notes string) error {
	requestID := obs.RequestIDFromContext(ctx)
	const q = `UPDATE patients SET full_name = $1, email = $2, phone = $3, medical_notes = $4
		WHERE id = $5 AND tenant_id = $6`
	tag, err := s.pool.Exec(ctx, q, fullName, email, phone, notes, id, tenantID)
	if err != nil {
		log.New(os.Stdout, "", 0).Printf(`{"event":"patient_update_error","request_id":"%s","tenant":"%s","patient_id":"%s"}`, requestID, tenantID, id)
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
