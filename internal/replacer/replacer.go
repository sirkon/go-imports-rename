package replacer

type Replacer interface {
	Replace(old string) Variant
}
