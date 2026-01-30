package tokenizer

type Tokenizer interface {
	CutForSearch(text string, hmm bool) []string
	GetSearchableTextFromSentences(sentences []string, hmm bool) string
}