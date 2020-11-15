<template>
<div class="container">
    <div class="row">
        <div class="col-12 text-center">
            <h3>
                {{project.name}}
            </h3>

            <h5 class="gray">
                {{stringDate}}
            </h5>

            <p>
                {{project.description}}
            </p>

        </div>


    </div>

    <div class="row">
        <div class="col-8 offset-2 layout">
            <button v-for="link in links" @click="goto(link.link, link.disabled)" v-bind:key="link.name" class="btn btn-primary">
                {{link.name}} {{l}}
            </button>
        </div>
    </div>
</div>
</template>

<script>
import moment from "moment";

export default {
    computed: {
        project() {
            return this.$store.getters.currentProject;
        },

        stringDate() {

            return moment().date(1).month(this.project.month).year(this.project.year).format("MMM YYYY");
        },
    },

    data() {
        return {
            links: [
                {name: "Upload Sources", link: "upload", disabled: false},
                {name: "Upload Demographics", link: "demo", disabled: true,},
                {name: "Code Sources", link: "code", disabled: false},
                {name: "View Coding Buckets", link: "buckets", disabled: false},
                {name: "View Insights", link: "insights", disabled: true},
                {name: "Export Insights", link: "insights", disabled: true}
            ]
        }
    },

    methods: {
        goto(destination, isDisabled) {
            if (isDisabled) {
                alert("not currently implemented!");
                return;
            }
            const maybeSlash = this.$route.fullPath.endsWith("/") ? "" : "/";

            this.$router.push(this.$route.fullPath + maybeSlash + destination);
        }
    },
name: "ProjectHub"

}
</script>

<style scoped>
.project-container {
    width: 100%;
}

.layout {
    display: grid;
    grid-template-rows: repeat(3, 3em);
    grid-template-columns: 1fr 1fr;
    grid-gap: 10px;
}
</style>