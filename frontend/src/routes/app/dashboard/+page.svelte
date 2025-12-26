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
			if ($authStore.user) {
				try {
					const summary = await apiClient.getCommissionSummary($authStore.user.id);
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

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
		<p class="mt-1 text-sm text-gray-500">Welcome back, {$authStore.user?.first_name}!</p>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
		</div>
	{:else}
		<!-- Stats -->
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<div class="rounded-lg bg-white p-6 shadow">
				<h3 class="text-sm font-medium text-gray-500">Total Registrations</h3>
				<p class="mt-2 text-3xl font-semibold text-gray-900">{recentRegistrations.length}</p>
			</div>

			<div class="rounded-lg bg-white p-6 shadow">
				<h3 class="text-sm font-medium text-gray-500">Pending Commissions</h3>
				<p class="mt-2 text-3xl font-semibold text-gray-900">
					{formatCurrency(commissionSummary.total_pending)}
				</p>
			</div>

			<div class="rounded-lg bg-white p-6 shadow">
				<h3 class="text-sm font-medium text-gray-500">Paid Commissions</h3>
				<p class="mt-2 text-3xl font-semibold text-gray-900">
					{formatCurrency(commissionSummary.total_paid)}
				</p>
			</div>
		</div>

		<!-- Recent Registrations -->
		<div class="rounded-lg bg-white shadow">
			<div class="border-b border-gray-200 px-6 py-4">
				<h2 class="text-lg font-medium text-gray-900">Recent Registrations</h2>
			</div>
			<div class="divide-y divide-gray-200">
				{#if recentRegistrations.length === 0}
					<div class="px-6 py-12 text-center">
						<p class="text-gray-500">No registrations yet</p>
						<a
							href="/app/registrations/new"
							class="mt-4 inline-flex items-center rounded-md bg-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90"
						>
							Create Registration
						</a>
					</div>
				{:else}
					{#each recentRegistrations as registration}
						<div class="px-6 py-4">
							<div class="flex items-center justify-between">
								<div>
									<h3 class="text-sm font-medium text-gray-900">{registration.company_name}</h3>
									<p class="mt-1 text-sm text-gray-500">
										{registration.registration_type} • {registration.jurisdiction}
									</p>
								</div>
								<div class="flex items-center space-x-4">
									<span
										class="inline-flex rounded-full px-2 py-1 text-xs font-semibold {getStatusColor(
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
				<div class="border-t border-gray-200 px-6 py-3">
					<a href="/app/registrations" class="text-sm font-medium text-primary hover:text-primary/80">
						View all registrations →
					</a>
				</div>
			{/if}
		</div>
	{/if}
</div>
