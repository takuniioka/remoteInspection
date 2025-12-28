// API Types (generated from OpenAPI)
export interface User {
    userId: string
    username: string
    email: string
    role: 'viewer' | 'capturer' | 'editor' | 'admin' | 'guestViewer'
    createdAt: string
}

export interface Inspection {
    inspectionId: string
    title: string
    place: string
    status: 'draft' | 'running' | 'closed'
    realtimeRef?: string
    createdAt: string
    startedAt?: string | null
    endedAt?: string | null
    createdBy: string
}

export interface ChecklistItem {
    itemId: string
    inspectionId: string
    label: string
    result: 'ok' | 'ng' | 'pending'
    comment: string
    updatedAt: string
    updatedBy: string
}

export interface Issue {
    issueId: string
    inspectionId: string
    title: string
    detail: string
    status: 'open' | 'resolved' | 'closed'
    assignee?: string | null
    createdAt: string
    createdBy: string
    updatedAt: string
    updatedBy: string
}

export interface CaptureRequest {
    inspectionId: string
    captureRequestId: string
    requestedAt: string
    requestedBy: string
    state: 'requested' | 'capturing' | 'uploaded' | 'failed' | 'timeout'
    capturedAt?: string | null
    uploadedAt?: string | null
    photoId?: string | null
    linkedType?: 'checklistItem' | 'issue' | null
    linkedId?: string | null
    failReason?: string | null
}

export interface EvidencePhoto {
    inspectionId: string
    photoId: string
    captureRequestId?: string | null
    s3OriginalKey: string
    s3AnnotatedKey?: string | null
    s3AnnotationJsonKey?: string | null
    capturedAt: string
    uploadedAt: string
    annotatedAt?: string | null
    linkedType?: 'checklistItem' | 'issue' | null
    linkedId?: string | null
    note: string
}

export interface AnnotationTemplate {
    templateId: string
    version: number
    name: string
    description: string
    isActive: boolean
    definitionJson: Record<string, any>
    createdAt: string
}

export interface ViewerLink {
    viewerAccessToken: string
    inspectionId: string
    expiresAt: string
    isRevoked: boolean
    maxUses?: number | null
    usedCount: number
    createdAt: string
    createdBy: string
}

export interface GuestSessionToken {
    guestSessionToken: string
    expiresIn: number
}

export interface VideoToken {
    token: string
    expiresAt: string
    participantId: string
}

export interface PresignedURLResponse {
    url: string
    expires: number
}

// API Response format
export interface ApiResponse<T> {
    code: string
    message: string
    data: T
}

export interface ApiErrorResponse {
    code: string
    message: string
    details?: Record<string, any>
}

// WebSocket message types
export type WSMessageType =
    | 'capture.request'
    | 'capture.ack'
    | 'capture.uploaded'
    | 'capture.failed'
    | 'capture.timeout'
    | 'checklist.updated'
    | 'issue.updated'
    | 'photo.annotation.updated'
    | 'connection.established'
    | 'error'

export interface WSMessage {
    type: WSMessageType
    [key: string]: any
}

export interface CaptureRequestMessage extends WSMessage {
    type: 'capture.request'
    inspectionId: string
    requestedBy: string
    linkedType?: 'checklistItem' | 'issue'
    linkedId?: string
}

export interface CaptureAckMessage extends WSMessage {
    type: 'capture.ack'
    inspectionId: string
    captureRequestId: string
    requestedAt: string
    requestedBy: string
    state: 'requested'
}

export interface CaptureUploadedMessage extends WSMessage {
    type: 'capture.uploaded'
    inspectionId: string
    captureRequestId: string
    photoId: string
    uploadedAt: string
    state: 'uploaded'
    linkedType?: 'checklistItem' | 'issue'
    linkedId?: string
}

export interface CaptureFailedMessage extends WSMessage {
    type: 'capture.failed'
    inspectionId: string
    captureRequestId: string
    state: 'failed'
    failReason: string
}

export interface CaptureTimeoutMessage extends WSMessage {
    type: 'capture.timeout'
    inspectionId: string
    captureRequestId: string
    state: 'timeout'
}

export interface ChecklistUpdatedMessage extends WSMessage {
    type: 'checklist.updated'
    inspectionId: string
    itemId: string
    label: string
    result: 'ok' | 'ng' | 'pending'
    comment: string
    updatedAt: string
    updatedBy: string
}

export interface IssueUpdatedMessage extends WSMessage {
    type: 'issue.updated'
    inspectionId: string
    issueId: string
    title: string
    detail: string
    status: 'open' | 'resolved' | 'closed'
    assignee?: string | null
    createdAt: string
    updatedAt: string
    updatedBy: string
}

export interface PhotoAnnotationUpdatedMessage extends WSMessage {
    type: 'photo.annotation.updated'
    inspectionId: string
    photoId: string
    s3AnnotatedKey: string
    annotatedAt: string
    updatedBy: string
}

export interface ConnectionEstablishedMessage extends WSMessage {
    type: 'connection.established'
    connectionId: string
    serverTime: string
}

export interface ErrorMessage extends WSMessage {
    type: 'error'
    code: string
    message: string
    details?: Record<string, any>
}
