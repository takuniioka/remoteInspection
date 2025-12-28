/**
 * API client wrapper
 *
 * This module provides a small wrapper around `axios` that centralizes
 * API requests, handles the JSON envelope used by the backend (`ApiResponse`),
 * and injects an Authorization header when a token is set.
 */

import {
    AnnotationTemplate,
    ApiErrorResponse,
    ApiResponse,
    CaptureRequest,
    ChecklistItem,
    EvidencePhoto,
    GuestSessionToken,
    Inspection,
    Issue,
    PresignedURLResponse,
    User,
    VideoToken,
    ViewerLink,
} from '@/types/api'
import axios, {
    AxiosInstance,
    AxiosRequestConfig
} from 'axios'

class ApiClient {
    private client: AxiosInstance
    private token: string | null = null

    /**
     * Create a new ApiClient.
     * @param baseURL base URL for API requests (e.g., VITE_API_BASE_URL)
     */
    constructor(baseURL: string) {
        this.client = axios.create({
            baseURL,
            timeout: 30000,
            headers: {
                'Content-Type': 'application/json',
            },
        })

        // Add request interceptor for auth token
        this.client.interceptors.request.use((config) => {
            if (this.token) {
                config.headers.Authorization = `Bearer ${this.token}`
            }
            return config
        })
    }

    /** Set the bearer token used for authenticated requests. */
    setToken(token: string) {
        this.token = token
    }

    /** Clear the stored authentication token. */
    clearToken() {
        this.token = null
    }

    /**
     * Internal request helper that unwraps the backend's `ApiResponse` envelope
     * and throws a JavaScript Error with the backend message on failure.
     */
    private async request<T>(
        method: string,
        url: string,
        data?: any,
        config?: AxiosRequestConfig
    ): Promise<T> {
        try {
            const response = await this.client({
                method,
                url,
                data,
                ...config,
            })
            const apiResponse: ApiResponse<T> = response.data
            return apiResponse.data
        } catch (error) {
            if (axios.isAxiosError(error)) {
                const apiError = error.response?.data as ApiErrorResponse
                throw new Error(apiError?.message || error.message)
            }
            throw error
        }
    }

    // Auth endpoints
    async getMe(): Promise<User> {
        return this.request<User>('GET', '/me')
    }

    // Inspection endpoints
    async createInspection(title: string, place: string): Promise<Inspection> {
        return this.request<Inspection>('POST', '/inspections', {
            title,
            place,
        })
    }

    async getInspection(inspectionId: string): Promise<Inspection> {
        return this.request<Inspection>(
            'GET',
            `/inspections/${inspectionId}`
        )
    }

    async listInspections(): Promise<Inspection[]> {
        return this.request<Inspection[]>('GET', '/inspections')
    }

    async startInspection(inspectionId: string): Promise<Inspection> {
        return this.request<Inspection>(
            'POST',
            `/inspections/${inspectionId}/start`
        )
    }

    async endInspection(inspectionId: string): Promise<Inspection> {
        return this.request<Inspection>(
            'POST',
            `/inspections/${inspectionId}/end`
        )
    }

    // Checklist endpoints
    async getChecklist(inspectionId: string): Promise<ChecklistItem[]> {
        return this.request<ChecklistItem[]>(
            'GET',
            `/inspections/${inspectionId}/checklist`
        )
    }

    async updateChecklistItem(
        inspectionId: string,
        itemId: string,
        result: string,
        comment: string
    ): Promise<ChecklistItem> {
        return this.request<ChecklistItem>(
            'PATCH',
            `/inspections/${inspectionId}/checklist/${itemId}`,
            {
                result,
                comment,
            }
        )
    }

    // Issue endpoints
    async getIssues(inspectionId: string): Promise<Issue[]> {
        return this.request<Issue[]>(
            'GET',
            `/inspections/${inspectionId}/issues`
        )
    }

    async createIssue(
        inspectionId: string,
        title: string,
        detail: string,
        linkedType?: string,
        linkedId?: string
    ): Promise<Issue> {
        return this.request<Issue>(
            'POST',
            `/inspections/${inspectionId}/issues`,
            {
                title,
                detail,
                linkedType,
                linkedId,
            }
        )
    }

