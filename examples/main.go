package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/remiges-sachin/sqld"
	"github.com/remiges-sachin/sqld/examples/db/sqlc-gen"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// EmployeeIDScanner handles scanning for EmployeeID type
type EmployeeIDScanner struct {
	valid bool
	value EmployeeID
}

func (s *EmployeeIDScanner) Scan(src interface{}) error {
	log.Printf("EmployeeIDScanner.Scan called with: %v of type %T", src, src) // Debug print
	if src == nil {
		s.valid = false
		return nil
	}

	switch v := src.(type) {
	case int64:
		s.value = EmployeeID(v)
		s.valid = true
	case int32:
		s.value = EmployeeID(v)
		s.valid = true
	case int:
		s.value = EmployeeID(v)
		s.valid = true
	default:
		s.valid = false
		return fmt.Errorf("cannot scan type %T into EmployeeID", src)
	}
	return nil
}

func (s *EmployeeIDScanner) Value() interface{} {
	if !s.valid {
		return nil
	}
	return s.value
}

// EmployeeID is a custom type for employee IDs
type EmployeeID int64

// Value implements driver.Valuer
func (id EmployeeID) Value() (driver.Value, error) {
	return int64(id), nil
}

// Employee represents our database model matching the employees table
type Employee struct {
	ID         EmployeeID `json:"id" db:"id"`
	FirstName  string     `json:"first_name" db:"first_name"`
	LastName   string     `json:"last_name" db:"last_name"`
	Email      string     `json:"email" db:"email"`
	Phone      string     `json:"phone" db:"phone"`
	HireDate   time.Time  `json:"hire_date" db:"hire_date"`
	Salary     float64    `json:"salary" db:"salary"`
	Department string     `json:"department" db:"department"`
	Position   string     `json:"position" db:"position"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

func (Employee) TableName() string {
	return "employees"
}

// Account represents our database model matching the accounts table
type Account struct {
	ID            int64     `json:"id" db:"id"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	AccountName   string    `json:"account_name" db:"account_name"`
	AccountType   string    `json:"account_type" db:"account_type"`
	Balance       float64   `json:"balance" db:"balance"`
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"`
	OwnerID       *int64    `json:"owner_id" db:"owner_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func (Account) TableName() string {
	return "accounts"
}

// AccountQueryParams represents parameters for account queries
type AccountQueryParams struct {
	MinBalance float64 `db:"min_balance" json:"min_balance"`
}

type Server struct {
	db *pgx.Conn
}

func NewServer(db *pgx.Conn) *Server {
	return &Server{db: db}
}

// DynamicQueryHandler demonstrates dynamic field selection and filtering
func (s *Server) DynamicQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req sqld.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := sqld.Execute[Employee](r.Context(), s.db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PaginatedQueryHandler demonstrates pagination with dynamic queries
func (s *Server) PaginatedQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req sqld.QueryRequest
	req.Pagination = &sqld.PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := sqld.Execute[Employee](r.Context(), s.db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CustomFilterHandler demonstrates using custom WHERE conditions
func (s *Server) CustomFilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req sqld.QueryRequest
	req.Select = []string{"id", "account_number", "balance", "status"}
	req.Where = map[string]interface{}{
		"status":  "active",
		"balance": 1000.00,
	}

	resp, err := sqld.Execute[Account](r.Context(), s.db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SQLCQueryHandler demonstrates integration with SQLC-generated types
func (s *Server) SQLCQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req sqld.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := sqld.Execute[sqlc.Employee](r.Context(), s.db, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) UCCQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Request structure to capture the parameter and the flag
	var reqData struct {
		UccParams     sqlc.UCCListParams `json:"params"`
		IncludeParent bool               `json:"include_parent_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	p := reqData.UccParams
	if !p.Limit.Valid {
		p.Limit.Int32 = 10
		p.Limit.Valid = true
	}
	if !p.Offset.Valid {
		p.Offset.Int32 = 0
		p.Offset.Valid = true
	}

	if p.SortBy == nil {
		p.SortBy = "client_code"
	}
	if p.SortOrder == nil {
		p.SortOrder = "A"
	}

	nullifyEmptyString := func(t *pgtype.Text) {
		if t.Valid && t.String == "" {
			t.Valid = false
		}
	}

	nullifyEmptyString(&p.ClientCode)
	nullifyEmptyString(&p.MemberCode)
	nullifyEmptyString(&p.ParentClientCode)
	nullifyEmptyString(&p.Search)

	paramMap := map[string]interface{}{
		"ucc_status":         p.UccStatus,
		"tax_status":         p.TaxStatus,
		"holding_nature":     p.HoldingNature,
		"is_client_physical": p.IsClientPhysical,
		"is_client_demat":    p.IsClientDemat,
		"client_code":        p.ClientCode,
		"member_code":        p.MemberCode,
		"parent_client_code": p.ParentClientCode,
		"search":             p.Search,
		"sort_by":            p.SortBy,
		"sort_order":         p.SortOrder,
		"offset":             p.Offset,
		"limit":              p.Limit,
	}

	// Base SELECT columns without parent_client_code
	selectCols := `
        u.ucc_id,
        u.client_code,
        u.member_code,
        ts.tax_name AS tax_status,
        crm.entity_key AS holding_nature,
        ucc_status.entity_key AS ucc_status,
        u.is_client_physical,
        u.is_client_demat,
        COALESCE(
            CASE
                WHEN pd.first_name IS NOT NULL THEN CONCAT(pd.first_name, ' ', pd.middle_name, ' ', pd.last_name)
                WHEN nid.org_name IS NOT NULL THEN nid.org_name
                ELSE ''
            END, 
        '') AS primary_holder_name
    `

	// If IncludeParent flag is true, add the parent_client_code column
	if reqData.IncludeParent {
		selectCols = `
            u.ucc_id,
            u.client_code,
            u.member_code,
            ts.tax_name AS tax_status,
            crm.entity_key AS holding_nature,
            ucc_status.entity_key AS ucc_status,
            u.is_client_physical,
            u.is_client_demat,
            u.parent_client_code,
            COALESCE(
                CASE
                    WHEN pd.first_name IS NOT NULL THEN CONCAT(pd.first_name, ' ', pd.middle_name, ' ', pd.last_name)
                    WHEN nid.org_name IS NOT NULL THEN nid.org_name
                    ELSE ''
                END, 
            '') AS primary_holder_name
        `
	}

	query := fmt.Sprintf(`
    SELECT
        %s
    FROM
        ucc u
    LEFT JOIN
        holder h ON u.ucc_id = h.ref_id 
        AND h.holder_rank = (
            SELECT id 
            FROM common_reference_master 
            WHERE entity = 'UCC_HOLDER_RANK' AND entity_key = 'FIRST'
        )
        AND h.deleted_at IS NULL 
        AND h.deleted_by IS NULL
    LEFT JOIN
        person_detail pd ON pd.ref_id = h.id AND h.ref_id = u.ucc_id AND pd.entity_type = 'HOLDER'
    LEFT JOIN 
        non_individual_detail nid ON nid.ref_id = h.id AND nid.entity_type = 'HOLDER'
    LEFT JOIN
        common_reference_master AS crm ON u.holding_nature = crm.id AND crm.entity = 'UCC_HOLDING_TYPE'
    LEFT JOIN
        common_reference_master AS ucc_status ON u.ucc_status = ucc_status.id AND ucc_status.entity = 'UCC_ACC_STATUS'
    LEFT JOIN
        tax_status_master AS ts ON u.tax_status = ts.id
    WHERE
        (u.ucc_status = {{ucc_status}}::bigint OR {{ucc_status}}::bigint IS NULL OR {{ucc_status}}::bigint = 0) AND 
        (u.tax_status = {{tax_status}}::bigint OR {{tax_status}}::bigint IS NULL OR {{tax_status}}::bigint = 0) AND
        (u.holding_nature = {{holding_nature}}::bigint OR {{holding_nature}}::bigint IS NULL OR {{holding_nature}}::bigint = 0) AND
        (
            (u.is_client_physical = {{is_client_physical}} OR {{is_client_physical}} IS NULL OR {{is_client_physical}} = false) AND 
            (u.is_client_demat = {{is_client_demat}} OR {{is_client_demat}} IS NULL OR {{is_client_demat}} = false)
        ) AND NOT (u.is_client_physical = false AND u.is_client_demat = false) AND 
        u.deleted_at IS NULL AND u.deleted_by IS NULL AND
        ({{client_code}}::text IS NULL OR {{client_code}}::text = '' OR u.client_code ILIKE {{client_code}}::text || '%%') AND
        ({{member_code}}::text IS NULL OR {{member_code}}::text = '' OR u.member_code ILIKE {{member_code}}::text || '%%') AND
        ({{parent_client_code}}::text IS NULL OR {{parent_client_code}}::text = '' OR u.parent_client_code ILIKE {{parent_client_code}}::text || '%%') AND
        (
            {{search}}::text IS NULL 
            OR {{search}}::text = '' 
            OR (
                CONCAT_WS(' ', 
                    u.client_code,  
                    u.member_code,
                    u.parent_client_code,
                    ts.tax_name, 
                    crm.entity_key,
                    COALESCE(
                        pd.first_name || ' ' || pd.middle_name || ' ' || pd.last_name,
                        nid.org_name
                    ),
                    ucc_status.entity_key,
                    CASE WHEN u.is_client_physical THEN 'true' ELSE 'false' END,
                    CASE WHEN u.is_client_demat THEN 'true' ELSE 'false' END,
                    COALESCE(pd.first_name, nid.org_name)
                ) ILIKE '%%' || {{search}}::text || '%%'
            )
        )
    ORDER BY
        CASE 
            WHEN {{sort_by}} = 'ucc_id' AND upper({{sort_order}}) = 'A' THEN u.ucc_id
        END ASC,
        CASE 
            WHEN {{sort_by}} = 'ucc_id' AND upper({{sort_order}}) = 'D' THEN u.ucc_id
        END DESC,
        CASE
            WHEN {{sort_by}} = 'client_code' AND upper({{sort_order}}) = 'A' THEN u.client_code
        END ASC,
        CASE 
            WHEN {{sort_by}} = 'client_code' AND upper({{sort_order}}) = 'D' THEN u.client_code
        END DESC,
        CASE
            WHEN {{sort_by}} = 'member_code' AND upper({{sort_order}}) = 'A' THEN u.member_code
        END ASC,
        CASE 
            WHEN {{sort_by}} = 'member_code' AND upper({{sort_order}}) = 'D' THEN u.member_code
        END DESC,
        CASE
            WHEN {{sort_by}} = 'tax_status' AND upper({{sort_order}}) = 'A' THEN ts.tax_name
        END ASC,
        CASE 
            WHEN {{sort_by}} = 'tax_status' AND upper({{sort_order}}) = 'D' THEN ts.tax_name
        END DESC,
        CASE
            WHEN {{sort_by}} = 'primary_holder_name' AND upper({{sort_order}}) = 'A' THEN COALESCE(pd.first_name, nid.org_name)
        END ASC,
        CASE 
            WHEN {{sort_by}} = 'primary_holder_name' AND upper({{sort_order}}) = 'D' THEN COALESCE(pd.first_name, nid.org_name)
        END DESC
    OFFSET {{offset}}
    LIMIT {{limit}}
    `, selectCols)

	results, err := sqld.ExecuteRaw[sqlc.UCCListParams, sqlc.UCCListRow](r.Context(), s.db, query, paramMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Data []map[string]interface{} `json:"data"`
	}{
		Data: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type SimpleQueryParams struct {
	Department string  `db:"department" json:"department"`
	MinSalary  float64 `db:"min_salary" json:"min_salary"`
}

type EmployeeRow struct {
	ID         int64   `db:"id" json:"id"`
	FirstName  string  `db:"first_name" json:"first_name"`
	LastName   string  `db:"last_name" json:"last_name"`
	Department string  `db:"department" json:"department"`
	Salary     float64 `db:"salary" json:"salary"`
}

// SimpleQueryHandler demonstrates a basic usage of ExecuteRaw
func (s *Server) RawSimpleQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params SimpleQueryParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple query with parameter placeholders
	query := `
		SELECT id, first_name, last_name, department, salary
		FROM employees
		WHERE department = {{department}}
		AND salary >= {{min_salary}}
		ORDER BY salary DESC
	`

	paramMap := map[string]interface{}{
		"department": params.Department,
		"min_salary": params.MinSalary,
	}

	results, err := sqld.ExecuteRaw[SimpleQueryParams, EmployeeRow](
		r.Context(),
		s.db,
		query,
		paramMap,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Data []map[string]interface{} `json:"data"`
	}{
		Data: results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RawQueryHandler demonstrates using ExecuteRaw for custom SQL queries
func (s *Server) RawQueryJoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT 
			a.id, 
			a.account_number, 
			a.balance, 
			COALESCE(e.first_name, 'Unknown') as owner_name
		FROM accounts a
		LEFT JOIN employees e ON e.id = a.owner_id
		WHERE a.balance > {{min_balance}}
		ORDER BY a.balance DESC
	`

	// Define reference struct for parameter validation
	type QueryParams struct {
		MinBalance float64 `db:"min_balance" json:"min_balance"`
	}

	params := map[string]interface{}{
		"min_balance": 1000.00,
	}

	// Define result struct that matches our query
	type Result struct {
		ID            int64   `db:"id" json:"id"`
		AccountNumber string  `db:"account_number" json:"account_number"`
		Balance       float64 `db:"balance" json:"balance"`
		OwnerName     string  `db:"owner_name" json:"owner_name"`
	}

	log.Printf("Executing query: %s with params: %v", query, params)
	results, err := sqld.ExecuteRaw[QueryParams, Result](r.Context(), s.db, query, params)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got %d results: %+v", len(results), results)
	if len(results) == 0 {
		log.Println("No results found")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func init() {
	sqld.RegisterScanner(reflect.TypeOf((*EmployeeID)(nil)).Elem(), func() sql.Scanner {
		return &EmployeeIDScanner{}
	})

	if err := sqld.Register(sqlc.Employee{}); err != nil {
		log.Fatalf("failed to register sqlc.Employee model: %v", err)
	}
}

func main() {
	ctx := context.Background()

	config, err := pgx.ParseConfig("postgres://alyatest:alyatest@localhost:5432/alyatest?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	if err := sqld.Register(Employee{}); err != nil {
		log.Fatal(err)
	}
	if err := sqld.Register(Account{}); err != nil {
		log.Fatal(err)
	}
	if err := sqld.Register(sqlc.Employee{}); err != nil {
		log.Fatal(err)
	}

	server := NewServer(conn)

	// Execute endpoints - Using dynamic query building
	http.HandleFunc("/api/dynamic", server.DynamicQueryHandler)
	http.HandleFunc("/api/paginated", server.PaginatedQueryHandler)
	http.HandleFunc("/api/filtered", server.CustomFilterHandler)
	http.HandleFunc("/api/sqlc", server.SQLCQueryHandler)

	// ExecuteRaw endpoints - Using raw SQL queries
	http.HandleFunc("/api/rawquery-join", server.RawQueryJoinHandler)     // Complex query with JOIN
	http.HandleFunc("/api/rawquery-simple", server.RawSimpleQueryHandler) // Simple single-table query
	http.HandleFunc("/api/ucc-queries", server.UCCQueryHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
