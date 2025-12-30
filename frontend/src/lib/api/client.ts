import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig } from 'axios';
import type {
	AuthResponse,
	LoginRequest,
	RegisterRequest,
	Registration,
	CreateRegistrationRequest,
	Document,
	Commission,
	CommissionSummary,
	PaginatedResponse,
	HealthCheck,
	Client
} from '$types';
import { AuthManager } from '$lib/auth/AuthManager';

class APIClient {
	public client: AxiosInstance;
	private tenantSubdomain: string | null = null;
	private authManager: AuthManager;

	// Rate limiting queue
	private pendingRequests = 0;
	private maxConcurrentRequests = 2; // Limit to 2 concurrent requests to avoid 429s
	private requestQueue: Array<() => void> = [];

	constructor(baseURL: string = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1') {
		// Get AuthManager instance
		this.authManager = AuthManager.getInstance();
		// Set this API client in AuthManager for login/logout operations
		this.authManager.setApiClient(this);
		this.client = axios.create({
			baseURL,
			headers: {
				'Content-Type': 'application/json'
			},
			withCredentials: false  // Changed to false for direct API calls
		});

		// Add rate limiting interceptor
		this.client.interceptors.request.use(async (config) => {
			await new Promise<void>((resolve) => {
				this.enqueueRequest(resolve);
			});
			return config;
		});

		// CRITICAL FIX: Ensure tenant ID has a fallback value
		// Use environment variable or hardcoded default for development
		this.tenantSubdomain = import.meta.env.VITE_DEFAULT_TENANT_ID || '00000000-0000-0000-0000-000000000001';
		console.log('[API Client] Initialized with tenant ID:', this.tenantSubdomain);

		// Request interceptor to add auth token and tenant ID
		this.client.interceptors.request.use((config) => {
			// Get token from AuthManager
			const token = this.authManager.getAccessToken();
			if (token && !this.authManager.isTokenExpired(token)) {
				config.headers.Authorization = `Bearer ${token}`;
				console.log('[API Client] Request interceptor added Authorization header for:', config.url);
			} else if (token) {
				console.warn('[API Client] Access token expired for request:', config.url);
			} else {
				console.warn('[API Client] No access token found for request:', config.url);
			}

			// Get tenant ID from AuthManager
			const tenantId = this.authManager.getTenantId();
			if (tenantId) {
				config.headers['X-Tenant-ID'] = tenantId;
				console.log('[API Client] Using tenant ID:', tenantId);
			} else {
				// CRITICAL FIX: Always provide a tenant ID, even if it's the default
				config.headers['X-Tenant-ID'] = '00000000-0000-0000-0000-000000000001';
				console.warn('[API Client] No tenant ID found, using default:', '00000000-0000-0000-0000-000000000001');
			}

			return config;
		});

		// Response interceptor for error handling
		this.client.interceptors.response.use(
			(response) => {
				this.dequeueRequest();
				console.log('[API Client] Response received:', {
					url: response.config.url,
					status: response.status
				});
				return response;
			},
			async (error) => {
				this.dequeueRequest();

				// Handle 429 Too Many Requests
				if (error.response?.status === 429) {
					const retryAfter = parseInt(error.response.headers['retry-after']) || 2;
					const delay = retryAfter * 1000;
					console.warn(`[API Client] 429 Too Many Requests. Retrying after ${delay}ms...`);
					
					await new Promise(resolve => setTimeout(resolve, delay));
					return this.client.request(error.config);
				}

				console.error('[API Client] Response error:', {
					url: error.config?.url,
					method: error.config?.method,
					status: error.response?.status,
					data: error.response?.data
				});

				if (error.response?.status === 401) {
					console.log('[API Client] 401 error - checking if should redirect');

					// CRITICAL FIX: Don't redirect if this is already a login/refresh request
					const isAuthEndpoint = error.config?.url?.includes('/auth/login') ||
					                       error.config?.url?.includes('/auth/refresh');

					if (isAuthEndpoint) {
						console.log('[API Client] 401 on auth endpoint, rejecting without redirect');
						return Promise.reject(error);
					}

					// Try to refresh token using AuthManager
					const refreshed = await this.authManager.refreshAccessToken();
					if (refreshed) {
						console.log('[API Client] Token refreshed, retrying request');
						// Retry original request with new token
						return this.client.request(error.config);
					} else {
						// Refresh failed - AuthManager will handle logout
						console.log('[API Client] Refresh failed, auth cleared by AuthManager');
						// Just reject the error, let Svelte routing handle the redirect
						return Promise.reject(error);
					}
				}
				return Promise.reject(error);
			}
		);
	}

	private enqueueRequest(resolve: () => void) {
		if (this.pendingRequests < this.maxConcurrentRequests) {
			this.pendingRequests++;
			resolve();
		} else {
			this.requestQueue.push(resolve);
		}
	}

