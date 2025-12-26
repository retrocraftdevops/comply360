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
		try {
			authStore.update((state) => ({ ...state, isLoading: true }));

			const response: AuthResponse = await apiClient.login({ email, password });

			// Extract roles from user object
			const roles = response.user.roles || response.roles || [];

			authStore.set({
				user: response.user,
				roles,
				isAuthenticated: true,
				isLoading: false
			});

			return { success: true };
		} catch (error: any) {
			authStore.update((state) => ({ ...state, isLoading: false }));
			console.error('Login error:', error);
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
		try {
			authStore.update((state) => ({ ...state, isLoading: true }));

			const response = await apiClient.getProfile();

			authStore.set({
				user: response.user,
				roles: response.roles,
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
