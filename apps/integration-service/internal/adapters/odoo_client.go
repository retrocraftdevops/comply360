package adapters

import (
	"fmt"
	"log"

	"github.com/kolo/xmlrpc"
)

// OdooConfig holds Odoo connection configuration
type OdooConfig struct {
	URL      string
	Database string
	Username string
	Password string
}

// OdooClient handles XML-RPC communication with Odoo
type OdooClient struct {
	config *OdooConfig
	uid    int
}

// NewOdooClient creates a new Odoo client and authenticates
func NewOdooClient(config *OdooConfig) (*OdooClient, error) {
	client := &OdooClient{
		config: config,
	}

	// Authenticate and get UID
	if err := client.authenticate(); err != nil {
		return nil, fmt.Errorf("failed to authenticate with Odoo: %w", err)
	}

	log.Printf("Successfully authenticated with Odoo as user: %s (UID: %d)", config.Username, client.uid)
	return client, nil
}

// authenticate authenticates with Odoo and stores the user ID
func (c *OdooClient) authenticate() error {
	commonClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/common", c.config.URL), nil)
	if err != nil {
		return fmt.Errorf("failed to create XML-RPC client: %w", err)
	}
	defer commonClient.Close()

	var uid int
	err = commonClient.Call("authenticate", []interface{}{
		c.config.Database,
		c.config.Username,
		c.config.Password,
		map[string]interface{}{},
	}, &uid)

	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if uid == 0 {
		return fmt.Errorf("authentication failed: invalid credentials")
	}

	c.uid = uid
	return nil
}

// Execute executes a method on an Odoo model
func (c *OdooClient) Execute(model string, method string, args []interface{}) (interface{}, error) {
	objectClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", c.config.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create object client: %w", err)
	}
	defer objectClient.Close()

	var result interface{}
	err = objectClient.Call("execute_kw", []interface{}{
		c.config.Database,
		c.uid,
		c.config.Password,
		model,
		method,
		args,
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("execute failed for %s.%s: %w", model, method, err)
	}

	return result, nil
}

// ExecuteKw executes a method with keyword arguments
func (c *OdooClient) ExecuteKw(model string, method string, args []interface{}, kwargs map[string]interface{}) (interface{}, error) {
	objectClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", c.config.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create object client: %w", err)
	}
	defer objectClient.Close()

	var result interface{}
	err = objectClient.Call("execute_kw", []interface{}{
		c.config.Database,
		c.uid,
		c.config.Password,
		model,
		method,
		args,
		kwargs,
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("execute_kw failed for %s.%s: %w", model, method, err)
	}

	return result, nil
}

// Search searches for records matching a domain
func (c *OdooClient) Search(model string, domain []interface{}, options map[string]interface{}) ([]int, error) {
	result, err := c.ExecuteKw(model, "search", []interface{}{domain}, options)
	if err != nil {
		return nil, err
	}

	// Convert result to []int
	if resultSlice, ok := result.([]interface{}); ok {
		ids := make([]int, len(resultSlice))
		for i, v := range resultSlice {
			if id, ok := v.(int); ok {
				ids[i] = id
			} else if id64, ok := v.(int64); ok {
				ids[i] = int(id64)
			}
		}
		return ids, nil
	}

	return nil, fmt.Errorf("unexpected search result type")
}

// Read reads records by IDs
func (c *OdooClient) Read(model string, ids []int, fields []string) ([]map[string]interface{}, error) {
	options := map[string]interface{}{
		"fields": fields,
	}

	result, err := c.ExecuteKw(model, "read", []interface{}{ids}, options)
	if err != nil {
		return nil, err
	}

	// Convert result to []map[string]interface{}
	if resultSlice, ok := result.([]interface{}); ok {
		records := make([]map[string]interface{}, len(resultSlice))
		for i, v := range resultSlice {
			if record, ok := v.(map[string]interface{}); ok {
				records[i] = record
			}
		}
		return records, nil
	}

	return nil, fmt.Errorf("unexpected read result type")
}

// SearchRead searches and reads records in one call
func (c *OdooClient) SearchRead(model string, domain []interface{}, fields []string, options map[string]interface{}) ([]map[string]interface{}, error) {
	kwargs := map[string]interface{}{
		"fields": fields,
	}

	// Merge additional options
	for k, v := range options {
		kwargs[k] = v
	}

	result, err := c.ExecuteKw(model, "search_read", []interface{}{domain}, kwargs)
	if err != nil {
		return nil, err
	}

	// Convert result to []map[string]interface{}
	if resultSlice, ok := result.([]interface{}); ok {
		records := make([]map[string]interface{}, len(resultSlice))
		for i, v := range resultSlice {
			if record, ok := v.(map[string]interface{}); ok {
				records[i] = record
			}
		}
		return records, nil
	}

	return nil, fmt.Errorf("unexpected search_read result type")
}

// Create creates a new record
func (c *OdooClient) Create(model string, values map[string]interface{}) (int, error) {
	result, err := c.Execute(model, "create", []interface{}{values})
	if err != nil {
		return 0, err
	}

	// Convert result to int
	if id, ok := result.(int); ok {
		return id, nil
	}
	if id64, ok := result.(int64); ok {
		return int(id64), nil
	}

	return 0, fmt.Errorf("unexpected create result type")
}

// Write updates existing records
func (c *OdooClient) Write(model string, ids []int, values map[string]interface{}) (bool, error) {
	result, err := c.Execute(model, "write", []interface{}{ids, values})
	if err != nil {
		return false, err
	}

	if success, ok := result.(bool); ok {
		return success, nil
	}

	return false, fmt.Errorf("unexpected write result type")
}

// Unlink deletes records
func (c *OdooClient) Unlink(model string, ids []int) (bool, error) {
	result, err := c.Execute(model, "unlink", []interface{}{ids})
	if err != nil {
		return false, err
	}

	if success, ok := result.(bool); ok {
		return success, nil
	}

	return false, fmt.Errorf("unexpected unlink result type")
}

// GetServerVersion returns the Odoo server version
func (c *OdooClient) GetServerVersion() (map[string]interface{}, error) {
	commonClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/common", c.config.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create XML-RPC client: %w", err)
	}
	defer commonClient.Close()

	var result map[string]interface{}
	err = commonClient.Call("version", []interface{}{}, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return result, nil
}

// CheckAccessRights checks if the user has specific access rights on a model
func (c *OdooClient) CheckAccessRights(model string, operation string) (bool, error) {
	result, err := c.Execute(model, "check_access_rights", []interface{}{operation, false})
	if err != nil {
		return false, err
	}

	if hasAccess, ok := result.(bool); ok {
		return hasAccess, nil
	}

	return false, nil
}
