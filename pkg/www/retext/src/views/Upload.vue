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


                <div class="source-file" v-for="file in slicedSourceFiles" v-bind:key="file.id">
                    <div class="mb-3">
                        <ToDocument :document-id="file.id"
                                    :document-name="file.name"
                                    :project-id="projectId"
                                    button-text="Code" path="/code"></ToDocument>
                    </div>
                </div>


                <div v-if="slicedSourceFiles.length > perPage">
                    <b-pagination v-model="currentPage"
                                  :total-rows="slicedSourceFiles.length"
                                  :per-page="perPage"></b-pagination>
                </div>

                <div>
                    <UploadFile file-type="Source"
                                :project-id="projectId"
                                accepted-files=".docx,.txt,.text" :multiple=true>Upload New Sources</UploadFile>
                </div>

            </div>

            <div class="col-6">
                <h4 class="upload-header"> Demographic Information </h4>
                <div class="source-file" v-for="file in demos" v-bind:key="file.id">
                    <div class="mb-3">
                        <ToDocument :document-id="file.id" :document-name="file.name"
                                    :project-id="projectId"
                                    button-text="Modify" path="/demo"></ToDocument>
                    </div>
                </div>

                <div>
                    <UploadFile file-type="Demographics" :project-id="projectId"
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
import {actions} from "@/store";

export default {
    components: {UploadFile, ToDocument},

    component: {
        UploadFile
    },

    name: "Upload",

    data() {
        return {
            currentPage: 1,
            perPage: 4,
        }
    },

    computed: {
        slicedSourceFiles() {
            return this.$store.getters.sources.slice(
                (this.currentPage - 1) * this.perPage,
                this.currentPage * this.perPage
            );
        },

        demos() {
            return this.$store.getters.demos;
        },

        projectId() {
            return parseInt(this.$route.params.projectId);
        }
    },

    mounted() {
        this.$store.dispatch(actions.file.getFiles, this.projectId);
    },
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