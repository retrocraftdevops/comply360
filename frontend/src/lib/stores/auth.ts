import { writable, derived } from 'svelte/store';
import { persisted } from 'svelte-persisted-store';
import type { User, AuthResponse } from '$types';
import { apiClient } from '$api/client';

interface AuthState {
	user: User | null;
	roles: string[];
	isAuthenticated: boolean;
	isLoading: boolean;
}

const initialState: AuthState = {
	user: null,
	roles: [],
	isAuthenticated: false,
	isLoading: true
};

// Create persisted store for user data
export const authStore = persisted<AuthState>('auth', initialState);

// Derived store for checking specific roles
export const hasRole = (role: string) => {
	return derived(authStore, ($auth) => $auth.roles.includes(role));
};

// Derived store for checking if user has any of the specified roles
export const hasAnyRole = (roles: string[]) => {
	return derived(authStore, ($auth) => roles.some((role) => $auth.roles.includes(role)));
};

// Derived store for checking if user is admin
export const isAdmin = derived(authStore, ($auth) =>
	$auth.roles.includes('tenant_admin') || $auth.roles.includes('system_admin')
);

// Derived store for checking if user is agent
export const isAgent = derived(authStore, ($auth) =>
	$auth.roles.includes('agent') || $auth.roles.includes('agent_assistant')
);

// Auth actions
export const authActions = {
	async login(email: string, password: string) {
		console.log('[Auth Store] Login called with:', { email, password: '***' });
		try {
			authStore.update((state) => ({ ...state, isLoading: true }));
			console.log('[Auth Store] Calling apiClient.login...');

			const response: AuthResponse = await apiClient.login({ email, password });
			console.log('[Auth Store] Received response:', response);

			// Extract roles from user object
			const roles = response.user?.roles || response.roles || [];
			console.log('[Auth Store] Extracted roles:', roles);

			authStore.set({
				user: response.user,
				roles,
				isAuthenticated: true,
				isLoading: false
			});
			console.log('[Auth Store] Auth state updated successfully');

			return { success: true };
		} catch (error: any) {
			authStore.update((state) => ({ ...state, isLoading: false }));
			console.error('[Auth Store] Login error:', error);
			console.error('[Auth Store] Error details:', {
				message: error.message,
				response: error.response?.data,
				status: error.response?.status
			});
			return {
				success: false,
				error: error.response?.data?.message || error.message || 'Login failed'
			};
		}
	},

	async register(email: string, password: string, firstName: string, lastName: string) {
		try {
			authStore.update((state) => ({ ...state, isLoading: true }));

			const response: AuthResponse = await apiClient.register({
				email,
				password,
				first_name: firstName,
				last_name: lastName
			});

			authStore.set({
				user: response.user,
				roles: response.roles,
				isAuthenticated: true,
				isLoading: false
			});

			return { success: true };
		} catch (error: any) {
			authStore.update((state) => ({ ...state, isLoading: false }));
			return {
				success: false,
				error: error.response?.data?.message || 'Registration failed'
			};
		}
	},

	async logout() {
		try {
			await apiClient.logout();
		} finally {
			authStore.set(initialState);
		}
	},

	async loadUser() {
		// Prevent multiple simultaneous load attempts
		if (authStore.get().isLoading) {
			return;
		}
		
		try {
			authStore.update((state) => ({ ...state, isLoading: true }));

			const response = await apiClient.getProfile();

			authStore.set({
				user: response.user,
				roles: response.roles || [],
				isAuthenticated: true,
				isLoading: false
			});
		} catch (error) {
			authStore.set({ ...initialState, isLoading: false });
		}
	},

	clearAuth() {
		authStore.set(initialState);
	}
};
