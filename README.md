# memoapp_docker_copose | 新人課題レビュー用メモアプリケーション

## ディレクトリ構成

| ディレクトリ/ ファイル名 | 用途                                                      |
| ------------------------ | --------------------------------------------------------- |
| db                       | docker-compose で立ち上がる MySQL で実行する SQL ファイル |
| handler                  | メモアプリハンドラー                                      |
| model                    | メモアプリのモデル郡                                      |
| repository               | DB 接続処理                                               |
| src                      | フロント HTML                                             |

## 必要なツール

- go 1.14 以上
- Docker 環境（docker-compose コマンドが使用できる）
- Docker for Mac

## 使い方

### メモアプリの起動

```console
cd ./
docker-compose up --build
```

### お片付け

```console
cd ./
docker-compose down
```
