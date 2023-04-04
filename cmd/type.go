package cmd

// パラメータを格納する構造体を定義
type paramSpec struct {
	name      string
	item_type string
	example   any
}

// メソッドを格納する構造体を定義
type baseSpec struct {
	method  string
	queries []paramSpec
	bodies  []paramSpec
}

// パスを格納する構造体を定義
type pathSpec struct {
	name    string
	path    string
	methods []baseSpec
}

// 以下はraxtestスキーマの構造を表す構造体である

// ルートを格納する構造体を定義
type rootRaxSpec struct {
	BaseUrl string        `yaml:"base_url"`
	Data    string        `yaml:"data"`
	Init    []stepRaxSpec `yaml:"init"`
	Steps   []stepRaxSpec `yaml:"steps"`
}

// ステップを格納する構造体を定義
type stepRaxSpec struct {
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	Method       string `yaml:"method"`
	Query        string `yaml:"query,omitempty"`
	Body         string `yaml:"body,omitempty"`
	ExpectStatus int    `yaml:"expect_status"`
}

// 以下はraxtestで参照するjsonの構造を表す構造体である

// データの構造体を定義
type dataRaxSpec struct {
	Bodies  map[string]any `json:"body,omitempty"`
	Queries map[string]any `json:"query,omitempty"`
}
