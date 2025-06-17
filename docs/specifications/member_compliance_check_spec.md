# メンバーコンプライアンスチェック機能仕様書

## 概要

このドキュメントは、`MemberService` によるDiscord上のメンバーのルール確認状況をチェックする機能について記述します。

本機能は、Discordの指定チャンネルからメンバーリストと、ルールへの同意（指定リアクション）を取得し、未確認メンバーの差分を検出・報告します。

---

## 1. クラス構成

### `MemberService`

| フィールド名 | 型 | 説明 |
|--------------|----|------|
| `discordToken` | `string` | Discord Botトークン |
| `memberChannelID` | `string` | メンバーを取得するチャンネルID |
| `ruleChannelID` | `string` | ルールリアクションを確認するチャンネルID |
| `masterUserID` | `string` | 管理者ユーザーID（未使用） |
| `goodReactionUrlEncodeString` | `string` | グッドリアクションのURLエンコード済み文字列 |

---

## 2. 主な構造体

### `ComplianceCheckResultResponse`

| フィールド名 | 型 | 説明 |
|--------------|----|------|
| `Message` | `string` | 差分の有無を示すメッセージ |
| `Diff` | `[]string` | 差分（未確認メンバー）リスト |

---

## 3. 主なメソッド

### `NewMemberService(config *appConfig.Config) *MemberService`

- **機能**: `MemberService` のコンストラクタ。設定ファイルからパラメータを取得。
- **戻り値**: `*MemberService`

---

### `CheckMemberComplianceWithContext(ctx context.Context) (*ComplianceCheckResultResponse, error)`

- **機能**:  
    - Discordセッションを作成し、メンバーリストとルール同意者を取得。
    - 差分を計算してレスポンス構造体にまとめて返す。
- **戻り値**: `*ComplianceCheckResultResponse`, `error`

---

### `getChannelMembers(dg *discordgo.Session, channelId string) ([]string, error)`

- **機能**: 指定チャンネルのGuild IDを使ってチャンネルメンバー一覧を取得。
- **戻り値**: メンバー名スライス, エラー

---

### `getRuleReaders(dg *discordgo.Session, channelId string) ([]string, error)`

- **機能**:  
    - 指定チャンネルの直近100件のメッセージにグッドリアクションをつけたユーザーを抽出。
    - Botユーザーは除外。
- **戻り値**: リアクション済みメンバー名スライス, エラー

---

### `isGoodReaction(emojiAPIName string, goodEmojiEncodeString string) bool`

- **機能**:  
    - リアクションが指定されたグッドマークかどうかを判定。
    - URLデコード比較も行う。
- **戻り値**: 真偽値

---

### `diffSlices(a, b []string) ([]string, bool)`

- **機能**: `a` に存在し、`b` に存在しない要素を抽出。
- **戻り値**: 差分スライス、完全一致かどうかのフラグ

---

## 4. 処理フロー

1. `NewMemberService` により設定をもとに初期化。
2. `CheckMemberComplianceWithContext` を実行。
3. メンバーチャンネルのメンバーリストを取得。
4. ルールチャンネルのリアクションから同意者リストを抽出。
5. 両者を比較し、未確認メンバーを特定。
6. メッセージと差分リストを返却。

---

## 5. 出力例

```json
{
  "message": "差分があります",
  "diff": ["userA", "userB"]
}
```