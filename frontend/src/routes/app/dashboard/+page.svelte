<script lang="ts">
	import { onMount } from 'svelte';
	import { authStore } from '$stores/auth';
	import apiClient from '$api/client';
	import type { Registration, Commission } from '$types';

	let isLoading = true;
	let recentRegistrations: Registration[] = [];
	let commissionSummary = {
		total_pending: 0,
		total_approved: 0,
		total_paid: 0
	};

	onMount(async () => {
		try {
			// Fetch recent registrations
			const regs = await apiClient.getRegistrations(1, 5);
			recentRegistrations = regs.data;

			// Fetch commission summary if user is an agent
			// Use get() to read store value non-reactively
			const auth = $authStore;
			if (auth.user) {
				try {
					const summary = await apiClient.getCommissionSummary(auth.user.id);
					commissionSummary = summary;
				} catch (e) {
					// User might not be an agent
					console.log('Not an agent user');
				}
			}
		} catch (error) {
			console.error('Failed to load dashboard data:', error);
		} finally {
			isLoading = false;
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-ZA', {
			style: 'currency',
			currency: 'ZAR'
		}).format(amount);
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'approved':
				return 'bg-green-100 text-green-800';
			case 'submitted':
			case 'in_review':
				return 'bg-yellow-100 text-yellow-800';
			case 'rejected':
				return 'bg-red-100 text-red-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="space-y-1">
		<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
			Dashboard
		</h1>
		<p class="text-sm text-gray-500 font-medium">
			Welcome back, <span class="font-semibold text-gray-700">{($authStore.user?.first_name) || 'User'}</span>! Here's what's happening today.
		</p>
	</div>

	{#if isLoading}
		<div class="card flex items-center justify-center py-16">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
				<p class="text-sm text-gray-500 font-medium">Loading dashboard...</p>
			</div>
		</div>
	{:else}
		<!-- Stats Cards -->
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<div class="card card-hover p-6">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-gray-500 uppercase tracking-wide">Total Registrations</p>
						<p class="mt-3 text-4xl font-bold text-gray-900">{recentRegistrations.length}</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-primary-100 flex items-center justify-center">
						<svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="card card-hover p-6 bg-gradient-to-br from-yellow-50 to-yellow-100/50 border-yellow-200/50">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-yellow-700 uppercase tracking-wide">Pending Commissions</p>
						<p class="mt-3 text-4xl font-bold text-yellow-900">
							{formatCurrency(commissionSummary.total_pending)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-yellow-200 flex items-center justify-center">
						<svg class="w-6 h-6 text-yellow-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>

			<div class="card card-hover p-6 bg-gradient-to-br from-green-50 to-green-100/50 border-green-200/50">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-sm font-semibold text-green-700 uppercase tracking-wide">Paid Commissions</p>
						<p class="mt-3 text-4xl font-bold text-green-900">
							{formatCurrency(commissionSummary.total_paid)}
						</p>
					</div>
					<div class="h-12 w-12 rounded-xl bg-green-200 flex items-center justify-center">
						<svg class="w-6 h-6 text-green-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
			</div>
		</div>

		<!-- Recent Registrations -->
		<div class="card card-hover overflow-hidden">
			<div class="border-b border-gray-200 bg-gray-50/50 px-6 py-4">
				<h2 class="text-lg font-semibold text-gray-900">Recent Registrations</h2>
			</div>
			<div class="divide-y divide-gray-100">
				{#if recentRegistrations.length === 0}
					<div class="px-6 py-16 text-center">
						<div class="flex flex-col items-center space-y-4">
							<div class="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
								<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
								</svg>
							</div>
							<div>
								<p class="text-sm font-medium text-gray-900">No registrations yet</p>
								<p class="text-sm text-gray-500 mt-1">Get started by creating your first registration</p>
							</div>
							<a
								href="/app/registrations/new"
								class="btn btn-primary mt-2"
							>
								Create Registration
							</a>
						</div>
					</div>
				{:else}
					{#each recentRegistrations as registration, index}
						<div class="px-6 py-4 table-row">
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-4">
									<div class="h-10 w-10 rounded-lg bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white text-sm font-semibold shadow-sm">
										{registration.company_name[0]}
									</div>
									<div>
										<h3 class="text-sm font-semibold text-gray-900">{registration.company_name}</h3>
										<p class="mt-1 text-xs text-gray-500">
											{registration.registration_type.replace(/_/g, ' ')} • {registration.jurisdiction}
										</p>
									</div>
								</div>
								<div class="flex items-center space-x-4">
									<span
										class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold shadow-sm {getStatusColor(
											registration.status
										)}"
									>
										{registration.status}
									</span>
									<span class="text-sm text-gray-500">{formatDate(registration.created_at)}</span>
								</div>
							</div>
						</div>
					{/each}
				{/if}
			</div>
			{#if recentRegistrations.length > 0}
				<div class="border-t border-gray-200 bg-gray-50/50 px-6 py-4">
					<a href="/app/registrations" class="text-sm font-semibold text-primary-600 hover:text-primary-700 transition-colors duration-200 inline-flex items-center">
						View all registrations
						<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
						</svg>
					</a>
				</div>
			{/if}
		</div>
	{/if}
</div>
