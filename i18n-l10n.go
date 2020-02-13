package i18n_l10n

import (
	// "fmt"
	"github.com/BurntSushi/toml"
)

/*
This is a somewhat opinionated golang package for simple and straightforward internationalization and localization of
web apps. This package assumes that you have one or more toml files marked in the format "active.{lang}.toml" in a
directory named "locales". It will crash if you don't. It will also crash if there are not corresponding string IDs
in each localization file.
Instantiate a Translator and use LoadStrings on launch to access translations. You can pass the Translator to a template
and use its .Translate function to translate strings. E.g., {{ .t.Translate "Hello, world!" }}
Alternatively, you could use GetStringsForPage to load the strings of each template into a one-dimensional map, and call
each string as such: {{ .strings.stringID }}.
 */

type Translator struct {
	Strings map[string]interface{} // strings will hold *all* translations for every page in every locale
	MasterLanguage string
	ActiveLanguage string
}

func loadStrings(lang string) (map[string]interface{}, error) {
	// Assume that the user has a file with required strings named "active.{language code}.toml" in a locales folder
	var activeStrings map[string]interface{}

	_, err := toml.DecodeFile("locales/active." + lang + ".toml", &activeStrings); if err != nil {
		return activeStrings, err
	}
	return activeStrings, nil
}

func (t *Translator) LoadStrings (langs []string) (map[string]interface{}, error) {
	// Assume that the user has a file with required strings named "active.{language code}.toml" in a locales folder
	strings := &t.Strings
	*strings = make(map[string]interface{})
	for _, lang := range langs {
		var activeStrings map[string]interface{}
		(*strings)[lang] = make(map[string]string)
		_, err := toml.DecodeFile("locales/active." + lang + ".toml", &activeStrings); if err != nil {
			return activeStrings, err
		}
		(*strings)[lang] = activeStrings
	}

	/*
	for k, v := range t.Strings {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
	 */

	return t.Strings, nil
}

func (t *Translator) GetStringsForPage (page string, lang string) map[string]interface{} {
	return (t.Strings)[lang].(map[string]interface{})[page].(map[string]interface{})
}

func (t *Translator) Translate(str string) (string) {
	// This should be able to be used like Gettext in Php, where you can surround a string with a function call _("str")
	// in this case it should be used in the template like {{ .t.Translate "whatever" }}
	// Go through relevant strings and find an equivalent value
	strings := t.Strings[t.MasterLanguage].(map[string]interface{})
	// go through all pages and try to find the string
	for k, _ := range strings {
		for l, j := range strings[k].(map[string]interface{}) {
			if j.(string) == str {
				// return j.(string)
				return t.Strings[t.ActiveLanguage].(map[string]interface{})[k].(map[string]interface{})[l].(string)
			}
		}
	}
	return "" // Return an empty string
}

func (t *Translator) Contains(str string) bool {
	for _, a := range t.Strings[t.MasterLanguage].(map[string]string) {
		if a == str {
			return true
		}
	}
	return false
}