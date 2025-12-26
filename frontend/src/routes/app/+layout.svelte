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

<div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-gray-50/50">
	<!-- Navigation -->
	<nav class="sticky top-0 z-50 glass border-b border-gray-200/50 backdrop-blur-xl">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center space-x-8">
					<!-- Logo -->
					<a href="/app/dashboard" class="flex items-center space-x-2 group">
						<div class="h-8 w-8 rounded-lg bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shadow-md group-hover:shadow-lg transition-shadow duration-300">
							<span class="text-white font-bold text-lg">C</span>
						</div>
						<h1 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
							Comply360
						</h1>
					</a>
					
					<!-- Navigation Links -->
					<div class="hidden md:flex items-center space-x-1">
						{#each navigation as item}
							<a
								href={item.href}
								class="relative px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200
									{$page.url.pathname.startsWith(item.href)
									? 'text-primary-600 bg-primary-50'
									: 'text-gray-600 hover:text-gray-900 hover:bg-gray-100/50'}"
							>
								{item.name}
								{#if $page.url.pathname.startsWith(item.href)}
									<span class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-500 rounded-full"></span>
								{/if}
							</a>
						{/each}
					</div>
				</div>
				
				<!-- User Menu -->
				<div class="flex items-center space-x-4">
					<div class="hidden sm:flex items-center space-x-3 px-4 py-2 rounded-lg bg-gray-50/50">
						<div class="h-8 w-8 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-sm font-semibold shadow-sm">
							{$authStore.user?.first_name?.[0] || 'U'}
						</div>
						<span class="text-sm font-medium text-gray-700">
							{$authStore.user?.first_name || ''} {$authStore.user?.last_name || ''}
						</span>
					</div>
					<button
						type="button"
						on:click={handleLogout}
						class="btn btn-ghost text-sm"
					>
						Logout
					</button>
				</div>
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	<main class="py-8 animate-fade-in">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<slot />
		</div>
	</main>
</div>
