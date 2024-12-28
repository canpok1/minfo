# minfo

Misskeyの情報を取得するツールです。

## 使い方

```
Usage:
  minfo <server_url> [flags]

Flags:
  -h, --help        help for minfo
  -l, --limit int   limit the number of notes (default 50)
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
make run server="https://xxxxx" option="-l 10"

# build
make build

# clean
make clean
```

## リリース方法

gitのタグをプッシュする
