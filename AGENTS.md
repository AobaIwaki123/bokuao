# bokuao

僕が見たかった青空 公式スケジュール取得・フィルタリング CLI

## 開発・PR運用ルール

- **PR-Only Workflow**:
  - `main` ブランチへの直接 Push は原則禁止とし、必ず feature ブランチを作成して Pull Request を起票・マージします。
  - コミット・Push 前には必ず `./scripts/verify-all.sh` を実行して 100% PASS することを確認します。
- **Go バージョン**:
  - `go 1.23.0` に固定します。
- **ドキュメント方針**:
  - 絵文字は使用しません。
