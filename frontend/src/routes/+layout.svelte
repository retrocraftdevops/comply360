<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';

	onMount(async () => {
		console.log('[Layout] onMount - checking authentication');
		// Try to load user if access token exists
		const accessToken = localStorage.getItem('access_token');
		console.log('[Layout] Access token exists:', !!accessToken);
		console.log('[Layout] Auth store authenticated:', $authStore.isAuthenticated);
		
		if (accessToken && !$authStore.isAuthenticated) {
			console.log('[Layout] Loading user from token...');
			try {
				await authActions.loadUser();
				console.log('[Layout] User loaded successfully');
			} catch (error) {
				console.error('[Layout] Failed to load user:', error);
				authStore.update((state) => ({ ...state, isLoading: false }));
			}
		} else {
			console.log('[Layout] No token or already authenticated, setting isLoading to false');
			authStore.update((state) => ({ ...state, isLoading: false }));
		}
	});

	// Redirect to login if not authenticated and not on auth pages
	// Only redirect if we've finished loading and there's no token
	$: if (
		!$authStore.isLoading &&
		!$authStore.isAuthenticated &&
		!$page.url.pathname.startsWith('/auth')
	) {
		if (typeof window !== 'undefined') {
			// Check if we have a token in localStorage (might be logged in but store not updated yet)
			const token = localStorage.getItem('access_token');
			if (!token) {
				console.log('[Layout] No token found, redirecting to login');
				window.location.href = '/auth/login';
			} else {
				console.log('[Layout] Token found but store not authenticated, waiting for store to update...');
				// Don't redirect - let the store update first
				// The store will update when the login completes
			}
		}
	}
</script>

{#if $authStore.isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-center">
			<div class="mb-4 h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading...</p>
		</div>
	</div>
{:else}
	<slot />
{/if}
