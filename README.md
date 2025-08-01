# Slackログ取得CLI

指定したSlackチャンネルのメッセージログを、指定した期間分取得し、JSONファイルとして保存するコマンドラインツールです。

## 概要

このツールは、Slackの会話履歴を分析やバックアップのために一括で取得したい場合に便利です。チャンネルIDと期間を指定するだけで、自動的に全メッセージを取得します。

## ✅ 要件

* Go言語の実行環境 (v1.18以上)
* Slack Botトークン

## ⚙️ セットアップ

### 1. SlackアプリとBotトークンの準備

1.  [Slack API: Your Apps](https://api.slack.com/apps) にアクセスし、`Create New App` から新しいアプリを作成します。
2.  `From scratch` を選択し、アプリ名とワークスペースを決定します。
3.  サイドバーの `OAuth & Permissions` をクリックします。
4.  `Scopes` > `Bot Token Scopes` セクションで、以下の権限（スコープ）を追加します。
    * `channels:history` (パブリックチャンネルの履歴)
    * `groups:history` (プライベートチャンネルの履歴)
    * `im:history` (DMの履歴)
    * `mpim:history` (グループDMの履歴)
5.  ページ上部の `Install to Workspace` をクリックして、アプリをワークスペースにインストールします。
6.  インストール後、`Bot User OAuth Token` (`xoxb-` で始まる文字列) が表示されるので、これをコピーします。このトークンをCLI実行時に使用します。

### 2. ツールのビルド

1.  このツール（`main.go`）を任意のフォルダに保存します。
2.  ターミナルでそのフォルダに移動し、以下のコマンドを実行して依存ライブラリをインストールし、実行ファイルをビルドします。

    ```bash
    # Goモジュールの初期化 (初回のみ)
    go mod init slack-log-fetcher

    # 依存ライブラリのインストール
    go get [github.com/slack-go/slack](https://github.com/slack-go/slack)

    # ビルドして実行ファイルを作成
    go build .
    ```

    成功すると、`slack-log-fetcher` (Windowsの場合は `slack-log-fetcher.exe`) という実行ファイルが生成されます。

## 🚀 使い方

以下の形式でコマンドを実行します。

```bash
./slack-to-json-exporter -token="YOUR_SLACK_TOKEN" -channels="C0123ABC,C4567DEF" -startDate="YYYY-MM-DD" -endDate="YYYY-MM-DD" -outputDir="./output"
```

### 引数の説明

| フラグ         | 説明                                                               | 必須 | デフォルト値     |
| :------------- | :----------------------------------------------------------------- | :--- | :--------------- |
| `-token`       | Slack Botトークン (`xoxb-...`)。                                   | **はい** | (なし)           |
| `-channels`    | ログを取得したいチャンネルIDのリスト。カンマ区切りで複数指定可能。 | **はい** | (なし)           |
| `-startDate`   | 取得を開始する日付 (`YYYY-MM-DD`形式)。                           | **はい** | (なし)           |
| `-endDate`     | 取得を終了する日付 (`YYYY-MM-DD`形式)。                           | **はい** | (なし)           |
| `-outputDir`   | JSONファイルを保存するフォルダのパス。                             | いいえ | `./slack_logs`   |

### 実行例

チャンネル `C0123ABC` と `C4567DEF` の、2025年7月1日から2025年7月31日までのログを `./my_slack_logs` フォルダに保存する場合：

```bash
./slack-to-json-exporter \
  -token="xoxb-xxxxxxxxxxxxxxxxxxxxxxxx" \
  -channels="C0123ABC,C4567DEF" \
  -startDate="2025-07-01" \
  -endDate="2025-07-31" \
  -outputDir="./my_slack_logs"
```

実行後、`./my_slack_logs` フォルダ内に `C0123ABC.json` と `C4567DEF.json` が作成されます。
