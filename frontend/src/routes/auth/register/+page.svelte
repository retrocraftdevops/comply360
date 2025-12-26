<script lang="ts">
	import { authActions } from '$stores/auth';
	import { goto } from '$app/navigation';

	let email = '';
	let password = '';
	let confirmPassword = '';
	let firstName = '';
	let lastName = '';
	let isLoading = false;
	let error = '';

	async function handleRegister() {
		error = '';

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (password.length < 8) {
			error = 'Password must be at least 8 characters';
			return;
		}

		isLoading = true;

		const result = await authActions.register(email, password, firstName, lastName);

		if (result.success) {
			goto('/app/dashboard');
		} else {
			error = result.error || 'Registration failed';
		}

		isLoading = false;
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h1 class="text-4xl font-bold text-primary">Comply360</h1>
			<h2 class="mt-6 text-3xl font-bold tracking-tight text-gray-900">Create your account</h2>
			<p class="mt-2 text-sm text-gray-600">
				Already have an account?
				<a href="/auth/login" class="font-medium text-primary hover:text-primary/80">
					Sign in here
				</a>
			</p>
		</div>

		<form class="mt-8 space-y-6" on:submit|preventDefault={handleRegister}>
			{#if error}
				<div class="rounded-md bg-red-50 p-4">
					<p class="text-sm text-red-800">{error}</p>
				</div>
			{/if}

			<div class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="first-name" class="block text-sm font-medium text-gray-700">
							First Name
						</label>
						<input
							id="first-name"
							name="first-name"
							type="text"
							required
							bind:value={firstName}
							class="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
						/>
					</div>
					<div>
						<label for="last-name" class="block text-sm font-medium text-gray-700">
							Last Name
						</label>
						<input
							id="last-name"
							name="last-name"
							type="text"
							required
							bind:value={lastName}
							class="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="email" class="block text-sm font-medium text-gray-700">Email address</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={email}
						class="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
					/>
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={password}
						class="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
					/>
					<p class="mt-1 text-xs text-gray-500">Must be at least 8 characters</p>
				</div>

				<div>
					<label for="confirm-password" class="block text-sm font-medium text-gray-700">
						Confirm Password
					</label>
					<input
						id="confirm-password"
						name="confirm-password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={confirmPassword}
						class="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary focus:outline-none focus:ring-primary sm:text-sm"
					/>
				</div>
			</div>

			<div>
				<button
					type="submit"
					disabled={isLoading}
					class="group relative flex w-full justify-center rounded-md border border-transparent bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 disabled:opacity-50"
				>
					{#if isLoading}
						<span class="mr-2">Creating account...</span>
						<div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
					{:else}
						Create Account
					{/if}
				</button>
			</div>

			<p class="text-center text-xs text-gray-500">
				By creating an account, you agree to our
				<a href="/terms" class="text-primary hover:text-primary/80">Terms of Service</a>
				and
				<a href="/privacy" class="text-primary hover:text-primary/80">Privacy Policy</a>
			</p>
		</form>
	</div>
</div>
