# Ness

Go 1.24.4 + Echoフレームワークを使用し、ドメイン駆動設計（DDD）に基づいて構築されたWebアプリケーション

## 技術スタック

- **言語**: Go 1.24.4
- **Webフレームワーク**: Echo v4.13.4
- **データベース**: PostgreSQL 15
- **ORM**: GORM
- **アーキテクチャ**: ドメイン駆動設計（DDD）
- **コンテナ**: Docker + Docker Compose
- **テスティング**: testify/assert
- **CI/CD**: GitHub Actions
- **ホットリロード**: Air

## 現在のプロジェクト構成

```
/
├── main.go                          # アプリケーションのエントリーポイント
├── internal/
│   └── presentation/                # プレゼンテーション層
│       └── handler/                 # HTTPハンドラー
│           ├── hello_handler.go     # Hello Worldハンドラー
│           └── hello_handler_test.go # ハンドラーのテスト
├── Dockerfile                       # 本番用（distroless）
├── Dockerfile-local                 # 開発用（ホットリロード対応）
├── docker-compose.yml               # ローカル開発環境
├── .air.toml                        # ホットリロード設定
└── .github/workflows/               # GitHub Actions
    └── test.yml                     # テスト実行ワークフロー
```

## 環境構築

### 必要なソフトウェア

- Docker
- Docker Compose
- Go 1.24.4（ローカル開発時）

### ローカル開発環境の起動

1. リポジトリをクローン
   ```bash
   git clone <repository-url>
   cd ness
   ```

2. Docker Composeで開発環境を起動
   ```bash
   docker-compose up -d
   ```

3. データベースの初期化を待つ
   ```bash
   # PostgreSQLコンテナが起動するまで少し待つ
   docker-compose logs -f postgres
   ```

4. アプリケーションの確認
   ```bash
   curl http://localhost:8080/
   # レスポンス: {"message":"Hello, World!"}
   ```

5. ログの確認
   ```bash
   docker-compose logs -f
   ```

6. 開発環境の停止
   ```bash
   docker-compose down
   ```

### ローカル開発（Go直接実行）

1. 依存関係のインストール
   ```bash
   go mod download
   ```

2. アプリケーションの実行
   ```bash
   go run main.go
   ```

3. テストの実行
   ```bash
   go test ./...
   ```

## コマンド

### Go開発
```bash
# プロジェクトをビルド
go build

# アプリケーションを実行
go run main.go

# 全てのテストを実行
go test ./...

# 特定のテストを実行
go test -run <TestName> ./...

# ドメインテストのみ実行
go test ./internal/domain/...

# コードフォーマット
go fmt ./...

# 依存関係を整理
go mod tidy
```

### Docker開発環境
```bash
# 開発環境を起動（アプリ + PostgreSQL）
docker-compose up -d

# ログを確認
docker-compose logs -f

# PostgreSQLのログのみ確認
docker-compose logs -f postgres

# アプリコンテナに入る
docker-compose exec app bash

# PostgreSQLコンテナに入る
docker-compose exec postgres psql -U ness_user -d ness_db

# 開発環境を停止
docker-compose down

# イメージを再ビルド
docker-compose up -d --build

# データベースデータも削除して完全リセット
docker-compose down -v
```

## データベース

### 環境変数
開発環境では以下の環境変数が設定されています：
- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=ness_user`
- `DB_PASSWORD=ness_password`
- `DB_NAME=ness_db`

### 接続情報
- **ホスト**: localhost（ホストマシンから接続時）
- **ポート**: 5432
- **ユーザー**: ness_user
- **パスワード**: ness_password
- **データベース名**: ness_db

## CI/CD

GitHub Actionsで以下を自動実行：
- テスト実行
- 静的解析（go vet、staticcheck）
- カバレッジ測定

## 開発ガイドライン

### DDDアーキテクチャ
- **ドメイン層**: 他の層に依存しない純粋なビジネスロジック
- **アプリケーション層**: ドメイン層のみに依存
- **インフラストラクチャ層**: ドメイン層のインターフェースを実装
- **プレゼンテーション層**: アプリケーション層とインフラストラクチャ層に依存

### 開発ルール
- エンティティはドメインロジックを持つ
- 値オブジェクトは不変性を保つ
- リポジトリはインターフェースと実装を分離
- ユースケースは単一の責任を持つ
- ドメインサービスは複数のエンティティにまたがるロジックを扱う

### テスト戦略
- ドメイン層は単体テストを重視
- アプリケーション層はモックを使用したテスト
- インフラストラクチャ層は結合テスト
- プレゼンテーション層はE2Eテスト