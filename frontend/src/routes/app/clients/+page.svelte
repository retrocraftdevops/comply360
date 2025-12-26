<script lang="ts">
	import { onMount } from 'svelte';
	import apiClient from '$api/client';
	import type { Client } from '$types';

	let isLoading = true;
	let clients: Client[] = [];
	let total = 0;
	let currentPage = 1;
	let limit = 20;
	let showCreateModal = false;
	let isCreating = false;
	let createError = '';

	// Form data
	let clientType: 'individual' | 'company' = 'individual';
	let fullName = '';
	let companyName = '';
	let email = '';
	let phone = '';
	let mobile = '';

	onMount(async () => {
		await loadClients();
	});

	async function loadClients() {
		isLoading = true;
		try {
			const response = await apiClient.getClients(currentPage, limit);
			clients = response.data;
			total = response.total;
		} catch (error) {
			console.error('Failed to load clients:', error);
		} finally {
			isLoading = false;
		}
	}

	async function handleCreate() {
		createError = '';

		if (!email) {
			createError = 'Email is required';
			return;
		}

		if (clientType === 'individual' && !fullName) {
			createError = 'Full name is required for individuals';
			return;
		}

		if (clientType === 'company' && !companyName) {
			createError = 'Company name is required for companies';
			return;
		}

		isCreating = true;

		try {
			await apiClient.createClient({
				client_type: clientType,
				full_name: clientType === 'individual' ? fullName : undefined,
				company_name: clientType === 'company' ? companyName : undefined,
				email,
				phone,
				mobile,
				status: 'active'
			});

			showCreateModal = false;
			resetForm();
			await loadClients();
		} catch (error: any) {
			createError = error.response?.data?.message || 'Failed to create client';
		} finally {
			isCreating = false;
		}
	}

	function resetForm() {
		clientType = 'individual';
		fullName = '';
		companyName = '';
		email = '';
		phone = '';
		mobile = '';
		createError = '';
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'active':
				return 'bg-green-100 text-green-800';
			case 'suspended':
				return 'bg-red-100 text-red-800';
			case 'inactive':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-blue-100 text-blue-800';
		}
	}
</script>

