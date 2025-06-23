# エラーレポート

環境構築時に発生したエラーと解決方法をまとめています。

## 発生したエラー

### 1. テスト実行時の依存関係エラー

**エラー内容:**
```
# ness
../../../go/pkg/mod/github.com/labstack/echo/v4@v4.13.4/middleware/rate_limiter.go:12:2: missing go.sum entry for module providing package golang.org/x/time/rate (imported by github.com/labstack/echo/v4/middleware); to add:
	go get github.com/labstack/echo/v4/middleware@v4.13.4

FAIL	ness [setup failed]
ok  	ness/internal/presentation/handler	0.828s
FAIL
```

**原因:**
- Echo v4.13.4の内部で使用している`golang.org/x/time/rate`パッケージがgo.sumに記録されていなかった
- `go mod download`実行時に、間接的な依存関係が完全に解決されていなかった

**解決方法:**
```bash
go mod tidy
```

**解決理由:**
- `go mod tidy`コマンドにより、不足していた間接依存関係が自動的に追加された
- 具体的には以下の依存関係が追加された：
  - `golang.org/x/time v0.11.0`
  - `gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`

**予防策:**
- 新しい依存関係を追加した後は必ず`go mod tidy`を実行する
- CI/CDパイプラインでも`go mod verify`を実行して依存関係の整合性を確認する

## 解決済み状況

- ✅ テスト実行: 正常に完了
- ✅ ビルド: 正常に完了
- ✅ Docker環境: 構築済み
- ✅ GitHub Actions: 設定済み

### 2. トレードノート機能実装時のモック関連エラー

**エラー内容:**
```
# ness/internal/application/usecase
../../../go/pkg/mod/github.com/stretchr/testify@v1.10.0/mock/mock.go:16:2: missing go.sum entry for module providing package github.com/stretchr/objx (imported by github.com/stretchr/testify/mock); to add:
	go get github.com/stretchr/testify/mock@v1.10.0
```

**原因:**
- testifyのmockパッケージが依存する`github.com/stretchr/objx`がgo.sumに記録されていなかった
- testifyを使用していたが、mockパッケージの依存関係が未解決だった

**解決方法:**
```bash
go get github.com/stretchr/testify/mock@v1.10.0
```

### 3. データベーステストでのSQL引数順序エラー

**エラー内容:**
```
Query 'INSERT INTO "trade_notes" ("trade_date","symbol","content","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"', arguments do not match: argument 3 expected [time.Time - 2024-06-23 10:30:00 +0000 UTC] does not match actual [time.Time - 2025-06-24 08:38:16.935254 +0900 JST]
```

**原因:**
- GORMが実際に生成するSQLの引数順序と、テストで期待している順序が異なっていた
- GORMは自動的にcreated_atとupdated_atを設定するため、引数の順序が変わった

**解決方法:**
テストでのモック期待値を実際のGORMのSQL生成順序に合わせて修正：
```go
// 修正前
mock.ExpectQuery(`INSERT INTO "trade_notes"`).
    WithArgs(
        sqlmock.AnyArg(), // created_at
        sqlmock.AnyArg(), // updated_at
        sqlmock.AnyArg(), // deleted_at
        tradeNote.TradeDate,
        tradeNote.Symbol,
        tradeNote.Content,
    )

// 修正後
mock.ExpectQuery(`INSERT INTO "trade_notes"`).
    WithArgs(
        tradeNote.TradeDate,
        tradeNote.Symbol,
        tradeNote.Content,
        sqlmock.AnyArg(), // created_at
        sqlmock.AnyArg(), // updated_at
        sqlmock.AnyArg(), // deleted_at
    )
```

### 4. テストでのインターフェース型不整合エラー

**エラー内容:**
```
cannot use mockUsecase (variable of type *MockTradeNoteUsecase) as *usecase.TradeNoteUsecase value in argument to NewTradeNoteHandler
```

**原因:**
- ハンドラーが具象型（*usecase.TradeNoteUsecase）を期待していたが、テスト用のモックは別の型だった
- インターフェースを使用していなかったため、モックオブジェクトを注入できなかった

**解決方法:**
1. ユースケースのインターフェースを定義：
```go
type TradeNoteUsecaseInterface interface {
    CreateTradeNote(ctx context.Context, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error)
    GetTradeNote(ctx context.Context, id uint) (*dto.TradeNoteResponse, error)
    GetAllTradeNotes(ctx context.Context) (*dto.TradeNoteListResponse, error)
    UpdateTradeNote(ctx context.Context, id uint, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error)
    DeleteTradeNote(ctx context.Context, id uint) error
}
```

2. ハンドラーでインターフェースを使用：
```go
type TradeNoteHandler struct {
    tradeNoteUsecase usecase.TradeNoteUsecaseInterface
}
```

### 5. テストでのContext型不整合エラー

**エラー内容:**
```
mock: Unexpected Method Call
GetAllTradeNotes(context.backgroundCtx)
The closest call I have is: 
GetAllTradeNotes(mock.anythingOfTypeArgument)
Diff: 0: FAIL: type *context.emptyCtx != type backgroundCtx
```

**原因:**
- モックの期待値で具体的なContext型を指定していたが、実際に渡されるContext型が異なっていた
- Echo経由で生成されるContextとテストで作成されるContextの型が違った

**解決方法:**
モックの期待値で`mock.Anything`を使用して任意のContext型を許可：
```go
// 修正前
mockUsecase.On("GetAllTradeNotes", mock.AnythingOfType("*context.emptyCtx")).Return(expectedResponse, nil)

// 修正後
mockUsecase.On("GetAllTradeNotes", mock.Anything).Return(expectedResponse, nil)
```

## 解決済み状況

- ✅ テスト実行: 正常に完了
- ✅ ビルド: 正常に完了
- ✅ Docker環境: 構築済み
- ✅ GitHub Actions: 設定済み
- ✅ トレードノート機能: 実装完了
- ✅ テストコード: 全層のテスト実装完了

## 今後の注意点

1. **依存関係管理**
   - 新しいパッケージ追加後は`go mod tidy`を実行
   - 定期的に`go mod verify`で依存関係を検証
   - testifyのmockパッケージを使用する場合は明示的にインストール

2. **テスト実行**
   - ローカルテスト実行前に`go mod download`を確実に実行
   - CI環境でもキャッシュが効かない場合に備えて依存関係の解決を確実に行う

3. **バージョン管理**
   - go.modとgo.sumファイルは必ずコミットに含める
   - 依存関係の更新は慎重に行い、テストで動作確認する

4. **テスト設計**
   - 依存性注入を前提とした設計にするため、インターフェースを活用する
   - データベーステストではGORMの実際のSQL生成順序を把握する
   - モックでのContext引数は`mock.Anything`を使用して柔軟に対応する
   - SQLMockを使用する場合は実際のクエリ形式を確認してから期待値を設定する