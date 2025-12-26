<script lang="ts">
	import { authActions } from '$stores/auth';
	import { goto } from '$app/navigation';

	let email = '';
	let password = '';
	let isLoading = false;
	let error = '';

	async function handleLogin() {
		error = '';
		isLoading = true;

		const result = await authActions.login(email, password);

		if (result.success) {
			goto('/app/dashboard');
		} else {
			error = result.error || 'Login failed';
		}

		isLoading = false;
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h1 class="text-4xl font-bold text-primary">Comply360</h1>
			<h2 class="mt-6 text-3xl font-bold tracking-tight text-gray-900">Sign in to your account</h2>
			<p class="mt-2 text-sm text-gray-600">
				Don't have an account?
				<a href="/auth/register" class="font-medium text-primary hover:text-primary/80">
					Register here
				</a>
			</p>
		</div>

		<form class="mt-8 space-y-6" on:submit|preventDefault={handleLogin}>
			{#if error}
				<div class="rounded-md bg-red-50 p-4">
					<p class="text-sm text-red-800">{error}</p>
				</div>
			{/if}

			<div class="space-y-4 rounded-md shadow-sm">
				<div>
					<label for="email" class="sr-only">Email address</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={email}
						class="relative block w-full appearance-none rounded-t-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
						placeholder="Email address"
					/>
				</div>
				<div>
					<label for="password" class="sr-only">Password</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="current-password"
						required
						bind:value={password}
						class="relative block w-full appearance-none rounded-b-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
						placeholder="Password"
					/>
				</div>
			</div>

			<div class="flex items-center justify-between">
				<div class="flex items-center">
					<input
						id="remember-me"
						name="remember-me"
						type="checkbox"
						class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
					/>
					<label for="remember-me" class="ml-2 block text-sm text-gray-900">Remember me</label>
				</div>

				<div class="text-sm">
					<a href="/auth/forgot-password" class="font-medium text-primary hover:text-primary/80">
						Forgot password?
					</a>
				</div>
			</div>

			<div>
				<button
					type="submit"
					disabled={isLoading}
					class="group relative flex w-full justify-center rounded-md border border-transparent bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 disabled:opacity-50"
				>
					{#if isLoading}
						<span class="mr-2">Signing in...</span>
						<div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
					{:else}
						Sign in
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
