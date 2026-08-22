package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/audit"
	"github.com/healthops/patient-service/internal/crypto"
	"github.com/healthops/patient-service/internal/middleware"
	"github.com/healthops/patient-service/internal/store"
)

type PatientAPI struct {
	Store *store.PatientStore
	Audit *audit.Logger
}

func (a *PatientAPI) Register(r *gin.RouterGroup) {
	r.GET("/patients/:id", a.Get)
}

func (a *PatientAPI) Get(c *gin.Context) {
	claims, ok := middleware.ClaimsFrom(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !claims.HasRole("patient:read") {
		a.emit(c, claims.Sub, claims.Tenant, c.Param("id"), "denied")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id := c.Param("id")
	if !middleware.ValidPatientID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}
	p, err := a.Store.GetByID(c.Request.Context(), claims.Tenant, id)
	if err != nil {
		a.emit(c, claims.Sub, claims.Tenant, id, "not_found")
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if !claims.HasRole("patient:phi-breakglass") {
		p.SSN = crypto.MaskSSN(p.SSN)
	}
	a.emit(c, claims.Sub, claims.Tenant, id, "ok")
	c.JSON(http.StatusOK, p)
}

func (a *PatientAPI) emit(c *gin.Context, actor, tenant, object, outcome string) {
	rid, _ := c.Get("request_id")
	ridStr, _ := rid.(string)
	if a.Audit != nil {
		a.Audit.Emit(audit.Event{
			Actor: actor, Tenant: tenant, Action: "patient.read", ObjectID: object, Outcome: outcome, RequestID: ridStr,
		})
	}
}
