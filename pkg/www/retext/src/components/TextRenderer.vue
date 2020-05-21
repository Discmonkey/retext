<template>
    <div class="large">
        <p class="paragraph" v-for="(paragraph, parIndex) in text.Paragraphs" v-bind:key="parIndex">
            <template v-for="(sentence, senIndex) in paragraph.Sentences">
                <span
                        v-for="(word, wordIndex) in sentence.Parts" v-bind:key="senIndex + '-' + wordIndex"
                        v-on:mousedown.stop="start(parIndex, senIndex, wordIndex, $event)"
                        v-on:mouseup.stop="stop(parIndex, senIndex, wordIndex, $event)"
                        v-on:mouseenter="dragged(parIndex, senIndex, wordIndex, $event)"
                        class="border-on-hover word non-selectable"
                >
                    {{word.Text}}
                </span>
            </template>
        </p>
    </div>
</template>

<script>
    export default {
        name: "TextRenderer",
        props: ["text"],
        methods: {
            start: function(paragraph, sentence, word, e) {
                console.log(paragraph, sentence, word,e );
            },
            stop: function(paragraph, sentence, word, e) {
                console.log(paragraph, sentence, word, e);
            },
            dragged: function(paragraph, sentence, word, e) {
                console.log(paragraph, sentence, word, e);
            }
        }
    }
</script>

<style scoped>
    .border-on-hover:hover {
        box-shadow: inset 0 0 5px #689268;
    }

    .non-selectable {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
    }

</style>