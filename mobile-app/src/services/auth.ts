/**
 * Authentication Service
 * Handles all authentication-related API calls
 */

import axios, { AxiosInstance } from 'axios';
import Config from 'react-native-config';

// API Configuration
const API_URL = Config.API_URL || 'http://localhost:8080/api/v1';
const API_TIMEOUT = 30000; // 30 seconds

// Types
export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  tenant_id: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  user: User;
  token: string;
  refreshToken: string;
  expiresIn: number;
}

export interface RefreshTokenResponse {
  token: string;
  refreshToken: string;
  expiresIn: number;
}

export interface PasswordResetResponse {
  message: string;
  success: boolean;
}

/**
 * Authentication Service Class
 */
class AuthenticationService {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: API_URL,
      timeout: API_TIMEOUT,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Response interceptor for error handling
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          // Server responded with error status
          const message = error.response.data?.message || error.response.data?.error || 'An error occurred';
          throw new Error(message);
        } else if (error.request) {
          // Request made but no response
          throw new Error('Network error. Please check your connection.');
        } else {
          // Error in request setup
          throw new Error(error.message || 'An unexpected error occurred');
        }
      }
    );
  }

  /**
   * Login with email and password
   * @param email User email
   * @param password User password
   * @returns Login response with user data and tokens
   */
  async login(email: string, password: string): Promise<LoginResponse> {
    try {
      const response = await this.axiosInstance.post<LoginResponse>('/auth/login', {
        email: email.toLowerCase().trim(),
        password,
      });

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Login error:', error.message);
      throw error;
    }
  }

  /**
   * Logout (invalidate token on server)
   * @param token Current access token
   */
  async logout(token: string): Promise<void> {
    try {
      await this.axiosInstance.post(
        '/auth/logout',
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    } catch (error: any) {
      console.error('[AuthService] Logout error:', error.message);
      // Don't throw - logout should succeed even if API call fails
    }
  }

  /**
   * Refresh access token
   * @param refreshToken Current refresh token
   * @returns New tokens
   */
  async refreshToken(refreshToken: string): Promise<RefreshTokenResponse> {
    try {
      const response = await this.axiosInstance.post<RefreshTokenResponse>('/auth/refresh', {
        refresh_token: refreshToken,
      });

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Refresh token error:', error.message);
      throw error;
    }
  }

  /**
   * Request password reset email
   * @param email User email
   * @returns Success response
   */
  async requestPasswordReset(email: string): Promise<PasswordResetResponse> {
    try {
      const response = await this.axiosInstance.post<PasswordResetResponse>(
        '/auth/forgot-password',
        {
          email: email.toLowerCase().trim(),
        }
      );

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Password reset request error:', error.message);
      throw error;
    }
  }

  /**
   * Verify user's current password
   * @param token Access token
   * @param password Password to verify
   * @returns True if password is correct
   */
  async verifyPassword(token: string, password: string): Promise<boolean> {
    try {
      const response = await this.axiosInstance.post(
        '/auth/verify-password',
        { password },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      return response.data.valid === true;
    } catch (error: any) {
      console.error('[AuthService] Password verification error:', error.message);
      return false;
    }
  }

  /**
   * Change user password
   * @param token Access token
   * @param currentPassword Current password
   * @param newPassword New password
   */
  async changePassword(
    token: string,
    currentPassword: string,
    newPassword: string
  ): Promise<void> {
    try {
      await this.axiosInstance.post(
        '/auth/change-password',
        {
          current_password: currentPassword,
          new_password: newPassword,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    } catch (error: any) {
      console.error('[AuthService] Password change error:', error.message);
      throw error;
    }
  }

  /**
   * Get current user profile
   * @param token Access token
   * @returns User data
   */
  async getCurrentUser(token: string): Promise<User> {
    try {
      const response = await this.axiosInstance.get<User>('/auth/me', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Get current user error:', error.message);
      throw error;
    }
  }

  /**
   * Update user profile
   * @param token Access token
   * @param userData Partial user data to update
   * @returns Updated user data
   */
  async updateProfile(
    token: string,
    userData: Partial<Pick<User, 'first_name' | 'last_name' | 'phone'>>
  ): Promise<User> {
    try {
      const response = await this.axiosInstance.patch<User>('/auth/profile', userData, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Update profile error:', error.message);
      throw error;
    }
  }

  /**
   * Validate token (check if still valid)
   * @param token Access token
   * @returns True if token is valid
   */
  async validateToken(token: string): Promise<boolean> {
    try {
      const response = await this.axiosInstance.get('/auth/validate', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      return response.data.valid === true;
    } catch (error: any) {
      console.error('[AuthService] Token validation error:', error.message);
      return false;
    }
  }

  /**
   * Register new user account
   * @param userData Registration data
   * @returns Login response with user data and tokens
   */
  async register(userData: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    phone?: string;
    company_name?: string;
  }): Promise<LoginResponse> {
    try {
      const response = await this.axiosInstance.post<LoginResponse>('/auth/register', {
        ...userData,
        email: userData.email.toLowerCase().trim(),
      });

      return response.data;
    } catch (error: any) {
      console.error('[AuthService] Registration error:', error.message);
      throw error;
    }
  }
}

// Export singleton instance
export const AuthService = new AuthenticationService();
