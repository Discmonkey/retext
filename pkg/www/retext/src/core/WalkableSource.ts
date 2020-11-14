import {Source} from "@/model/source";


export interface WordWithAttributes {
    text: string;
    attributes: {
        [key: string]: boolean;
    }
}

export type coord = {
    paragraph: number;
    sentence: number;
    word: number;
};

type maybeCoord = [coord, boolean];

/**
 *
 * @param a
 * @param b
 *
 * return -1 if a < b, 1 if a > b, o otherwise
 */
function comp(a: coord, b: coord) {
    if (a.paragraph < b.paragraph) {
        return -1;
    }

    if (a.paragraph > b.paragraph) {
        return 1;
    }

    if (a.sentence < b.sentence) {
        return -1;
    }

    if (a.sentence > b.sentence) {
        return 1;
    }

    if (a.word < b.word) {
        return -1;
    }

    if (a.word > b.word) {
        return 1;
    }

    return 0;
}

export class WalkableSource {
    // Array<Paragraph<Sentence<WordWithAttributes>>>
    source: Array<Array<Array<WordWithAttributes>>>;

    constructor(source: Source) {

        // flatten the data structure and remove any empty paragraphs or sentences
        this.source = source.paragraphs.filter(par => par.sentences.length > 0).map(
            paragraph => paragraph.sentences.filter(sen => sen.words.length > 0).map(
                sentence => sentence.words.map(
                    word => {
                        return {
                            text: word,
                            attributes: {}
                        }
                    })))
    }

    walk(a: coord, b: coord, callback: (word: WordWithAttributes) => void) {
        const order = comp(a, b);

        switch (order) {
            case 0:
            case -1:
                this.iterate(a, b, callback);
                break;
            case 1:
                this.iterate(b, a, callback);
        }

    }

    // a is always less than b
    private iterate(a: coord, b: coord, callback: (word: WordWithAttributes) => void) {
        let valid = true
        let current = a;
        while (valid && comp(current, b) != 1) {
            callback(this.word(current));

            [current, valid] = this.next(current)

        }
    }

    private next(coord: coord): maybeCoord {
        let paragraph = coord.paragraph;
        let sentence = coord.sentence;
        let word = coord.word + 1;

        if (word < this.source[paragraph][sentence].length) {
            return [{paragraph, sentence, word}, true];
        }

        word = 0;
        sentence += 1;

        if (sentence < this.source[paragraph].length) {
            return [{paragraph, sentence, word}, true];
        }

        sentence = 0;
        paragraph += 1;

        return [{paragraph, sentence, word}, paragraph < this.source.length]
    }

    private word(coord: coord): WordWithAttributes {
        return this.source[coord.paragraph][coord.sentence][coord.word];
    }

}