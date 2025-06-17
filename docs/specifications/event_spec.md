# イベント取得機能仕様書

## 概要

このドキュメントは、`EventService` による外部Pythonスクリプトを介したイベントデータ取得・整形処理の仕様について記述します。

本機能は、Pythonスクリプトを実行してイベント情報を取得し、構造体としてパース、保存、検証を行います。

---

## 1. クラス構成

### `EventService`

| フィールド名 | 型 | 説明 |
|--------------|----|------|
| `scriptPath` | `string` | Pythonスクリプトのパス |
| `pythonCommand` | `string` | Python実行コマンド（例: `"python3"`） |
| `outputFilePath` | `string` | 出力ファイルパス（ハッシュ保存用） |
| `defaultTimeout` | `time.Duration` | スクリプト実行のデフォルトタイムアウト |

---

## 2. 主な構造体

### `ScheduleEvent`

| フィールド名 | 型 | 説明 |
|--------------|----|------|
| `Title` | `string` | イベントタイトル |
| `Level` | `string` | 難易度などの分類 |
| `Period` | `string` | 開催期間 |
| `Description` | `string` | 詳細説明 |
| `ImageURL` | `string` | イベント画像のURL |

### `EventFetchResponse`

| フィールド名 | 型 | 説明 |
|--------------|----|------|
| `Events` | `[]ScheduleEvent` | 取得されたイベント一覧 |
| `Error` | `error` | エラー情報 |

---

## 3. 主なメソッド

### `NewEventService(config *appConfig.Config) *EventService`

- **機能**: `EventService` を構成ファイルから初期化する。
- **戻り値**: `*EventService`

---

### `FetchEvents() ([]ScheduleEvent, error)`

- **機能**:  
    - デフォルトのタイムアウト付きで `FetchEventsWithContext` を実行。
    - バックワードコンパチ対応。
- **戻り値**: イベント一覧, エラー

---

### `FetchEventsWithContext(ctx context.Context) ([]ScheduleEvent, error)`

- **機能**:  
    - 指定されたPythonスクリプトを実行し、JSON出力をパース。
    - 整形された出力をファイルに保存。
- **戻り値**: イベント一覧, エラー

---

### `cleanJSONOutput(output []byte) []byte`

- **機能**:  
    - JSON出力の前後に付くゴミ文字を除去。
    - `[` または `{` から始まる正しい形式を探してトリム。
- **戻り値**: 整形済みJSON

---

### `parseEvents(output []byte) ([]ScheduleEvent, error)`

- **機能**:  
    - 出力をクリーンアップし、構造体へパース。
    - JSON構造の妥当性をチェック。
- **戻り値**: イベント一覧, エラー

---

### `validateEventData(events []ScheduleEvent) error`

- **機能**:  
    - 各イベントの `Title`, `Level` が空でないかを検証。
    - 今後さらに厳格なチェックが可能。
- **戻り値**: エラー（なければ nil）

---

## 4. 補助・拡張用メソッド

| メソッド名 | 説明 |
|------------|------|
| `executeScript(ctx context.Context)` | Pythonスクリプトを実行し出力を返却（将来的な拡張用） |
| `processOutput(output []byte)` | 出力をファイルに保存（将来的な拡張用） |
| `isValidJSON(data []byte)` | JSONの整合性を確認 |

---

## 5. レガシー関数

これらは今後非推奨となる可能性がある古いインタフェースです。

| メソッド名 | 概要 |
|------------|------|
| `FetchEvents()` | タイムアウト10秒の固定でスクリプトを実行 |
| `FetchEventsWithContext(ctx context.Context)` | グローバル関数版のイベント取得 |

---

## 6. 処理フロー

1. `NewEventService` により初期化。
2. `FetchEventsWithContext` を実行。
3. Pythonスクリプトを実行し、標準出力からJSONを取得。
4. ゴミを取り除き、構造体へパース。
5. 保存（`PersistJSONToHash256`）を試行。
6. 結果を返却。

---

## 7. 出力例（JSON）

```json
[
  {
    "title": "サクラフェスティバル",
    "level": "中級",
    "period": "2025-03-01 ～ 2025-03-31",
    "description": "春の祭典！",
    "image_url": "https://example.com/event.jpg"
  }
]
