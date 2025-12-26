<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { authStore, isAgent, isAdmin } from '$stores/auth';
	import apiClient from '$api/client';
	import type { Commission, CommissionSummary } from '$types';

	let isLoading = true;
	let commissions: Commission[] = [];
	let summary: CommissionSummary | null = null;
	let total = 0;
	let currentPage = 1;
	let limit = 20;
	let showPaymentModal = false;
	let selectedCommission: Commission | null = null;
	let paymentReference = '';
	let isProcessing = false;

	onMount(async () => {
		await loadCommissions();
		if ($authStore.user && $isAgent) {
			await loadSummary();
		}
	});

	async function loadCommissions() {
		isLoading = true;
		try {
			const response = await apiClient.getCommissions(currentPage, limit);
			commissions = response.data;
			total = response.total;
		} catch (error) {
			console.error('Failed to load commissions:', error);
		} finally {
			isLoading = false;
		}
	}

	async function loadSummary() {
		if (!$authStore.user) return;
		try {
			summary = await apiClient.getCommissionSummary($authStore.user.id);
		} catch (error) {
			console.error('Failed to load summary:', error);
		}
	}

	async function handleApprove(commission: Commission) {
		if (!confirm('Approve this commission?')) return;

		try {
			await apiClient.approveCommission(commission.id);
			await loadCommissions();
			if (summary) await loadSummary();
		} catch (error) {
			console.error('Failed to approve commission:', error);
		}
	}

	function openPaymentModal(commission: Commission) {
		selectedCommission = commission;
		paymentReference = `PAY-${new Date().getFullYear()}${String(new Date().getMonth() + 1).padStart(2, '0')}${String(new Date().getDate()).padStart(2, '0')}-${Math.random().toString(36).substring(2, 8).toUpperCase()}`;
		showPaymentModal = true;
	}

	async function handlePayment() {
		if (!selectedCommission || !paymentReference) return;

		isProcessing = true;
		try {
			await apiClient.payCommission(selectedCommission.id, paymentReference);
			showPaymentModal = false;
			selectedCommission = null;
			paymentReference = '';
			await loadCommissions();
			if (summary) await loadSummary();
		} catch (error) {
			console.error('Failed to process payment:', error);
		} finally {
			isProcessing = false;
		}
	}

	function formatCurrency(amount: number, currency: string = 'ZAR'): string {
		return new Intl.NumberFormat('en-ZA', {
			style: 'currency',
			currency: currency
		}).format(amount);
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
			case 'paid':
				return 'bg-green-100 text-green-800';
			case 'approved':
				return 'bg-blue-100 text-blue-800';
			case 'pending':
				return 'bg-yellow-100 text-yellow-800';
			case 'cancelled':
				return 'bg-red-100 text-red-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}
</script>

<div class="space-y-8 animate-fade-in">
	<!-- Header -->
	<div class="space-y-1">
		<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
			Commissions
		</h1>
		<p class="text-sm text-gray-500 font-medium">
			Track and manage agent commissions
		</p>
	</div>

	<!-- Summary Cards (for agents) -->
	{#if $isAgent && summary}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
			<div class="card card-hover p-6 animate-fade-in" style="animation-delay: 0ms">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-gray-500 uppercase tracking-wide">Total Commissions</p>
						<p class="mt-3 text-3xl font-bold text-gray-900">
							{formatCurrency(summary.total_commissions, summary.currency)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-primary-100 flex items-center justify-center">
						<svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="card card-hover p-6 animate-fade-in bg-gradient-to-br from-yellow-50 to-yellow-100/50 border-yellow-200/50" style="animation-delay: 100ms">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-yellow-700 uppercase tracking-wide">Pending</p>
						<p class="mt-3 text-3xl font-bold text-yellow-900">
							{formatCurrency(summary.total_pending, summary.currency)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-yellow-200 flex items-center justify-center">
						<svg class="w-6 h-6 text-yellow-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="card card-hover p-6 animate-fade-in bg-gradient-to-br from-blue-50 to-blue-100/50 border-blue-200/50" style="animation-delay: 200ms">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-blue-700 uppercase tracking-wide">Approved</p>
						<p class="mt-3 text-3xl font-bold text-blue-900">
							{formatCurrency(summary.total_approved, summary.currency)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-blue-200 flex items-center justify-center">
						<svg class="w-6 h-6 text-blue-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="card card-hover p-6 animate-fade-in bg-gradient-to-br from-green-50 to-green-100/50 border-green-200/50" style="animation-delay: 300ms">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-green-700 uppercase tracking-wide">Paid</p>
						<p class="mt-3 text-3xl font-bold text-green-900">
							{formatCurrency(summary.total_paid, summary.currency)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-green-200 flex items-center justify-center">
						<svg class="w-6 h-6 text-green-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
						</svg>
					</div>
				</div>
			</div>
		</div>
	{/if}

	{#if isLoading}
		<div class="card flex items-center justify-center py-16">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
				<p class="text-sm text-gray-500 font-medium">Loading commissions...</p>
			</div>
		</div>
	{:else}
		<div class="card card-hover overflow-hidden">
			<div class="overflow-x-auto">
				<table class="table">
					<thead class="table-header">
						<tr>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Registration
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Fee
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Rate
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Commission
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Status
							</th>
							<th class="px-6 py-4 text-left text-xs font-semibold uppercase tracking-wider text-gray-600">
								Date
							</th>
							{#if $isAdmin}
								<th class="px-6 py-4 text-right text-xs font-semibold uppercase tracking-wider text-gray-600">
									Actions
								</th>
							{/if}
						</tr>
					</thead>
					<tbody class="bg-white divide-y divide-gray-100">
						{#if commissions.length === 0}
							<tr>
								<td colspan="7" class="px-6 py-16 text-center">
									<div class="flex flex-col items-center space-y-4">
										<div class="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
											<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
											</svg>
										</div>
										<div>
											<p class="text-sm font-medium text-gray-900">No commissions found</p>
											<p class="text-sm text-gray-500 mt-1">Commissions will appear here once registrations are completed</p>
										</div>
									</div>
								</td>
							</tr>
						{:else}
							{#each commissions as commission, index}
								<tr class="table-row animate-fade-in" style="animation-delay: {index * 50}ms">
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="flex items-center space-x-3">
											<div class="h-10 w-10 rounded-lg bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-xs font-semibold shadow-sm">
												{commission.registration_id.substring(0, 2).toUpperCase()}
											</div>
											<div class="text-sm font-semibold text-gray-900 font-mono">{commission.registration_id.substring(0, 8)}...</div>
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
										{formatCurrency(commission.registration_fee, commission.currency)}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{commission.commission_rate}%
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm font-bold text-primary-600">
										{formatCurrency(commission.commission_amount, commission.currency)}
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span
											class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold shadow-sm {getStatusColor(
												commission.status
											)}"
										>
											{commission.status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
										{formatDate(commission.created_at)}
									</td>
									{#if $isAdmin}
										<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
											{#if commission.status === 'pending'}
												<button
													on:click={() => handleApprove(commission)}
													class="text-primary-600 hover:text-primary-700 font-medium transition-colors duration-200"
												>
													Approve
												</button>
											{:else if commission.status === 'approved'}
												<button
													on:click={() => openPaymentModal(commission)}
													class="text-primary-600 hover:text-primary-700 font-medium transition-colors duration-200"
												>
													Mark Paid
												</button>
											{:else if commission.status === 'paid'}
												<span class="text-gray-400">Completed</span>
											{/if}
										</td>
									{/if}
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
							<span class="font-semibold text-gray-900">{total}</span> commissions
						</p>
						<div class="flex space-x-2">
							<button
								disabled={currentPage === 1}
								on:click={() => {
									currentPage--;
									loadCommissions();
								}}
								class="btn btn-secondary text-sm disabled:opacity-50"
							>
								Previous
							</button>
							<button
								disabled={currentPage * limit >= total}
								on:click={() => {
									currentPage++;
									loadCommissions();
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

<!-- Payment Modal -->
{#if showPaymentModal && selectedCommission}
	<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
		<div class="flex min-h-screen items-center justify-center px-4 py-4 text-center sm:p-0">
			<div class="fixed inset-0 bg-gray-900/50 backdrop-blur-sm" on:click={() => showPaymentModal = false} transition:fade={{ duration: 200 }}></div>

			<div class="relative inline-block transform overflow-hidden rounded-2xl bg-white text-left align-bottom shadow-2xl transition-all animate-scale-in sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-6 pb-6 pt-6 sm:p-8">
					<div class="flex items-center justify-between mb-6">
						<h3 class="text-2xl font-bold text-gray-900" id="modal-title">
							Mark Commission as Paid
						</h3>
						<button
							on:click={() => showPaymentModal = false}
							class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<div class="space-y-5">
						<div class="rounded-xl bg-gradient-to-br from-primary-50 to-primary-100/50 border border-primary-200/50 p-5">
							<div class="grid grid-cols-2 gap-4">
								<div>
									<span class="text-xs font-semibold text-primary-700 uppercase tracking-wide">Amount</span>
									<p class="mt-2 text-xl font-bold text-primary-900">{formatCurrency(selectedCommission.commission_amount, selectedCommission.currency)}</p>
								</div>
								<div>
									<span class="text-xs font-semibold text-primary-700 uppercase tracking-wide">Agent ID</span>
									<p class="mt-2 text-sm font-mono font-semibold text-primary-900">{selectedCommission.agent_id.substring(0, 8)}...</p>
								</div>
							</div>
						</div>

						<div>
							<label for="payment-reference" class="block text-sm font-semibold text-gray-700 mb-2">
								Payment Reference <span class="text-red-500">*</span>
							</label>
							<input
								id="payment-reference"
								type="text"
								bind:value={paymentReference}
								required
								class="input"
								placeholder="BANK-TXN-2025-001"
							/>
							<p class="mt-2 text-xs text-gray-500">Enter the bank transaction reference</p>
						</div>

						<div class="flex items-center justify-end space-x-3 pt-4 border-t border-gray-200">
							<button
								type="button"
								on:click={() => showPaymentModal = false}
								class="btn btn-secondary"
							>
								Cancel
							</button>
							<button
								type="button"
								on:click={handlePayment}
								disabled={isProcessing || !paymentReference}
								class="btn btn-primary shadow-md hover:shadow-lg"
							>
								{#if isProcessing}
									<span class="flex items-center">
										<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
										Processing...
									</span>
								{:else}
									Confirm Payment
								{/if}
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