	private dequeueRequest() {
		this.pendingRequests--;
		if (this.requestQueue.length > 0) {
			const nextResolve = this.requestQueue.shift();
			if (nextResolve) {
				this.pendingRequests++;
				nextResolve();
			}
		}
	}

	setTenantSubdomain(subdomain: string) {
		this.tenantSubdomain = subdomain;
	}

	private getServiceURL(service: 'analytics' | 'agent' | 'integration' | 'ml' | 'chatbot' | 'compliance' | 'commission-agent'): string {
		const urls = {
			analytics: import.meta.env.VITE_ANALYTICS_API_URL || 'http://localhost:8088/api/v1',
			agent: import.meta.env.VITE_AGENT_API_URL || 'http://localhost:8092/api/v1',
			integration: import.meta.env.VITE_INTEGRATION_API_URL || 'http://localhost:8084/api/v1',
			ml: import.meta.env.VITE_ML_API_URL || 'http://localhost:8095/api/v1',
			chatbot: import.meta.env.VITE_CHATBOT_API_URL || 'http://localhost:8094/api/v1',
			compliance: import.meta.env.VITE_COMPLIANCE_API_URL || 'http://localhost:8096/api/v1',
			'commission-agent': import.meta.env.VITE_COMMISSION_AGENT_API_URL || 'http://localhost:8097/api/v1'
		};
		return urls[service];
	}

	private getTenantId(): string | null {
		// Get tenant ID from AuthManager first
		const tenantId = this.authManager.getTenantId();
		if (tenantId) return tenantId;

		// Fall back to subdomain or env variable
		return this.tenantSubdomain || import.meta.env.VITE_DEFAULT_TENANT_ID || '00000000-0000-0000-0000-000000000001';
	}

	private getTenantHeader(): string {
		return this.getTenantId() || '00000000-0000-0000-0000-000000000001';
	}

