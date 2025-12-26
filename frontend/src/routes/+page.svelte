<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		// Redirect to login page on mount
		goto('/auth/login');
	});
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-50">
	<div class="text-center">
		<h1 class="text-4xl font-bold text-gray-900 mb-4">Comply360</h1>
		<p class="text-gray-600">Redirecting to login...</p>
	</div>
</div>
