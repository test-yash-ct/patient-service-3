CREATE TABLE IF NOT EXISTS patients (
    id TEXT NOT NULL,
    tenant_id TEXT NOT NULL,
    full_name TEXT NOT NULL,
    ssn TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    medical_notes TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (id, tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_patients_tenant ON patients (tenant_id, id);

CREATE TABLE IF NOT EXISTS care_team_assignments (
    tenant_id TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    patient_id TEXT NOT NULL,
    PRIMARY KEY (tenant_id, provider_id, patient_id)
);

CREATE INDEX IF NOT EXISTS idx_care_team_patient ON care_team_assignments (tenant_id, patient_id);
