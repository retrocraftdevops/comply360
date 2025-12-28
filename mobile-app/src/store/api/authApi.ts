/**
 * Auth API Slice
 * RTK Query API for authentication operations
 */

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import Config from 'react-native-config';
import type { User } from '../slices/authSlice';

// API Configuration
const API_URL = Config.API_URL || 'http://localhost:8080/api/v1';

// Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
  refreshToken: string;
  expiresIn: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  company_name?: string;
}

export interface PasswordResetRequest {
  email: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  phone?: string;
}

/**
 * Auth API Slice
 */
export const authApi = createApi({
  reducerPath: 'authApi',
  baseQuery: fetchBaseQuery({
    baseUrl: `${API_URL}/auth`,
  }),
  tagTypes: ['User'],
  endpoints: (builder) => ({
    /**
     * Login with email and password
     */
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (credentials) => ({
        url: '/login',
        method: 'POST',
        body: credentials,
      }),
    }),

    /**
     * Register new user
     */
    register: builder.mutation<LoginResponse, RegisterRequest>({
      query: (userData) => ({
        url: '/register',
        method: 'POST',
        body: userData,
      }),
    }),

    /**
     * Logout
     */
    logout: builder.mutation<void, string>({
      query: (token) => ({
        url: '/logout',
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),

    /**
     * Refresh access token
     */
    refreshToken: builder.mutation<{ token: string; refreshToken: string }, string>({
      query: (refreshToken) => ({
        url: '/refresh',
        method: 'POST',
        body: { refresh_token: refreshToken },
      }),
    }),

    /**
     * Request password reset email
     */
    requestPasswordReset: builder.mutation<{ message: string }, PasswordResetRequest>({
      query: (data) => ({
        url: '/forgot-password',
        method: 'POST',
        body: data,
      }),
    }),

    /**
     * Change password (authenticated)
     */
    changePassword: builder.mutation<void, { token: string; data: ChangePasswordRequest }>({
      query: ({ token, data }) => ({
        url: '/change-password',
        method: 'POST',
        body: data,
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),

    /**
     * Get current user profile
     */
    getCurrentUser: builder.query<User, string>({
      query: (token) => ({
        url: '/me',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
      providesTags: ['User'],
    }),

    /**
     * Update user profile
     */
    updateProfile: builder.mutation<User, { token: string; data: UpdateProfileRequest }>({
      query: ({ token, data }) => ({
        url: '/profile',
        method: 'PATCH',
        body: data,
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
      invalidatesTags: ['User'],
    }),

    /**
     * Validate token
     */
    validateToken: builder.query<{ valid: boolean }, string>({
      query: (token) => ({
        url: '/validate',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),

    /**
     * Verify password
     */
    verifyPassword: builder.mutation<{ valid: boolean }, { token: string; password: string }>({
      query: ({ token, password }) => ({
        url: '/verify-password',
        method: 'POST',
        body: { password },
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),
  }),
});

export const {
  useLoginMutation,
  useRegisterMutation,
  useLogoutMutation,
  useRefreshTokenMutation,
  useRequestPasswordResetMutation,
  useChangePasswordMutation,
  useGetCurrentUserQuery,
  useUpdateProfileMutation,
  useValidateTokenQuery,
  useVerifyPasswordMutation,
} = authApi;
