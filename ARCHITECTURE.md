# Remote Inspection Capture Tool - Repository Structure

```
.
├── README.md                      # Main documentation (this file)
├── openapi.yaml                   # REST API specification (OpenAPI 3.0)
├── ws-schema.json                 # WebSocket message schema (JSON Schema)
├── docker-compose.yaml            # Local development environment (DynamoDB Local + MinIO)
├── .gitignore                     # Git ignore rules
│
├── backend/                       # Go backend application
│   ├── go.mod                     # Go module definition
│   ├── go.sum                     # Go dependency checksums
│   ├── (environment variables are set via OS or .env not included)
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go            # REST API server (Lambda-ready)
│   │   ├── ws/
│   │   │   └── main.go            # WebSocket server (Lambda-ready)
│   │   └── (no migrate script included)
│   │
│   └── internal/
│       ├── config/
│       │   └── config.go          # Configuration management
│       │
│       ├── domain/                # Domain models (entities)
│       │   ├── models.go          # Entity definitions (User, Inspection, Issue, etc.)
│       │   ├── errors.go          # Error types and codes
│       │   ├── video.go           # VideoProvider interface definition
│       │   └── repositories.go    # Repository interfaces
│       │
│       ├── service/
│       │   └── services.go        # Business logic services
│       │       - InspectionService
│       │       - PhotoService
│       │       - CaptureRequestService
│       │       - TemplateService
│       │       - ViewerLinkService
│       │
│       ├── handler/
│       │   ├── inspection.go      # Inspection HTTP handlers (TODO)
│       │   ├── photo.go           # Photo/capture HTTP handlers (TODO)
│       │   └── public.go          # Public endpoint handlers (TODO)
│       │
│       ├── repository/
│       │   └── dynamodb.go        # DynamoDB adapter implementation
│       │       - InspectionRepositoryImpl
│       │       - ChecklistItemRepositoryImpl
│       │       - IssueRepositoryImpl
│       │       - CaptureRequestRepositoryImpl
│       │       - EvidencePhotoRepositoryImpl
│       │       - AnnotationTemplateRepositoryImpl
│       │       - ViewerLinkRepositoryImpl
│       │       - AuditLogRepositoryImpl
│       │
│       ├── middleware/
│       │   └── auth.go            # JWT authentication/authorization middleware
│       │
│       └── infrastructure/
│           ├── storage/
│           │   └── s3.go          # S3 operations (presigned URLs, etc.)
│           │
│           └── video/
│               └── mock.go        # Mock VideoProvider implementation
│
│
├── frontend/                      # React + TypeScript frontend application
│   ├── package.json               # NPM dependencies and scripts
│   ├── package-lock.json          # NPM lock file
│   ├── tsconfig.json              # TypeScript configuration
│   ├── vite.config.ts             # Vite configuration
│   ├── index.html                 # HTML entry point
│   ├── .env                       # Frontend environment variables
│   │
│   └── src/
│       ├── main.tsx               # React entry point
│       ├── App.tsx                # Main App component with routing
│       ├── index.css              # Global styles
│       │
│       ├── pages/                 # Page components
│       │   ├── LoginPage.tsx      # Login page (Cognito integration)
│       │   ├── InspectionListPage.tsx  # Inspection list
│       │   ├── RemoteInspectionPage.tsx # Remote inspection view (video + checklist + issues)
│       │   ├── FieldCameraPage.tsx      # Field camera/broadcaster
│       │   ├── AnnotatePhotoPage.tsx    # Photo annotation with Konva canvas
│       │   ├── AdminTemplatesPage.tsx   # Template management
│       │   └── GuestViewPage.tsx        # Guest viewer page (no login required)
│       │
│       ├── components/            # Reusable UI components
│       │   ├── VideoViewer.tsx    # IVS Real-Time video viewer
│       │   ├── CameraCapture.tsx  # Camera capture (ImageCapture API)
│       │   ├── AnnotationCanvas.tsx # Konva canvas for drawing
│       │   ├── ChecklistPanel.tsx  # Checklist UI
│       │   ├── IssuePanel.tsx      # Issues UI
│       │   └── PhotoTimeline.tsx   # Photo evidence timeline
│       │
│       ├── hooks/                 # Custom React hooks
│       │   ├── useInspection.ts   # Inspection data management
│       │   ├── useWebSocket.ts    # WebSocket connection management
│       │   ├── useAnnotation.ts   # Annotation state management
│       │   └── useAuth.ts         # Authentication state
│       │
│       ├── services/              # API and integration services
│       │   ├── apiClient.ts       # HTTP API client (axios-based)
│       │   └── wsClient.ts        # WebSocket client
│       │
│       ├── types/                 # TypeScript type definitions
│       │   └── api.ts             # All API response/request types and WebSocket messages
│       │
│       ├── context/               # React Context for global state
│       │   ├── AuthContext.tsx    # Authentication context
│       │   └── InspectionContext.tsx # Inspection data context
│       │
│       └── utils/                 # Utility functions
│           ├── validators.ts      # Form validation helpers
│           ├── formatters.ts      # Date/time and data formatters
│           └── helpers.ts         # General utility functions
│
│
└── infra/                         # AWS CDK infrastructure
    ├── package.json               # CDK npm dependencies
    ├── tsconfig.json              # TypeScript configuration
    ├── cdk.json                   # CDK configuration
    ├── bin/
    │   └── app.ts                 # CDK app entrypoint
    │
    └── lib/
        └── inspection-stack.ts    # Main CDK stack definition
            - DynamoDB tables (Inspections, ChecklistItems, Issues, CaptureRequests, EvidencePhotos, AnnotationTemplates, ViewerLinks)
            - S3 buckets (Evidence, UI)
            - CloudFront distribution for UI
            - Cognito user pool
            - API Gateway (HTTP API + WebSocket) - future
            - Lambda functions - future
```

