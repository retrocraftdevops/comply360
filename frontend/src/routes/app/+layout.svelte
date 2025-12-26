<script lang="ts">
	import { authStore, authActions } from '$stores/auth';
	import { page } from '$app/stores';

	const navigation = [
		{ name: 'Dashboard', href: '/app/dashboard', icon: 'home' },
		{ name: 'Registrations', href: '/app/registrations', icon: 'file-text' },
		{ name: 'Documents', href: '/app/documents', icon: 'upload' },
		{ name: 'Commissions', href: '/app/commissions', icon: 'dollar-sign' },
		{ name: 'Clients', href: '/app/clients', icon: 'users' }
	];

	async function handleLogout() {
		await authActions.logout();
		window.location.href = '/auth/login';
	}
</script>

<div class="min-h-screen bg-gray-100">
	<!-- Navigation -->
	<nav class="bg-white shadow-sm">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 justify-between">
				<div class="flex">
					<div class="flex flex-shrink-0 items-center">
						<h1 class="text-2xl font-bold text-primary">Comply360</h1>
					</div>
					<div class="hidden sm:ml-6 sm:flex sm:space-x-8">
						{#each navigation as item}
							<a
								href={item.href}
								class="inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium
									{$page.url.pathname.startsWith(item.href)
									? 'border-primary text-gray-900'
									: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'}"
							>
								{item.name}
							</a>
						{/each}
					</div>
				</div>
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<span class="text-sm text-gray-700">
							{$authStore.user?.first_name || ''} {$authStore.user?.last_name || ''}
						</span>
					</div>
					<button
						type="button"
						on:click={handleLogout}
						class="ml-4 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
					>
						Logout
					</button>
				</div>
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	<main class="py-10">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<slot />
		</div>
	</main>
</div>
