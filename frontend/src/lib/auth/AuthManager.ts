/**
 * AuthManager - Single Source of Truth for Authentication State
 *
 * This singleton class manages ALL authentication state in the application.
 * It provides atomic operations for login, logout, token management, and session handling.
 *
 * Key Principles:
 * - Singleton pattern ensures only ONE instance exists
 * - In-memory token storage prevents XSS attacks
 * - Atomic operations prevent state desynchronization
 * - All auth operations go through this manager
 *
 * @module AuthManager
 */

import type { User, AuthResponse, LoginRequest } from '$lib/types';
import { jwtDecode } from 'jwt-decode';

interface TokenClaims {
    user_id: string;
    tenant_id: string;
    email: string;
    roles: string[];
    exp: number;
    iat: number;
}

export class AuthManager {
    private static instance: AuthManager;

    // Core auth state (in-memory)
    private user: User | null = null;
    private accessToken: string | null = null;
    private refreshToken: string | null = null;
    private permissions: string[] = [];
    private tenantId: string | null = null;

    // Session management
    private sessionTimeoutId: number | null = null;
    private readonly SESSION_TIMEOUT: number;

    // API client reference (set externally to avoid circular dependency)
    private apiClient: any = null;

    /**
     * Private constructor enforces singleton pattern
     */
    private constructor() {
        // Get session timeout from env (in minutes) and convert to milliseconds
        const timeoutMinutes = Number(import.meta.env.VITE_SESSION_TIMEOUT) || 60;
        this.SESSION_TIMEOUT = timeoutMinutes * 60 * 1000;

        console.log('[AuthManager] Initialized with session timeout:', timeoutMinutes, 'minutes');

        // Try to restore session from localStorage on initialization
        this.restoreSession();
    }

    /**
     * Get singleton instance
     */
    public static getInstance(): AuthManager {
        if (!AuthManager.instance) {
            AuthManager.instance = new AuthManager();
        }
        return AuthManager.instance;
    }

    /**
     * Set API client reference (called from API client initialization)
     */
    public setApiClient(client: any): void {
        this.apiClient = client;
    }

    /**
     * Restore session from localStorage (if refresh token exists)
     */
    private restoreSession(): void {
        if (typeof window === 'undefined') return;

        try {
            const storedRefreshToken = localStorage.getItem('refresh_token');
            if (storedRefreshToken) {
                this.refreshToken = storedRefreshToken;
                console.log('[AuthManager] Session restoration: refresh token found');
            }
        } catch (error) {
            console.error('[AuthManager] Session restoration failed:', error);
        }
    }

    /**
     * Login user with email and password
     * This is an ATOMIC operation - either all state is set or none
     */
    public async login(email: string, password: string): Promise<{success: boolean, error?: string}> {
        console.log('[AuthManager] Login attempt for:', email);

        try {
            // Validate API client is set
            if (!this.apiClient) {
                throw new Error('API client not initialized');
            }

            // Call login API
            const response: AuthResponse = await this.apiClient.login({ email, password });

            // Validate response
            if (!response.user || !response.access_token) {
                throw new Error('Invalid login response - missing user or token');
            }

            // ATOMIC STATE UPDATE - all or nothing
            this.user = response.user;
            this.accessToken = response.access_token;
            this.refreshToken = response.refresh_token;
            this.tenantId = response.user.tenant_id;

            // Store refresh token in localStorage for session persistence
            if (typeof window !== 'undefined') {
                localStorage.setItem('refresh_token', response.refresh_token);
            }

            // IMMEDIATELY load permissions (blocking operation)
            await this.loadPermissions(response.user.id);

            // Start session timeout timer
            this.startSessionTimer();

            console.log('[AuthManager] Login successful:', {
                userId: this.user.id,
                email: this.user.email,
                tenantId: this.tenantId,
                permissions: this.permissions.length,
                roles: response.roles?.length || 0
            });

            return { success: true };

        } catch (error: any) {
            // On ANY error, clear all state (atomic failure)
            console.error('[AuthManager] Login failed:', error);
            this.clearAllState();

            return {
                success: false,
                error: error.response?.data?.message || error.message || 'Login failed'
            };
        }
    }

    /**
     * Logout user and clear ALL auth state
     * This performs comprehensive cleanup
     */
    public async logout(): Promise<void> {
        console.log('[AuthManager] Starting logout cleanup');

        // 1. Call backend logout endpoint (best effort)
        try {
            if (this.apiClient && this.accessToken) {
                await this.apiClient.logout();
                console.log('[AuthManager] Backend logout successful');
            }
        } catch (error) {
            console.error('[AuthManager] Backend logout failed (continuing cleanup):', error);
        }

        // 2. Clear session timer
        this.clearSessionTimer();

        // 3. Clear ALL auth state
        this.clearAllState();

        // 4. Clear localStorage (selective - keep user preferences)
        if (typeof window !== 'undefined') {
            const keysToKeep = ['theme', 'language', 'currency'];
            const allKeys = Object.keys(localStorage);

            allKeys.forEach(key => {
                if (!keysToKeep.includes(key)) {
                    localStorage.removeItem(key);
                }
            });

            console.log('[AuthManager] localStorage cleared (except preferences)');
        }

        console.log('[AuthManager] Logout cleanup complete');
    }

