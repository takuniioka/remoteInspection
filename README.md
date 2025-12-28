# Remote Inspection Capture Tool - MVP

高速リモート検収ツール。現場がカメラ映像を配信し、リモート（最大5人）が1秒未満で視聴しながら検収。リモートシャッターで現場が高画質撮影→自動紐付け。

**最優先**：導入のしやすさ（ブラウザのみ）+ リモートシャッター→高画質撮影→自動紐付けの動作

## 特徴

- ✅ **ブラウザのみで動作** - インストール不要
- ✅ **リモートシャッター** - 最大5人のリモートから撮影要求、現場は自動（またはワンタップ）で撮影
- ✅ **WebRTC低遅延配信** - Amazon IVS Real-Time ベース（モック対応で即座に検証可能）
- ✅ **権限制御** - viewer/capturer/editor/admin + guestViewer（招待リンク）
- ✅ **注釈・看板テンプレ** - Konva.js で画像にお絵かき、看板テンプレ入力→焼き込み保存
- ✅ **報告書生成** - HTML→S3（PDF は設計・TODO）
- ✅ **監査ログ** - 全アクション記録

## システムアーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│ Client Layer (React + TypeScript + Vite + Konva)                │
├─────────────────────────────────────────────────────────────────┤
│ ・/login - Cognito ログイン                                      │
│ ・/inspections/:id/remote - リモート検収画面                    │
│ ・/inspections/:id/field - 現場配信・撮影画面                  │
│ ・/inspections/:id/view - ゲスト閲覧画面                        │
│ ・/photos/:photoId/annotate - 注釈編集                          │
│ ・/admin/templates - テンプレ管理                               │
└─────────────────────────────────────────────────────────────────┘
         ↕ HTTP REST API + WebSocket
