<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';

	onMount(async () => {
		// Try to load user if access token exists
		const accessToken = localStorage.getItem('access_token');
		if (accessToken && !$authStore.isAuthenticated) {
			await authActions.loadUser();
		} else {
			authStore.update((state) => ({ ...state, isLoading: false }));
		}
	});

	// Redirect to login if not authenticated and not on auth pages
	$: if (
		!$authStore.isLoading &&
		!$authStore.isAuthenticated &&
		!$page.url.pathname.startsWith('/auth')
	) {
		if (typeof window !== 'undefined') {
			window.location.href = '/auth/login';
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
