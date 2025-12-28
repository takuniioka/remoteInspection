# Testing Guide

このドキュメントでは、プロジェクト全体のテスト戦略と実行方法について説明します。

## テスト構成

### 1. バックエンド (Go) - ユニットテスト

**場所**: `backend/internal/*/`

**対象**: サービス層、リポジトリ層、ミドルウェア層、ドメインモデル

#### テストファイル一覧

- `backend/internal/service/services_test.go` - InspectionService, PhotoService, CaptureRequestService, ViewerLinkService
- `backend/internal/middleware/auth_test.go` - JWT認証、ロールチェック、エラーレスポンス
- `backend/internal/domain/models_test.go` - ドメインモデル検証、ステータス遷移

#### 実行方法

```bash
# すべてのGoテストを実行
cd backend
go test ./...

# 特定のパッケージをテスト
go test ./internal/service/...
go test ./internal/middleware/...
go test ./internal/domain/...

# 詳細な出力を表示
go test -v ./...

# テストカバレッジを確認
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### テストケース例

```go
// サービスのユニットテスト
func TestInspectionService_CreateInspection(t *testing.T) {
  // テストロジック
}

// ミドルウェアのテスト
func TestAuthMiddleware_ValidToken(t *testing.T) {
  // 有効なJWTトークンのテスト
}

// ドメインモデルのテスト
func TestInspectionStatusTransitions(t *testing.T) {
  // ステータス遷移の検証
}
```

### 2. フロントエンド (React) - ユニット・インテグレーションテスト

**フレームワーク**: Vitest + React Testing Library

**場所**: `frontend/src/**/*.test.ts[x]`

#### テストファイル一覧

**サービステスト**:
- `frontend/src/services/apiClient.test.ts` - API通信のモック、エラーハンドリング
- `frontend/src/services/wsClient.test.ts` - WebSocket接続、イベントルーティング

**ページコンポーネントテスト**:
- `frontend/src/pages/LoginPage.test.tsx` - ログインフォーム
- `frontend/src/pages/InspectionListPage.test.tsx` - 検査一覧
- `frontend/src/pages/RemoteInspectionPage.test.tsx` - リモート検査ビュー
- `frontend/src/pages/FieldCameraPage.test.tsx` - フィールドカメラ
- `frontend/src/pages/AnnotatePhotoPage.test.tsx` - 写真注釈
- `frontend/src/pages/AdminTemplatesPage.test.tsx` - テンプレート管理
- `frontend/src/pages/GuestViewPage.test.tsx` - ゲストビューアー

#### 実行方法

```bash
cd frontend

# すべてのテストを実行
npm test

# UIモードでテストを実行（ブラウザで確認）
npm run test:ui

# カバレッジレポートを生成
npm run test:coverage

# 特定のテストファイルのみ実行
npm test -- apiClient.test.ts

# ウォッチモードで実行
npm test -- --watch
```

#### テストカバレッジの確認

```bash
# カバレッジレポートを生成
npm run test:coverage

# HTML形式のレポートを開く
open coverage/index.html
```

#### テスト設定

**設定ファイル**: `frontend/vitest.config.ts`

- テスト環境: jsdom (ブラウザ互換)
- セットアップファイル: `src/test/setup.ts`
  - localStorage/sessionStorage自動クリア
  - WebSocketモック
  - IntersectionObserverモック
  - navigator.mediaDevicesモック

**テストユーティリティ**: `frontend/src/test/test-utils.tsx`
- BrowserRouterで囲まれたカスタムレンダー関数
- @testing-library/reactの全エクスポート

### 3. E2Eテスト (Playwright)

**フレームワーク**: Playwright (複数ブラウザ対応)

**場所**: `e2e/tests/`

#### テストシナリオ

**inspection-flow.spec.ts** (メインフロー):
- 認証フロー (ログイン、ログアウト、エラーハンドリング)
- 検査管理 (作成、表示、開始、終了)
- チェックリスト管理 (項目確認、問題作成)
- 写真キャプチャ・アップロード
- 写真注釈 (描画、テキスト追加、保存)
- ビューアリンク・ゲストアクセス
- WebSocket リアルタイム更新
- エラーハンドリング
- パフォーマンス

**advanced-features.spec.ts** (高度な機能):
- 管理者機能 (テンプレート管理、ユーザー管理)
- 監査ログ (フィルタリング、エクスポート)
- 報告書生成 (HTML/PDF、プレビュー、ダウンロード)
- アクセシビリティ (ARIAラベル、キーボード操作)
- モバイルレスポンシブネス

#### 実行方法

```bash
cd e2e

# 必要なパッケージをインストール
npm install

# 全E2Eテストを実行
npm test

# 特定のブラウザで実行
npm test -- --project=chromium
npm test -- --project=firefox
npm test -- --project=webkit

# UIモードで実行（対話的にテストを確認）
npm run test:ui

# デバッグモードで実行
npm run test:debug

# ヘッドレスモードで実行（ブラウザUI表示）
npm run test:headed

