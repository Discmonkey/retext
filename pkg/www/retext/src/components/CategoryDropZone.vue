<template>
    <div v-bind:style="{backgroundColor: category.bgColor, color: category.color}"
         class="category-drop-zone"
         @mouseup="$emit('category-drop', category.id)">
        <span class="category-drop-zone-display">{{category.name}}</span>
    </div>
</template>

<script>
    import {gen_html, invertColor} from "@/components/Colors"

    export default {
        name: "CategoryDropZone",
        props: ["category"],
        beforeMount() {
            if (!this.category.bgColor) {
                this.category.bgColor = gen_html();
                this.category.color = invertColor(this.category.bgColor, true);
            }
        }
    }
</script>

<style>
    .category-drop-zone {
        border-radius: 2px;
        min-height: 3em;
        text-align: center;
        /* https://css-tricks.com/centering-css-complete-guide/#both-flexbox */
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .category-drop-zone:hover {
        box-shadow: inset 0 0 0 99999px rgba(255, 255, 255, 0.3);
    }
</style>