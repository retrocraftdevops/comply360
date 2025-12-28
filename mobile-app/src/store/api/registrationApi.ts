/**
 * Registration API Slice
 * RTK Query API for company registration operations
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import Config from 'react-native-config';
import type { RootState } from '../store';

// API Configuration
const API_URL = Config.API_URL || 'http://localhost:8080/api/v1';

// Types
export interface Registration {
  id: string;
  tenant_id: string;
  client_id: string;
  client_name: string;
  company_name: string;
  registration_type: 'COMPANY' | 'CLOSE_CORPORATION' | 'TRUST' | 'NPO';
  country: string;
  status: 'DRAFT' | 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'REJECTED';
  cipc_registration_number?: string;
  tax_number?: string;
  vat_number?: string;
  total_amount: number;
  paid_amount: number;
  commission_amount: number;
  commission_paid: boolean;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface CreateRegistrationRequest {
  client_id: string;
  company_name: string;
  registration_type: string;
  country: string;
  business_activity?: string;
  directors?: Director[];
  shareholders?: Shareholder[];
}

export interface Director {
  first_name: string;
  last_name: string;
  id_number: string;
  nationality: string;
  residential_address: string;
}

export interface Shareholder {
  name: string;
  id_number?: string;
  shares: number;
  share_percentage: number;
}

export interface UpdateRegistrationRequest {
  status?: string;
  cipc_registration_number?: string;
  tax_number?: string;
  vat_number?: string;
  notes?: string;
}

export interface RegistrationListResponse {
  registrations: Registration[];
  total: number;
  page: number;
  limit: number;
}

/**
 * Registration API Slice
 */
export const registrationApi = createApi({
  reducerPath: 'registrationApi',
  baseQuery: fetchBaseQuery({
    baseUrl: `${API_URL}/registrations`,
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
  tagTypes: ['Registration'],
  endpoints: (builder) => ({
    /**
     * Get list of registrations
     */
    getRegistrations: builder.query<RegistrationListResponse, { page?: number; limit?: number; status?: string }>({
      query: ({ page = 1, limit = 20, status }) => ({
        url: '',
        params: { page, limit, ...(status && { status }) },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.registrations.map(({ id }) => ({ type: 'Registration' as const, id })),
              { type: 'Registration', id: 'LIST' },
            ]
          : [{ type: 'Registration', id: 'LIST' }],
    }),

    /**
     * Get single registration by ID
     */
    getRegistration: builder.query<Registration, string>({
      query: (id) => `/${id}`,
      providesTags: (result, error, id) => [{ type: 'Registration', id }],
    }),

    /**
     * Create new registration
     */
    createRegistration: builder.mutation<Registration, CreateRegistrationRequest>({
      query: (data) => ({
        url: '',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [{ type: 'Registration', id: 'LIST' }],
    }),

    /**
     * Update registration
     */
    updateRegistration: builder.mutation<Registration, { id: string; data: UpdateRegistrationRequest }>({
      query: ({ id, data }) => ({
        url: `/${id}`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Registration', id },
        { type: 'Registration', id: 'LIST' },
      ],
    }),

    /**
     * Delete registration
     */
    deleteRegistration: builder.mutation<void, string>({
      query: (id) => ({
        url: `/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Registration', id },
        { type: 'Registration', id: 'LIST' },
      ],
    }),

    /**
     * Get registration statistics
     */
    getRegistrationStats: builder.query<{
      total: number;
      pending: number;
      in_progress: number;
      completed: number;
      total_revenue: number;
      total_commission: number;
    }, void>({
      query: () => '/stats',
    }),
  }),
});

export const {
  useGetRegistrationsQuery,
  useGetRegistrationQuery,
  useCreateRegistrationMutation,
  useUpdateRegistrationMutation,
  useDeleteRegistrationMutation,
  useGetRegistrationStatsQuery,
} = registrationApi;
