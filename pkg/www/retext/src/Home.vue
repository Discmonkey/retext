<template><router-view/></template>

<script>
import {ProjectActions} from "@/store/modules/project";

export default {
    name: "Home",
    mounted() {
        this.$store.dispatch(ProjectActions.SELECT_PROJECT, this.$route.params.projectId)
    }
}

</script>
