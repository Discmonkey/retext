<template>
    <div class="border--blue layout">
        <h5 class="link"> {{projectName}}</h5>

        <h5 class="date"> {{stringDate}}</h5>

        <p>  {{projectDescription}}</p>
    </div>
</template>

<script>
import moment from 'moment'

export default {
    name: "ToProject",
    props: [
        "projectName",
        "projectDate",
        "projectDescription"
    ],

    computed: {
        stringDate() {
            return moment(this.projectDate).format("MMM YYYY");
        }
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