<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';

	let isLoadingUser = false;

	onMount(async () => {
		console.log('[Layout] onMount - checking authentication');
		// Try to load user if access token exists
		const accessToken = localStorage.getItem('access_token');
		console.log('[Layout] Access token exists:', !!accessToken);
		console.log('[Layout] Auth store authenticated:', $authStore.isAuthenticated);
		
		if (accessToken && !$authStore.isAuthenticated && !isLoadingUser) {
			console.log('[Layout] Loading user from token...');
			isLoadingUser = true;
			try {
				await authActions.loadUser();
				console.log('[Layout] User loaded successfully');
			} catch (error) {
				console.error('[Layout] Failed to load user:', error);
				authStore.update((state) => ({ ...state, isLoading: false }));
			} finally {
				isLoadingUser = false;
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
		!$page.url.pathname.startsWith('/auth') &&
		!isLoadingUser
	) {
		if (typeof window !== 'undefined') {
			// Check if we have a token in localStorage (might be logged in but store not updated yet)
			const token = localStorage.getItem('access_token');
			if (!token) {
				console.log('[Layout] No token found, redirecting to login');
				window.location.href = '/auth/login';
			} else {
				console.log('[Layout] Token found but store not authenticated, will wait for store to update');
				// Don't redirect - the store will update when login completes or onMount will load user
			}
		}
	}
</script>

{#if $authStore.isLoading}
	<div class="flex h-screen items-center justify-center bg-gradient-to-br from-gray-50 via-white to-primary-50/20">
		<div class="text-center space-y-4">
			<div class="flex justify-center">
				<div class="h-16 w-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shadow-lg shadow-primary-500/20 animate-pulse">
					<span class="text-white font-bold text-2xl">C</span>
				</div>
			</div>
			<div class="h-8 w-8 mx-auto animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
			<p class="text-sm font-medium text-gray-600">Loading Comply360...</p>
		</div>
	</div>
{:else}
	<slot />
{/if}
