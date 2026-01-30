package tokenizer

import (
	"strings"

	"github.com/yanyiwu/gojieba"
)

type jiebaTokenizer struct {
	jieba *gojieba.Jieba
}

func NewJiebaTokenizer(dicPath ...string) Tokenizer {
	return &jiebaTokenizer{
		jieba: gojieba.NewJieba(dicPath...),
	}
}

func (j *jiebaTokenizer) CutForSearch(text string, hmm bool) []string {
	return j.jieba.CutForSearch(text, hmm)
}

func (j *jiebaTokenizer) GetSearchableTextFromSentences(sentences []string, hmm bool) string {
	searchText := strings.Join(sentences, " ")
	words := j.jieba.CutForSearch(searchText, hmm)
	var filteredWords []string
	for _, word := range words {
		if len([]rune(word)) > 1 {
			filteredWords = append(filteredWords, word)
		}
	}
	searchableText := strings.Join(filteredWords, " ")
	return searchableText
}
