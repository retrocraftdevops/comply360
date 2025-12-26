<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		// Redirect to login page on mount
		goto('/auth/login');
	});
</script>

<div class="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-50 via-white to-primary-50/20">
	<div class="text-center space-y-6 animate-fade-in">
		<div class="flex justify-center">
			<div class="h-20 w-20 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center shadow-lg shadow-primary-500/20">
				<span class="text-white font-bold text-3xl">C</span>
			</div>
		</div>
		<div>
			<h1 class="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-2">
				Comply360
			</h1>
			<p class="text-sm font-medium text-gray-600">Redirecting to login...</p>
		</div>
		<div class="flex justify-center">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
		</div>
	</div>
</div>
