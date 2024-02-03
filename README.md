# 技術記事ショート
Qiitaの記事をchatGPTで要約、voicevox_engineで読み上げするツール
## 初期設定
 - このツールはdocker及びdocker-composeを使い、別途インストールが必要。
 - 基礎ディレクトリ内にある.example.envファイルを参考に.envファイルを作成。
 - ディレクトリgo/app/config以下にある例を参考にconfig.xmlを作成。
 - 以下のコマンドを入力する。
```
$ docker network create ictshort_network
```
## 実行方法
 - 以下のコマンドを入力する。
```
$ docker-compose up
```
