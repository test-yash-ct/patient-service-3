package service

import (
	"context"
	"errors"
	"testing"

	"github.com/healthops/patient-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type fakePatients struct {
	gotTenant string
	gotID     string
	row       models.Patient
	getErr    error
}

func (f *fakePatients) GetByID(_ context.Context, tenantID, id string) (models.Patient, error) {
	f.gotTenant = tenantID
	f.gotID = id
	if f.getErr != nil {
		return models.Patient{}, f.getErr
	}
	return f.row, nil
}

func (f *fakePatients) ListByTenant(_ context.Context, tenantID string) ([]models.Patient, error) {
	f.gotTenant = tenantID
	return []models.Patient{f.row}, nil
}

func (f *fakePatients) Update(_ context.Context, tenantID, id, _, _, _, _ string) error {
	f.gotTenant = tenantID
	f.gotID = id
	return f.getErr
}

func TestResolveTenantMismatch(t *testing.T) {
	if _, err := ResolveTenant("t1", "t2"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("got %v", err)
	}
}

func TestGetRequiresTenant(t *testing.T) {
	svc := NewPatients(&fakePatients{})
	_, err := svc.Get(context.Background(), "", "", "p-1", "sub")
	if !errors.Is(err, ErrMissingTenant) {
		t.Fatalf("got %v", err)
	}
}

func TestGetScopesByResolvedTenant(t *testing.T) {
	repo := &fakePatients{row: models.Patient{ID: "p-1", TenantID: "acme"}}
	svc := NewPatients(repo)
	p, err := svc.Get(context.Background(), "acme", "", "p-1", "nurse")
	if err != nil {
		t.Fatal(err)
	}
	if repo.gotTenant != "acme" || p.ID != "p-1" {
		t.Fatalf("tenant=%q id=%q", repo.gotTenant, p.ID)
	}
}

func TestGetMapsNotFound(t *testing.T) {
	svc := NewPatients(&fakePatients{getErr: pgx.ErrNoRows})
	_, err := svc.Get(context.Background(), "acme", "acme", "missing", "sub")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}
