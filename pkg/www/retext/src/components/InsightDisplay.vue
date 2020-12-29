<template>
<div>
    <div v-for="i in insights" :key="i.id" style="border: 1px solid blue">
        <div>{{ i.value }}</div>
        <ul>
            <li v-for="t in i.text_ids" :key="t">{{t}}</li>
        </ul>
    </div>
</div>
</template>

<script>
import {API} from "@/core/API.ts";

export default {
    name: "InsightDisplay",
    data: function() {
        return {
            insights: []
        };
    },

    async mounted() {
        this.insights = await API.insights.get(this.$store.getters.currentProject.id);
    },
}
</script>

<style scoped>

</style>