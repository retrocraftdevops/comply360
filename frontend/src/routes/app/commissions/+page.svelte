<script lang="ts">
	import { onMount } from 'svelte';
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

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-gray-900">Commissions</h1>
		<p class="mt-1 text-sm text-gray-500">
			Track and manage agent commissions
		</p>
	</div>

	<!-- Summary Cards (for agents) -->
	{#if $isAgent && summary}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
			<div class="rounded-lg bg-white p-6 shadow">
				<h3 class="text-sm font-medium text-gray-500">Total Commissions</h3>
				<p class="mt-2 text-2xl font-semibold text-gray-900">
					{formatCurrency(summary.total_commissions, summary.currency)}
				</p>
			</div>

			<div class="rounded-lg bg-yellow-50 p-6 shadow">
				<h3 class="text-sm font-medium text-yellow-800">Pending</h3>
				<p class="mt-2 text-2xl font-semibold text-yellow-900">
					{formatCurrency(summary.total_pending, summary.currency)}
				</p>
			</div>

			<div class="rounded-lg bg-blue-50 p-6 shadow">
				<h3 class="text-sm font-medium text-blue-800">Approved</h3>
				<p class="mt-2 text-2xl font-semibold text-blue-900">
					{formatCurrency(summary.total_approved, summary.currency)}
				</p>
			</div>

			<div class="rounded-lg bg-green-50 p-6 shadow">
				<h3 class="text-sm font-medium text-green-800">Paid</h3>
				<p class="mt-2 text-2xl font-semibold text-green-900">
					{formatCurrency(summary.total_paid, summary.currency)}
				</p>
			</div>
		</div>
	{/if}

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
								Registration
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Fee
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Rate
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Commission
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Status
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
								Date
							</th>
							{#if $isAdmin}
								<th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
									Actions
								</th>
							{/if}
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200 bg-white">
						{#if commissions.length === 0}
							<tr>
								<td colspan="7" class="px-6 py-12 text-center">
									<p class="text-gray-500">No commissions found</p>
								</td>
							</tr>
						{:else}
							{#each commissions as commission}
								<tr class="hover:bg-gray-50">
									<td class="whitespace-nowrap px-6 py-4">
										<div class="text-sm font-medium text-gray-900">{commission.registration_id.substring(0, 8)}...</div>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
										{formatCurrency(commission.registration_fee, commission.currency)}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
										{commission.commission_rate}%
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
										{formatCurrency(commission.commission_amount, commission.currency)}
									</td>
									<td class="whitespace-nowrap px-6 py-4">
										<span
											class="inline-flex rounded-full px-2 py-1 text-xs font-semibold {getStatusColor(
												commission.status
											)}"
										>
											{commission.status}
										</span>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{formatDate(commission.created_at)}
									</td>
									{#if $isAdmin}
										<td class="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
											{#if commission.status === 'pending'}
												<button
													on:click={() => handleApprove(commission)}
													class="text-primary hover:text-primary/80"
												>
													Approve
												</button>
											{:else if commission.status === 'approved'}
												<button
													on:click={() => openPaymentModal(commission)}
													class="text-primary hover:text-primary/80"
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
				<div class="border-t border-gray-200 px-6 py-3">
					<div class="flex items-center justify-between">
						<p class="text-sm text-gray-700">
							Showing <span class="font-medium">{(currentPage - 1) * limit + 1}</span> to
							<span class="font-medium">{Math.min(currentPage * limit, total)}</span> of
							<span class="font-medium">{total}</span> commissions
						</p>
						<div class="flex space-x-2">
							<button
								disabled={currentPage === 1}
								on:click={() => {
									currentPage--;
									loadCommissions();
								}}
								class="rounded-md border border-gray-300 bg-white px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
							>
								Previous
							</button>
							<button
								disabled={currentPage * limit >= total}
								on:click={() => {
									currentPage++;
									loadCommissions();
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

<!-- Payment Modal -->
{#if showPaymentModal && selectedCommission}
	<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
		<div class="flex min-h-screen items-end justify-center px-4 pb-20 pt-4 text-center sm:block sm:p-0">
			<div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" on:click={() => showPaymentModal = false}></div>

			<div class="inline-block transform overflow-hidden rounded-lg bg-white text-left align-bottom shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">
						Mark Commission as Paid
					</h3>

					<div class="mt-4 space-y-4">
						<div class="rounded-md bg-gray-50 p-4">
							<div class="grid grid-cols-2 gap-4 text-sm">
								<div>
									<span class="text-gray-500">Amount:</span>
									<p class="font-medium">{formatCurrency(selectedCommission.commission_amount, selectedCommission.currency)}</p>
								</div>
								<div>
									<span class="text-gray-500">Agent ID:</span>
									<p class="font-medium">{selectedCommission.agent_id.substring(0, 8)}...</p>
								</div>
							</div>
						</div>

						<div>
							<label for="payment-reference" class="block text-sm font-medium text-gray-700">
								Payment Reference
							</label>
							<input
								id="payment-reference"
								type="text"
								bind:value={paymentReference}
								required
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
							/>
							<p class="mt-1 text-xs text-gray-500">Enter the bank transaction reference</p>
						</div>

						<div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
							<button
								type="button"
								on:click={handlePayment}
								disabled={isProcessing || !paymentReference}
								class="inline-flex w-full justify-center rounded-md bg-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90 disabled:opacity-50 sm:ml-3 sm:w-auto"
							>
								{#if isProcessing}
									Processing...
								{:else}
									Confirm Payment
								{/if}
							</button>
							<button
								type="button"
								on:click={() => showPaymentModal = false}
								class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
							>
								Cancel
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
