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

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Clients</h1>
			<p class="mt-1 text-sm text-gray-500">
				Manage client information and contacts
			</p>
		</div>
		<button
			on:click={() => showCreateModal = true}
			class="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90"
		>
			Add Client
		</button>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
		</div>
	{:else}
		<div class="rounded-lg bg-white shadow">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Name
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Type
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Email
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Phone
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Status
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Created
							</th>
							<th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 bg-white">
						{#if clients.length === 0}
							<tr>
								<td colspan="7" class="px-6 py-12 text-center">
									<p class="text-gray-500">No clients found</p>
									<button
										on:click={() => showCreateModal = true}
										class="mt-4 inline-flex items-center text-sm font-medium text-primary hover:text-primary/80"
									>
										Add your first client
									</button>
								</td>
							</tr>
						{:else}
							{#each clients as client}
								<tr class="hover:bg-gray-50">
									<td class="whitespace-nowrap px-6 py-4">
										<div class="text-sm font-medium text-gray-900">
											{client.client_type === 'individual' ? client.full_name : client.company_name}
										</div>
									</td>
									<td class="whitespace-nowrap px-6 py-4">
										<span class="inline-flex rounded-full px-2 py-1 text-xs font-medium capitalize {client.client_type === 'company' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'}">
											{client.client_type}
										</span>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{client.email}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{client.mobile || client.phone || '-'}
									</td>
									<td class="whitespace-nowrap px-6 py-4">
										<span
											class="inline-flex rounded-full px-2 py-1 text-xs font-semibold {getStatusColor(
												client.status
											)}"
										>
											{client.status}
										</span>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{formatDate(client.created_at)}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
										<a
											href="/app/clients/{client.id}"
											class="text-primary hover:text-primary/80"
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
		<div class="flex min-h-screen items-end justify-center px-4 pb-20 pt-4 text-center sm:block sm:p-0">
			<div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" on:click={() => showCreateModal = false}></div>

			<div class="inline-block transform overflow-hidden rounded-lg bg-white text-left align-bottom shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">
						Add New Client
					</h3>

					<form on:submit|preventDefault={handleCreate} class="mt-4 space-y-4">
						{#if createError}
							<div class="rounded-md bg-red-50 p-3">
								<p class="text-sm text-red-800">{createError}</p>
							</div>
						{/if}

						<div>
							<label class="block text-sm font-medium text-gray-700">Client Type</label>
							<div class="mt-2 flex space-x-4">
								<label class="flex items-center">
									<input
										type="radio"
										bind:group={clientType}
										value="individual"
										class="h-4 w-4 border-gray-300 text-primary focus:ring-primary"
									/>
									<span class="ml-2 text-sm text-gray-700">Individual</span>
								</label>
								<label class="flex items-center">
									<input
										type="radio"
										bind:group={clientType}
										value="company"
										class="h-4 w-4 border-gray-300 text-primary focus:ring-primary"
									/>
									<span class="ml-2 text-sm text-gray-700">Company</span>
								</label>
							</div>
						</div>

						{#if clientType === 'individual'}
							<div>
								<label for="full-name" class="block text-sm font-medium text-gray-700">
									Full Name <span class="text-red-500">*</span>
								</label>
								<input
									id="full-name"
									type="text"
									bind:value={fullName}
									required
									class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
								/>
							</div>
						{:else}
							<div>
								<label for="company-name" class="block text-sm font-medium text-gray-700">
									Company Name <span class="text-red-500">*</span>
								</label>
								<input
									id="company-name"
									type="text"
									bind:value={companyName}
									required
									class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
								/>
							</div>
						{/if}

						<div>
							<label for="email" class="block text-sm font-medium text-gray-700">
								Email <span class="text-red-500">*</span>
							</label>
							<input
								id="email"
								type="email"
								bind:value={email}
								required
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
							/>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div>
								<label for="phone" class="block text-sm font-medium text-gray-700">
									Phone
								</label>
								<input
									id="phone"
									type="tel"
									bind:value={phone}
									class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
								/>
							</div>
							<div>
								<label for="mobile" class="block text-sm font-medium text-gray-700">
									Mobile
								</label>
								<input
									id="mobile"
									type="tel"
									bind:value={mobile}
									class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
								/>
							</div>
						</div>

						<div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
							<button
								type="submit"
								disabled={isCreating}
								class="inline-flex w-full justify-center rounded-md bg-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90 disabled:opacity-50 sm:ml-3 sm:w-auto"
							>
								{#if isCreating}
									Creating...
								{:else}
									Create Client
								{/if}
							</button>
							<button
								type="button"
								on:click={() => {
									showCreateModal = false;
									resetForm();
								}}
								class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
							>
								Cancel
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
{/if}
