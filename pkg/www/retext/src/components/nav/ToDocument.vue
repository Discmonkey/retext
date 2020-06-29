

<template>
    <div class="doc-nav">
        <div class="row">
            <div class="col-md-2">
                <h4> {{documentName}}</h4>
            </div>

            <div class="offset-7 col-3">
                <button class="btn btn-success bg-blue float-right" @click="goto()">
                    Code
                    <i class="fa fa-arrow-right"></i>
                </button>
            </div>
        </div>

        <div class="row">
            <div class="col-9">
                ID: {{documentId}}
            </div>

        </div>
    </div>
</template>

<script>

    export default {
        name: "ToDocument",
        props: ["documentName", "documentId"],
        methods: {
            goto() {
                this.$router.push(`/code/${this.documentId}`)
            }
        }
    }
</script>

<style scoped>
    .doc-nav {
        border: 2px solid lightblue;
        border-radius: 5px;
        padding: 10px;
    }
</style>