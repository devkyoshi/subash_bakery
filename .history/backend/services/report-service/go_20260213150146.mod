module github.com/yourusername/erp-system/services/report-service

go 1.23.0

replace github.com/yourusername/erp-system/shared => ../../shared

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/xuri/excelize/v2 v2.9.0
	github.com/yourusername/erp-system/shared v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.17.9
)
