// User and Auth Types
export interface User {
	id: string;
	tenant_id: string;
	email: string;
	first_name?: string;
	last_name?: string;
	status: 'active' | 'suspended' | 'locked' | 'deleted';
	email_verified: boolean;
	mfa_enabled: boolean;
	created_at: string;
	updated_at: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	first_name: string;
	last_name: string;
}

export interface AuthResponse {
	access_token: string;
	refresh_token: string;
	token_type: string;
	expires_in: number;
	user: User;
	roles: string[];
}

export interface TokenClaims {
	user_id: string;
	tenant_id: string;
	email: string;
	roles: string[];
	exp: number;
	iat: number;
}

// Registration Types
export interface Registration {
	id: string;
	tenant_id: string;
	client_id: string;
	registration_type: 'pty_ltd' | 'close_corporation' | 'business_name' | 'vat_registration';
	company_name: string;
	registration_number?: string;
	jurisdiction: 'ZA' | 'ZW';
	status: 'draft' | 'submitted' | 'in_review' | 'approved' | 'rejected' | 'cancelled';
	submitted_at?: string;
	approved_at?: string;
	rejected_at?: string;
	rejection_reason?: string;
	assigned_to?: string;
	form_data?: Record<string, any>;
	created_at: string;
	updated_at: string;
}

export interface CreateRegistrationRequest {
	client_id: string;
	registration_type: string;
	company_name: string;
	jurisdiction: string;
	form_data?: Record<string, any>;
}

// Document Types
export interface Document {
	id: string;
	tenant_id: string;
	registration_id?: string;
	client_id?: string;
	document_type: string;
	file_name: string;
	file_size: number;
	mime_type: string;
	storage_path: string;
	status: 'pending' | 'verified' | 'rejected' | 'expired';
	verified_at?: string;
	ocr_processed: boolean;
	ai_verified: boolean;
	ai_verification_score?: number;
	created_at: string;
	updated_at: string;
}

// Commission Types
export interface Commission {
	id: string;
	tenant_id: string;
	registration_id: string;
	agent_id: string;
	registration_fee: number;
	commission_rate: number;
	commission_amount: number;
	currency: string;
	status: 'pending' | 'approved' | 'paid' | 'cancelled';
	approved_at?: string;
	paid_at?: string;
	payment_reference?: string;
	created_at: string;
	updated_at: string;
}

export interface CommissionSummary {
	total_commissions: number;
	total_pending: number;
	total_approved: number;
	total_paid: number;
	currency: string;
}

// API Response Types
export interface APIError {
	code: string;
	message: string;
	details?: any;
}

export interface PaginatedResponse<T> {
	data: T[];
	total: number;
	page: number;
	limit: number;
	has_next: boolean;
	has_prev: boolean;
}

export interface HealthCheck {
	service: string;
	version: string;
	status: 'healthy' | 'degraded' | 'unhealthy';
	checks: CheckResult[];
	timestamp: string;
}

export interface CheckResult {
	name: string;
	status: 'healthy' | 'degraded' | 'unhealthy';
	message?: string;
	latency?: string;
	timestamp: string;
}

// Client Types
export interface Client {
	id: string;
	tenant_id: string;
	client_type: 'individual' | 'company';
	full_name?: string;
	company_name?: string;
	email: string;
	phone?: string;
	mobile?: string;
	status: 'active' | 'inactive' | 'suspended';
	created_at: string;
	updated_at: string;
}

// Tenant Types
export interface Tenant {
	id: string;
	name: string;
	subdomain: string;
	status: 'active' | 'suspended' | 'deleted';
	subscription_tier: 'starter' | 'professional' | 'enterprise';
	created_at: string;
}
