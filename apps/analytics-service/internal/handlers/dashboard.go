package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// DashboardMetrics represents real-time dashboard metrics
type DashboardMetrics struct {
	TotalMRR           float64 `json:"total_mrr"`
	TotalARR           float64 `json:"total_arr"`
	TotalCustomers     int     `json:"total_customers"`
	ActiveAgents       int     `json:"active_agents"`
	MonthlyGrowthRate  float64 `json:"monthly_growth_rate"`
	ChurnRate          float64 `json:"churn_rate"`
	AverageARPU        float64 `json:"average_arpu"`
	TotalCommission    float64 `json:"total_commission"`
	CommissionPercentage float64 `json:"commission_percentage"`
}

// ExecutiveDashboardResponse represents executive dashboard data
type ExecutiveDashboardResponse struct {
	Metrics      DashboardMetrics            `json:"metrics"`
	TrendData    []TrendPoint                `json:"trend_data"`
	TopAgents    []AgentPerformance          `json:"top_agents"`
	RevenueByTerritory []TerritoryRevenue    `json:"revenue_by_territory"`
	LastUpdated  string                      `json:"last_updated"`
}

// TrendPoint represents a point in time-series data
type TrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// AgentPerformance represents agent performance summary
type AgentPerformance struct {
	AgentID    string  `json:"agent_id"`
	AgentName  string  `json:"agent_name"`
	MRR        float64 `json:"mrr"`
	Customers  int     `json:"customers"`
	Commission float64 `json:"commission"`
	GrowthRate float64 `json:"growth_rate"`
}

// TerritoryRevenue represents revenue by territory
type TerritoryRevenue struct {
	Territory  string  `json:"territory"`
	MRR        float64 `json:"mrr"`
	ARR        float64 `json:"arr"`
	Agents     int     `json:"agents"`
	Customers  int     `json:"customers"`
}

// GetExecutiveDashboard returns executive dashboard data
func GetExecutiveDashboard(c *gin.Context) {
	_ = c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id
	// TODO: Filter data by tenantID when implementing database queries
	
	// TODO: Fetch actual data from database
	// For now, return mock data
	response := ExecutiveDashboardResponse{
		Metrics: DashboardMetrics{
			TotalMRR:          125000.00,
			TotalARR:          1500000.00,
			TotalCustomers:    104,
			ActiveAgents:      8,
			MonthlyGrowthRate: 0.08,
			ChurnRate:         0.03,
			AverageARPU:       1200.00,
			TotalCommission:   31250.00,
			CommissionPercentage: 25.0,
		},
		TrendData: []TrendPoint{
			{Date: "2024-07", Value: 100000},
			{Date: "2024-08", Value: 110000},
			{Date: "2024-09", Value: 115000},
			{Date: "2024-10", Value: 120000},
			{Date: "2024-11", Value: 123000},
			{Date: "2024-12", Value: 125000},
		},
		TopAgents: []AgentPerformance{
			{AgentID: "agent-1", AgentName: "John Smith", MRR: 35000, Customers: 29, Commission: 8750, GrowthRate: 0.12},
			{AgentID: "agent-2", AgentName: "Jane Doe", MRR: 28000, Customers: 23, Commission: 7000, GrowthRate: 0.10},
			{AgentID: "agent-3", AgentName: "Mike Johnson", MRR: 22000, Customers: 18, Commission: 5500, GrowthRate: 0.08},
		},
		RevenueByTerritory: []TerritoryRevenue{
			{Territory: "Gauteng", MRR: 60000, ARR: 720000, Agents: 3, Customers: 50},
			{Territory: "Western Cape", MRR: 40000, ARR: 480000, Agents: 3, Customers: 33},
			{Territory: "KwaZulu-Natal", MRR: 25000, ARR: 300000, Agents: 2, Customers: 21},
		},
		LastUpdated: "2024-12-27T10:30:00Z",
	}
	
	c.JSON(http.StatusOK, response)
}

