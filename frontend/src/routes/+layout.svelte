<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { beforeNavigate } from '$app/navigation';

	let isLoadingUser = false;
	let hasCheckedAuth = false;
	let redirecting = false;
	let currentPath = '';
	let showLoading = false;

	// Use beforeNavigate to check auth before navigation
	beforeNavigate(({ to, cancel }) => {
		if (!to) return;
		
		const auth = get(authStore);
		const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
		
		// Allow navigation to auth pages
		if (to.url.pathname.startsWith('/auth') || to.url.pathname === '/') {
			return;
		}
		
		// Redirect to login if not authenticated
		if (!auth.isAuthenticated && !token && !redirecting) {
			redirecting = true;
			cancel();
			goto('/auth/login', { replaceState: true });
		}
	});

	onMount(async () => {
		if (hasCheckedAuth) return;
		
		hasCheckedAuth = true;
		currentPath = $page.url.pathname;
		
		// On auth pages, set loading to false immediately to prevent flash
		if (currentPath.startsWith('/auth')) {
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
	$: {
		currentPath = $page.url.pathname;
		showLoading = (isLoadingUser || $authStore.isLoading) && !currentPath.startsWith('/auth');
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
