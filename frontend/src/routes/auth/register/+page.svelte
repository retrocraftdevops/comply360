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
			<h2 class="text-2xl font-semibold text-gray-800">Create your account</h2>
			<p class="text-sm text-gray-500">
				Already have an account?
				<a href="/auth/login" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">
					Sign in here
				</a>
			</p>
		</div>

		<!-- Registration Form -->
		<div class="card card-hover p-8">
			<form class="space-y-5" on:submit|preventDefault={handleRegister}>
				{#if error}
					<div class="rounded-xl bg-red-50 border border-red-200 p-4 animate-fade-in">
						<p class="text-sm font-medium text-red-800">{error}</p>
					</div>
				{/if}

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="first-name" class="block text-sm font-semibold text-gray-700 mb-2">
							First Name <span class="text-red-500">*</span>
						</label>
						<input
							id="first-name"
							name="first-name"
							type="text"
							required
							bind:value={firstName}
							class="input"
							placeholder="John"
						/>
					</div>
					<div>
						<label for="last-name" class="block text-sm font-semibold text-gray-700 mb-2">
							Last Name <span class="text-red-500">*</span>
						</label>
						<input
							id="last-name"
							name="last-name"
							type="text"
							required
							bind:value={lastName}
							class="input"
							placeholder="Doe"
						/>
					</div>
				</div>

				<div>
					<label for="email" class="block text-sm font-semibold text-gray-700 mb-2">
						Email address <span class="text-red-500">*</span>
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
						Password <span class="text-red-500">*</span>
					</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={password}
						class="input"
						placeholder="••••••••"
					/>
					<p class="mt-2 text-xs text-gray-500">Must be at least 8 characters</p>
				</div>

				<div>
					<label for="confirm-password" class="block text-sm font-semibold text-gray-700 mb-2">
						Confirm Password <span class="text-red-500">*</span>
					</label>
					<input
						id="confirm-password"
						name="confirm-password"
						type="password"
						autocomplete="new-password"
						required
						bind:value={confirmPassword}
						class="input"
						placeholder="••••••••"
					/>
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
								Creating account...
							</span>
						{:else}
							Create Account
						{/if}
					</button>
				</div>

				<p class="text-center text-xs text-gray-500 pt-4 border-t border-gray-200">
					By creating an account, you agree to our
					<a href="/terms" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">Terms of Service</a>
					and
					<a href="/privacy" class="font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200">Privacy Policy</a>
				</p>
			</form>
		</div>
	</div>
</div>
