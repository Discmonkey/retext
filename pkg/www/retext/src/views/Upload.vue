<template>
    <div class="container">
        <div class="row">
            <div class="col-12 text-center">
                <h3>Kidney Cancer Study</h3>
                <p>Ten interviews with patients and cargivers for concept testing</p>
            </div>
        </div>

        <div class="row text-center space-bottom">
            <div class="col-4">
                <button class="btn btn-primary">
                    Upload Directory
                </button>
            </div>
            <div class="col-4">
                <UploadFile v-on:success="add($event)"></UploadFile>
            </div>

            <div class="col-4">
                <button class="btn btn-primary">
                    Upload Demographics
                </button>
            </div>

        </div>

        <div class="row">
            <div class="col-12 text-center">
                <h3>Uploaded So Far</h3>
            </div>
        </div>
        <div class="row">

            <div class="col-6">
                <h5>Source Files </h5>

                <div class="source-file" v-for="file in uploadedSourceFiles" v-bind:key="file">
                    <div class="mb-2">
                        <ToDocument :document-id="file" :document-name="file"></ToDocument>
                    </div>
                </div>
            </div>



            <div class="col-6">
                <h5> Demographic Information </h5>

            </div>
        </div>
    </div>


</template>

<script>
    import UploadFile from "../components/files/UploadFile";
    import ToDocument from "../components/nav/ToDocument";

    export default {
        components: {UploadFile, ToDocument},

        component: {
            UploadFile
        },

        name: "Upload",

        data() {
            return {
                uploadedSourceFiles: [],
                uploadedDemoFiles: [],
            }
        },

        mounted() {
            this.axios.get("/file/list").then((res) => {
                for (let f of res.data.Files) {
                    this.uploadedSourceFiles.push(f);
                }
            })
        },

        methods: {
            add(item) {
                this.uploadedSourceFiles.push(item.Key);
            }
        }
    }
</script>

<style scoped>
    .row {
        margin-top: 40px;
    }

    .space-bottom {
        margin-bottom: 5em;
    }

    h3, h4, h5 {
        font-weight: bold;
    }

</style>