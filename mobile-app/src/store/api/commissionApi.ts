/**
 * Commission API Slice
 * RTK Query API for commission tracking operations
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import Config from 'react-native-config';
import type { RootState } from '../store';

// API Configuration
const API_URL = Config.API_URL || 'http://localhost:8080/api/v1';

// Types
export interface Commission {
  id: string;
  tenant_id: string;
  registration_id: string;
  agent_id: string;
  agent_name: string;
  client_name: string;
  company_name: string;
  amount: number;
  percentage: number;
  status: 'PENDING' | 'APPROVED' | 'PAID' | 'DISPUTED';
  payment_method?: string;
  payment_reference?: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCommissionRequest {
  registration_id: string;
  amount: number;
  percentage?: number;
  notes?: string;
}

export interface UpdateCommissionRequest {
  status?: string;
  payment_method?: string;
  payment_reference?: string;
  notes?: string;
}

export interface CommissionListResponse {
  commissions: Commission[];
  total: number;
  page: number;
  limit: number;
}

export interface CommissionStats {
  total_earned: number;
  total_pending: number;
  total_paid: number;
  total_disputed: number;
  count_total: number;
  count_pending: number;
  count_paid: number;
  count_disputed: number;
  average_commission: number;
  this_month_earned: number;
  this_month_paid: number;
}

export interface CommissionEarnings {
  month: string;
  earned: number;
  paid: number;
}

/**
 * Commission API Slice
 */
export const commissionApi = createApi({
  reducerPath: 'commissionApi',
  baseQuery: fetchBaseQuery({
    baseUrl: `${API_URL}/commissions`,
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
  tagTypes: ['Commission'],
  endpoints: (builder) => ({
    /**
     * Get list of commissions
     */
    getCommissions: builder.query<CommissionListResponse, {
      page?: number;
      limit?: number;
      status?: string;
      start_date?: string;
      end_date?: string;
    }>({
      query: ({ page = 1, limit = 20, status, start_date, end_date }) => ({
        url: '',
        params: {
          page,
          limit,
          ...(status && { status }),
          ...(start_date && { start_date }),
          ...(end_date && { end_date }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.commissions.map(({ id }) => ({ type: 'Commission' as const, id })),
              { type: 'Commission', id: 'LIST' },
            ]
          : [{ type: 'Commission', id: 'LIST' }],
    }),

    /**
     * Get single commission by ID
     */
    getCommission: builder.query<Commission, string>({
      query: (id) => `/${id}`,
      providesTags: (result, error, id) => [{ type: 'Commission', id }],
    }),

    /**
     * Create new commission
     */
    createCommission: builder.mutation<Commission, CreateCommissionRequest>({
      query: (data) => ({
        url: '',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [{ type: 'Commission', id: 'LIST' }],
    }),

    /**
     * Update commission
     */
    updateCommission: builder.mutation<Commission, { id: string; data: UpdateCommissionRequest }>({
      query: ({ id, data }) => ({
        url: `/${id}`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Commission', id },
        { type: 'Commission', id: 'LIST' },
      ],
    }),

    /**
     * Mark commission as paid
     */
    markCommissionAsPaid: builder.mutation<Commission, {
      id: string;
      payment_method: string;
      payment_reference: string;
    }>({
      query: ({ id, payment_method, payment_reference }) => ({
        url: `/${id}/mark-paid`,
        method: 'POST',
        body: { payment_method, payment_reference },
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Commission', id },
        { type: 'Commission', id: 'LIST' },
      ],
    }),

    /**
     * Delete commission
     */
    deleteCommission: builder.mutation<void, string>({
      query: (id) => ({
        url: `/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Commission', id },
        { type: 'Commission', id: 'LIST' },
      ],
    }),

    /**
     * Get commission statistics
     */
    getCommissionStats: builder.query<CommissionStats, void>({
      query: () => '/stats',
    }),

    /**
     * Get monthly commission earnings
     */
    getCommissionEarnings: builder.query<{ earnings: CommissionEarnings[] }, {
      start_date?: string;
      end_date?: string;
    }>({
      query: ({ start_date, end_date }) => ({
        url: '/earnings',
        params: {
          ...(start_date && { start_date }),
          ...(end_date && { end_date }),
        },
      }),
    }),

    /**
     * Request commission payout
     */
    requestPayout: builder.mutation<{
      message: string;
      payout_request_id: string;
    }, {
      commission_ids: string[];
      payment_method: string;
      notes?: string;
    }>({
      query: (data) => ({
        url: '/request-payout',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [{ type: 'Commission', id: 'LIST' }],
    }),
  }),
});

export const {
  useGetCommissionsQuery,
  useGetCommissionQuery,
  useCreateCommissionMutation,
  useUpdateCommissionMutation,
  useMarkCommissionAsPaidMutation,
  useDeleteCommissionMutation,
  useGetCommissionStatsQuery,
  useGetCommissionEarningsQuery,
  useRequestPayoutMutation,
} = commissionApi;
