# go-rest-clean-plane-chi
base project for go rest api with chi

## 使用ライブラリ
- wire
- viper
- swag
- gomock
- go-playground validator
- chi
- go.opentelemetry.io
- 

## 実装内容
- 標準ライブラリでのルーディング
- Middleware
  - context value
    - request id
    - request info
  - CORS
  - logging
  - エラーハンドリング（recover）
  - タイムアウト
  - 認証
- カスタムロガー
- バリデーター
- wire ジェネレート
- swagger定義のジェネレート
- mock ジェネレート
- CI
  - lint
  - test
  - Dockerfile/docker-compose

### README に載せたい情報
- swagger の使い方
- mock generator の使い方
- wire の使い方
- datadog container の使い方


## このプロジェクトは以下のディレクトリ構造に基づいています：
```
.
├── cmd
│     └── api : アプリケーションのエントリーポイント
├── docs
│     └── swagger
├── internal : プロジェクト固有のパッケージ
│     ├── adapters : 外部システムとのインターフェース
│     │     ├── primary
│     │     │     └── http
│     │     │         ├── handlers
│     │     │         ├── middleware
│     │     │         └── routes
│     │     └── secondary
│     │         ├── aws
│     │         ├── db
│     │         └── graphql
│     ├── core : ビジネスロジック
│     │     ├── domain
│     │     ├── services
│     │     └── usecases
│     ├── errors
│     ├── infrastructure : 技術的な実装（ロギングなど）
│     │     ├── auth
│     │     ├── config
│     │     └── logger
│     └── utils
├── pkg
└── scripts
```

### direnv
1. direnv をインストール
```bash
brew install direnv
```
- シェルの設定を追加（~/.zshrc や ~/.bashrc）
```
eval "$(direnv hook zsh)" 
```

2. 環境変数テンプレートをコピー
```bash
cp .envrc.example .envrc
```

3. `.envrc` を編集

4. direnv の許可
```bash
direnv allow
```
