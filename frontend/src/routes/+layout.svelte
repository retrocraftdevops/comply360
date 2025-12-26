<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	let isLoadingUser = false;
	let hasCheckedAuth = false;
	let lastRedirectPath = '';
	let currentPath = '';
	let showLoading = false;

	onMount(async () => {
		if (hasCheckedAuth) return;
		
		hasCheckedAuth = true;
		currentPath = $page.url.pathname;
		
		// On auth pages, set loading to false immediately to prevent flash
		if (currentPath.startsWith('/auth') || currentPath === '/') {
			authStore.update((state) => ({ ...state, isLoading: false }));
			return;
		}
		
		// Try to load user if access token exists
		const accessToken = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
		const auth = get(authStore);
		
		if (accessToken && !auth.isAuthenticated && !isLoadingUser) {
			isLoadingUser = true;
			try {
				await authActions.loadUser();
			} catch (error) {
				// Clear invalid token
				if (typeof window !== 'undefined') {
					localStorage.removeItem('access_token');
					localStorage.removeItem('refresh_token');
				}
				authStore.update((state) => ({ ...state, isLoading: false, isAuthenticated: false }));
			} finally {
				isLoadingUser = false;
			}
		} else {
			authStore.update((state) => ({ ...state, isLoading: false }));
		}
	});

	// Update current path on navigation and loading state
	// Only redirect if we're not already on auth page and haven't just redirected
	$: {
		const newPath = $page.url.pathname;
		const auth = get(authStore);
		const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
		
		currentPath = newPath;
		showLoading = (isLoadingUser || $authStore.isLoading) && !newPath.startsWith('/auth');
		
		// Only redirect if:
		// 1. We're not on an auth page
		// 2. We're not on the root page
		// 3. User is not authenticated
		// 4. No token exists
		// 5. We haven't already redirected to this path
		// 6. Auth check is complete
		if (
			typeof window !== 'undefined' &&
			hasCheckedAuth &&
			!isLoadingUser &&
			!$authStore.isLoading &&
			!newPath.startsWith('/auth') &&
			newPath !== '/' &&
			!auth.isAuthenticated &&
			!token &&
			lastRedirectPath !== newPath
		) {
			lastRedirectPath = newPath;
			// Use setTimeout to break the reactive cycle
			setTimeout(() => {
				goto('/auth/login', { replaceState: true });
			}, 0);
		}
		
		// Reset redirect path if we're on auth page
		if (newPath.startsWith('/auth')) {
			lastRedirectPath = '';
		}
	}
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
