package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/erp-system/services/api-gateway/config"
)

func main() {
	cfg := config.LoadConfig()

	router := gin.Default()

	// Middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-organization-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// Proxy handler
	router.Any("/api/v1/*path", func(c *gin.Context) {
		path := c.Param("path")

		var targetURL string

		// Procurement Service Routes
		if strings.HasPrefix(path, "/suppliers") ||
			strings.HasPrefix(path, "/purchase-orders") ||
			strings.HasPrefix(path, "/grns") {
			targetURL = cfg.ProcurementServiceURL
		} else if strings.HasPrefix(path, "/organizations") {
			// Check for procurement sub-resources under organizations
			if strings.Contains(path, "/suppliers") ||
				strings.Contains(path, "/purchase-orders") ||
				strings.Contains(path, "/grns") {
				targetURL = cfg.ProcurementServiceURL
			} else if strings.Contains(path, "/stock-movements") ||
				strings.Contains(path, "/batches") ||
				strings.Contains(path, "/stock-adjustments") ||
				strings.Contains(path, "/inventory-counts") ||
				strings.Contains(path, "/serial-numbers") {
				targetURL = cfg.InventoryServiceURL
			} else {
				// Default to Org Service for other /organizations routes
				targetURL = cfg.OrgServiceURL
			}
		} else if strings.HasPrefix(path, "/auth") {
			targetURL = cfg.AuthServiceURL
		} else if strings.HasPrefix(path, "/companies") || strings.HasPrefix(path, "/locations") || strings.HasPrefix(path, "/users") {
			targetURL = cfg.OrgServiceURL
		} else if strings.Contains(path, "/stock") ||
			strings.HasPrefix(path, "/stock-levels") ||
			strings.HasPrefix(path, "/stock-movements") ||
			strings.HasPrefix(path, "/batches") ||
			strings.HasPrefix(path, "/stock-adjustments") ||
			strings.HasPrefix(path, "/inventory-counts") ||
			strings.HasPrefix(path, "/serial-numbers") {
			targetURL = cfg.InventoryServiceURL
		} else if strings.HasPrefix(path, "/products") {
			// Check for inventory sub-resources under products
			if strings.Contains(path, "/stock") ||
				strings.Contains(path, "/movements") ||
				strings.Contains(path, "/batches") ||
				strings.Contains(path, "/serial-numbers") {
				targetURL = cfg.InventoryServiceURL
			} else {
				targetURL = cfg.ProductServiceURL
			}
		} else if strings.HasPrefix(path, "/categories") || strings.HasPrefix(path, "/brands") || strings.HasPrefix(path, "/attributes") || strings.HasPrefix(path, "/units") {
			targetURL = cfg.ProductServiceURL
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}

		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Update the director to keeping the original path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = remote.Host
			// We can strip prefix if needed, but here we forward the full /api/v1/... path
			// if the services expect /api/v1/...
			// The services in this project seem to define groups like router.Group("/api/v1") in their main.go
			// So forwarding the full path is correct.
		}

		// Modify response to strip CORS headers from backend to avoid duplication
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Credentials")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Allow-Methods")
			return nil
		}

		// Error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Bad Gateway"))
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("API Gateway starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
