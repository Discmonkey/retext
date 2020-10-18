<template>
    <div class="container">
        <div class="row">
            <div class="col-8">
                <h3 class="text-center font-weight-bold mb-4">Existing Projects</h3>

                <ToProject v-for="project in projects"  class="mb-1" :project="project" v-bind:key="project.id"></ToProject>
            </div>

            <div class="col-4">
                <div class="form-group">
                    <h3 class="text-center font-weight-bold mb-4">Start a new Project</h3>

                    <input type="text" placeholder="Project Name" class="mb-4 w-100 d-inline-block form-control" v-model="name">

                    <input type="month" class="mb-4 d-inline-block form-control" v-model="date">

                    <textarea class="mb-4 w-100 form-control" placeholder="Project Description" v-model="description" rows="10">

                    </textarea>

                    <div class="w-100 text-center">
                    <button class="btn btn-primary d-inline-block" @click="createProject()"> Create Project </button>
                    </div>
                </div>

            </div>
        </div>
    </div>
</template>

<script>
import ToProject from "@/components/nav/ToProject";
import {ProjectActions} from "@/store/modules/project";
export default {
    name: "Project",
    data() {
        const current = new Date();
        return {
            name: "",
            description: "",
            date: `${current.getFullYear()}-${(current.getMonth() + 1).toString().padStart(2, "0")}`,
        }
    },

    mounted() {
        this.$store.dispatch(ProjectActions.LOAD_PROJECTS);
    },

    methods: {
        createProject() {
            const [year, month] = this.date.split("-").map(i => parseInt(i))
            this.$store.dispatch(ProjectActions.ADD_PROJECT, ProjectActions.makeAddProjectPayload(this.name, this.description, year, month));
        },
    },

    computed: {
        projects() {
            return this.$store.getters.projects;
        }
    },
    components: {ToProject}
}
</script>

<style scoped>
input {
    padding: 10px;
    border-radius: 5px;
}

.mb-1 {
    margin-bottom: 1em;
}
.mr-4 {
    margin-right: 4%;
}
.w-48 {
    width: 48%;
}

button {
    font-weight: bolder;
    padding: 10px 20px 10px 20px;
}
</style>