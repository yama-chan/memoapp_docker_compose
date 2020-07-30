# memoapp_docker_copose | 新人課題レビュー用メモアプリケーション

## ディレクトリ構成

| ディレクトリ/ ファイル名 | 用途                                                      |
| ------------------------ | --------------------------------------------------------- |
| db                       | docker-compose で立ち上がる MySQL で実行する SQL ファイル |
| images                   | アイコン画像など                                          |
| internal                 | サーバー処理                                              |
| model                    | メモアプリのモデル郡                                      |
| src                      | フロント HTML                                             |

## 必要なツール

- go 1.14 以上
- Docker 環境（docker-compose コマンドが使用できる）
- Docker for Mac

## 使い方

### メモアプリの起動

```console
make up
```

### お片付け

```console
make down
```
