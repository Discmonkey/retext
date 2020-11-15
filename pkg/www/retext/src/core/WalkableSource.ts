import {Source} from "@/model/source";
import {Word} from "@/model/word";


export interface WordWithAttributes {
    text: string;
    index: number;
    attributes: {
        [key: string]: boolean;
    }
}

export type coord = {
    paragraph: number;
    sentence: number;
    word: number;
};


type WordWithPointers = {
    coord: coord;
    word: WordWithAttributes;
}


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
    private flatSource: Array<WordWithPointers> = [];

    constructor(source: Source) {

        // flatten the data structure and remove any empty paragraphs or sentences
        this.source = source.paragraphs.filter(par => par.sentences.length > 0).map(
            paragraph => paragraph.sentences.filter(sen => sen.words.length > 0).map(
                sentence => sentence.words.map(
                    word => {
                        return {
                            text: word,
                            attributes: {},
                            index: 0,
                        }
                    })))

        for (let paragraph = 0; paragraph < this.source.length; paragraph++) {
            for (let sentence = 0; sentence < this.source[paragraph].length; sentence++) {
                for (let word = 0; word < this.source[paragraph][sentence].length; word++) {
                    this.source[paragraph][sentence][word].index = this.flatSource.length;

                    this.flatSource.push({
                        coord: {paragraph, sentence, word},
                        word: this.source[paragraph][sentence][word],
                    })
                }
            }
        }
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

    word(coord: coord): WordWithAttributes {
        return this.source[coord.paragraph][coord.sentence][coord.word];
    }

    /**
     * Search returns an ordered list of text which satisfy the given method, the returned words are contiguous
     * and contain the given coordinate
     */
    search(coord: coord, f: (word: WordWithAttributes) => boolean): [Array<WordWithAttributes>, coord, coord] {
        const anchor = this.word(coord).index;
        const ret = [];
        let start = anchor;
        let idx = anchor;
        let end = anchor;

        // search back while the condition is true
        while (idx > 0 && f(this.flatSource[idx].word)) {
            start = idx;
            idx--;
        }

        idx = anchor + 1;
        //search forward while the condition is true
        while (idx < this.flatSource.length && f(this.flatSource[idx].word)) {
            end = idx;
            idx++;
        }

        while (start <= end) {
            ret.push(this.flatSource[start].word);
            start++;
        }

        return [ret, this.flatSource[start].coord, this.flatSource[end].coord];

    }

    // a is always less than b
    private iterate(a: coord, b: coord, callback: (word: WordWithAttributes) => void) {
        const end = this.word(b).index;

        for (let start = this.word(a).index; start <= end; start++) {
            callback(this.flatSource[start].word);
        }
    }

}