<script lang="ts">
	import { onMount } from 'svelte';
	import apiClient from '$api/client';
	import type { Registration } from '$types';

	let isLoading = true;
	let registrations: Registration[] = [];
	let total = 0;
	let currentPage = 1;
	let limit = 20;

	onMount(async () => {
		await loadRegistrations();
	});

	async function loadRegistrations() {
		isLoading = true;
		try {
			const response = await apiClient.getRegistrations(currentPage, limit);
			registrations = response.data;
			total = response.total;
		} catch (error) {
			console.error('Failed to load registrations:', error);
		} finally {
			isLoading = false;
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'approved':
				return 'bg-green-100 text-green-800';
			case 'submitted':
			case 'in_review':
				return 'bg-yellow-100 text-yellow-800';
			case 'rejected':
			case 'cancelled':
				return 'bg-red-100 text-red-800';
			case 'draft':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-blue-100 text-blue-800';
		}
	}

	function formatRegistrationType(type: string): string {
		const types = {
			pty_ltd: 'Pty Ltd',
			close_corporation: 'Close Corporation',
			business_name: 'Business Name',
			vat_registration: 'VAT Registration'
		};
		return types[type as keyof typeof types] || type;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Registrations</h1>
			<p class="mt-1 text-sm text-gray-500">
				Manage company registrations and track their status
			</p>
		</div>
		<a
			href="/app/registrations/new"
			class="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90"
		>
			New Registration
		</a>
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
								Company Name
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Type
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Jurisdiction
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
						{#if registrations.length === 0}
							<tr>
								<td colspan="6" class="px-6 py-12 text-center">
									<p class="text-gray-500">No registrations found</p>
									<a
										href="/app/registrations/new"
										class="mt-4 inline-flex items-center text-sm font-medium text-primary hover:text-primary/80"
									>
										Create your first registration
									</a>
								</td>
							</tr>
						{:else}
							{#each registrations as registration}
								<tr class="hover:bg-gray-50">
									<td class="whitespace-nowrap px-6 py-4">
										<div class="text-sm font-medium text-gray-900">{registration.company_name}</div>
										{#if registration.registration_number}
											<div class="text-sm text-gray-500">{registration.registration_number}</div>
										{/if}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{formatRegistrationType(registration.registration_type)}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{registration.jurisdiction}
									</td>
									<td class="whitespace-nowrap px-6 py-4">
										<span
											class="inline-flex rounded-full px-2 py-1 text-xs font-semibold {getStatusColor(
												registration.status
											)}"
										>
											{registration.status}
										</span>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{formatDate(registration.created_at)}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
										<a
											href="/app/registrations/{registration.id}"
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
							<span class="font-medium">{total}</span> registrations
						</p>
						<div class="flex space-x-2">
							<button
								disabled={currentPage === 1}
								on:click={() => {
									currentPage--;
									loadRegistrations();
								}}
								class="rounded-md border border-gray-300 bg-white px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
							>
								Previous
							</button>
							<button
								disabled={currentPage * limit >= total}
								on:click={() => {
									currentPage++;
									loadRegistrations();
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
