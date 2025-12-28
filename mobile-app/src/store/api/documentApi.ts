/**
 * Document API Slice
 * RTK Query API for document management operations
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import Config from 'react-native-config';
import type { RootState } from '../store';

// API Configuration
const API_URL = Config.API_URL || 'http://localhost:8080/api/v1';

// Types
export interface Document {
  id: string;
  tenant_id: string;
  registration_id?: string;
  client_id?: string;
  name: string;
  description?: string;
  file_name: string;
  file_size: number;
  file_type: string;
  category: string;
  storage_path: string;
  url?: string;
  status: 'PENDING' | 'VERIFIED' | 'REJECTED';
  uploaded_by: string;
  created_at: string;
  updated_at: string;
}

export interface UploadDocumentRequest {
  registration_id?: string;
  client_id?: string;
  name: string;
  description?: string;
  category: string;
  file: {
    uri: string;
    type: string;
    name: string;
  };
}

export interface UpdateDocumentRequest {
  name?: string;
  description?: string;
  category?: string;
  status?: string;
}

export interface DocumentListResponse {
  documents: Document[];
  total: number;
  page: number;
  limit: number;
}

/**
 * Document API Slice
 */
export const documentApi = createApi({
  reducerPath: 'documentApi',
  baseQuery: fetchBaseQuery({
    baseUrl: `${API_URL}/documents`,
    prepareHeaders: (headers, { getState }) => {
      const state = getState() as RootState;
      const token = state.auth.token;
      const tenantId = state.auth.user?.tenant_id;

      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }

      if (tenantId) {
        headers.set('X-Tenant-ID', tenantId);
      }

      return headers;
    },
  }),
  tagTypes: ['Document'],
  endpoints: (builder) => ({
    /**
     * Get list of documents
     */
    getDocuments: builder.query<DocumentListResponse, {
      page?: number;
      limit?: number;
      category?: string;
      registration_id?: string;
    }>({
      query: ({ page = 1, limit = 20, category, registration_id }) => ({
        url: '',
        params: {
          page,
          limit,
          ...(category && { category }),
          ...(registration_id && { registration_id }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.documents.map(({ id }) => ({ type: 'Document' as const, id })),
              { type: 'Document', id: 'LIST' },
            ]
          : [{ type: 'Document', id: 'LIST' }],
    }),

    /**
     * Get single document by ID
     */
    getDocument: builder.query<Document, string>({
      query: (id) => `/${id}`,
      providesTags: (result, error, id) => [{ type: 'Document', id }],
    }),

    /**
     * Upload new document
     * Note: This uses multipart/form-data for file upload
     */
    uploadDocument: builder.mutation<Document, UploadDocumentRequest>({
      query: (data) => {
        const formData = new FormData();

        if (data.registration_id) {
          formData.append('registration_id', data.registration_id);
        }
        if (data.client_id) {
          formData.append('client_id', data.client_id);
        }
        formData.append('name', data.name);
        if (data.description) {
          formData.append('description', data.description);
        }
        formData.append('category', data.category);
        formData.append('file', data.file as any);

        return {
          url: '/upload',
          method: 'POST',
          body: formData,
        };
      },
      invalidatesTags: [{ type: 'Document', id: 'LIST' }],
    }),

    /**
     * Update document metadata
     */
    updateDocument: builder.mutation<Document, { id: string; data: UpdateDocumentRequest }>({
      query: ({ id, data }) => ({
        url: `/${id}`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Document', id },
        { type: 'Document', id: 'LIST' },
      ],
    }),

    /**
     * Delete document
     */
    deleteDocument: builder.mutation<void, string>({
      query: (id) => ({
        url: `/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Document', id },
        { type: 'Document', id: 'LIST' },
      ],
    }),

    /**
     * Download document (get download URL)
     */
    getDocumentDownloadUrl: builder.query<{ url: string }, string>({
      query: (id) => `/${id}/download`,
    }),

    /**
     * Get document categories
     */
    getDocumentCategories: builder.query<{ categories: string[] }, void>({
      query: () => '/categories',
    }),

    /**
     * Get document statistics
     */
    getDocumentStats: builder.query<{
      total: number;
      pending: number;
      verified: number;
      rejected: number;
      total_size_mb: number;
    }, void>({
      query: () => '/stats',
    }),
  }),
});

export const {
  useGetDocumentsQuery,
  useGetDocumentQuery,
  useUploadDocumentMutation,
  useUpdateDocumentMutation,
  useDeleteDocumentMutation,
  useGetDocumentDownloadUrlQuery,
  useGetDocumentCategoriesQuery,
  useGetDocumentStatsQuery,
} = documentApi;
