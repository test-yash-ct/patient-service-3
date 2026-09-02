package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/auth"
	"github.com/healthops/patient-service/internal/obs"
	"github.com/healthops/patient-service/internal/store"
)

type PatientAPI struct {
	Store  *store.PatientStore
	Secret string
}

func (a *PatientAPI) Register(r *gin.RouterGroup) {
	r.GET("/patients/:id", a.Get)
}

func (a *PatientAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := obs.RequestIDFromContext(ctx)

	claims, err := auth.ParseBearer(ctx, c.GetHeader("Authorization"), a.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tenantID := claims.Tenant
	if headerTenant := c.GetHeader(obs.HeaderTenantID); headerTenant != "" {
		if tenantID != "" && headerTenant != tenantID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant_mismatch"})
			return
		}
		if tenantID == "" {
			tenantID = headerTenant
		}
	}
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_tenant"})
		return
	}
	ctx = obs.WithTenant(ctx, tenantID)
	c.Request = c.Request.WithContext(ctx)

	id := c.Param("id")
	p, err := a.Store.GetByID(ctx, tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	logPatientAccess(requestID, tenantID, id, claims.Sub)
	c.JSON(http.StatusOK, p)
}

func logPatientAccess(requestID, tenant, patientID, subject string) {
	entry := map[string]string{
		"event":      "patient_access",
		"request_id": requestID,
		"tenant":     tenant,
		"patient_id": patientID,
		"subject":    subject,
	}
	b, _ := json.Marshal(entry)
	log.New(os.Stdout, "", 0).Println(string(b))
}