<div class="space-y-6 animate-fade-in">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="space-y-1">
			<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
				Clients
			</h1>
			<p class="text-sm text-gray-500 font-medium">
				Manage client information and contacts
			</p>
		</div>
		<button
			on:click={() => showCreateModal = true}
			class="btn btn-primary shadow-md hover:shadow-lg transition-all duration-200 group"
		>
			<svg class="w-4 h-4 mr-2 group-hover:rotate-90 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			Add Client
		</button>
	</div>

	{#if isLoading}
		<div class="card flex items-center justify-center py-16">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
				<p class="text-sm text-gray-500 font-medium">Loading clients...</p>
			</div>
		</div>
	{:else}
		<div class="card card-hover overflow-hidden">
			<div class="overflow-x-auto">
				<table class="table">
					<thead class="table-header">
						<tr>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Name
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Type
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Email
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Phone
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Status
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Created
							</th>
							<th class="px-6 py-4 text-right text-xs font-semibold uppercase tracking-wider text-gray-600">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-100">
						{#if clients.length === 0}
							<tr>
								<td colspan="7" class="px-6 py-16 text-center">
									<div class="flex flex-col items-center space-y-4">
										<div class="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
											<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
											</svg>
										</div>
										<div>
											<p class="text-sm font-medium text-gray-900">No clients found</p>
											<p class="text-sm text-gray-500 mt-1">Get started by adding your first client</p>
										</div>
										<button
											on:click={() => showCreateModal = true}
											class="btn btn-primary mt-2"
										>
											Add your first client
										</button>
									</div>
								</td>
							</tr>
						{:else}
							{#each clients as client, index}
								<tr class="table-row animate-fade-in" style="animation-delay: {index * 50}ms">
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="flex items-center space-x-3">
											<div class="h-10 w-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-sm font-semibold shadow-sm">
												{(client.client_type === 'individual' ? client.full_name : client.company_name)?.[0] || 'C'}
											</div>
											<div class="text-sm font-semibold text-gray-900">
												{client.client_type === 'individual' ? client.full_name : client.company_name}
											</div>
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold capitalize shadow-sm
											{client.client_type === 'company' 
												? 'bg-blue-50 text-blue-700 border border-blue-200' 
												: 'bg-purple-50 text-purple-700 border border-purple-200'}">
											{client.client_type}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{client.email}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{client.mobile || client.phone || '-'}
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span
											class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold shadow-sm {getStatusColor(
												client.status
											)}"
										>
											{client.status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{formatDate(client.created_at)}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
										<a
											href="/app/clients/{client.id}"
											class="text-primary-600 hover:text-primary-700 font-medium transition-colors duration-200"
										>
											View
										</a>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>

			{#if total > limit}
				<div class="border-t border-gray-200 px-6 py-3">
					<div class="flex items-center justify-between">
						<p class="text-sm text-gray-700">
							Showing <span class="font-medium">{(currentPage - 1) * limit + 1}</span> to
							<span class="font-medium">{Math.min(currentPage * limit, total)}</span> of
							<span class="font-medium">{total}</span> clients
						</p>
						<div class="flex space-x-2">
							<button
								disabled={currentPage === 1}
								on:click={() => {
									currentPage--;
									loadClients();
								}}
								class="rounded-md border border-gray-300 bg-white px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
							>
								Previous
							</button>
							<button
								disabled={currentPage * limit >= total}
								on:click={() => {
									currentPage++;
									loadClients();
								}}
								class="rounded-md border border-gray-300 bg-white px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
							>
								Next
							</button>
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Create Client Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
		<div class="flex min-h-screen items-center justify-center px-4 py-4 text-center sm:p-0">
			<!-- Backdrop -->
			<div 
				class="fixed inset-0 bg-gray-900/50 backdrop-blur-sm transition-opacity animate-fade-in" 
				on:click={() => showCreateModal = false}
				transition:fade={{ duration: 200 }}
			></div>

			<!-- Modal -->
			<div class="relative inline-block transform overflow-hidden rounded-2xl bg-white text-left align-bottom shadow-2xl transition-all animate-scale-in sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-6 pb-6 pt-6 sm:p-8">
					<div class="flex items-center justify-between mb-6">
						<h3 class="text-2xl font-bold text-gray-900" id="modal-title">
							Add New Client
						</h3>
						<button
							on:click={() => {
								showCreateModal = false;
								resetForm();
							}}
							class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<form on:submit|preventDefault={handleCreate} class="space-y-5">
						{#if createError}
							<div class="rounded-xl bg-red-50 border border-red-200 p-4 animate-fade-in">
								<p class="text-sm font-medium text-red-800">{createError}</p>
							</div>
						{/if}

						<div>
							<label class="block text-sm font-semibold text-gray-700 mb-3">Client Type</label>
							<div class="grid grid-cols-2 gap-3">
								<label class="relative flex cursor-pointer rounded-xl border-2 p-4 transition-all duration-200 {clientType === 'individual' ? 'border-primary-500 bg-primary-50 shadow-sm' : 'border-gray-200 hover:border-gray-300'}">
									<input
										type="radio"
										bind:group={clientType}
										value="individual"
										class="sr-only"
									/>
									<div class="flex items-center space-x-3">
										<div class="h-5 w-5 rounded-full border-2 flex items-center justify-center {clientType === 'individual' ? 'border-primary-500' : 'border-gray-300'}">
											{#if clientType === 'individual'}
												<div class="h-2.5 w-2.5 rounded-full bg-primary-500"></div>
											{/if}
										</div>
										<span class="text-sm font-medium text-gray-700">Individual</span>
									</div>
								</label>
								<label class="relative flex cursor-pointer rounded-xl border-2 p-4 transition-all duration-200 {clientType === 'company' ? 'border-primary-500 bg-primary-50 shadow-sm' : 'border-gray-200 hover:border-gray-300'}">
									<input
										type="radio"
										bind:group={clientType}
										value="company"
										class="sr-only"
									/>
									<div class="flex items-center space-x-3">
										<div class="h-5 w-5 rounded-full border-2 flex items-center justify-center {clientType === 'company' ? 'border-primary-500' : 'border-gray-300'}">
											{#if clientType === 'company'}
												<div class="h-2.5 w-2.5 rounded-full bg-primary-500"></div>
											{/if}
										</div>
										<span class="text-sm font-medium text-gray-700">Company</span>
									</div>
								</label>
							</div>
						</div>

						{#if clientType === 'individual'}
							<div>
								<label for="full-name" class="block text-sm font-semibold text-gray-700 mb-2">
									Full Name <span class="text-red-500">*</span>
								</label>
								<input
									id="full-name"
									type="text"
									bind:value={fullName}
									required
									class="input"
									placeholder="John Doe"
								/>
							</div>
						{:else}
							<div>
								<label for="company-name" class="block text-sm font-semibold text-gray-700 mb-2">
									Company Name <span class="text-red-500">*</span>
								</label>
								<input
									id="company-name"
									type="text"
									bind:value={companyName}
									required
									class="input"
									placeholder="Acme Corporation"
								/>
							</div>
						{/if}

						<div>
							<label for="email" class="block text-sm font-semibold text-gray-700 mb-2">
								Email <span class="text-red-500">*</span>
							</label>
							<input
								id="email"
								type="email"
								bind:value={email}
								required
								class="input"
								placeholder="client@example.com"
							/>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div>
								<label for="phone" class="block text-sm font-semibold text-gray-700 mb-2">
									Phone
								</label>
								<input
									id="phone"
									type="tel"
									bind:value={phone}
									class="input"
									placeholder="+27 11 123 4567"
								/>
							</div>
							<div>
								<label for="mobile" class="block text-sm font-semibold text-gray-700 mb-2">
									Mobile
								</label>
								<input
									id="mobile"
									type="tel"
									bind:value={mobile}
									class="input"
									placeholder="+27 82 123 4567"
								/>
							</div>
						</div>

						<div class="flex items-center justify-end space-x-3 pt-4 border-t border-gray-200">
							<button
								type="button"
								on:click={() => {
									showCreateModal = false;
									resetForm();
								}}
								class="btn btn-secondary"
							>
								Cancel
							</button>
							<button
								type="submit"
								disabled={isCreating}
								class="btn btn-primary shadow-md hover:shadow-lg"
							>
								{#if isCreating}
									<span class="flex items-center">
										<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
										Creating...
									</span>
								{:else}
									Create Client
								{/if}
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
{/if}
