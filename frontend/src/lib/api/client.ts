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

class APIClient {
	private client: AxiosInstance;
	private tenantSubdomain: string | null = null;

	constructor(baseURL: string = '/api/v1') {
		this.client = axios.create({
			baseURL,
			headers: {
				'Content-Type': 'application/json'
			},
			withCredentials: true
		});

		// Set default tenant ID for development
		this.tenantSubdomain = '9ac5aa3e-91cd-451f-b182-563b0d751dc7';

		// Request interceptor to add auth token and tenant ID
		this.client.interceptors.request.use((config) => {
			const token = this.getAccessToken();
			if (token) {
				config.headers.Authorization = `Bearer ${token}`;
			}

			if (this.tenantSubdomain) {
				config.headers['X-Tenant-ID'] = this.tenantSubdomain;
			}

			return config;
		});

		// Response interceptor for error handling
		this.client.interceptors.response.use(
			(response) => response,
			async (error) => {
				if (error.response?.status === 401) {
					// Token expired, try to refresh
					const refreshed = await this.refreshToken();
					if (refreshed) {
						// Retry original request
						return this.client.request(error.config);
					} else {
						// Redirect to login
						window.location.href = '/auth/login';
					}
				}
				return Promise.reject(error);
			}
		);
	}

	setTenantSubdomain(subdomain: string) {
		this.tenantSubdomain = subdomain;
	}

	private getAccessToken(): string | null {
		return localStorage.getItem('access_token');
	}

	private setAccessToken(token: string) {
		localStorage.setItem('access_token', token);
	}

	private getRefreshToken(): string | null {
		return localStorage.getItem('refresh_token');
	}

	private setRefreshToken(token: string) {
		localStorage.setItem('refresh_token', token);
	}

	private clearTokens() {
		localStorage.removeItem('access_token');
		localStorage.removeItem('refresh_token');
	}

	// Auth endpoints
	async login(data: LoginRequest): Promise<AuthResponse> {
		const response = await this.client.post<AuthResponse>('/auth/login', data);
		this.setAccessToken(response.data.access_token);
		this.setRefreshToken(response.data.refresh_token);
		return response.data;
	}

	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await this.client.post<AuthResponse>('/auth/register', data);
		this.setAccessToken(response.data.access_token);
		this.setRefreshToken(response.data.refresh_token);
		return response.data;
	}

	async refreshToken(): Promise<boolean> {
		try {
			const refreshToken = this.getRefreshToken();
			if (!refreshToken) return false;

			const response = await this.client.post<AuthResponse>('/auth/refresh', {
				refresh_token: refreshToken
			});

			this.setAccessToken(response.data.access_token);
			this.setRefreshToken(response.data.refresh_token);
			return true;
		} catch {
			this.clearTokens();
			return false;
		}
	}

	async logout(): Promise<void> {
		try {
			await this.client.post('/auth/logout');
		} finally {
			this.clearTokens();
		}
	}

	async getProfile() {
		const response = await this.client.get('/auth/me');
		return response.data;
	}

	// Registration endpoints
	async getRegistrations(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Registration>> {
		const response = await this.client.get<PaginatedResponse<Registration>>('/registrations', {
			params: { page, limit }
		});
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

	// Document endpoints
	async getDocuments(registrationId?: string): Promise<Document[]> {
		const params = registrationId ? { registration_id: registrationId } : {};
		const response = await this.client.get<Document[]>('/documents', { params });
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

	async getDocumentDownloadURL(id: string): Promise<{ download_url: string }> {
		const response = await this.client.get<{ download_url: string }>(`/documents/${id}/download`);
		return response.data;
	}

	async verifyDocument(id: string): Promise<Document> {
		const response = await this.client.post<Document>(`/documents/${id}/verify`);
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

	async createClient(data: Partial<Client>): Promise<Client> {
		const response = await this.client.post<Client>('/clients', data);
		return response.data;
	}

	// Health check
	async healthCheck(): Promise<HealthCheck> {
		const response = await this.client.get<HealthCheck>('/health');
		return response.data;
	}
}

// Export singleton instance
export const apiClient = new APIClient();
export default apiClient;
