<template>
    <div class="container-fluid">

        <div class="row">
            <div class="col-md-12 document-display">
                <TextRenderer :document-id="documentId"></TextRenderer>
            </div>
        </div>
    </div>
</template>

<script>
    import TextRenderer from "./TextRenderer";
    import {actions} from "@/store";

    export default {

        computed: {
            documentId() {
                return parseInt(this.$route.params.documentId);
            },
        },

        name: "DocumentDisplay",

        components: {
            TextRenderer
        },
        data: () => {
            return {
                files: [],
                uploadedFiles: [],
                currentText: "",
                selected: "",
            }
        },

        mounted() {
            this.$store.dispatch(actions.file.selectSource, {
                fileId: this.documentId
            });
        },

    }
</script>

<style scoped>
    .document-display {
        max-height: 90%;
        overflow-y:hidden;
    }

    .row {
        padding-top: 10px;
    }
</style>