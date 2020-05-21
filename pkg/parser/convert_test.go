package parser

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {

	text := []byte("This is a short sentence. \n\n\r\t And another one. Two in this one... oh no? \n\tWith different spacing!!!\n")

	expected := Document{
		Paragraphs: []Paragraph{
			{
				Sentences: []Sentence{
					{
						Parts: []Word{
							{Text: "This"},
							{Text: "is"},
							{Text: "a"},
							{Text: "short"},
							{Text: "sentence."},
						},
					},
				},
			}, {
				Sentences: []Sentence{
					{
						Parts: []Word{
							{Text: "And"},
							{Text: "another"},
							{Text: "one."},
						},
					},
					{
						Parts: []Word{
							{Text: "Two"},
							{Text: "in"},
							{Text: "this"},
							{Text: "one."},
						},
					},
					{
						Parts: []Word{
							{Text: "."},
						},
					},
					{
						Parts: []Word{
							{Text: "."},
						},
					},
					{
						Parts: []Word{
							{Text: "oh"},
							{Text: "no?"},
						},
					},
				},
			}, {
				Sentences: []Sentence{
					{
						Parts: []Word{
							{Text: "With"},
							{Text: "different"},
							{Text: "spacing!"},
						},
					},
					{
						Parts: []Word{
							{Text: "!"},
						},
					},
					{
						Parts: []Word{
							{Text: "!"},
						},
					},
				},
			},
		},
	}

	converted := Convert(text)

	if len(converted.Paragraphs) != len(expected.Paragraphs) {
		t.Fatal("paragraph length does not match")
	}
	for numP, paragraph := range expected.Paragraphs {

		paragraphConverted := converted.Paragraphs[numP]

		if len(paragraphConverted.Sentences) != len(paragraph.Sentences) {
			t.Fatal("sentence length in paragraph ", numP, "does not match")
		}

		for numS, sentence := range paragraph.Sentences {

			sentenceConverted := paragraphConverted.Sentences[numS]

			if len(sentenceConverted.Parts) != len(sentence.Parts) {
				t.Fatal("word length does not match in sentence")
			}

			for numW, word := range sentence.Parts {
				textConverted := sentenceConverted.Parts[numW].Text

				if textConverted != word.Text {
					t.Fatal(textConverted, "does not match", word.Text)
				}
			}
		}
	}
	fmt.Println(converted)
}