## Key Files and Their Purposes

### Backend

- **internal/domain/models.go**: Core entity definitions (User, Inspection, ChecklistItem, Issue, CaptureRequest, EvidencePhoto, AnnotationTemplate, ViewerLink)
- **internal/domain/repositories.go**: Interface definitions for data access layer
- **internal/service/services.go**: Business logic implementations
- **internal/middleware/auth.go**: JWT validation and authorization
- **internal/repository/dynamodb.go**: DynamoDB adapter (repository pattern)
- **internal/infrastructure/video/mock.go**: Mock video provider for local development

### Frontend

- **src/types/api.ts**: Complete TypeScript definitions for all API responses and WebSocket messages
- **src/services/apiClient.ts**: HTTP API client with all endpoints
- **src/services/wsClient.ts**: WebSocket client for real-time messaging
- **src/pages/*.tsx**: Page-level components with business logic

### Infrastructure

- **infra/lib/inspection-stack.ts**: AWS CDK stack defining all cloud resources

## Database Schema (DynamoDB)

| Table Name          | Partition Key     | Sort Key         | Purpose                        |
| ------------------- | ----------------- | ---------------- | ------------------------------ |
| Inspections         | inspectionId      | -                | Inspection projects            |
| ChecklistItems      | inspectionId      | itemId           | Checklist items per inspection |
| Issues              | inspectionId      | issueId          | Issues per inspection          |
| CaptureRequests     | inspectionId      | captureRequestId | Capture request history        |
| EvidencePhotos      | inspectionId      | photoId          | Photo evidence per inspection  |
| AnnotationTemplates | templateId        | version          | Annotation templates           |
| ViewerLinks         | viewerAccessToken | -                | Guest access tokens (with TTL) |

## Development Workflow

1. **Setup**: `docker-compose up -d` (DynamoDB + MinIO)
2. **Backend**: `cd backend && go run ./cmd/api/main.go`
3. **Frontend**: `cd frontend && npm run dev`
4. **Access**: http://localhost:5173

## Deployment

1. **AWS CDK**: `cd infra && cdk deploy`
2. **Backend**: Build and push Lambda functions
3. **Frontend**: Build and upload to S3 + CloudFront

## TODO Items

- [ ] Complete HTTP handler implementations
- [ ] WebSocket server implementation (gorilla/websocket)
- [ ] IVS Real-Time video provider implementation
- [ ] Image processing for annotation (drawing-to-image conversion)
- [ ] PDF report generation
- [ ] Unit tests and integration tests
- [ ] E2E testing setup
- [x] Unit tests and integration tests (backend + frontend unit tests added)
- [x] E2E testing setup (Playwright skeleton added under `e2e/`)
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Monitoring and alerting (CloudWatch)
- [ ] More template types and customization
