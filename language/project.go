package language

type Project struct {
	Language        string
	LanguageVersion string
	PackageManager  string
	RunBuild        bool
	RelativeDir     string
	BuildOutput     string
	StartCommand    string
	Port            string
}
