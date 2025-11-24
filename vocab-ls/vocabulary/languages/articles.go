package languages

import "strings"

var articeln = func() map[string]struct{} {
	arr := []string{
		"der ",
		"den ",
		"dem ",
		"des ",
		"die ",
		"das ",
		"ein ",
		"einem ",
		"eines ",
		"eine ",
		"einer ",
	}
	set := make(map[string]struct{})
	for _, item := range arr {
		set[item] = struct{}{}
	}
	return set
}()

var article = func() map[string]struct{} {
	arr := []string{
		"la",
		"le",
		"l'",
		"une",
		"un",
		"des",
		"du",
	}
	set := make(map[string]struct{})
	for _, item := range arr {
		set[item] = struct{}{}
	}
	return set
}()

var articoli = func() map[string]struct{} {
	arr := []string{
		"il ",
		"la ",
		"l'",
		"lo ",
		"i ",
		"gli ",
		"le ",
		"una ",
		"uno ",
		"un'",
	}
	set := make(map[string]struct{})
	for _, item := range arr {
		set[item] = struct{}{}
	}
	return set
}()

func strip(set map[string]struct{}, word string, checkFor string) (string, bool) {
	splitted := strings.Split(word, checkFor)
	if len(splitted) == 1 {
		return word, false
	}

	possibleArticle := splitted[0]
	if _, exists := set[possibleArticle+checkFor]; exists {
		return strings.Join(splitted[1:], checkFor), true
	}
	return word, false
}

func StripGermanArticleFromWord(word string) string {
	maybeStripped, _ := strip(articeln, word, " ")
	return maybeStripped
}

func StripFrenchArticleFromWord(word string) string {
	if maybeStripped, handled := strip(article, word, " "); handled {
		return maybeStripped
	}

	maybeStripped, _ := strip(article, word, "'")
	return maybeStripped
}

func StripItalianArticleFromWord(word string) string {
	if maybeStripped, handled := strip(articoli, word, " "); handled {
		return maybeStripped
	}

	maybeStripped, _ := strip(articoli, word, "'")
	return maybeStripped
}
