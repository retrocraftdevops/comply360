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

	// Update current path and loading state
	$: {
		currentPath = $page.url.pathname;
		showLoading = (isLoadingUser || $authStore.isLoading) && !currentPath.startsWith('/auth');
	}

	// Separate reactive statement for redirect - only runs when path changes
	$: if (
		typeof window !== 'undefined' &&
		hasCheckedAuth &&
		!isLoadingUser &&
		!$authStore.isLoading &&
		currentPath &&
		!currentPath.startsWith('/auth') &&
		currentPath !== '/' &&
		!get(authStore).isAuthenticated &&
		!localStorage.getItem('access_token') &&
		lastRedirectPath !== currentPath
	) {
		lastRedirectPath = currentPath;
		// Use requestAnimationFrame to break reactive cycle and prevent loops
		requestAnimationFrame(() => {
			if (lastRedirectPath === currentPath) {
				goto('/auth/login', { replaceState: true });
			}
		});
	}

	// Reset redirect path when on auth page
	$: if (currentPath.startsWith('/auth')) {
		lastRedirectPath = '';
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
