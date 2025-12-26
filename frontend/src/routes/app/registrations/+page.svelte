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

<div class="space-y-6 animate-fade-in">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="space-y-1">
			<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
				Registrations
			</h1>
			<p class="text-sm text-gray-500 font-medium">
				Manage company registrations and track their status
			</p>
		</div>
		<a
			href="/app/registrations/new"
			class="btn btn-primary shadow-md hover:shadow-lg transition-all duration-200 group"
		>
			<svg class="w-4 h-4 mr-2 group-hover:rotate-90 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Registration
		</a>
	</div>

	{#if isLoading}
		<div class="card flex items-center justify-center py-16">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
				<p class="text-sm text-gray-500 font-medium">Loading registrations...</p>
			</div>
		</div>
	{:else}
		<div class="card card-hover overflow-hidden">
			<div class="overflow-x-auto">
				<table class="table">
					<thead class="table-header">
						<tr>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Company Name
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Type
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Jurisdiction
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
						{#if registrations.length === 0}
							<tr>
								<td colspan="6" class="px-6 py-16 text-center">
									<div class="flex flex-col items-center space-y-4">
										<div class="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
											<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
											</svg>
										</div>
										<div>
											<p class="text-sm font-medium text-gray-900">No registrations found</p>
											<p class="text-sm text-gray-500 mt-1">Get started by creating your first registration</p>
										</div>
										<a
											href="/app/registrations/new"
											class="btn btn-primary mt-2"
										>
											Create your first registration
										</a>
									</div>
								</td>
							</tr>
						{:else}
							{#each registrations as registration, index}
								<tr class="table-row animate-fade-in" style="animation-delay: {index * 50}ms">
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="flex items-center space-x-3">
											<div class="h-10 w-10 rounded-lg bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-sm font-semibold shadow-sm">
												{registration.company_name[0]}
											</div>
											<div>
												<div class="text-sm font-semibold text-gray-900">{registration.company_name}</div>
												{#if registration.registration_number}
													<div class="text-xs text-gray-500 mt-0.5">{registration.registration_number}</div>
												{/if}
											</div>
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{formatRegistrationType(registration.registration_type)}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{registration.jurisdiction}
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span
											class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold shadow-sm {getStatusColor(
												registration.status
											)}"
										>
											{registration.status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{formatDate(registration.created_at)}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
										<a
											href="/app/registrations/{registration.id}"
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
				<div class="border-t border-gray-200 bg-gray-50/50 px-6 py-4">
					<div class="flex items-center justify-between">
						<p class="text-sm text-gray-600 font-medium">
							Showing <span class="font-semibold text-gray-900">{(currentPage - 1) * limit + 1}</span> to
							<span class="font-semibold text-gray-900">{Math.min(currentPage * limit, total)}</span> of
							<span class="font-semibold text-gray-900">{total}</span> registrations
						</p>
						<div class="flex space-x-2">
							<button
								disabled={currentPage === 1}
								on:click={() => {
									currentPage--;
									loadRegistrations();
								}}
								class="btn btn-secondary text-sm disabled:opacity-50"
							>
								Previous
							</button>
							<button
								disabled={currentPage * limit >= total}
								on:click={() => {
									currentPage++;
									loadRegistrations();
								}}
								class="btn btn-secondary text-sm disabled:opacity-50"
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
