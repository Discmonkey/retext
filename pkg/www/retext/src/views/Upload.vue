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


                <div class="source-file" v-for="file in slicedSourceFiles" v-bind:key="file.Id">
                    <div class="mb-3">
                        <ToDocument :document-id="file.Id"
                                    :document-name="file.Name"
                                    :project-id="projectId"
                                    button-text="Code" path="/code"></ToDocument>
                    </div>
                </div>


                <div v-if="uploadedSourceFiles.length > perPage">
                    <b-pagination v-model="currentPage"
                                  :total-rows="uploadedSourceFiles.length"
                                  :per-page="perPage"></b-pagination>
                </div>

                <div>
                    <UploadFile file-type="Source"
                                v-on:success="addSource($event)" :project-id="projectId"
                                accepted-files=".docx,.txt,.text" :multiple=true>Upload New Sources</UploadFile>
                </div>

            </div>

            <div class="col-6">
                <h4 class="upload-header"> Demographic Information </h4>
                <div class="source-file" v-for="file in uploadedDemoFiles" v-bind:key="file.Id">
                    <div class="mb-3">
                        <ToDocument :document-id="file.Id" :document-name="file.Name"
                                    :project-id="projectId"
                                    button-text="Modify" path="/demo"></ToDocument>
                    </div>
                </div>

                <div>
                    <UploadFile file-type="Demographics" v-on:success="addDemo($event)" :project-id="projectId"
                                tooltip="For demographic information, please upload a .xlsx or .csv file in which each participant is a different row (a header row is required)."
                                accepted-files=".xlsx,.csv">Upload New Demographics</UploadFile>
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
            currentPage: 1,
            perPage: 4,
        }
    },

    computed: {
        slicedSourceFiles() {
            return this.uploadedSourceFiles.slice(
                (this.currentPage - 1) * this.perPage,
                this.currentPage * this.perPage
            );
        },

        projectId() {
            return parseInt(this.$route.params.projectId);
        }
    },

    mounted() {
        this.axios.get(`/file/list?projectId=${this.$route.params.projectId}`).then((res) => {
            for (let f of res.data.Files) {
                if (f.Type === "SourceFile") {
                    this.uploadedSourceFiles.push(f)
                } else if (f.Type === "DemoFile") {
                    this.uploadedDemoFiles.push(f)
                }
            }
        })
    },

    methods: {
        addSource(items) {
            items.forEach((item) => {
                this.uploadedSourceFiles.push(item.File);
            });
        },

        addDemo(items) {
            items.forEach((item) => {
                this.uploadedDemoFiles.push(item.File);
            });
        },

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