<template>
    <span @text-drop="textDrop($event)" class="code-drop-zone">
        <slot></slot>
    </span>
</template>

<script>

    export default {
        name: "CodeDropZone",
        props: ["code"],

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
    .code-drop-zone {
        cursor: default;
    }
</style>