    /**
     * Refresh access token using refresh token
     */
    public async refreshAccessToken(): Promise<boolean> {
        console.log('[AuthManager] Attempting token refresh');

        try {
            if (!this.refreshToken) {
                console.warn('[AuthManager] No refresh token available');
                return false;
            }

            if (!this.apiClient) {
                console.error('[AuthManager] API client not initialized');
                return false;
            }

            // Call refresh endpoint
            const response: AuthResponse = await this.apiClient.post('/auth/refresh', {
                refresh_token: this.refreshToken
            });

            // Update tokens
            this.accessToken = response.data.access_token;
            this.refreshToken = response.data.refresh_token;

            // Update refresh token in localStorage
            if (typeof window !== 'undefined') {
                localStorage.setItem('refresh_token', response.data.refresh_token);
            }

            // Reset session timer
            this.resetSessionTimer();

            console.log('[AuthManager] Token refresh successful');
            return true;

        } catch (error) {
            console.error('[AuthManager] Token refresh failed:', error);
            // Clear state on refresh failure
            await this.logout();
            return false;
        }
    }

    /**
     * Load user permissions from backend
     */
    public async loadPermissions(userId: string): Promise<void> {
        console.log('[AuthManager] Loading permissions for user:', userId);

        try {
            if (!this.apiClient) {
                console.warn('[AuthManager] Cannot load permissions - API client not set');
                return;
            }

            // Call permissions endpoint
            const response = await this.apiClient.get(`/admin/users/${userId}/effective-permissions`);
            this.permissions = response.data.permissions || [];

            console.log('[AuthManager] Loaded permissions:', this.permissions.length);

        } catch (error) {
            console.error('[AuthManager] Failed to load permissions:', error);
            // Don't fail login if permissions can't be loaded - set empty array
            this.permissions = [];
        }
    }

    /**
     * Check if user has specific permission
     */
    public hasPermission(permission: string): boolean {
        return this.permissions.includes(permission);
    }

    /**
     * Check if token is expired
     */
    public isTokenExpired(token: string): boolean {
        try {
            const decoded = jwtDecode<TokenClaims>(token);
            const now = Date.now() / 1000;
            // Consider token expired 30 seconds before actual expiry
            return decoded.exp < (now + 30);
        } catch (error) {
            console.error('[AuthManager] Token decode failed:', error);
            return true;
        }
    }

    /**
     * Start session timeout timer
     */
    public startSessionTimer(): void {
        this.clearSessionTimer();

        this.sessionTimeoutId = window.setTimeout(() => {
            console.log('[AuthManager] Session timeout - logging out');
            this.logout();

            // Redirect to login if in browser
            if (typeof window !== 'undefined') {
                window.location.href = '/auth/login?reason=session_timeout';
            }
        }, this.SESSION_TIMEOUT);

        console.log('[AuthManager] Session timer started:', this.SESSION_TIMEOUT / 1000, 'seconds');
    }

    /**
     * Reset session timeout timer (call on user activity)
     */
    public resetSessionTimer(): void {
        if (this.isAuthenticated()) {
            this.startSessionTimer();
        }
    }

    /**
     * Clear session timeout timer
     */
    public clearSessionTimer(): void {
        if (this.sessionTimeoutId) {
            clearTimeout(this.sessionTimeoutId);
            this.sessionTimeoutId = null;
        }
    }

    /**
     * Clear all authentication state
     */
    public clearAllState(): void {
        this.user = null;
        this.accessToken = null;
        this.refreshToken = null;
        this.permissions = [];
        this.tenantId = null;
        this.clearSessionTimer();

        console.log('[AuthManager] All state cleared');
    }

    // ==================== GETTERS ====================

    /**
     * Get current user
     */
    public getUser(): User | null {
        return this.user;
    }

    /**
     * Get access token (checks expiry)
     */
    public getAccessToken(): string | null {
        if (!this.accessToken) return null;

        // Return null if token is expired (triggers refresh)
        if (this.isTokenExpired(this.accessToken)) {
            console.warn('[AuthManager] Access token expired');
            return null;
        }

        return this.accessToken;
    }

    /**
     * Get refresh token
     */
    public getRefreshToken(): string | null {
        return this.refreshToken;
    }

    /**
     * Get user permissions
     */
    public getPermissions(): string[] {
        return [...this.permissions]; // Return copy to prevent external mutation
    }

    /**
     * Get tenant ID
     */
    public getTenantId(): string | null {
        return this.tenantId;
    }

    /**
     * Check if user is authenticated
     */
    public isAuthenticated(): boolean {
        return this.user !== null && this.accessToken !== null;
    }

    /**
     * Set tokens directly (used by API client after login)
     */
    public setTokens(access: string, refresh: string): void {
        this.accessToken = access;
        this.refreshToken = refresh;

        // Store refresh token in localStorage
        if (typeof window !== 'undefined') {
            localStorage.setItem('refresh_token', refresh);
        }
    }

    /**
     * Clear tokens
     */
    public clearTokens(): void {
        this.accessToken = null;
        this.refreshToken = null;

        if (typeof window !== 'undefined') {
            localStorage.removeItem('refresh_token');
        }
    }
}

// Export singleton instance
export const authManager = AuthManager.getInstance();
export default authManager;
