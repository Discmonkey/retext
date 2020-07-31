<template>
    <div class="container">

        <div class="row">
            <div class="col-12 text-center">
                <h3>Uploads</h3>
            </div>
        </div>
        <div class="row">

            <div class="col-6">
                <h4 class="upload-header">Source Files </h4>

                <div class="source-file" v-for="file in uploadedSourceFiles" v-bind:key="file">
                    <div class="mb-3">
                        <ToDocument :document-id="file" :document-name="file" button-text="Code" path="/code"></ToDocument>
                    </div>
                </div>

                <div>
                    <UploadFile file-type="Source" v-on:success="addSource($event)" accepted-files=".docx,.txt,.text" :multiple=true></UploadFile>
                </div>
            </div>


            <div class="col-6">
                <h4 class="upload-header"> Demographic Information </h4>
                <div class="source-file" v-for="file in uploadedDemoFiles" v-bind:key="file">
                    <div class="mb-3">
                        <ToDocument :document-id="file" :document-name="file" button-text="Modify" path="/demo"></ToDocument>
                    </div>
                </div>

                <div>
                    <UploadFile file-type="Demographics" v-on:success="addDemo($event)" accepted-files=".xlsx"></UploadFile>
                </div>
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
                if (f.Type === "SourceFile") {
                    this.uploadedSourceFiles.push(f.ID)
                } else if (f.Type === "DemoFile") {
                    this.uploadedDemoFiles.push(f.ID)
                }
            }
        })
    },

    methods: {
        addSource(items) {
            items.forEach((item) => {
                this.uploadedSourceFiles.push(item.Key);
            });
        },

        addDemo(items) {
            items.forEach((item) => {
                this.uploadedDemoFiles.push(item.Key);
            });
        }
    }
}
</script>

<style scoped>
.row {
    margin-top: 40px;
}

.mb-3 {
    margin-bottom: 1em;
}


h3, h4, h5 {
    font-weight: bold;
}

.upload-header {
    padding-bottom: 10px;
}
</style>