	private getAuthHeaders(): Record<string, string> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json'
		};

		const tenantId = this.getTenantId();
		if (tenantId) {
			headers['X-Tenant-ID'] = tenantId;
		}

		const token = this.authManager.getAccessToken();
		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
			console.log('[API Client] Auth headers include Authorization token');
		} else {
			console.warn('[API Client] No access token found for auth headers');
		}
		return headers;
	}

	// Auth endpoints
	async login(data: LoginRequest): Promise<AuthResponse> {
		console.log('[API Client] Sending login request to:', this.client.defaults.baseURL + '/auth/login');
		console.log('[API Client] Request data:', { email: data.email, password: '***' });
		console.log('[API Client] Headers:', {
			'Content-Type': 'application/json',
			'X-Tenant-ID': this.getTenantId()
		});

		try {
			const response = await this.client.post<AuthResponse>('/auth/login', data);
			console.log('[API Client] Login response received:', {
				status: response.status,
				hasToken: !!response.data.access_token,
				user: response.data.user?.email
			});

			// Store tokens in AuthManager
			this.authManager.setTokens(response.data.access_token, response.data.refresh_token);
			return response.data;

		} catch (error: any) {
			console.error('[API Client] Login error:', {
				message: error.message,
				response: error.response?.data,
				status: error.response?.status,
				statusText: error.response?.statusText,
				headers: error.response?.headers,
				cors: error.code === 'ERR_NETWORK' || error.message?.includes('CORS'),
				config: {
					url: error.config?.url,
					method: error.config?.method,
					baseURL: error.config?.baseURL,
					headers: error.config?.headers
				}
			});
			throw error;
		}
	}

	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await this.client.post<AuthResponse>('/auth/register', data);
		// Store tokens in AuthManager
		this.authManager.setTokens(response.data.access_token, response.data.refresh_token);
		return response.data;
	}

	async logout(): Promise<void> {
		try {
			await this.client.post('/auth/logout', {}, {
				headers: {
					'Content-Type': 'application/json'
				}
			});
		} catch (error) {
			console.error('[Auth] Logout API call failed:', error);
		}
		// Note: Token clearing is handled by AuthManager.logout()
	}

	async getProfile() {
		const response = await this.client.get('/auth/me');
		return response.data;
	}

	async verifyEmail(token: string): Promise<{ message: string }> {
		const response = await this.client.post('/auth/verify-email', { token });
		return response.data;
	}

	async resendVerification(email: string): Promise<{ message: string }> {
		const response = await this.client.post('/auth/resend-verification', { email });
		return response.data;
	}

	// Registration endpoints
	async getRegistrations(
		page: number = 1,
		limit: number = 20,
		filters?: {
			search?: string;
			status?: string;
			registration_type?: string;
			jurisdiction?: string;
			agent_id?: string;
			date_from?: string;
			date_to?: string;
			sort_by?: string;
			sort_order?: string;
		}
	): Promise<PaginatedResponse<Registration>> {
		const params: any = { page, limit };

		// Add filters if provided
		if (filters) {
			Object.keys(filters).forEach(key => {
				const value = (filters as any)[key];
				if (value !== undefined && value !== null && value !== '') {
					params[key] = value;
				}
			});
		}

		const response = await this.client.get<PaginatedResponse<Registration>>('/registrations', { params });
		return response.data;
	}

	async getRegistration(id: string): Promise<Registration> {
		const response = await this.client.get<Registration>(`/registrations/${id}`);
		return response.data;
	}

	async createRegistration(data: CreateRegistrationRequest): Promise<Registration> {
		const response = await this.client.post<Registration>('/registrations', data);
		return response.data;
	}

	async updateRegistration(id: string, data: Partial<Registration>): Promise<Registration> {
		const response = await this.client.put<Registration>(`/registrations/${id}`, data);
		return response.data;
	}

	async deleteRegistration(id: string): Promise<void> {
		await this.client.delete(`/registrations/${id}`);
	}

	// Registration workflow actions
	async submitRegistration(id: string): Promise<{ message: string }> {
		const response = await this.client.post<{ message: string }>(`/registrations/${id}/submit`);
		return response.data;
	}

	async approveRegistration(id: string): Promise<{ message: string }> {
		const response = await this.client.post<{ message: string }>(`/registrations/${id}/approve`);
		return response.data;
	}

	async rejectRegistration(id: string, reason: string): Promise<{ message: string }> {
		const response = await this.client.post<{ message: string }>(`/registrations/${id}/reject`, { reason });
		return response.data;
	}

	async cancelRegistration(id: string, reason?: string): Promise<{ message: string }> {
		const response = await this.client.post<{ message: string }>(`/registrations/${id}/cancel`, { reason });
		return response.data;
	}

	// Document endpoints
	async getDocuments(page: number = 1, limit: number = 20, registrationId?: string): Promise<PaginatedResponse<Document>> {
		const params: any = { page, limit };
		if (registrationId) {
			params.registration_id = registrationId;
		}
		const response = await this.client.get<PaginatedResponse<Document>>('/documents', { params });
		return response.data;
	}

	async getDocument(id: string): Promise<Document> {
		const response = await this.client.get<Document>(`/documents/${id}`);
		return response.data;
	}

	async uploadDocument(file: File, documentType: string, registrationId?: string): Promise<Document> {
		const formData = new FormData();
		formData.append('file', file);
		formData.append('document_type', documentType);
		if (registrationId) {
			formData.append('registration_id', registrationId);
		}

		const response = await this.client.post<Document>('/documents', formData, {
			headers: {
				'Content-Type': 'multipart/form-data'
			}
		});
		return response.data;
	}

	async updateDocument(id: string, data: Partial<Document>): Promise<Document> {
		const response = await this.client.put<Document>(`/documents/${id}`, data);
		return response.data;
	}

	async deleteDocument(id: string): Promise<void> {
		await this.client.delete(`/documents/${id}`);
	}

	async getDocumentDownloadURL(id: string): Promise<{ download_url: string }> {
		const response = await this.client.get<{ download_url: string }>(`/documents/${id}/download`);
		return response.data;
	}

	async verifyDocument(id: string): Promise<Document> {
		const response = await this.client.post<Document>(`/documents/${id}/verify`);
		return response.data;
	}

	async rejectDocument(id: string, reason: string): Promise<Document> {
		const response = await this.client.post<Document>(`/documents/${id}/reject`, { reason });
		return response.data;
	}

	// Commission endpoints
	async getCommissions(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Commission>> {
		const response = await this.client.get<PaginatedResponse<Commission>>('/commissions', {
			params: { page, limit }
		});
		return response.data;
	}

	async getCommission(id: string): Promise<Commission> {
		const response = await this.client.get<Commission>(`/commissions/${id}`);
		return response.data;
	}

	async createCommission(data: {
		registration_id: string;
		agent_id: string;
		registration_fee: number;
		commission_rate: number;
		currency: string;
	}): Promise<Commission> {
		const response = await this.client.post<Commission>('/commissions', data);
		return response.data;
	}

	async approveCommission(id: string): Promise<Commission> {
		const response = await this.client.post<Commission>(`/commissions/${id}/approve`);
		return response.data;
	}

	async payCommission(id: string, paymentReference: string): Promise<Commission> {
		const response = await this.client.post<Commission>(`/commissions/${id}/pay`, {
			payment_reference: paymentReference
		});
		return response.data;
	}

	async getCommissionSummary(agentId: string): Promise<CommissionSummary> {
		const response = await this.client.get<CommissionSummary>(`/commissions/summary/${agentId}`);
		return response.data;
	}

	// Client endpoints
	async getClients(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Client>> {
		const response = await this.client.get<PaginatedResponse<Client>>('/clients', {
			params: { page, limit }
		});
		return response.data;
	}

	async getClient(id: string): Promise<Client> {
		const response = await this.client.get<Client>(`/clients/${id}`);
		return response.data;
	}

	async createClient(data: Partial<Client>): Promise<Client> {
		const response = await this.client.post<Client>('/clients', data);
		return response.data;
	}

	async updateClient(id: string, data: Partial<Client>): Promise<Client> {
		const response = await this.client.put<Client>(`/clients/${id}`, data);
		return response.data;
	}

	async deleteClient(id: string): Promise<void> {
		await this.client.delete(`/clients/${id}`);
	}

	async getClientStats(clientId: string): Promise<{
		registrations_count: number;
		documents_count: number;
		total_revenue: number;
		total_paid: number;
	}> {
		const response = await this.client.get(`/clients/${clientId}/stats`);
		return response.data;
	}

	async getClientActivity(clientId: string, limit: number = 5): Promise<any[]> {
		const response = await this.client.get(`/clients/${clientId}/activity`, {
			params: { limit }
		});
		return response.data;
	}

	// Client Groups endpoints
	async getClientGroups(): Promise<any> {
		const response = await this.client.get('/clients/groups');
		return response.data;
	}

	async getClientGroup(id: string): Promise<any> {
		const response = await this.client.get(`/clients/groups/${id}`);
		return response.data;
	}

	async createClientGroup(data: {
		name: string;
		description?: string;
		color?: string;
	}): Promise<any> {
		const response = await this.client.post('/clients/groups', data);
		return response.data;
	}

	async updateClientGroup(id: string, data: {
		name?: string;
		description?: string;
		color?: string;
	}): Promise<any> {
		const response = await this.client.put(`/clients/groups/${id}`, data);
		return response.data;
	}

	async deleteClientGroup(id: string): Promise<void> {
		await this.client.delete(`/clients/groups/${id}`);
	}

	// Health check
	async healthCheck(): Promise<HealthCheck> {
		const response = await this.client.get<HealthCheck>('/health');
		return response.data;
	}

	// Analytics endpoints
	async getExecutiveDashboard(): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/dashboards/executive`, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async getAnalyticsTrends(period: 'daily' | 'weekly' | 'monthly' = 'monthly'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/trends`, {
			headers: this.getAuthHeaders(),
			params: { period }
		});
		return response.data;
	}

	async getTopPerformers(limit: number = 5): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/top-performers`, {
			headers: this.getAuthHeaders(),
			params: { limit }
		});
		return response.data;
	}

	// Agent endpoints (via API Gateway)
	async getAgents(page: number = 1, limit: number = 50): Promise<any> {
		const response = await this.client.get('/agents/performance', {
			params: { page, limit }
		});
		return response.data;
	}

	async getAgent(agentId: string): Promise<any> {
		const response = await this.client.get(`/agents/${agentId}`);
		return response.data;
	}

	async createAgent(agentData: any): Promise<any> {
		const response = await this.client.post('/agents', agentData);
		return response.data;
	}

	async updateAgent(agentId: string, agentData: any): Promise<any> {
		const response = await this.client.put(`/agents/${agentId}`, agentData);
		return response.data;
	}

	async deleteAgent(agentId: string): Promise<any> {
		const response = await this.client.delete(`/agents/${agentId}`);
		return response.data;
	}

	async getAgentPerformance(agentId: string): Promise<any> {
		const response = await this.client.get(`/agents/${agentId}/performance`);
		return response.data;
	}

	async getAgentPerformanceHistory(agentId: string): Promise<any> {
		const response = await this.client.get(`/agents/${agentId}/performance/history`);
		return response.data;
	}

	// Advanced Analytics & Reports Endpoints
	async getAnalyticsReports(period: 'daily' | 'weekly' | 'monthly' = 'monthly'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/reports`, {
			headers: this.getAuthHeaders(),
			params: { period }
		});
		return response.data;
	}

	async getRegistrationTrends(period: 'daily' | 'weekly' | 'monthly' = 'monthly'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/registration-trends`, {
			headers: this.getAuthHeaders(),
			params: { period }
		});
		return response.data;
	}

	async getCommissionAnalytics(period: 'daily' | 'weekly' | 'monthly' = 'monthly'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/commissions`, {
			headers: this.getAuthHeaders(),
			params: { period }
		});
		return response.data;
	}

	async getClientAcquisitionData(period: 'daily' | 'weekly' | 'monthly' = 'monthly'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/client-acquisition`, {
			headers: this.getAuthHeaders(),
			params: { period }
		});
		return response.data;
	}

	async exportReport(reportType: string, format: 'pdf' | 'csv' | 'excel'): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('analytics')}/analytics/export/${reportType}`, {
			headers: this.getAuthHeaders(),
			params: { format },
			responseType: format === 'pdf' ? 'blob' : 'json'
		});
		return response.data;
	}

	// Integration Health & Monitoring Endpoints
	async getIntegrationHealth(): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('integration')}/integrations/health`, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async getIntegrationStatus(integrationId: string): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('integration')}/integrations/${integrationId}/status`, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async testIntegration(integrationId: string): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('integration')}/integrations/${integrationId}/test`, {}, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async getIntegrationLogs(integrationId: string, limit: number = 50): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('integration')}/integrations/${integrationId}/logs`, {
			headers: this.getAuthHeaders(),
			params: { limit }
		});
		return response.data;
	}

	// Agent Performance Endpoints
	async getAllAgentPerformance(page: number = 1, limit: number = 20): Promise<any> {
		const response = await axios.get(`${this.getServiceURL('agent')}/agents/performance`, {
			headers: this.getAuthHeaders(),
			params: { page, limit }
		});
		return response.data;
	}

	async getAgentCommissionBreakdown(agentId: string, period: 'monthly' | 'quarterly' | 'yearly' = 'monthly'): Promise<any> {
		const response = await this.client.get(`/agents/${agentId}/breakdown`, {
			params: { period }
		});
		return response.data;
	}

	async getAgentComparisonMetrics(): Promise<any> {
		const response = await this.client.get('/agents/comparison');
		return response.data;
	}

	async getAgentTierManagement(): Promise<any> {
		const response = await this.client.get('/agents/tiers');
		return response.data;
	}

	// AI Agent endpoints
	async getAIAgents(): Promise<any[]> {
		const aiAgentURL = import.meta.env.VITE_AI_AGENT_API_URL || 'http://localhost:8098/api/v1';
		const response = await axios.get(`${aiAgentURL}/agents`, {
			headers: this.getAuthHeaders()
		});
		return response.data.agents || [];
	}

	async getAIAgentHealth(endpointUrl: string): Promise<any> {
		const healthUrl = endpointUrl.replace(/\/$/, '') + '/health';
		const response = await axios.get(healthUrl, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	// Analytics Agent (ML Service)
	async predictChurn(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('ml')}/analytics/predict/churn`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async forecastRevenue(input: any): Promise<any> {
		try {
			const response = await axios.post(`${this.getServiceURL('ml')}/analytics/forecast/revenue`, input, {
				headers: this.getAuthHeaders()
			});
			return response.data;
		} catch (error: any) {
			// If endpoint doesn't exist, return mock data
			if (error.response?.status === 404) {
				console.warn('[API Client] Revenue forecast endpoint not found, returning mock data');
				return {
					forecast: [],
					confidence: 0,
					message: 'Revenue forecasting is not yet available'
				};
			}
			throw error;
		}
	}

	/**
	 * Calculate Customer Lifetime Value using both Simple and DCF methods
	 *
	 * @param input - LTV calculation parameters
	 * @returns Object containing simple_ltv, dcf_ltv, and supporting metrics
	 *
	 * Backend endpoint: POST /api/v1/calculators/ltv (via analytics service)
	 */
	async calculateLTV(input: {
		average_mrr: number;
		gross_margin: number;
		monthly_churn_rate: number;
		discount_rate?: number;
		time_horizon_months?: number;
	}): Promise<any> {
		const response = await this.client.post('/calculators/ltv', input);
		return response.data;
	}

	async detectAnomalies(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('ml')}/analytics/detect/anomalies`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	// Chatbot Agent (via API Gateway)
	async sendChatMessage(params: {
		message: string;
		conversation_id?: string;
		context?: Record<string, any>;
	}): Promise<any> {
		const response = await this.client.post('/chat/message', params);
		return response.data;
	}

	async getConversationHistory(conversationId: string): Promise<any> {
		const response = await this.client.get(`/chat/conversations/${conversationId}`);
		return response.data;
	}

	async getAllConversations(): Promise<any> {
		const response = await this.client.get('/chat/conversations');
		return response.data;
	}

	async deleteConversation(conversationId: string): Promise<any> {
		const response = await this.client.delete(`/chat/conversations/${conversationId}`);
		return response.data;
	}

	async sendChatFeedback(params: {
		conversation_id: string;
		message_id: string;
		rating: number;
		comment?: string;
	}): Promise<any> {
		const response = await this.client.post('/chat/feedback', params);
		return response.data;
	}

	async checkChatbotHealth(): Promise<any> {
		const response = await this.client.get('/chat/health');
		return response.data;
	}

	// Commission Calculator Agent
	async calculateCommission(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('commission-agent')}/commission/calculate`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async forecastCommission(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('commission-agent')}/commission/forecast`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	// Compliance Validator Agent
	async validateCompanyName(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('compliance')}/validate/company-name`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async validateDocuments(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('compliance')}/validate/documents`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async validateCIPCRequirements(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('compliance')}/validate/cipc`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async validateSARSRequirements(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('compliance')}/validate/sars`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async comprehensiveValidation(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('compliance')}/validate/comprehensive`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	// Document Processor Agent
	async analyzeDocument(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('agent')}/documents/analyze`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	async compareDocuments(input: any): Promise<any> {
		const response = await axios.post(`${this.getServiceURL('agent')}/documents/compare`, input, {
			headers: this.getAuthHeaders()
		});
		return response.data;
	}

	// Notification endpoints
	async getNotifications(params?: any): Promise<any> {
		const response = await this.client.get('/notifications', { params });
		return response.data;
	}

	async getUnreadNotifications(page: number = 1, limit: number = 20): Promise<any> {
		const response = await this.client.get('/notifications/unread', {
			params: { page, limit }
		});
		return response.data;
	}

	async getUnreadCount(): Promise<number> {
		const response = await this.client.get('/notifications/unread/count');
		return response.data.count || 0;
	}

	async markNotificationAsRead(notificationId: string): Promise<void> {
		await this.client.post(`/notifications/${notificationId}/read`);
	}

	async markAllNotificationsAsRead(): Promise<void> {
		await this.client.post('/notifications/read-all');
	}

	async dismissNotification(notificationId: string): Promise<void> {
		await this.client.delete(`/notifications/${notificationId}`);
	}

	async clearAllNotifications(): Promise<void> {
		await this.client.delete('/notifications/clear-all');
	}

	async getNotificationPreferences(): Promise<any> {
		const response = await this.client.get('/notifications/preferences');
		return response.data;
	}

	async updateNotificationPreferences(preferences: any): Promise<any> {
		const response = await this.client.put('/notifications/preferences', preferences);
		return response.data;
	}

	async sendTestNotification(): Promise<void> {
		await this.client.post('/notifications/test');
	}

	// =====================================================
	// Profile Management Endpoints
	// =====================================================

	/**
	 * Get user profile
	 */
	async getUserProfile(): Promise<any> {
		const response = await this.client.get('/profile');
		return response.data;
	}

	/**
	 * Update user profile
	 */
	async updateUserProfile(data: {
		first_name?: string;
		last_name?: string;
		phone?: string;
		timezone?: string;
		language?: string;
		avatar_url?: string;
		bio?: string;
		address?: string;
	}): Promise<any> {
		const response = await this.client.put('/profile', data);
		return response.data;
	}

	/**
	 * Get user preferences
	 */
	async getUserPreferences(): Promise<{
		language: string;
		timezone: string;
		date_format: string;
		time_format: string;
		theme: 'light' | 'dark' | 'auto';
		currency: string;
		items_per_page: number;
	}> {
		const response = await this.client.get('/preferences');
		return response.data;
	}

	/**
	 * Update user preferences
	 */
	async updateUserPreferences(preferences: {
		language?: string;
		timezone?: string;
		date_format?: string;
		time_format?: string;
		theme?: 'light' | 'dark' | 'auto';
		currency?: string;
		items_per_page?: number;
	}): Promise<{ message: string }> {
		const response = await this.client.put('/preferences', preferences);
		return response.data;
	}

	// =====================================================
	// Security Settings Endpoints
	// =====================================================

	/**
	 * Change password
	 */
	async changePassword(data: {
		current_password: string;
		new_password: string;
	}): Promise<any> {
		const response = await this.client.post('/security/password/change', data);
		return response.data;
	}

	/**
	 * Setup 2FA - Get QR code and secret
	 */
	async setup2FA(method: 'totp' | 'sms' | 'email' = 'totp'): Promise<{
		qr_code: string;
		secret: string;
		backup_codes: string[];
	}> {
		const response = await this.client.post('/security/2fa/setup', { method });
		return response.data;
	}

	/**
	 * Verify 2FA code
	 */
	async verify2FA(code: string): Promise<{ valid: boolean }> {
		const response = await this.client.post('/security/2fa/verify', { code });
		return response.data;
	}

	/**
	 * Enable 2FA
	 */
	async enable2FA(code: string): Promise<{ message: string }> {
		const response = await this.client.post('/security/2fa/enable', { code });
		return response.data;
	}

	/**
	 * Disable 2FA
	 */
	async disable2FA(password: string): Promise<{ message: string }> {
		const response = await this.client.post('/security/2fa/disable', { password });
		return response.data;
	}

	/**
	 * Generate new backup codes for 2FA
	 */
	async generate2FABackupCodes(): Promise<{ backup_codes: string[] }> {
		const response = await this.client.post('/security/2fa/backup-codes');
		return response.data;
	}

	/**
	 * Update security preferences
	 */
	async updateSecurityPreferences(preferences: {
		sessionTimeout?: number;
		requirePasswordChange?: boolean;
		allowMultipleSessions?: boolean;
		emailOnNewLogin?: boolean;
	}): Promise<{ message: string }> {
		const response = await this.client.put('/security/preferences', preferences);
		return response.data;
	}

	/**
	 * Get security preferences
	 */
	async getSecurityPreferences(): Promise<{
		sessionTimeout: number;
		requirePasswordChange: boolean;
		allowMultipleSessions: boolean;
		emailOnNewLogin: boolean;
	}> {
		const response = await this.client.get('/security/preferences');
		return response.data;
	}

	/**
	 * Get active sessions
	 */
	async getActiveSessions(): Promise<any[]> {
		const response = await this.client.get('/security/sessions');
		return response.data.sessions || [];
	}

	/**
	 * Revoke specific session
	 */
	async revokeSession(sessionId: string): Promise<{ message: string }> {
		const response = await this.client.delete(`/security/sessions/${sessionId}`);
		return response.data;
	}

	/**
	 * Revoke all sessions except current
	 */
	async revokeAllSessions(): Promise<{ message: string; revoked_count: number }> {
		const response = await this.client.delete('/security/sessions/all');
		return response.data;
	}

	// =====================================================
	// Billing & Subscription Endpoints
	// =====================================================

	/**
	 * Get available subscription plans
	 */
	async getSubscriptionPlans(): Promise<any[]> {
		const response = await this.client.get('/billing/plans');
		return response.data.plans || [];
	}

	/**
	 * Get current tenant subscription
	 */
	async getCurrentSubscription(): Promise<any> {
		const response = await this.client.get('/billing/subscription');
		return response.data;
	}

	/**
	 * Update subscription plan
	 */
	async updateSubscription(data: { plan_id: string }): Promise<{ message: string }> {
		const response = await this.client.put('/billing/subscription', data);
		return response.data;
	}

	/**
	 * Cancel subscription
	 */
	async cancelSubscription(): Promise<{ message: string }> {
		const response = await this.client.post('/billing/subscription/cancel');
		return response.data;
	}

	/**
	 * Set default payment method
	 */
	async setDefaultPaymentMethod(methodId: string): Promise<{ message: string }> {
		const response = await this.client.put(`/billing/payment-methods/${methodId}/default`);
		return response.data;
	}

	/**
	 * Delete payment method
	 */
	async deletePaymentMethod(methodId: string): Promise<{ message: string }> {
		const response = await this.client.delete(`/billing/payment-methods/${methodId}`);
		return response.data;
	}

	/**
	 * Get usage metrics for current subscription
	 */
	async getUsageMetrics(): Promise<{
		registrations_used: number;
		registrations_limit: number | null;
		registrations_percentage: number;
		storage_used_mb: number;
		storage_limit_mb: number | null;
		storage_percentage: number;
		users_count: number;
		users_limit: number | null;
		users_percentage: number;
	}> {
		const response = await this.client.get('/billing/usage');
		return response.data;
	}

	/**
	 * Get invoice history
	 */
	async getInvoices(page: number = 1, limit: number = 20): Promise<PaginatedResponse<any>> {
		const response = await this.client.get('/billing/invoices', {
			params: { page, limit }
		});
		return response.data;
	}

	/**
	 * Get payment methods
	 */
	async getPaymentMethods(): Promise<any[]> {
		const response = await this.client.get('/billing/payment-methods');
		return response.data.payment_methods || [];
	}

	/**
	 * Update payment method
	 */
	async updatePaymentMethod(data: {
		type: 'card' | 'bank_account';
		card_number?: string;
		card_expiry?: string;
		card_cvv?: string;
		cardholder_name?: string;
		billing_address?: string;
	}): Promise<any> {
		const response = await this.client.put('/billing/payment-methods', data);
		return response.data;
	}

	// =====================================================
	// Analytics & Insights Endpoints
	// =====================================================

	/**
	 * Get automated daily insights from Analytics Agent
	 */
	async getDailyInsights(): Promise<any> {
		try {
			const response = await this.client.get('/insights/daily', {
				headers: {
					'X-Tenant-ID': this.getTenantHeader()
				}
			});
			return response.data;
		} catch (error: any) {
			// Return empty array if endpoint doesn't exist or fails
			console.warn('[API Client] getDailyInsights failed:', error.message);
			return { data: [] };
		}
	}

	/**
	 * Get current alerts and anomalies
	 */
	async getAlerts(): Promise<any> {
		const response = await this.client.get('/insights/alerts');
		return response.data;
	}

	/**
	 * Get recommendations for specific entity
	 */
	async getRecommendations(entityType: string, entityId: string): Promise<any> {
		const response = await this.client.get(`/insights/recommendations/${entityType}/${entityId}`);
		return response.data;
	}

	/**
	 * Get commission scenarios by tier
	 */
	async getCommissionScenarios(input: {
		monthly_registrations: number;
		avg_registration_value: number;
		service_type: string;
	}): Promise<any> {
		const response = await this.client.post('/calculators/commission/scenarios', input);
		return response.data;
	}

	/**
	 * Get revenue forecast scenarios
	 */
	async getForecastScenarios(input: {
		current_mrr: number;
		growth_rates: number[];
		churn_rates: number[];
		forecast_months: number;
	}): Promise<any> {
		const response = await this.client.post('/calculators/forecast/scenarios', input);
		return response.data;
	}

	/**
	 * Detect anomaly in metric
	 */
	async detectAnomaly(input: {
		metric_name: string;
		current_value: number;
		historical_values: number[];
		sensitivity: number;
	}): Promise<any> {
		const response = await this.client.post('/insights/detect-anomaly', input);
		return response.data;
	}

	// =====================================================
	// API Keys Management Endpoints
	// =====================================================

	/**
	 * Get all API keys for current user/tenant
	 */
	async getAPIKeys(): Promise<any[]> {
		const response = await this.client.get('/api-keys');
		return response.data.keys || [];
	}

	/**
	 * Create new API key
	 */
	async createAPIKey(data: {
		name: string;
		permissions: string[];
		expires_at?: string;
	}): Promise<{
		key: string;
		id: string;
		name: string;
		permissions: string[];
		created_at: string;
	}> {
		const response = await this.client.post('/api-keys', data);
		return response.data;
	}

	/**
	 * Revoke an API key
	 */
	async revokeAPIKey(keyId: string): Promise<{ message: string }> {
		const response = await this.client.delete(`/api-keys/${keyId}`);
		return response.data;
	}

	/**
	 * Update API key permissions
	 */
	async updateAPIKey(keyId: string, data: {
		name?: string;
		permissions?: string[];
	}): Promise<{ message: string }> {
		const response = await this.client.put(`/api-keys/${keyId}`, data);
		return response.data;
	}

	// =====================================================
	// Generic HTTP Methods for Stores (Public Access)
	// =====================================================

	/**
	 * Generic GET method for use by stores and other modules
	 * Includes auth headers and tenant context automatically via interceptors
	 */
	async get(url: string, config?: AxiosRequestConfig): Promise<any> {
		return this.client.get(url, config);
	}

	/**
	 * Generic POST method for use by stores and other modules
	 */
	async post(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
		return this.client.post(url, data, config);
	}

	/**
	 * Generic PUT method for use by stores and other modules
	 */
	async put(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
		return this.client.put(url, data, config);
	}

	/**
	 * Generic DELETE method for use by stores and other modules
	 */
	async delete(url: string, config?: AxiosRequestConfig): Promise<any> {
		return this.client.delete(url, config);
	}

	// =====================================================
	// Admin Management Endpoints
	// =====================================================

	/**
	 * Get all users (admin only)
	 */
	async getUsers(page: number = 1, limit: number = 20, filters?: any): Promise<PaginatedResponse<any>> {
		const response = await this.client.get('/admin/users', {
			params: { page, limit, ...filters }
		});
		return response.data;
	}

	/**
	 * Get user by ID (admin only)
	 */
	async getUser(userId: string): Promise<any> {
		const response = await this.client.get(`/admin/users/${userId}`);
		return response.data;
	}

	/**
	 * Create new user (admin only)
	 */
	async createUser(data: {
		email: string;
		password: string;
		first_name: string;
		last_name: string;
		roles?: string[];
	}): Promise<any> {
		const response = await this.client.post('/admin/users', data);
		return response.data;
	}

	/**
	 * Update user (admin only)
	 */
	async updateUser(userId: string, data: any): Promise<any> {
		const response = await this.client.put(`/admin/users/${userId}`, data);
		return response.data;
	}

	/**
	 * Delete user (admin only)
	 */
	async deleteUser(userId: string): Promise<void> {
		await this.client.delete(`/admin/users/${userId}`);
	}

	/**
	 * Get user's effective permissions (admin)
	 */
	async getUserEffectivePermissions(userId: string): Promise<{ permissions: string[] }> {
		const response = await this.client.get(`/admin/users/${userId}/effective-permissions`);
		return response.data;
	}

	/**
	 * Get enabled features for tenant
	 */
	async getEnabledFeatures(): Promise<{ features: any[]; feature_codes: string[] }> {
		const response = await this.client.get('/admin/features/enabled');
		return response.data;
	}

	/**
	 * Get tenant information
	 */
	async getTenantInfo(): Promise<any> {
		const response = await this.client.get('/tenants/me');
		return response.data;
	}
}

// Export singleton instance
export const apiClient = new APIClient();
export default apiClient;
