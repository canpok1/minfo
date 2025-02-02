# minfo

Misskeyの情報を取得するツールです。

## 使い方

```
Usage:
  minfo <server_url> [flags]

Flags:
  -h, --help                       help for minfo
      --ignore-usernames strings   list of usernames to ignore (comma separated). when the display name is xxxxx@yyyy, username is yyyy.
  -l, --limit uint                 limit the number of notes (default 50)
```


## 開発環境
### 必要条件

- golang: 1.23

### 各種操作

基本的にmakeコマンドで操作可能

```
# run
# server : misskey server url origin
# option : tool option
make run server="https://xxxxx" option="-l 10 --ignore-user-names xxxx,yyyy,zzzz"

# build
make build

# clean
make clean
```

## リリース方法

gitのタグをプッシュする