┌─────────────────────────────────────────────────────────────────┐
│ API Gateway + Lambda (Go)                                        │
│ ・HTTP API: /api/* (REST)                                        │
│ ・WebSocket: /ws (inspectionId ルーム)                           │
└─────────────────────────────────────────────────────────────────┘
         ↕
┌─────────────────────────────────────────────────────────────────┐
│ Backend Services (Go 1.22+)                                      │
├─────────────────────────────────────────────────────────────────┤
│ Domain Layer:                                                    │
│ ・Inspection, ChecklistItem, Issue, CaptureRequest, Photo, etc. │
│                                                                  │
│ Service Layer:                                                   │
│ ・InspectionService, PhotoService, AnnotationService, etc.      │
│                                                                  │
│ Repository Layer:                                                │
│ ・DynamoDB Adapter                                               │
│ ・S3 Adapter (presigned URLs)                                    │
│ ・Cognito Adapter                                                │
│                                                                  │
│ Infrastructure:                                                  │
│ ・VideoProvider (IVS Real-Time + Mock)                           │
│ ・Storage (S3)                                                   │
│ ・Image Processor (Annotated 生成)                               │
└─────────────────────────────────────────────────────────────────┘
         ↕
┌─────────────────────────────────────────────────────────────────┐
│ Data & Storage                                                   │
├─────────────────────────────────────────────────────────────────┤
│ ・DynamoDB (Inspections, Photos, Issues, etc.)                   │
│ ・S3 (Original images, Annotated images, Annotation JSON, HTML)  │
│ ・Cognito (User auth, Sessions)                                  │
│ ・CloudWatch Logs (Audit logs)                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

### 前提条件

- Node.js 18+ (フロント)
- Go 1.22+ (バック)
- Docker / Docker Compose (ローカルDB/S3)
- AWS CLI (本番デプロイ)

### ローカル開発セットアップ

#### 1. リポジトリクローン

```bash
cd e:\新岡屋プロダクト\リモート検収＋報告書出力ツール
```

#### 2. 環境変数設定

**バックエンド** (`backend/.env`):
```env
# ローカル開発用
AWS_REGION=ap-northeast-1
DYNAMODB_ENDPOINT=http://localhost:8000
S3_ENDPOINT=http://localhost:9000
S3_BUCKET_NAME=inspection-evidence
S3_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin

# Cognito (ローカルでも同じ)
COGNITO_REGION=ap-northeast-1
COGNITO_USER_POOL_ID=ap-northeast-1_xxxxxx
COGNITO_CLIENT_ID=xxxxx

# Video
VIDEO_PROVIDER=mock  # 'ivs' for production

# JWT signing
JWT_SECRET=your-secret-key-change-in-production

# API
API_PORT=8080
WS_PORT=8081
```

**フロント** (`frontend/.env`):
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8081
VITE_COGNITO_DOMAIN=https://your-cognito-domain.auth.ap-northeast-1.amazoncognito.com
VITE_COGNITO_CLIENT_ID=xxxxx
VITE_COGNITO_REDIRECT_URI=http://localhost:5173/login/callback
```

#### 3. ローカルDB/S3 起動

```bash
docker-compose up -d
```

このコマンドで以下が起動：
- DynamoDB Local (`localhost:8000`)
- MinIO (S3互換) (`localhost:9000`)

#### 4. バックエンド - テーブル初期化 & 起動

```bash
cd backend

# DynamoDB Local テーブル作成
go run ./cmd/migrate/main.go

# REST API サーバ起動
go run ./cmd/api/main.go

# WebSocket サーバ起動（別ターミナル）
go run ./cmd/ws/main.go
```

#### 5. フロント起動

```bash
cd frontend

npm install
npm run dev
```

ブラウザで `http://localhost:5173` を開く。

## ローカル検証フロー

### 1. ゲスト閲覧（ログインなし）

1. リモート側でログイン（仮ユーザー）
2. 検収案件を作成
3. 「閲覧リンクを発行」を押す
4. 生成されたリンクをゲストで開く
5. ログインなしで案件を閲覧できる

### 2. リモートシャッター → 撮影

1. 現場側：`/inspections/:id/field` へアクセス
2. カメラ許可を与え、配信開始
3. リモート側：`/inspections/:id/remote` へアクセス
4. 「撮影」ボタンを押す
5. 現場側で自動撮影（ImageCapture API）、またはダイアログが表示
6. S3へアップロード → 証跡写真として反映

### 3. 注釈編集 → Annotated生成

1. 証跡写真をクリック
2. `/photos/:photoId/annotate` で注釈編集
3. ペン・図形・テキスト・看板テンプレを使用
4. 「保存」を押す
5. Annotated画像 + Annotation JSON が S3へ保存される
6. 証跡として Annotated画像が表示される

## デプロイ

### AWS へのデプロイ（概要）

1. **インフラ（CDK）デプロイ**

```bash
cd infra
npm install
cdk bootstrap  # 初回のみ
cdk deploy
```

CDK により以下が自動作成：
- DynamoDB テーブル
- S3 バケット
- API Gateway (HTTP API + WebSocket)
- Lambda (Go 関数)
- Cognito ユーザープール
- CloudWatch Logs
- IAM ロール/ポリシー

2. **バックエンド - Lambda へのビルド・デプロイ**

```bash
cd backend

# ビルド
GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go

# CDK または AWS CLI で Lambda 更新
cd ../infra
cdk deploy
```

3. **フロント - CloudFront + S3 へのデプロイ**

```bash
cd frontend

# ビルド
npm run build

# S3 へアップロード（CDK が S3 バケットを用意）
aws s3 sync dist/ s3://inspection-ui-bucket/ --delete
```

## OpenAPI / WebSocket スキーマ

- [OpenAPI 定義](./openapi.yaml)
- [WebSocket Message Schema](./ws-schema.json)

これらから自動生成されたクライアント型定義は `frontend/src/types/api.ts` に格納。

## データモデル

### DynamoDB テーブル設計

| テーブル              | PK                  | SK                 | 主要属性                                                                                                                            |
| --------------------- | ------------------- | ------------------ | ----------------------------------------------------------------------------------------------------------------------------------- |
| `Inspections`         | `inspectionId`      | -                  | title, place, status (draft/running/closed), realtimeRef, createdAt, startedAt, endedAt                                             |
| `ChecklistItems`      | `inspectionId`      | `itemId`           | label, result (ok/ng/pending), comment, updatedAt, updatedBy                                                                        |
| `Issues`              | `inspectionId`      | `issueId`          | title, detail, status (open/resolved/closed), assignee, createdAt, updatedAt                                                        |
| `CaptureRequests`     | `inspectionId`      | `captureRequestId` | requestedAt, requestedBy, state (requested/capturing/uploaded/failed/timeout), capturedAt, uploadedAt, photoId, linkedType/linkedId |
| `EvidencePhotos`      | `inspectionId`      | `photoId`          | captureRequestId, s3OriginalKey, s3AnnotatedKey, s3AnnotationJsonKey, annotatedAt, linkedType/linkedId                              |
| `AnnotationTemplates` | `templateId`        | `version`          | name, isActive, definitionJson (Konva オブジェクト定義), createdAt                                                                  |
| `ViewerLinks`         | `viewerAccessToken` | -                  | inspectionId, expiresAt, isRevoked, maxUses, usedCount, createdAt, createdBy                                                        |
| `AuditLogs`           | `logId`             | `timestamp`        | userId, action, inspectionId, details, timestamp                                                                                    |

## ロール・権限

| ロール        | 説明           | 機能                                                  |
| ------------- | -------------- | ----------------------------------------------------- |
| `admin`       | 管理者         | 全機能 + ユーザー管理 + テンプレ管理                  |
| `editor`      | 検収編集者     | チェック/Issue 編集 + 注釈編集 + 看板テンプレ置き換え |
| `capturer`    | 現場撮影者     | 高画質写真撮影・アップロード                          |
| `viewer`      | 閲覧者         | 案件・写真・Issue 閲覧のみ                            |
| `guestViewer` | 招待リンク閲覧 | 招待リンク経由で閲覧のみ（ログイン不要）              |

**認可**はサーバ側で強制実装。UI では機能を非表示にするが、API 呼び出しで権限チェック。

## 主な実装のポイント

### 1. リモートシャッター → 高画質撮影の流れ

```
Remote (capturer)
  ↓ WebSocket: capture.request (inspectionId, requestedBy)
Backend
  ↓ captureRequestId 採番、state = "requested"
  ↓ DB 保存
  ↓ WebSocket: capture.ack (captureRequestId, state) → 全員
Field (現場)
  ↓ WebSocket で capture.request 受信
  ↓ 通知画面表示
  ↓ ImageCapture API で撮影（自動）or ワンタップダイアログ表示
  ↓ presigned URL でアップロード
Backend (presigned URL検証)
  ↓ state = "uploaded"、photoId 割り当て
  ↓ WebSocket: capture.uploaded → 全員
Remote/Field/All
  ↓ 証跡タイムラインに反映
```

### 2. 注釈・看板テンプレの流れ

```
Editor
  ↓ 写真クリック → /photos/:photoId/annotate へ遷移
Annotation Canvas
  ↓ Konva で描画
  ↓ 看板テンプレ選択（ドロップダウン）
  ↓ テンプレ配置、パラメータ入力（タイトル/場所/日時/担当者）
  ↓ 「保存」 → presigned URL で Annotated画像 + Annotation JSON をアップロード
Backend
  ↓ state = "annotated"、s3AnnotatedKey / s3AnnotationJsonKey 保存
  ↓ WebSocket: photo.annotation.updated → 全員
Timeline
  ↓ Annotated画像が証跡として表示される
```

### 3. ゲスト閲覧の流れ

```
Admin/Editor/Capturer
  ↓ 案件の詳細ページで「閲覧リンクを発行」
Backend
  ↓ viewerAccessToken (推測困難な UUID/署名) 生成
  ↓ ViewerLinks テーブル保存（expiresAt = 7日後など、maxUses = 無制限 or 制限）
  ↓ URL 生成: https://...?viewerAccessToken=xxxxx
Frontend
  ↓ リンク表示 (コピー可能)

Guest
  ↓ リンク開く
  ↓ /inspections/:id/view へリダイレクト
Backend
  ↓ viewerAccessToken → guestSessionToken (JWT、1時間) 発行
  ↓ クライアント側で guestSessionToken を Authorization ヘッダで送信
Frontend
  ↓ 閲覧系 API のみ可 (isPublic または guestViewer で認可)
```

### 4. VideoProvider インターフェース（拡張性）

```go
type VideoProvider interface {
    // リアルタイム配信用 token 発行
    IssuePubToken(ctx context.Context, stageId string) (*VideoToken, error)
    IssuePubToken(ctx context.Context, stageId string) (*VideoToken, error)

    // Participant 管理など
    ListParticipants(ctx context.Context, stageId string) ([]*Participant, error)
}

// 実装 A: IVS Real-Time
type IVSProvider struct { ... }

// 実装 B: モック
type MockProvider struct {
    tokens map[string]*VideoToken
}
```

ローカル開発時は `VIDEO_PROVIDER=mock` で MockProvider を使用。本番では `VIDEO_PROVIDER=ivs` で IVSProvider を使用。

## テスト & 検証

### ユニットテスト (Go)

```bash
cd backend
go test ./...
```

### E2E フロー（手動）

1. **ゲスト閲覧**
   - ログインなしで招待リンクを開く → 閲覧できる ✓

2. **リモートシャッター**
   - リモート側：capturer ロールでログイン
   - 現場側：field 画面で配信開始
   - リモート側：「撮影」ボタン → 現場側が自動/ワンタップ撮影 ✓

3. **注釈編集**
   - 証跡写真をクリック → 注釈編集画面へ
   - ペン・図形・テキストで描画 → 看板テンプレ配置 → 保存 ✓
   - Annotated 画像が証跡として表示される ✓

4. **権限制御**
   - viewer ロール：撮影・注釈編集機能が UI に表示されない ✓
   - API 直呼び出し → 401/403 を返す ✓

## トラブルシューティング

### DynamoDB Local が起動しない

```bash
docker-compose ps
# ステータスを確認。"Up" なら正常
```

### S3 (MinIO) が起動しない

```bash
docker-compose logs minio
# MinIO のログを確認
# デフォルトアクセスキー: minioadmin / minioadmin
```

### Cognito 認証に失敗する

- `.env` の `COGNITO_USER_POOL_ID` / `COGNITO_CLIENT_ID` が正しいか確認
- AWS CLI で設定を確認

```bash
aws cognito-idp describe-user-pool --user-pool-id <YOUR_POOL_ID> --region ap-northeast-1
```

### WebSocket 接続に失敗

- フロント側 `.env` の `VITE_WS_URL` が正しいか確認
- バック側でログレベルを DEBUG に設定

```go
// backend/cmd/ws/main.go で
log.SetLevel("DEBUG")
```

## MVP 実装範囲 & 今後の拡張

### MVP（このリリース）

- ✅ 認証（Cognito） + 権限制御
- ✅ ゲスト閲覧（招待リンク）
- ✅ 検収案件/チェック/Issue CRUD
- ✅ WebSocket 基盤 + capture イベント
- ✅ 写真アップロード → 証跡一覧
- ✅ 注釈（ペン/図形/テキスト） → Annotated 生成・保存
- ✅ 看板テンプレ（基本機能）
- ✅ 報告書（HTML → S3）

### 将来の拡張（Phase 2+）

- [ ] 報告書 PDF 生成（ヘッドレスChrome）
- [ ] マルチテナント対応
- [ ] リアルタイム映像の自動録画機能
- [ ] AI による自動チェックリスト判定
- [ ] モバイルアプリ（React Native）
- [ ] オフライン機能
- [ ] 複数現場の並行監視

## ファイル構成

```
.
├── README.md (this file)
├── openapi.yaml (REST API 定義)
├── ws-schema.json (WebSocket message schema)
├── docker-compose.yaml (ローカル開発)
├── backend/
│   ├── go.mod / go.sum
│   ├── .env
│   ├── cmd/
│   │   ├── api/main.go (REST API Lambda)
│   │   ├── ws/main.go (WebSocket Lambda)
│   │   └── migrate/main.go (DB 初期化)
│   ├── internal/
│   │   ├── domain/ (Entity definitions)
│   │   ├── service/ (Business logic)
│   │   ├── handler/ (HTTP handlers)
│   │   ├── repository/ (DB access)
│   │   ├── middleware/ (Auth, logging)
│   │   └── infrastructure/ (S3, Video, etc.)
│   └── tests/
├── frontend/
│   ├── package.json / package-lock.json
│   ├── .env
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── index.html
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── pages/ (Login, InspectionList, Remote, Field, View, Annotate, AdminTemplates)
│   │   ├── components/ (UI components)
│   │   ├── hooks/ (useInspection, useWebSocket, useAnnotation, etc.)
│   │   ├── services/ (apiClient, wsClient)
│   │   ├── types/ (API types, Domain types)
│   │   ├── context/ (AuthContext, InspectionContext)
│   │   └── utils/ (helpers, validators)
│   └── public/
├── infra/
│   ├── package.json
│   ├── cdk.json
│   ├── tsconfig.json
│   └── lib/
│       └── inspection-stack.ts (CDK 定義)
└── docker-compose.yaml
```

## 本リポジトリの開発ガイドライン

### Code Style & Quality

**Go:**
- `gofmt` で自動フォーマット
- `golangci-lint` でチェック
- エラーハンドリングは明示的（`if err != nil`）
- ログは `log/slog` で構造化

**TypeScript/React:**
- `prettier` で自動フォーマット
- `eslint` でチェック
- 型定義を完全に（`any` は禁止）
- コンポーネントは関数型、フック使用

### Git Workflow

1. Feature ブランチで開発
2. Pull Request でレビュー
3. CI/CD（GitHub Actions）でテスト・デプロイ
4. Main へマージ

### コミットメッセージ

```
<type>(<scope>): <subject>

<body>

<footer>
```

例:
```
feat(capture): リモートシャッター機能実装

WebSocket でシャッター要求を受信し、
presigned URL で写真をアップロードする。

Fixes #123
```

## Support & License

- **Issue / Bug Report**: GitHub Issues
- **License**: (別途指定)

---

**最終更新**: 2025-12-28
**MVP ステータス**: 実装中

---

### 次のステップ

1. ローカルセットアップを実行
2. `docker-compose up -d` で DB/S3 起動
3. バック/フロント起動
4. http://localhost:5173 でアプリ確認
5. 本 README の E2E フロー実施

問題が発生した場合は、issue を作成またはサポートチャネルへ。
