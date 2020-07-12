

<template>
    <div class="doc-nav">
        <div class="row">
            <div class="col-md-9">
                <h5> {{documentName}}</h5>
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
        methods: {
            goto() {
                this.$router.push(`${this.path}/${this.documentId}`)
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
    h5 {
        font-weight: bolder;
        color: gray;
    }
</style>