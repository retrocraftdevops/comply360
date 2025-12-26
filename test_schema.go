package main

import (
	"fmt"
	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

func main() {
	testUUID := uuid.MustParse("9ac5aa3e-91cd-451f-b182-563b0d751dc7")
	tenant := &models.Tenant{
		ID: testUUID,
	}
	fmt.Println("Schema name:", tenant.TenantSchema())
}
