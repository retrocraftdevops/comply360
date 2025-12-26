<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { afterNavigate } from '$app/navigation';

	let isLoadingUser = false;
	let showLoading = false;

	// Check auth and redirect if needed - NO reactive statements
	async function checkAuth() {
		console.log('[Layout v2.0] checkAuth called - path:', $page.url.pathname);
		const path = $page.url.pathname;
		
		// Allow auth pages and root
		if (path.startsWith('/auth') || path === '/') {
			console.log('[Layout v2.0] Auth page, skipping check');
			authStore.update((state) => ({ ...state, isLoading: false }));
			return;
		}
		
		// Check if authenticated
		const auth = get(authStore);
		const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
		console.log('[Layout v2.0] Auth check:', { authenticated: auth.isAuthenticated, hasToken: !!token });
		
		// If not authenticated and no token, redirect to login
		if (!auth.isAuthenticated && !token) {
			console.log('[Layout v2.0] Not authenticated, redirecting to login');
			authStore.update((state) => ({ ...state, isLoading: false }));
			goto('/auth/login', { replaceState: true });
			return;
		}
		
		// If token exists but not authenticated, try to load user
		if (token && !auth.isAuthenticated && !isLoadingUser) {
			console.log('[Layout v2.0] Token exists, loading user...');
			isLoadingUser = true;
			showLoading = true;
			try {
				await authActions.loadUser();
				console.log('[Layout v2.0] User loaded successfully');
			} catch (error) {
				console.error('[Layout v2.0] Failed to load user:', error);
				// Clear invalid token
				if (typeof window !== 'undefined') {
					localStorage.removeItem('access_token');
					localStorage.removeItem('refresh_token');
				}
				authStore.update((state) => ({ ...state, isLoading: false, isAuthenticated: false }));
				goto('/auth/login', { replaceState: true });
			} finally {
				isLoadingUser = false;
				showLoading = false;
			}
		} else {
			console.log('[Layout v2.0] Already authenticated or loading');
			authStore.update((state) => ({ ...state, isLoading: false }));
			showLoading = false;
		}
	}

	onMount(() => {
		console.log('[Layout v2.0] onMount - NEW VERSION LOADED');
		checkAuth();
	});

	// Check auth after navigation - this is stable and won't cause loops
	afterNavigate(() => {
		console.log('[Layout v2.0] afterNavigate - checking auth');
		checkAuth();
	});

	// Only update loading state reactively - no redirects here
	$: showLoading = (isLoadingUser || $authStore.isLoading) && !$page.url.pathname.startsWith('/auth');
</script>

{#if showLoading}
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