    async updateIssue(
        inspectionId: string,
        issueId: string,
        status?: string,
        assignee?: string,
        detail?: string
    ): Promise<Issue> {
        return this.request<Issue>(
            'PATCH',
            `/inspections/${inspectionId}/issues/${issueId}`,
            {
                status,
                assignee,
                detail,
            }
        )
    }

    // Video endpoints
    async getRealtimeToken(inspectionId: string): Promise<VideoToken> {
        return this.request<VideoToken>(
            'POST',
            `/inspections/${inspectionId}/realtime/token`
        )
    }

    // Capture endpoints
    async createCaptureRequest(
        inspectionId: string,
        linkedType?: string,
        linkedId?: string
    ): Promise<CaptureRequest> {
        return this.request<CaptureRequest>(
            'POST',
            `/inspections/${inspectionId}/capture-requests`,
            {
                linkedType,
                linkedId,
            }
        )
    }

    // Photo endpoints
    async generateUploadPresignedURL(
        inspectionId: string,
        captureRequestId: string,
        fileName: string,
        contentType: string
    ): Promise<PresignedURLResponse> {
        return this.request<PresignedURLResponse>(
            'POST',
            `/inspections/${inspectionId}/photos/presign`,
            {
                captureRequestId,
                fileName,
                contentType,
            }
        )
    }

    async completePhotoUpload(
        inspectionId: string,
        photoId: string,
        captureRequestId: string
    ): Promise<EvidencePhoto> {
        return this.request<EvidencePhoto>(
            'POST',
            `/inspections/${inspectionId}/photos/${photoId}/complete`,
            {
                captureRequestId,
            }
        )
    }

    // Annotation endpoints
    async generateAnnotationPresignedURL(
        inspectionId: string,
        photoId: string,
        fileName: string,
        contentType: string
    ): Promise<PresignedURLResponse> {
        return this.request<PresignedURLResponse>(
            'POST',
            `/inspections/${inspectionId}/photos/${photoId}/annotation/presign`,
            {
                fileName,
                contentType,
            }
        )
    }

    async completeAnnotation(
        inspectionId: string,
        photoId: string,
        s3AnnotatedKey: string,
        annotationJson: any
    ): Promise<EvidencePhoto> {
        return this.request<EvidencePhoto>(
            'POST',
            `/inspections/${inspectionId}/photos/${photoId}/annotation/complete`,
            {
                s3AnnotatedKey,
                annotationJson,
            }
        )
    }

    // Template endpoints
    async getTemplates(): Promise<AnnotationTemplate[]> {
        return this.request<AnnotationTemplate[]>('GET', '/annotation-templates')
    }

    async createTemplate(
        name: string,
        definitionJson: any,
        description?: string
    ): Promise<AnnotationTemplate> {
        return this.request<AnnotationTemplate>(
            'POST',
            '/annotation-templates',
            {
                name,
                description,
                definitionJson,
            }
        )
    }

    async updateTemplate(
        templateId: string,
        name?: string,
        description?: string,
        definitionJson?: any,
        isActive?: boolean
    ): Promise<AnnotationTemplate> {
        return this.request<AnnotationTemplate>(
            'PUT',
            `/annotation-templates/${templateId}`,
            {
                name,
                description,
                definitionJson,
                isActive,
            }
        )
    }

    // Viewer Link endpoints
    async createViewerLink(
        inspectionId: string,
        expiresInDays?: number,
        maxUses?: number
    ): Promise<ViewerLink> {
        return this.request<ViewerLink>(
            'POST',
            `/inspections/${inspectionId}/viewer-links`,
            {
                expiresInDays: expiresInDays || 7,
                maxUses,
            }
        )
    }

    // Public endpoints
    async createGuestSession(
        viewerAccessToken: string
    ): Promise<GuestSessionToken> {
        return this.request<GuestSessionToken>(
            'POST',
            '/public/viewer-sessions',
            {
                viewerAccessToken,
            }
        )
    }

    async getPublicInspection(inspectionId: string): Promise<Inspection> {
        return this.request<Inspection>(
            'GET',
            `/public/inspections/${inspectionId}`
        )
    }

    async getPublicPhotos(inspectionId: string): Promise<EvidencePhoto[]> {
        return this.request<EvidencePhoto[]>(
            'GET',
            `/public/inspections/${inspectionId}/photos`
        )
    }
}

const apiClient = new ApiClient(
    import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
)

export default apiClient