# 特定のテストスイートのみ実行
npm test -- inspection-flow.spec.ts

# 特定のテストケースのみ実行
npm test -- --grep "should render login page"
```

#### テストレポート

テスト実行後、自動生成されるレポート:

```bash
# HTMLレポートを開く
npx playwright show-report
```

### 4. 統合テスト

#### ローカル環境でのテスト

1. **Dockerで依存サービスを起動**:

```bash
docker-compose up -d
```

2. **バックエンドを起動**:

```bash
cd backend
go run ./cmd/api/main.go
go run ./cmd/ws/main.go
```

3. **フロントエンドを起動**:

```bash
cd frontend
npm install
npm run dev
```

4. **テストを実行**:

```bash
# Go テスト
cd backend && go test ./...

# React テスト
cd frontend && npm test

# E2E テスト
cd e2e && npm test
```

## テストベストプラクティス

### Go テスト

```go
// モックリポジトリを使用
type MockInspectionRepository struct {
  inspections map[string]*domain.Inspection
  shouldFail  bool
}

// テストでのセットアップ
func TestService(t *testing.T) {
  repo := &MockInspectionRepository{}
  service := NewInspectionService(repo, auditRepo)

  // テストロジック
}
```

### React テスト

```typescript
// コンポーネントテンプレート
import { render, screen } from '../test/test-utils';

describe('LoginPage', () => {
  it('should render login form', () => {
    render(<LoginPage />);
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });
});

// サービステンプレート
vi.mock('axios');
const mockedAxios = axios as any;

describe('ApiClient', () => {
  it('should fetch user', async () => {
    mockedAxios.get.mockResolvedValue({ data: mockUser });
    const user = await apiClient.getMe();
    expect(user).toEqual(mockUser.data);
  });
});
```

### E2E テスト

```typescript
// ページ遷移テスト
test('should navigate after login', async ({ page }) => {
  await page.goto('/login');
  await page.fill('input[type="email"]', 'test@example.com');
  await page.click('button:has-text("Login")');
  await expect(page).toHaveURL(/.*inspections/);
});

// リアルタイムイベントテスト
test('should receive real-time update', async ({ page }) => {
  await page.goto('/inspections/test-id');
  // WebSocketメッセージをシミュレート
  await expect(page.locator('[data-test="photo-item"]')).toBeVisible();
});
```

## テストカバレッジ目標

| レイヤー         | 目標 | 現在 |
| ---------------- | ---- | ---- |
| Go Domain Models | 100% | 95%  |
| Go Services      | 90%  | 85%  |
| Go Middleware    | 95%  | 90%  |
| React Components | 80%  | 60%  |
| React Services   | 90%  | 85%  |
| E2E Scenarios    | 80%  | 70%  |

## CI/CDパイプラインでのテスト実行

GitHub Actions (推奨):

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.22
      - run: cd backend && go test -cover ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: 18
      - run: cd frontend && npm install && npm test

  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: 18
      - run: cd e2e && npm install && npm test
```

※ リポジトリには CI ワークフローのサンプルを追加しています: `.github/workflows/ci-tests.yml`。
  - `backend-tests`: Go ユニットテストを実行
  - `frontend-tests`: Vitest を使ったフロントのユニットテスト + カバレッジ
  - `e2e-tests`: Playwright を使った E2E（必要に応じて有効化）

## トラブルシューティング

### Go テストが失敗する場合

```bash
# キャッシュをクリア
go clean -testcache

# 詳細ログを確認
go test -v -run TestName ./...

# タイムアウト値を増やす
go test -timeout 30s ./...
```

### React テストが失敗する場合

```bash
# node_modulesを再インストール
rm -rf node_modules package-lock.json
npm install

# テストを詳細表示で実行
npm test -- --reporter=verbose

# 特定のテストのみ実行
npm test -- --testNamePattern="LoginPage"
```

### E2E テストが失敗する場合

```bash
# ブラウザドライバを更新
npx playwright install

# スクリーンショットとトレースを確認
npm test -- --screenshot=only-on-failure --trace=on

# 特定のブラウザのみテスト
npm test -- --project=chromium
```

## テスト実行時間の目安

| テストタイプ          | 実行時間 |
| --------------------- | -------- |
| Go Unit Tests         | ~5秒     |
| React Component Tests | ~10秒    |
| React Coverage Report | ~20秒    |
| E2E Full Suite        | ~3分     |
| E2E Single Test       | ~30秒    |

## 参考資料

- **Go Testing**: https://golang.org/doc/effective_go#testing
- **Vitest**: https://vitest.dev/
- **React Testing Library**: https://testing-library.com/react
- **Playwright**: https://playwright.dev/
- **Testing Best Practices**: https://testing-library.com/docs/guiding-principles

## 質問・問題報告

テスト実行時に問題が発生した場合:

1. エラーメッセージを記録
2. 失敗したテスト名を確認
3. ローカル環境で単体テストを実行
4. キャッシュとnode_modulesをクリア
5. リポジトリのissueを確認
