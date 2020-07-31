

<template>
    <div class="doc-nav">
        <div class="row">
            <div class="col-md-9">
                <h5 class="doc-name" :title="documentNameNoExtension">
                    {{documentNameNoExtension}}
                </h5>
            </div>

            <div class="col-3">
                <button class="btn btn-primary float-right" @click="goto()">
                    {{ buttonText }}
                </button>
            </div>
        </div>
    </div>
</template>

<script>

    export default {
        name: "ToDocument",
        props: ["documentName", "documentId", "path", "buttonText"],
        computed: {
            documentNameNoExtension: function() {
                return this.documentName.split('.').slice(0, -1).join('.');
            }
        },
        methods: {
            goto() {
                this.$router.push(`${this.path}/${this.documentId}`)
            }
        }
    }
</script>

<style scoped>
    .doc-nav {
        border: 2px solid var(--primary);
        border-radius: 5px;
        padding: 10px;
    }
    .doc-name {
        font-weight: bolder;
        color: gray;
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>