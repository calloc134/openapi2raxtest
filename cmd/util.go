package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	"github.com/getkin/kin-openapi/openapi3"
)

// ディレクトリ名を生成する関数
func genDirName(str, sep string) string {

	re := regexp.MustCompile("[@{}]+")
	str_noSP := re.ReplaceAllString(str, "")
	parts := strings.Split(str_noSP, sep)
	for i, part := range parts {
		parts[i] = cases.Title(language.Und, cases.NoLower).String(part)
	}
	return strings.Join(parts, "")
}

// テストデータとなるJSONを生成する関数
func genJson(paramSpecs *[]paramSpec) map[string]any {

	jsonBodyMap := map[string]any{}

	// パラメータ毎に型を判定し、テストデータを生成
	for _, param := range *paramSpecs {

		// パラメータの型がstringの場合
		if param.item_type == "string" {
			// exampleが設定されている場合はexampleを
			// 設定されていない場合はダミーデータを設定
			if param.example != nil {
				jsonBodyMap[param.name] = param.example
			} else {
				jsonBodyMap[param.name] = "dummy"
			}
			// パラメータの型がnumberの場合
		} else if param.item_type == "number" {
			// 0を設定
			jsonBodyMap[param.name] = 0
			// パラメータの型がそれ以外の場合
		} else {
			// 空文字列を設定
			jsonBodyMap[param.name] = ""
		}
	}

	return jsonBodyMap

}

// openapiから採取したデータを格納する構造体
func genItem(inputFileName string) (*[]pathSpec, *[]pathSpec, *[]pathSpec, error) {
	// パス毎の構造体を格納するスライスを定義
	var pathSpecs []pathSpec

	// ログイン用とログアウト用のパス構造体を定義
	var loginSpecs []pathSpec
	var logoutSpecs []pathSpec

	// OpenAPIのYAMLファイルを読み込み
	doc, err := openapi3.NewLoader().LoadFromFile(inputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, nil, err
	}

	// パス毎に処理
	for _, path := range doc.Paths.InMatchingOrder() {

		// それぞれのパスに対するメソッドの一覧を取得
		obj := doc.Paths.Find(path).Operations()

		// メソッド毎の構造体を格納するスライスを定義
		var baseSpecs []baseSpec

		// メソッド毎に処理
		for method, op := range obj {

			// クエリとボディに当たるパラメータ構造体を格納するスライスを定義
			var queries []paramSpec
			var bodies []paramSpec

			// 元データにクエリパラメータがある場合
			if op.Parameters != nil {
				for _, q := range op.Parameters {

					// クエリ毎にクエリパラメータ構造体を生成
					queries = append(queries, paramSpec{
						// フィールドの名前
						name: q.Value.Name,
						// フィールドの型
						item_type: q.Value.Schema.Value.Type,
						// フィールドのサンプル値
						example: q.Value.Example,
					})
				}
			}

			// 元データにボディパラメータがある場合
			if op.RequestBody != nil {
				for name, b := range op.RequestBody.Value.Content["application/json"].Schema.Value.Properties {

					// ボディ毎にボディパラメータ構造体を生成
					bodies = append(bodies, paramSpec{
						// フィールドの名前
						name: name,
						// フィールドの型
						item_type: b.Value.Type,
						// フィールドのサンプル値
						example: b.Value.Example,
					})
				}
			}

			// メソッド毎にメソッド構造体を生成して末尾に追加
			baseSpecs = append(baseSpecs, baseSpec{
				// メソッド名
				method: method,
				// ボディパラメータ構造体のスライス
				bodies: bodies,
				// クエリパラメータ構造体のスライス
				queries: queries,
			})

		}

		if strings.Contains(path, "login") {
			// ログイン用のパス構造体を生成
			loginSpecs = append(loginSpecs, pathSpec{
				// ディレクトリ名
				name: genDirName(path, "/"),
				// パス
				path: path,
				// メソッド構造体のスライス
				methods: baseSpecs,
			})
		} else if strings.Contains(path, "logout") {
			logoutSpecs = append(logoutSpecs, pathSpec{
				// ディレクトリ名
				name: genDirName(path, "/"),
				// パス
				path: path,
				// メソッド構造体のスライス
				methods: baseSpecs,
			})
		} else {

			// パス毎にパス構造体を生成して末尾に追加
			pathSpecs = append(pathSpecs, pathSpec{
				// ディレクトリ名
				name: genDirName(path, "/"),
				// パス
				path: path,
				// メソッド構造体のスライス
				methods: baseSpecs,
			})
		}
	}

	return &pathSpecs, &loginSpecs, &logoutSpecs, nil
}

