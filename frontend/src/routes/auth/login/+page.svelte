<script lang="ts">
	import { authActions, authStore } from '$stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let email = '';
	let password = '';
	let isLoading = false;
	let error = '';

	// For debugging - check if we're already logged in
	onMount(() => {
		console.log('Login page mounted');
		console.log('[Login Page] Auth state:', {
			isAuthenticated: $authStore.isAuthenticated,
			hasToken: !!localStorage.getItem('access_token'),
			user: $authStore.user
		});

		// Only redirect if both token exists AND store is authenticated
		// This prevents redirect loops
		if ($authStore.isAuthenticated && $authStore.user) {
			console.log('[Login Page] Already authenticated with user, redirecting to dashboard');
			goto('/app/dashboard');
		} else if (localStorage.getItem('access_token') && !$authStore.isAuthenticated) {
			console.log('[Login Page] Token exists but store not updated, waiting...');
			// Don't redirect - let the store update first
		}
	});

	async function handleLogin() {
		console.log('[Login Page] handleLogin called');
		console.log('[Login Page] Form values:', { email, password: password ? '***' : '(empty)' });
		
		if (!email || !password) {
			console.error('[Login Page] Missing email or password');
			error = 'Please enter both email and password';
			return;
		}

		error = '';
		isLoading = true;

		console.log('[Login Page] Attempting login with:', { email, password: '***' });

		try {
			console.log('[Login Page] Calling authActions.login...');
			const result = await authActions.login(email, password);

			console.log('[Login Page] Login result:', result);

			if (result.success) {
				console.log('[Login Page] Login successful, redirecting to dashboard');
				await goto('/app/dashboard');
			} else {
				error = result.error || 'Login failed';
				console.error('[Login Page] Login failed:', error);
			}
		} catch (err) {
			console.error('[Login Page] Login exception:', err);
			error = 'An unexpected error occurred. Please try again.';
		} finally {
			isLoading = false;
			console.log('[Login Page] handleLogin completed, isLoading:', isLoading);
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-50 via-white to-primary-50/20 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8 animate-fade-in">
		<!-- Logo and Header -->
		<div class="text-center space-y-4">
			<div class="flex justify-center">
				<div class="h-16 w-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shadow-lg shadow-primary-500/20">
					<span class="text-white font-bold text-2xl">C</span>
				</div>
			</div>
			<h1 class="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
				Comply360
			</h1>
			<h2 class="text-2xl font-semibold text-gray-800">Sign in to your account</h2>
			<p class="text-sm text-gray-500">
				Don't have an account?
				<a href="/auth/register" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">
					Register here
				</a>
			</p>
		</div>

		<!-- Login Form -->
		<div class="card card-hover p-8">
			<form class="space-y-5" on:submit|preventDefault={handleLogin}>
				{#if error}
					<div class="rounded-xl bg-red-50 border border-red-200 p-4 animate-fade-in">
						<p class="text-sm font-medium text-red-800">{error}</p>
					</div>
				{/if}

				<div class="space-y-4">
					<div>
						<label for="email" class="block text-sm font-semibold text-gray-700 mb-2">
							Email address
						</label>
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							bind:value={email}
							class="input"
							placeholder="you@example.com"
						/>
					</div>
					<div>
						<label for="password" class="block text-sm font-semibold text-gray-700 mb-2">
							Password
						</label>
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="current-password"
							required
							bind:value={password}
							class="input"
							placeholder="••••••••"
						/>
					</div>
				</div>

				<div class="flex items-center justify-between">
					<div class="flex items-center">
						<input
							id="remember-me"
							name="remember-me"
							type="checkbox"
							class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
						/>
						<label for="remember-me" class="ml-2 block text-sm text-gray-700">Remember me</label>
					</div>

					<div class="text-sm">
						<a href="/auth/forgot-password" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">
							Forgot password?
						</a>
					</div>
				</div>

				<div class="pt-2">
					<button
						type="submit"
						disabled={isLoading}
						class="btn btn-primary w-full shadow-md hover:shadow-lg transition-all duration-200"
					>
						{#if isLoading}
							<span class="flex items-center justify-center">
								<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Signing in...
							</span>
						{:else}
							Sign in
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
</div>
