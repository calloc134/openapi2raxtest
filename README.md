```
 _____  _____  _____  _____  _____  _____  ___  _____  _____  _____  __  __  ____  _____  _____  ____ 
/  _  \/  _  \/   __\/  _  \/  _  \/  _  \/___\<___  \/  _  \/  _  \/  \/  \/    \/   __\/  ___>/    \
|  |  ||   __/|   __||  |  ||  _  ||   __/|   | /  __/|  _  <|  _  |>-    -<\-  -/|   __||___  |\-  -/
\_____/\__/   \_____/\__|__/\__|__/\__/   \___/<_____|\__|\_/\__|__/\__/\__/ |__| \_____/<_____/ |__| 

```
# openapi2raxtest

## 概要 / general

このリポジトリは、[raxtest](https://github.com/calloc134/raxtest) のテストスキーマを、OpenAPI定義から生成するためのユーティリティです。

## インストール
```bash
$ go install github.com/calloc134/openapi2raxtest@latest
```

## 使い方

例:
```bash
$ openapi2raxtest -i openapi.yaml -o runn.yaml -d data.json -s http://localhost:8080
```

 - `-i` : 入力ファイル名。OpenAPIスキーマを指定する
 - `-o` : raxtest構成ファイルの出力ファイル名  
yaml形式として、raxtestスキーマが出力される
 - `-d` : JSONデータファイル名  
json形式として、raxtestのテストデータが出力される
 - `-s` : サーバのホスト

## 注意事項 / caution
このプログラムは現在開発中のため、バグが含まれている可能性があります。  
また、バグを発見した場合は、PRを送っていただけると幸いです。

## 姉妹プロジェクト / sister projects
 - [raxtest](https://github.com/calloc134/raxtest) : 高速に動作するAPIのテストツール
 - [openapi2runn](https://github.com/calloc134/openapi2runn) : OpenAPI定義からrunn構成ファイルを生成するユーティリティ