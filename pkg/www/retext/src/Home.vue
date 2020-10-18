<template><router-view v-if="projectLoaded"/></template>

<script>
import {ProjectActions} from "@/store/modules/project";

export default {
    name: "Home",
    mounted() {
        this.$store.dispatch(ProjectActions.SELECT_PROJECT, this.$route.params.projectId)
    },

    computed: {
        projectLoaded() {
            return this.$store.getters.currentProject !== null &&
                this.$store.getters.currentProject.id === parseInt(this.$route.params.projectId);
        }
    }
}

</script>
