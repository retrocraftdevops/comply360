/**
 * Auth Store - Reactive Wrapper around AuthManager
 *
 * This store provides a Svelte-reactive interface to the AuthManager singleton.
 * All actual auth state lives in AuthManager - this store just wraps it for reactivity.
 *
 * IMPORTANT: This store does NOT persist auth state directly.
 * Only the refresh token is persisted (in localStorage by AuthManager).
 *
 * @module AuthStore
 */

import { writable, derived } from 'svelte/store';
import { AuthManager } from '$lib/auth/AuthManager';
import type { User } from '$lib/types';

const authManager = AuthManager.getInstance();

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    permissions: string[];
    roles: string[];
    isLoading: boolean;
}

/**
 * Create auth store that wraps AuthManager
 */
function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>({
        user: authManager.getUser(),
        isAuthenticated: authManager.isAuthenticated(),
        permissions: authManager.getPermissions(),
        roles: authManager.getUser()?.roles || [],
        isLoading: false
    });

    return {
        subscribe,

        /**
         * Refresh store from AuthManager state
         * Call this after any AuthManager operation
         */
        refresh: () => {
            const user = authManager.getUser();
            set({
                user,
                isAuthenticated: authManager.isAuthenticated(),
                permissions: authManager.getPermissions(),
                roles: user?.roles || [],
                isLoading: false
            });
        },

        /**
         * Set loading state
         */
        setLoading: (loading: boolean) => {
            update(state => ({ ...state, isLoading: loading }));
        }
    };
}

export const authStore = createAuthStore();

/**
 * Derived store for checking specific roles
 */
export const hasRole = (role: string) => {
    return derived(authStore, ($auth) => $auth.roles.includes(role));
};

/**
 * Derived store for checking if user has any of the specified roles
 */
export const hasAnyRole = (roles: string[]) => {
    return derived(authStore, ($auth) => roles.some((role) => $auth.roles.includes(role)));
};

/**
 * Derived store for checking if user is admin
 */
export const isAdmin = derived(authStore, ($auth) =>
    $auth.roles.includes('tenant_admin') || $auth.roles.includes('system_admin')
);

/**
 * Derived store for checking if user is agent
 */
export const isAgent = derived(authStore, ($auth) =>
    $auth.roles.includes('agent') || $auth.roles.includes('agent_assistant')
);

/**
 * Auth Actions - All operations go through AuthManager
 */
export const authActions = {
    /**
     * Login with email and password
     */
    async login(email: string, password: string) {
        console.log('[Auth Store] Login called with:', { email, password: '***' });
        authStore.setLoading(true);

        try {
            const result = await authManager.login(email, password);
            authStore.refresh(); // Update reactive store
            return result;
        } catch (error: any) {
            console.error('[Auth Store] Login error:', error);
            authStore.setLoading(false);
            return {
                success: false,
                error: error.message || 'Login failed'
            };
        }
    },

    /**
     * Register new user
     * Note: This doesn't authenticate - user needs to verify email first
     */
    async register(email: string, password: string, firstName: string, lastName: string) {
        console.log('[Auth Store] Register called with:', { email, firstName, lastName });
        authStore.setLoading(true);

        try {
            // Import API client dynamically to avoid circular dependency
            const { apiClient } = await import('$lib/api/client');

            const response = await apiClient.register({
                email,
                password,
                first_name: firstName,
                last_name: lastName
            });

            authStore.setLoading(false);
            return { success: true, user: response.user };

        } catch (error: any) {
            console.error('[Auth Store] Register error:', error);
            authStore.setLoading(false);
            return {
                success: false,
                error: error.response?.data?.message || error.message || 'Registration failed'
            };
        }
    },

    /**
     * Logout user
     */
    async logout() {
        console.log('[Auth Store] Logout called');
        authStore.setLoading(true);

        try {
            await authManager.logout();
            authStore.refresh(); // Update reactive store (will show logged out state)
        } catch (error) {
            console.error('[Auth Store] Logout error:', error);
            // Even if logout fails, clear local state
            authStore.refresh();
        }
    },

    /**
     * Load current user profile
     */
    async loadUser() {
        console.log('[Auth Store] Loading user profile');
        authStore.setLoading(true);

        try {
            // Check if we have a refresh token to restore session
            const refreshToken = authManager.getRefreshToken();
            if (!refreshToken) {
                console.log('[Auth Store] No refresh token - cannot restore session');
                authStore.setLoading(false);
                return;
            }

            // Try to refresh access token
            const refreshed = await authManager.refreshAccessToken();
            if (!refreshed) {
                console.log('[Auth Store] Token refresh failed');
                authStore.setLoading(false);
                return;
            }

            // Load user profile
            const { apiClient } = await import('$lib/api/client');
            const response = await apiClient.getProfile();

            // Update AuthManager with user data
            authManager.setTokens(
                authManager.getAccessToken() || '',
                authManager.getRefreshToken() || ''
            );

            authStore.refresh();
            console.log('[Auth Store] User profile loaded');

        } catch (error) {
            console.error('[Auth Store] Load user failed:', error);
            await authManager.logout();
            authStore.refresh();
        }
    },

    /**
     * Clear auth state (for logout or error recovery)
     */
    clearAuth() {
        console.log('[Auth Store] Clearing auth state');
        authManager.clearAllState();
        authStore.refresh();
    },

    /**
     * Check if user has specific permission
     */
    hasPermission(permission: string): boolean {
        return authManager.hasPermission(permission);
    },

    /**
     * Refresh store state from AuthManager
     */
    refresh() {
        authStore.refresh();
    }
};
