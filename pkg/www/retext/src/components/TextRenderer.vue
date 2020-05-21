<template>
    <div class="large">
        <p class="paragraph" v-for="(paragraph, index) in text.Paragraphs" v-bind:key="index">
            <span class="sentence" v-for="(sentence, index) in paragraph.Sentences" v-bind:key="index">
                <span v-for="(word, index) in sentence.Parts" v-bind:key="index">{{word.Text}} </span>
            </span>
        </p>
    </div>
</template>

<script>
    export default {
        name: "TextRenderer",
        props: ["text"]
    }
</script>

<style scoped>

</style>