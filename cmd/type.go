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
	BaseUrl        string                     `yaml:"base_url"`
	Data           string                     `yaml:"data"`
	Init           []stepRaxSpec              `yaml:"init"`
	StepCategories map[string]categoryRaxSpec `yaml:"categories"`
}

// カテゴリーを格納する構造体を定義
type categoryRaxSpec struct {
	Login string         `yaml:"login,omitempty"`
	Steps *[]stepRaxSpec `yaml:"steps"`
}

// ステップを格納する構造体を定義
type stepRaxSpec struct {
	Name    string        `yaml:"name"`
	Path    string        `yaml:"path"`
	Method  string        `yaml:"method"`
	RefData string        `yaml:"ref_data"`
	Option  optionRaxSpec `yaml:"option"`
}

// オプションを格納する構造体を定義
type optionRaxSpec struct {
	Query bool `yaml:"query"`
	Body  bool `yaml:"body"`
}

// 以下はraxtestで参照するjsonの構造を表す構造体である

// データの構造体を定義
type dataRaxSpec struct {
	Bodies       map[string]any `json:"body,omitempty"`
	Queries      map[string]any `json:"query,omitempty"`
	ExpectStatus int            `json:"expect_status"`
}
