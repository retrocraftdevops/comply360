<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
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

<div class="space-y-6 animate-fade-in">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="space-y-1">
			<h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
				Documents
			</h1>
			<p class="text-sm text-gray-500 font-medium">
				Upload and manage registration documents
			</p>
		</div>
		<button
			on:click={() => showUploadModal = true}
			class="btn btn-primary shadow-md hover:shadow-lg transition-all duration-200 group"
		>
			<svg class="w-4 h-4 mr-2 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
			</svg>
			Upload Document
		</button>
	</div>

	{#if isLoading}
		<div class="card flex items-center justify-center py-16">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary-200 border-t-primary-500"></div>
				<p class="text-sm text-gray-500 font-medium">Loading documents...</p>
			</div>
		</div>
	{:else}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#if documents.length === 0}
				<div class="col-span-full card p-16 text-center">
					<div class="flex flex-col items-center space-y-4">
						<div class="h-16 w-16 rounded-full bg-gray-100 flex items-center justify-center">
							<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
							</svg>
						</div>
						<div>
							<p class="text-sm font-medium text-gray-900">No documents uploaded yet</p>
							<p class="text-sm text-gray-500 mt-1">Get started by uploading your first document</p>
						</div>
						<button
							on:click={() => showUploadModal = true}
							class="btn btn-primary mt-2"
						>
							Upload your first document
						</button>
					</div>
				</div>
			{:else}
				{#each documents as doc, index}
					<div class="card card-hover p-6 animate-fade-in" style="animation-delay: {index * 50}ms">
						<div class="flex items-start justify-between mb-4">
							<div class="flex-1 min-w-0">
								<div class="flex items-center space-x-3">
									<div class="h-10 w-10 rounded-lg bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white shadow-sm flex-shrink-0">
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
										</svg>
									</div>
									<div class="flex-1 min-w-0">
										<h3 class="text-sm font-semibold text-gray-900 truncate">{doc.file_name}</h3>
										<p class="mt-1 text-xs text-gray-500 capitalize">
											{doc.document_type.replace(/_/g, ' ')}
										</p>
									</div>
								</div>
							</div>
							<span
								class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold shadow-sm flex-shrink-0 ml-2 {getStatusColor(
									doc.status
								)}"
							>
								{doc.status}
							</span>
						</div>

						<div class="space-y-2 text-xs text-gray-600 mb-4 pb-4 border-b border-gray-100">
							<div class="flex justify-between">
								<span class="text-gray-500">Size:</span>
								<span class="font-semibold text-gray-700">{formatFileSize(doc.file_size)}</span>
							</div>
							<div class="flex justify-between">
								<span class="text-gray-500">Uploaded:</span>
								<span class="font-semibold text-gray-700">{formatDate(doc.created_at)}</span>
							</div>
							{#if doc.verified_at}
								<div class="flex justify-between">
									<span class="text-gray-500">Verified:</span>
									<span class="font-semibold text-gray-700">{formatDate(doc.verified_at)}</span>
								</div>
							{/if}
							{#if doc.ai_verified}
								<div class="flex justify-between">
									<span class="text-gray-500">AI Score:</span>
									<span class="font-semibold text-green-600">{doc.ai_verification_score?.toFixed(1)}%</span>
								</div>
							{/if}
						</div>

						<div class="flex space-x-2">
							<button
								on:click={() => handleDownload(doc)}
								class="btn btn-secondary flex-1 text-xs"
							>
								Download
							</button>
							{#if doc.status === 'pending'}
								<button
									on:click={() => handleVerify(doc)}
									class="btn btn-primary flex-1 text-xs"
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
		<div class="flex min-h-screen items-center justify-center px-4 py-4 text-center sm:p-0">
			<div class="fixed inset-0 bg-gray-900/50 backdrop-blur-sm" on:click={() => showUploadModal = false} transition:fade={{ duration: 200 }}></div>

			<div class="relative inline-block transform overflow-hidden rounded-2xl bg-white text-left align-bottom shadow-2xl transition-all animate-scale-in sm:my-8 sm:w-full sm:max-w-lg sm:align-middle">
				<div class="bg-white px-6 pb-6 pt-6 sm:p-8">
					<div class="flex items-center justify-between mb-6">
						<h3 class="text-2xl font-bold text-gray-900" id="modal-title">
							Upload Document
						</h3>
						<button
							on:click={() => showUploadModal = false}
							class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<form on:submit|preventDefault={handleUpload} class="space-y-5">
						{#if uploadError}
							<div class="rounded-xl bg-red-50 border border-red-200 p-4 animate-fade-in">
								<p class="text-sm font-medium text-red-800">{uploadError}</p>
							</div>
						{/if}

						<div>
							<label for="document-type" class="block text-sm font-semibold text-gray-700 mb-2">
								Document Type
							</label>
							<select
								id="document-type"
								bind:value={documentType}
								class="input"
							>
								{#each documentTypes as type}
									<option value={type.value}>{type.label}</option>
								{/each}
							</select>
						</div>

						<div>
							<label for="registration-id" class="block text-sm font-semibold text-gray-700 mb-2">
								Registration ID <span class="text-gray-400 font-normal">(Optional)</span>
							</label>
							<input
								id="registration-id"
								type="text"
								bind:value={registrationId}
								placeholder="Link to a registration"
								class="input"
							/>
						</div>

						<div>
							<label for="file" class="block text-sm font-semibold text-gray-700 mb-2">
								File
							</label>
							<div class="mt-1 flex justify-center rounded-xl border-2 border-dashed border-gray-300 px-6 pt-5 pb-6 hover:border-primary-400 transition-colors duration-200">
								<div class="space-y-1 text-center">
									<svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
										<path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28.828m-12 0l-4 4m0 0l4 4m4-4l4-4m12 0l-4-4m4 4v-4m0 4h-4" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" />
									</svg>
									<div class="flex text-sm text-gray-600">
										<label for="file" class="relative cursor-pointer rounded-md font-semibold text-primary-600 hover:text-primary-500">
											<span>Upload a file</span>
											<input
												id="file"
												type="file"
												on:change={handleFileChange}
												accept=".pdf,.jpg,.jpeg,.png"
												required
												class="sr-only"
											/>
										</label>
										<p class="pl-1">or drag and drop</p>
									</div>
									<p class="text-xs text-gray-500">PDF, JPG, PNG (max 50MB)</p>
									{#if uploadFile}
										<p class="mt-2 text-xs font-medium text-gray-700">
											✓ {uploadFile.name} ({formatFileSize(uploadFile.size)})
										</p>
									{/if}
								</div>
							</div>
						</div>

						<div class="flex items-center justify-end space-x-3 pt-4 border-t border-gray-200">
							<button
								type="button"
								on:click={() => showUploadModal = false}
								class="btn btn-secondary"
							>
								Cancel
							</button>
							<button
								type="submit"
								disabled={isUploading}
								class="btn btn-primary shadow-md hover:shadow-lg"
							>
								{#if isUploading}
									<span class="flex items-center">
										<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
											<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
											<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
										</svg>
										Uploading...
									</span>
								{:else}
									Upload
								{/if}
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
{/if}
