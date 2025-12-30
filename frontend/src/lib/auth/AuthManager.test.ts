/**
 * AuthManager Tests
 *
 * These tests verify the AuthManager singleton implementation
 * and its core authentication functionality.
 */

import { describe, it, expect, beforeEach, afterEach, vi, beforeAll } from 'vitest';
import { AuthManager } from './AuthManager';
import type { AuthResponse } from '$lib/types';

// Setup localStorage mock before all tests
beforeAll(() => {
    const localStorageMock = (() => {
        let store: Record<string, string> = {};

        return {
            getItem: (key: string) => store[key] || null,
            setItem: (key: string, value: string) => {
                store[key] = value;
            },
            removeItem: (key: string) => {
                delete store[key];
            },
            clear: () => {
                store = {};
            },
            get length() {
                return Object.keys(store).length;
            },
            key: (index: number) => Object.keys(store)[index] || null
        };
    })();

    Object.defineProperty(global, 'localStorage', {
        value: localStorageMock,
        writable: true
    });
});

describe('AuthManager', () => {
    let authManager: AuthManager;
    let mockApiClient: any;

    beforeEach(() => {
        // Clear localStorage
        if (typeof window !== 'undefined' && window.localStorage) {
            window.localStorage.clear();
        }

        // Get fresh instance
        authManager = AuthManager.getInstance();

        // Clear any existing state
        authManager.clearAllState();

        // Create mock API client
        mockApiClient = {
            login: vi.fn(),
            logout: vi.fn(),
            post: vi.fn(),
            get: vi.fn()
        };

        // Set mock API client
        authManager.setApiClient(mockApiClient);
    });

    afterEach(() => {
        vi.clearAllMocks();
        authManager.clearAllState();
    });

    describe('Singleton Pattern', () => {
        it('should return the same instance', () => {
            const instance1 = AuthManager.getInstance();
            const instance2 = AuthManager.getInstance();

            expect(instance1).toBe(instance2);
        });
    });

    describe('Login', () => {
        // Create a valid JWT token (expires in 1 hour)
        const futureExp = Math.floor(Date.now() / 1000) + 3600;
        const validAccessToken = `header.${btoa(JSON.stringify({ exp: futureExp, user_id: 'user-123' }))}.signature`;

        const mockAuthResponse: AuthResponse = {
            access_token: validAccessToken,
            refresh_token: 'mock-refresh-token',
            token_type: 'Bearer',
            expires_in: 3600,
            user: {
                id: 'user-123',
                tenant_id: 'tenant-456',
                email: 'test@example.com',
                first_name: 'Test',
                last_name: 'User',
                status: 'active',
                email_verified: true,
                mfa_enabled: false,
                roles: ['agent'],
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z'
            },
            roles: ['agent']
        };

        it('should login successfully and store all state atomically', async () => {
            // Mock API responses
            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockResolvedValue({
                data: { permissions: ['read:clients', 'write:registrations'] }
            });

            const result = await authManager.login('test@example.com', 'password123');

            expect(result.success).toBe(true);
            expect(authManager.isAuthenticated()).toBe(true);
            expect(authManager.getUser()).toEqual(mockAuthResponse.user);
            expect(authManager.getAccessToken()).toBe(validAccessToken);
            expect(authManager.getTenantId()).toBe('tenant-456');
            expect(authManager.getPermissions()).toEqual(['read:clients', 'write:registrations']);
        });

        it('should clear all state on login failure (atomic failure)', async () => {
            mockApiClient.login.mockRejectedValue(new Error('Invalid credentials'));

            const result = await authManager.login('test@example.com', 'wrongpassword');

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
            expect(authManager.isAuthenticated()).toBe(false);
            expect(authManager.getUser()).toBeNull();
            expect(authManager.getAccessToken()).toBeNull();
        });

        it('should load permissions immediately after login', async () => {
            const mockPermissions = ['read:clients', 'write:registrations', 'delete:documents'];
            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockResolvedValue({
                data: { permissions: mockPermissions }
            });

            await authManager.login('test@example.com', 'password123');

            expect(mockApiClient.get).toHaveBeenCalledWith(
                expect.stringContaining('/admin/users/user-123/effective-permissions')
            );
            expect(authManager.getPermissions()).toEqual(mockPermissions);
        });

        it('should continue login even if permission loading fails', async () => {
            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockRejectedValue(new Error('Permissions API failed'));

            const result = await authManager.login('test@example.com', 'password123');

            expect(result.success).toBe(true);
            expect(authManager.getPermissions()).toEqual([]);
        });
    });

    describe('Logout', () => {
        beforeEach(async () => {
            // Setup logged in state
            const mockAuthResponse: AuthResponse = {
                access_token: 'mock-access-token',
                refresh_token: 'mock-refresh-token',
                token_type: 'Bearer',
                expires_in: 3600,
                user: {
                    id: 'user-123',
                    tenant_id: 'tenant-456',
                    email: 'test@example.com',
                    status: 'active',
                    email_verified: true,
                    mfa_enabled: false,
                    created_at: '2024-01-01T00:00:00Z',
                    updated_at: '2024-01-01T00:00:00Z'
                },
                roles: []
            };

            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockResolvedValue({ data: { permissions: ['read:clients'] } });

            await authManager.login('test@example.com', 'password123');
        });

        it('should clear all auth state on logout', async () => {
            await authManager.logout();

            expect(authManager.isAuthenticated()).toBe(false);
            expect(authManager.getUser()).toBeNull();
            expect(authManager.getAccessToken()).toBeNull();
            expect(authManager.getRefreshToken()).toBeNull();
            expect(authManager.getTenantId()).toBeNull();
            expect(authManager.getPermissions()).toEqual([]);
        });

        it('should call backend logout endpoint', async () => {
            mockApiClient.logout.mockResolvedValue({});

            await authManager.logout();

            expect(mockApiClient.logout).toHaveBeenCalled();
        });

        it('should still clear state even if backend logout fails', async () => {
            mockApiClient.logout.mockRejectedValue(new Error('Network error'));

            await authManager.logout();

            expect(authManager.isAuthenticated()).toBe(false);
            expect(authManager.getUser()).toBeNull();
        });
    });

    describe('Token Management', () => {
        it('should validate token expiry correctly', () => {
            // Create a valid token (expires in 1 hour)
            const futureExp = Math.floor(Date.now() / 1000) + 3600;
            const validToken = `header.${btoa(JSON.stringify({ exp: futureExp }))}.signature`;

            expect(authManager.isTokenExpired(validToken)).toBe(false);
        });

        it('should detect expired tokens', () => {
            // Create an expired token
            const pastExp = Math.floor(Date.now() / 1000) - 3600;
            const expiredToken = `header.${btoa(JSON.stringify({ exp: pastExp }))}.signature`;

            expect(authManager.isTokenExpired(expiredToken)).toBe(true);
        });

        it('should return null for expired access token', () => {
            const pastExp = Math.floor(Date.now() / 1000) - 3600;
            const expiredToken = `header.${btoa(JSON.stringify({ exp: pastExp }))}.signature`;

            authManager.setTokens(expiredToken, 'refresh-token');

            expect(authManager.getAccessToken()).toBeNull();
        });

        it('should refresh access token successfully', async () => {
            // Set initial refresh token
            authManager.setTokens('old-access-token', 'refresh-token');

            // Mock refresh response
            mockApiClient.post.mockResolvedValue({
                data: {
                    access_token: 'new-access-token',
                    refresh_token: 'new-refresh-token'
                }
            });

            const result = await authManager.refreshAccessToken();

            expect(result).toBe(true);
            expect(mockApiClient.post).toHaveBeenCalledWith('/auth/refresh', {
                refresh_token: 'refresh-token'
            });
        });
    });

    describe('Permission Management', () => {
        it('should check permissions correctly', async () => {
            const mockAuthResponse: AuthResponse = {
                access_token: 'token',
                refresh_token: 'refresh',
                token_type: 'Bearer',
                expires_in: 3600,
                user: {
                    id: 'user-123',
                    tenant_id: 'tenant-456',
                    email: 'test@example.com',
                    status: 'active',
                    email_verified: true,
                    mfa_enabled: false,
                    created_at: '2024-01-01T00:00:00Z',
                    updated_at: '2024-01-01T00:00:00Z'
                },
                roles: []
            };

            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockResolvedValue({
                data: { permissions: ['read:clients', 'write:documents'] }
            });

            await authManager.login('test@example.com', 'password123');

            expect(authManager.hasPermission('read:clients')).toBe(true);
            expect(authManager.hasPermission('delete:users')).toBe(false);
        });

        it('should return immutable copy of permissions', () => {
            const permissions = authManager.getPermissions();
            permissions.push('new:permission');

            expect(authManager.getPermissions()).not.toContain('new:permission');
        });
    });

    describe('State Getters', () => {
        it('should return correct authentication status', () => {
            expect(authManager.isAuthenticated()).toBe(false);

            authManager.setTokens('access-token', 'refresh-token');
            // Still false because user is null
            expect(authManager.isAuthenticated()).toBe(false);
        });

        it('should return tenant ID correctly', async () => {
            const mockAuthResponse: AuthResponse = {
                access_token: 'token',
                refresh_token: 'refresh',
                token_type: 'Bearer',
                expires_in: 3600,
                user: {
                    id: 'user-123',
                    tenant_id: 'tenant-789',
                    email: 'test@example.com',
                    status: 'active',
                    email_verified: true,
                    mfa_enabled: false,
                    created_at: '2024-01-01T00:00:00Z',
                    updated_at: '2024-01-01T00:00:00Z'
                },
                roles: []
            };

            mockApiClient.login.mockResolvedValue(mockAuthResponse);
            mockApiClient.get.mockResolvedValue({ data: { permissions: [] } });

            await authManager.login('test@example.com', 'password123');

            expect(authManager.getTenantId()).toBe('tenant-789');
        });
    });
});
