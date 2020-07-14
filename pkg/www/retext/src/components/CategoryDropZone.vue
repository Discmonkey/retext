<template>
    <span @text-drop="textDrop($event)" class="category-drop-zone">
        <slot></slot>
    </span>
</template>

<script>
    import {getColor, invertColor} from "@/core/Colors"

    export default {
        name: "CategoryDropZone",
        props: ["category"],
        beforeMount() {
            if (!this.category.bgColor) {
                this.category.bgColor = getColor(this.category.id);
                this.category.color = invertColor(this.category.bgColor, true);
            }
        },
        methods: {
            textDrop: function(e) {
                e.detail.data.category = this.category;
            },
            getTextsLength() {
                let length = this.category.texts.length;

                if(this.category.subcategories) {
                    for(let subCat of this.category.subcategories) {
                        length += subCat.texts.length;
                    }
                }

                return length;
            }
        }
    }
</script>

<style scoped>
    /* apply to empty element + a background-color to get a circle */
    .category-color-drop {
        border-radius: 50%;
        padding: .5em;
        margin-right: 5px;
    }

    .category-drop-zone {
        cursor: default;
    }
</style>