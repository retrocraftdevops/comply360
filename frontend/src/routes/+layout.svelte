<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	let isLoadingUser = false;
	let hasCheckedAuth = false;

	onMount(async () => {
		if (hasCheckedAuth) return;
		
		console.log('[Layout] onMount - checking authentication');
		hasCheckedAuth = true;
		
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
				// Clear invalid token
				localStorage.removeItem('access_token');
				localStorage.removeItem('refresh_token');
				authStore.update((state) => ({ ...state, isLoading: false, isAuthenticated: false }));
			} finally {
				isLoadingUser = false;
			}
		} else {
			console.log('[Layout] No token or already authenticated, setting isLoading to false');
			authStore.update((state) => ({ ...state, isLoading: false }));
		}
	});

	// Redirect to login if not authenticated and not on auth pages
	// Use a more stable reactive statement that doesn't cause loops
	$: if (
		typeof window !== 'undefined' &&
		hasCheckedAuth &&
		!isLoadingUser &&
		!$authStore.isLoading &&
		!$page.url.pathname.startsWith('/auth') &&
		!$authStore.isAuthenticated &&
		!localStorage.getItem('access_token')
	) {
		console.log('[Layout] No auth, redirecting to login');
		goto('/auth/login', { replaceState: true });
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