// GetAgentDashboard returns agent personal dashboard data
func GetAgentDashboard(c *gin.Context) {
	agentID := c.Param("agent_id")
	_ = c.GetHeader("X-Tenant-ID")

	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agent ID required"})
		return
	}
	// System users (system_admin, global_admin) may have empty tenant_id
	// TODO: Filter data by tenantID when implementing database queries
	
	// TODO: Fetch actual agent data
	response := gin.H{
		"agent_id":   agentID,
		"agent_name": "John Smith",
		"metrics": gin.H{
			"current_mrr":         35000.00,
			"current_arr":         420000.00,
			"total_customers":     29,
			"new_customers_mtd":   3,
			"commission_mtd":      8750.00,
			"commission_rate":     0.25,
			"effective_rate":      0.28,
			"next_tier_mrr":       40000.00,
			"next_tier_rate":      0.30,
			"progress_to_tier":    0.875,
		},
		"performance": gin.H{
			"mom_growth":    0.12,
			"yoy_growth":    0.45,
			"churn_rate":    0.02,
			"nrr":           1.15,
		},
		"recent_customers": []gin.H{
			{
				"customer_id":   "cust-1",
				"company_name":  "ABC Company",
				"mrr":           1500,
				"health_score":  85,
				"days_active":   45,
			},
		},
		"last_updated": "2024-12-27T10:30:00Z",
	}
	
	c.JSON(http.StatusOK, response)
}

// GetBusinessManagerDashboard returns business manager dashboard
func GetBusinessManagerDashboard(c *gin.Context) {
	_ = c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id
	// TODO: Filter data by tenantID when implementing database queries
	
	response := gin.H{
		"program_health": gin.H{
			"total_mrr":           125000.00,
			"total_arr":           1500000.00,
			"active_agents":       8,
			"total_customers":     104,
			"monthly_growth":      0.08,
			"quarterly_growth":    0.25,
			"churn_rate":          0.03,
			"nrr":                 1.12,
			"ltv_to_cac":          5.2,
			"payback_months":      8.5,
		},
		"agent_performance": []gin.H{
			{
				"agent_id":        "agent-1",
				"agent_name":      "John Smith",
				"mrr":             35000,
				"customers":       29,
				"growth_rate":     0.12,
				"churn_rate":      0.02,
				"tier":            "High Performer",
			},
		},
		"customer_health": gin.H{
			"healthy":    75,
			"at_risk":    20,
			"critical":   9,
		},
		"revenue_forecast": gin.H{
			"next_month_mrr":    135000,
			"next_quarter_mrr":  155000,
			"eoy_mrr":           180000,
			"confidence":        0.85,
		},
		"last_updated": "2024-12-27T10:30:00Z",
	}
	
	c.JSON(http.StatusOK, response)
}

// GetAnalystDashboard returns analyst dashboard with detailed metrics
func GetAnalystDashboard(c *gin.Context) {
	_ = c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id
	// TODO: Filter data by tenantID when implementing database queries
	
	response := gin.H{
		"acquisition_metrics": gin.H{
			"new_customers_mtd":     12,
			"new_customers_qtd":     35,
			"acquisition_cost":      1500.00,
			"conversion_rate":       0.25,
			"avg_sales_cycle_days":  45,
			"lead_velocity":         25,
		},
		"monetization_metrics": gin.H{
			"total_mrr":            125000.00,
			"total_arr":            1500000.00,
			"average_arpu":         1200.00,
			"expansion_mrr":        5000.00,
			"contraction_mrr":      1000.00,
			"nrr":                  1.12,
			"quick_ratio":          4.5,
		},
		"retention_metrics": gin.H{
			"gross_churn_rate":     0.03,
			"net_churn_rate":       -0.01,
			"customer_lifetime":    33.3,
			"retention_rate":       0.97,
		},
		"cohort_analysis": []gin.H{
			{
				"cohort":            "2024-Q1",
				"initial_customers": 20,
				"current_customers": 18,
				"retention_rate":    0.90,
				"ltv":               45000,
			},
		},
		"last_updated": "2024-12-27T10:30:00Z",
	}
	
	c.JSON(http.StatusOK, response)
}

// GetRealtimeMetrics returns real-time streaming metrics
func GetRealtimeMetrics(c *gin.Context) {
	_ = c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id
	// TODO: Filter data by tenantID when implementing database queries
	
	// Return current real-time snapshot
	response := gin.H{
		"timestamp": "2024-12-27T10:30:00Z",
		"metrics": gin.H{
			"active_users_now":    45,
			"registrations_today": 3,
			"mrr_added_today":     3600,
			"logins_today":        127,
			"support_tickets_open": 8,
		},
		"alerts": []gin.H{
			{
				"severity":    "warning",
				"message":     "Customer health score dropped below 50",
				"customer_id": "cust-123",
			},
		},
	}
	
	c.JSON(http.StatusOK, response)
}

