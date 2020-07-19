<template>
    <span @text-drop="textDrop($event)" class="code-drop-zone">
        <slot></slot>
    </span>
</template>

<script>
    import {getColor, invertColor} from "@/core/Colors"

    export default {
        name: "CodeDropZone",
        props: ["code"],
        beforeMount() {
            if (!this.code.bgColor) {
                this.code.bgColor = getColor(this.code.id);
                this.code.color = invertColor(this.code.bgColor, true);
            }
        },
        methods: {
            textDrop: function(e) {
                e.detail.data.code = this.code;
            },
            getTextsLength() {
                let length = this.code.texts.length;

                if(this.code.subcodes) {
                    for(let subCode of this.code.subcodes) {
                        length += subCode.texts.length;
                    }
                }

                return length;
            }
        }
    }
</script>

<style scoped>
    /* apply to empty element + a background-color to get a circle */
    /*noinspection CssUnusedSymbol*/
    .code-color-drop {
        border-radius: 50%;
        padding: .5em;
        margin-right: 5px;
    }

    .code-drop-zone {
        cursor: default;
    }
</style>