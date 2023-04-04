/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openapi2raxtest",
	Short: "Create test scenario files for runn based on OpenAPI documentation",
	Long: `The process of creating test scenario files for runn from OpenAPI documentation involves creating scenarios for each API method endpoint,
and placing data in JSON files in the same directory as the scenarios. By increasing the number of arrays in the JSON file, multiple test data can be included in the scenarios.`,

	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
	_____  _____  _____  _____  _____  _____  ___  _____  _____  __ __  _____  _____ 
	/  _  \/  _  \/   __\/  _  \/  _  \/  _  \/___\<___  \/  _  \/  |  \/  _  \/  _  \
	|  |  ||   __/|   __||  |  ||  _  ||   __/|   | /  __/|  _  <|  |  ||  |  ||  |  |
	\_____/\__/   \_____/\__|__/\__|__/\__/   \___/<_____|\__|\_/\_____/\__|__/\__|__/
	`)
	},

	Run: func(cmd *cobra.Command, args []string) {

		// OpenAPIのYAMLファイルを読み込みしてオブジェクトを生成
		flags := *cmd.Flags()

		// フラグからOpenAPIスキーマの入力ファイル名を取得
		input, err := flags.GetString("input")
		if err != nil {
			fmt.Println(err)
			input = "openapi.yml"
		}

		// フラグからymlの出力先パスを取得
		output_path, err := flags.GetString("output")
		if err != nil {
			fmt.Println(err)
			output_path = "index.yml"
		}

		// フラグからデータを格納するjsonの出力先パスを取得
		data_path, err := flags.GetString("data")
		if err != nil {
			fmt.Println(err)
			data_path = "json://data.json"
		}

		// フラグからアクセスするサーバのホストを取得
		host, err := flags.GetString("server")
		if err != nil {
			fmt.Println(err)
			host = "http://localhost:8080"
		}

		println("[*] input file name: " + input)
		println("[*] output file name: " + output_path)
		println("[*] data file name: " + data_path)
		println("[*] host server host: " + host)
		println("")

		println("[*] Reading OpenAPI file...")

		// OpenAPIのYAMLファイルを読み込みしてオブジェクトを生成
		pathSpecs, loginSpec, _, err := genItem(input)
		if err != nil {
			fmt.Println(err)
			return
		}

		println("[*] Generating test raxtest config files...")
		// raxtest構造体を生成
		rootRaxSpec, dataRaxSpec, err := genRaxtestStruct(&host, &data_path, pathSpecs, loginSpec)

		if err != nil {
			fmt.Println(err)
			return
		}

		println("[*] Rendering raxtest config files...")
		// raxtest構造体をyamlファイル、jsonファイルにしてそれぞれ出力
		err = renderRaxTestStruct(&output_path, &data_path, rootRaxSpec, dataRaxSpec)
		if err != nil {
			fmt.Println(err)
			return
		}

		println("[#] Done! Please check the output files.")

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "Input file name")
	rootCmd.Flags().StringP("output", "o", "", "Output index yaml file name")
	rootCmd.Flags().StringP("data", "d", "", "Output data json file name")
	rootCmd.Flags().StringP("server", "s", "", "Host server")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
	rootCmd.MarkFlagRequired("data")
	rootCmd.MarkFlagRequired("server")

}
