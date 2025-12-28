import axios from 'axios';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiClient } from '../../services/apiClient';

vi.mock('axios');
const mockedAxios = axios as any;

describe('ApiClient', () => {
    let apiClient: ApiClient;

    beforeEach(() => {
        apiClient = new ApiClient('http://localhost:8080');
        vi.clearAllMocks();
    });

    describe('setToken', () => {
        it('should set the authorization token', () => {
            const token = 'test-jwt-token';
            apiClient.setToken(token);
            expect(apiClient['token']).toBe(token);
        });
    });

    describe('clearToken', () => {
        it('should clear the authorization token', () => {
            apiClient.setToken('test-token');
            apiClient.clearToken();
            expect(apiClient['token']).toBe('');
        });
    });

    describe('getMe', () => {
        it('should fetch current user', async () => {
            const mockUser = {
                code: 'SUCCESS',
                data: {
                    userId: 'user1',
                    username: 'testuser',
                    email: 'test@example.com',
                    role: 'editor',
                },
            };

            mockedAxios.get.mockResolvedValue({ data: mockUser });

            const user = await apiClient.getMe();
            expect(user).toEqual(mockUser.data);
            expect(mockedAxios.get).toHaveBeenCalledWith('/api/me');
        });

        it('should handle error when fetching user', async () => {
            const mockError = new Error('Network error');
            mockedAxios.get.mockRejectedValue(mockError);

            await expect(apiClient.getMe()).rejects.toThrow('Network error');
        });
    });

    describe('createInspection', () => {
        it('should create a new inspection', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    inspectionId: 'insp-123',
                    title: 'Test Inspection',
                    place: 'Tokyo',
                    status: 'draft',
                },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.createInspection('Test Inspection', 'Tokyo');
            expect(result.inspectionId).toBe('insp-123');
            expect(mockedAxios.post).toHaveBeenCalledWith('/api/inspections', {
                title: 'Test Inspection',
                place: 'Tokyo',
            });
        });
    });

    describe('getInspection', () => {
        it('should fetch an inspection by ID', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    inspectionId: 'insp-123',
                    title: 'Test Inspection',
                    status: 'draft',
                },
            };

            mockedAxios.get.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.getInspection('insp-123');
            expect(result.inspectionId).toBe('insp-123');
            expect(mockedAxios.get).toHaveBeenCalledWith('/api/inspections/insp-123');
        });
    });

    describe('listInspections', () => {
        it('should list all inspections', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: [
                    { inspectionId: 'insp-1', title: 'Test 1' },
                    { inspectionId: 'insp-2', title: 'Test 2' },
                ],
            };

            mockedAxios.get.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.listInspections(10);
            expect(result).toHaveLength(2);
            expect(mockedAxios.get).toHaveBeenCalledWith('/api/inspections', {
                params: { limit: 10 },
            });
        });
    });

    describe('startInspection', () => {
        it('should start an inspection', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: { inspectionId: 'insp-123', status: 'running' },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.startInspection('insp-123');
            expect(result.status).toBe('running');
            expect(mockedAxios.post).toHaveBeenCalledWith(
                '/api/inspections/insp-123/start',
                {}
            );
        });
    });

    describe('createCaptureRequest', () => {
        it('should create a capture request', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    captureRequestId: 'cap-123',
                    inspectionId: 'insp-123',
                    state: 'requested',
                },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.createCaptureRequest('insp-123', 'device1');
            expect(result.captureRequestId).toBe('cap-123');
            expect(result.state).toBe('requested');
        });
    });

    describe('generateUploadPresignedURL', () => {
        it('should generate a presigned upload URL', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    uploadUrl: 'https://s3.amazonaws.com/bucket/object',
                    objectKey: 'photos/photo1.jpg',
                },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.generateUploadPresignedURL(
                'insp-123',
                'photo1.jpg'
            );
            expect(result.uploadUrl).toContain('s3.amazonaws.com');
        });
    });

    describe('createViewerLink', () => {
        it('should create a viewer link', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    viewerAccessToken: 'viewer-token-123',
                    inspectionId: 'insp-123',
                },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.createViewerLink('insp-123');
            expect(result.viewerAccessToken).toBe('viewer-token-123');
        });
    });

    describe('createGuestSession', () => {
        it('should create a guest session', async () => {
            const mockResponse = {
                code: 'SUCCESS',
                data: {
                    guestSessionToken: 'guest-123',
                    inspectionId: 'insp-123',
                },
            };

            mockedAxios.post.mockResolvedValue({ data: mockResponse });

            const result = await apiClient.createGuestSession('viewer-token-123');
            expect(result.guestSessionToken).toBe('guest-123');
        });
    });

    describe('request error handling', () => {
        it('should handle 401 unauthorized error', async () => {
            const mockError = {
                response: {
                    status: 401,
                    data: { code: 'ERR_UNAUTHORIZED', message: 'Token invalid' },
                },
            };

            mockedAxios.get.mockRejectedValue(mockError);

            await expect(apiClient.getMe()).rejects.toThrow();
        });

        it('should handle 403 forbidden error', async () => {
            const mockError = {
                response: {
                    status: 403,
                    data: { code: 'ERR_FORBIDDEN', message: 'Access denied' },
                },
            };

            mockedAxios.get.mockRejectedValue(mockError);

            await expect(apiClient.getMe()).rejects.toThrow();
        });

        it('should handle 500 server error', async () => {
            const mockError = {
                response: {
                    status: 500,
                    data: { code: 'ERR_INTERNAL_ERROR', message: 'Internal server error' },
                },
            };

            mockedAxios.get.mockRejectedValue(mockError);

            await expect(apiClient.getMe()).rejects.toThrow();
        });
    });
});