// openapi構造体よりraxtestのデータ構造体を生成する関数
func genRaxtestStruct(base_url *string, data_path *string, pathSpecs *[]pathSpec, loginSpecs *[]pathSpec) (*rootRaxSpec, *map[string]dataRaxSpec, error) {

	// 使う構造体を定義。処理後は返り値として返す
	rootRaxSpec := rootRaxSpec{
		BaseUrl: *base_url,
		Data:    *data_path,
		Init:    []stepRaxSpec{},
		Steps:   []stepRaxSpec{},
	}

	// JSONになるデータ構造体の連想配列を定義
	dataRaxSpecs := make(map[string]dataRaxSpec)
	// 引数で受け取ったopenapi構造体を読み込み

	// まずはログイン用のステップをinitとして処理
	for _, loginSpec := range *loginSpecs {
		// メソッド毎に処理
		for _, method_item := range loginSpec.methods {

			// ステップ名を生成
			step_name := loginSpec.name + "_" + method_item.method

			// クエリとボディが両方ある場合
			if method_item.bodies != nil && method_item.queries != nil {

				// データ構造体の連想配列にステップ名をキーにしてクエリとボディのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Bodies:  genJson(&method_item.bodies),
					Queries: genJson(&method_item.queries),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: loginSpec.path,
					// ボディパラメータ
					Body: step_name,
					// クエリパラメータ
					Query: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// ボディだけある場合
			} else if method_item.bodies != nil {

				// データ構造体の連想配列にステップ名をキーにしてボディのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Bodies: genJson(&method_item.bodies),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: loginSpec.path,
					// ボディパラメータ
					Body: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// クエリだけある場合
			} else if method_item.queries != nil {

				// データ構造体の連想配列にステップ名をキーにしてクエリのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Queries: genJson(&method_item.queries),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: loginSpec.path,
					// クエリパラメータ
					Query: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// クエリとボディが両方ない場合
			} else {

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: loginSpec.path,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})
			}
		}

	}

	// 次に通常のステップを処理
	for _, pathSpec := range *pathSpecs {
		// メソッド毎に処理
		for _, method_item := range pathSpec.methods {

			// ステップ名を生成
			step_name := pathSpec.name + "_" + method_item.method

			// クエリとボディが両方ある場合
			if method_item.bodies != nil && method_item.queries != nil {

				// データ構造体の連想配列にステップ名をキーにしてクエリとボディのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Bodies:  genJson(&method_item.bodies),
					Queries: genJson(&method_item.queries),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: pathSpec.path,
					// ボディパラメータ
					Body: step_name,
					// クエリパラメータ
					Query: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// ボディだけある場合
			} else if method_item.bodies != nil {

				// データ構造体の連想配列にステップ名をキーにしてボディのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Bodies: genJson(&method_item.bodies),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: pathSpec.path,
					// ボディパラメータ
					Body: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// クエリだけある場合
			} else if method_item.queries != nil {

				// データ構造体の連想配列にステップ名をキーにしてクエリのデータを格納
				dataRaxSpecs[step_name] = dataRaxSpec{
					Queries: genJson(&method_item.queries),
				}

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Init = append(rootRaxSpec.Init, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: pathSpec.path,
					// クエリパラメータ
					Query: step_name,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})

				// クエリとボディが両方ない場合
			} else {

				// ステップ構造体を生成してraxtest構造体内の配列の末尾に追加
				rootRaxSpec.Steps = append(rootRaxSpec.Steps, stepRaxSpec{
					// 名前
					Name: step_name,
					// メソッド名
					Method: method_item.method,
					// パス
					Path: pathSpec.path,
					// 予期しているステータスコード
					ExpectStatus: 200,
				})
			}
		}

	}

	return &rootRaxSpec, &dataRaxSpecs, nil
}

// raxtest構造体を受け取って、JSONに変換して指定されたパスに出力する関数
func renderRaxTestStruct(output_path *string, json_path *string, rootRaxSpec *rootRaxSpec, dataRaxSpecs *map[string]dataRaxSpec) error {

	// データをJSONに変換
	json, err := json.MarshalIndent(dataRaxSpecs, "", "  ")
	if err != nil {
		return err
	}

	// JSONをファイルに出力
	err = ioutil.WriteFile(*json_path, json, 0644)
	if err != nil {
		return err
	}

	// データをyamlに変換
	yaml, err := yaml.Marshal(rootRaxSpec)
	if err != nil {
		return err
	}

	// yamlをファイルに出力
	err = ioutil.WriteFile(*output_path, yaml, 0644)
	if err != nil {
		return err
	}

	return nil
}
