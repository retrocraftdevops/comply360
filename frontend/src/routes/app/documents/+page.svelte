<script lang="ts">
	import { onMount } from 'svelte';
	import apiClient from '$api/client';
	import type { Document } from '$types';

	let isLoading = true;
	let documents: Document[] = [];
	let showUploadModal = false;
	let uploadFile: File | null = null;
	let documentType = 'id_document';
	let registrationId = '';
	let isUploading = false;
	let uploadError = '';

	const documentTypes = [
		{ value: 'id_document', label: 'ID Document' },
		{ value: 'proof_of_address', label: 'Proof of Address' },
		{ value: 'company_constitution', label: 'Company Constitution' },
		{ value: 'tax_certificate', label: 'Tax Certificate' },
		{ value: 'banking_details', label: 'Banking Details' },
		{ value: 'other', label: 'Other' }
	];

	onMount(async () => {
		await loadDocuments();
	});

	async function loadDocuments() {
		isLoading = true;
		try {
			documents = await apiClient.getDocuments();
		} catch (error) {
			console.error('Failed to load documents:', error);
		} finally {
			isLoading = false;
		}
	}

	async function handleUpload() {
		if (!uploadFile) {
			uploadError = 'Please select a file';
			return;
		}

		isUploading = true;
		uploadError = '';

		try {
			await apiClient.uploadDocument(
				uploadFile,
				documentType,
				registrationId || undefined
			);

			showUploadModal = false;
			uploadFile = null;
			documentType = 'id_document';
			registrationId = '';
			await loadDocuments();
		} catch (error: any) {
			uploadError = error.response?.data?.message || 'Upload failed';
		} finally {
			isUploading = false;
		}
	}

	async function handleDownload(doc: Document) {
		try {
			const response = await apiClient.getDocumentDownloadURL(doc.id);
			window.open(response.download_url, '_blank');
		} catch (error) {
			console.error('Failed to download document:', error);
		}
	}

	async function handleVerify(doc: Document) {
		try {
			await apiClient.verifyDocument(doc.id);
			await loadDocuments();
		} catch (error) {
			console.error('Failed to verify document:', error);
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
		return (bytes / 1048576).toFixed(1) + ' MB';
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
			case 'verified':
				return 'bg-green-100 text-green-800';
			case 'rejected':
				return 'bg-red-100 text-red-800';
			case 'expired':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-yellow-100 text-yellow-800';
		}
	}

	function handleFileChange(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			uploadFile = target.files[0];
		}
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900">Documents</h1>
			<p class="mt-1 text-sm text-gray-500">
				Upload and manage registration documents
			</p>
		</div>
		<button
			on:click={() => showUploadModal = true}
			class="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90"
		>
			Upload Document
		</button>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
		</div>
	{:else}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#if documents.length === 0}
				<div class="col-span-full rounded-lg bg-white p-12 text-center shadow">
					<p class="text-gray-500">No documents uploaded yet</p>
					<button
						on:click={() => showUploadModal = true}
						class="mt-4 inline-flex items-center text-sm font-medium text-primary hover:text-primary/80"
					>
						Upload your first document
					</button>
				</div>
			{:else}
				{#each documents as doc}
					<div class="rounded-lg bg-white p-6 shadow">
						<div class="flex items-start justify-between">
							<div class="flex-1">
								<h3 class="text-sm font-medium text-gray-900 truncate">{doc.file_name}</h3>
								<p class="mt-1 text-xs text-gray-500">
									{doc.document_type.replace(/_/g, ' ')}
								</p>
							</div>
							<span
								class="inline-flex rounded-full px-2 py-1 text-xs font-semibold {getStatusColor(
									doc.status
								)}"
							>
								{doc.status}
							</span>
						</div>

						<div class="mt-4 space-y-2 text-xs text-gray-500">
							<div class="flex justify-between">
								<span>Size:</span>
								<span class="font-medium">{formatFileSize(doc.file_size)}</span>
							</div>
							<div class="flex justify-between">
								<span>Uploaded:</span>
								<span class="font-medium">{formatDate(doc.created_at)}</span>
							</div>
							{#if doc.verified_at}
								<div class="flex justify-between">
									<span>Verified:</span>
									<span class="font-medium">{formatDate(doc.verified_at)}</span>
								</div>
							{/if}
							{#if doc.ai_verified}
								<div class="flex justify-between">
									<span>AI Score:</span>
									<span class="font-medium">{doc.ai_verification_score?.toFixed(1)}%</span>
								</div>
							{/if}
						</div>

						<div class="mt-4 flex space-x-2">
							<button
								on:click={() => handleDownload(doc)}
								class="flex-1 rounded-md border border-gray-300 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50"
							>
								Download
							</button>
							{#if doc.status === 'pending'}
								<button
									on:click={() => handleVerify(doc)}
									class="flex-1 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-white hover:bg-primary/90"
								>
									Verify
								</button>
							{/if}
						</div>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
</div>

<!-- Upload Modal -->
{#if showUploadModal}
	<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
		<div class="flex min-h-screen items-end justify-center px-4 pb-20 pt-4 text-center sm:block sm:p-0">
			<div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" on:click={() => showUploadModal = false}></div>

			<div class="inline-block transform overflow-hidden rounded-lg bg-white text-left align-bottom shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">
						Upload Document
					</h3>

					<form on:submit|preventDefault={handleUpload} class="mt-4 space-y-4">
						{#if uploadError}
							<div class="rounded-md bg-red-50 p-3">
								<p class="text-sm text-red-800">{uploadError}</p>
							</div>
						{/if}

						<div>
							<label for="document-type" class="block text-sm font-medium text-gray-700">
								Document Type
							</label>
							<select
								id="document-type"
								bind:value={documentType}
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
							>
								{#each documentTypes as type}
									<option value={type.value}>{type.label}</option>
								{/each}
							</select>
						</div>

						<div>
							<label for="registration-id" class="block text-sm font-medium text-gray-700">
								Registration ID (Optional)
							</label>
							<input
								id="registration-id"
								type="text"
								bind:value={registrationId}
								placeholder="Link to a registration"
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary sm:text-sm"
							/>
						</div>

						<div>
							<label for="file" class="block text-sm font-medium text-gray-700">
								File
							</label>
							<input
								id="file"
								type="file"
								on:change={handleFileChange}
								accept=".pdf,.jpg,.jpeg,.png"
								required
								class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:rounded-md file:border-0 file:bg-primary file:px-4 file:py-2 file:text-sm file:font-semibold file:text-white hover:file:bg-primary/90"
							/>
							<p class="mt-1 text-xs text-gray-500">PDF, JPG, PNG (max 50MB)</p>
							{#if uploadFile}
								<p class="mt-1 text-xs text-gray-700">
									Selected: {uploadFile.name} ({formatFileSize(uploadFile.size)})
								</p>
							{/if}
						</div>

						<div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
							<button
								type="submit"
								disabled={isUploading}
								class="inline-flex w-full justify-center rounded-md bg-primary px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90 disabled:opacity-50 sm:ml-3 sm:w-auto"
							>
								{#if isUploading}
									Uploading...
								{:else}
									Upload
								{/if}
							</button>
							<button
								type="button"
								on:click={() => showUploadModal = false}
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
