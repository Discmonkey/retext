package parser

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {

	text := []byte("This is a short sentence. \n\n\r\t And another one. Two in this one... oh no? \n\tWith different spacing!!!\n")

	expected := Source{
		Paragraphs: []Paragraph{
			{
				Sentences: []Sentence{
					{
						Words: []string{
							"This",
							"is",
							"a",
							"short",
							"sentence.",
						},
					},
				},
			}, {
				Sentences: []Sentence{
					{
						Words: []string{
							"And",
							"another",
							"one.",
						},
					},
					{
						Words: []string{
							"Two",
							"in",
							"this",
							"one...",
						},
					},
					{
						Words: []string{
							"oh",
							"no?",
						},
					},
				},
			}, {
				Sentences: []Sentence{
					{
						Words: []string{
							"With",
							"different",
							"spacing!",
						},
					},
					{
						Words: []string{
							"!",
						},
					},
					{
						Words: []string{
							"!",
						},
					},
				},
			},
		},
	}

	converted, err := ConvertSource(text, Text)

	if err != nil {
		t.Fatal(err.Error())
	}

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

			if len(sentenceConverted.Words) != len(sentence.Words) {
				t.Fatal("word length does not match in sentence")
			}

			for numW, word := range sentence.Words {
				textConverted := sentenceConverted.Words[numW]

				if textConverted != word {
					t.Fatal(textConverted, " does not match ", word)
				}
			}
		}
	}
	fmt.Println(converted)
}

func TestIsLastWord(t *testing.T) {
	withNumber := []byte("1.2 is a float")
	asSentence := []byte("the end.")
	ellipsis := []byte("... Thought continues")
	inBetween := []byte(". is end")
	aFile := []byte("t.doc")

	if isTerminating(withNumber, 1) {
		t.Fatalf("failed on number")
	}

	if !isTerminating(asSentence, len(asSentence)-1) {
		t.Fatalf("failed on actual sentence")
	}

	if isTerminating(ellipsis, 0) {
		t.Fatalf("failed on ellipsis first")
	}

	if isTerminating(ellipsis, 1) {
		t.Fatalf("failed on ellipsis middle")
	}

	if !isTerminating(ellipsis, 2) {
		t.Fatalf("failed on ellipsis last")
	}

	if !isTerminating(inBetween, 0) {
		t.Fatalf("failed on mid sentence period")
	}

	if isTerminating(aFile, 1) {
		t.Fatalf("failed on file extension")
	}
}
