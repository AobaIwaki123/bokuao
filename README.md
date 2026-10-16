# bokuao

アイドルグループ「僕が見たかった青空」（通称: 僕青）の公式スケジュールを取得・フィルタリングする高速・軽量なGo製CLIツールです。

公式サイト（SSR）からリアルタイムにイベント一覧を取得し、接触イベント（握手会・オンラインお話し会等）と純粋な音楽ライブ・フェス・公演を自動識別して抽出できます。

## 主な機能

- **超高速・依存ゼロ**: Go標準ライブラリのみで実装されており、外部依存なしで即座に動作します（実行時間約0.1秒）。
- **ライブ・公演フィルタ (`-live-only`)**: 「個別握手会」「オンラインお話し会」「サイン会」「バラエティ収録」を自動で除外し、生パフォーマンスがある現場（フリーライブ、大型フェス、対バン、単独・雲組公演）のみを瞬時に抽出します。
- **年・月指定対応**: 過去や翌月以降のスケジュール一覧を取得可能。
- **JSON出力対応 (`-json`)**: Bot連携、RSS生成、外部パイプライン処理への統合が容易。

## インストール

### Go でのインストール

```bash
go install github.com/AobaIwaki123/bokuao/cmd/bokuao@latest
```

### ソースからビルド

```bash
git clone https://github.com/AobaIwaki123/bokuao.git
cd bokuao
./scripts/verify-all.sh
# bin/bokuao が生成されます
```

## 使い方

### 1. 当月の全スケジュールを取得

```bash
bokuao
```

### 2. 生ライブ・公演のみに絞り込んで取得 (`-live-only`)

```bash
bokuao -live-only
```

### 3. 特定の年月を指定

```bash
# 2026年11月のライブのみを取得
bokuao -year 2026 -month 11 -live-only
```

### 4. JSON形式で出力 (`-json`)

```bash
bokuao -month 11 -live-only -json
```

## CLI オプション一覧

| フラグ | 型 | デフォルト値 | 説明 |
| :--- | :--- | :--- | :--- |
| `-year` | int | 現在の年 | 取得対象年 |
| `-month` | int | 現在の月 | 取得対象月 (1-12) |
| `-live-only` | bool | `false` | ライブ・フェス・公演パフォーマンスのみに絞り込む |
| `-json` | bool | `false` | JSON形式で出力 |

## アーキテクチャ

- `cmd/bokuao/`: CLIのエントリーポイント。引数のパースおよびフォーマット出力。
- `pkg/schedule/`: 公式サイトとの通信、HTMLパース、ライブイベント判定ロジック。
- `scripts/verify-all.sh`: ローカル検証用スクリプト（Linter、競合検出テスト、ビルド）。

## ライセンス

MIT License
