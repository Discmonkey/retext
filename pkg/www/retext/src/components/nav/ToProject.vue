<template>
    <div class="border--blue layout">
        <h5 class="link" >
            <router-link :to="link"> {{project.name}}</router-link>

        </h5>

        <h5 class="date"> {{stringDate}}</h5>

        <p>  {{project.description}}</p>
    </div>
</template>

<script>
import moment from 'moment'

export default {
    name: "ToProject",
    props: [
        "project",
    ],

    computed: {
        stringDate() {
            return moment(this.project.timeTag).add(1, "days").format("MMM YYYY");
        },

        link() {
            return `/project/${this.project.id}/upload`
        },
    }
}
</script>

<style scoped>

    .layout {
        width: 100%;
        height: 120px;
        padding: 10px;
        border-width: 2px;
        border-style: solid;
        border-radius: 5px;
        box-sizing: border-box;
    }

    .link {
        color: #74a9f8
    }

    .date {
        font-weight: bold;
    }
</style>