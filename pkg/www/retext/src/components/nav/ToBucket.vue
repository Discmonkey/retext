<template>
    <div class="bucket-container">
        <div class="grid-header">
            {{codeContainer.main.name}}
        </div>

        <div class="grid-view">
            {{codeContainer.id}}
        </div>

        <div class="grid-progress">
            {{codeContainer.percentage}}
        </div>
    </div>
</template>

<script>
export default {
    name: "ToBucket",
    props: ["codeContainer"],
}
</script>

<style scoped>
.bucket-container {
    display: grid;
    grid-template-columns: 5fr 2fr 2fr;
    grid-template-rows:
            1fr
            1fr;

    grid-template-areas:
            "header        .  view"
            "progress-bar  .  view";

}

.grid-header {
    grid-area: header;
}

.grid-view {
    grid-area: view;
}

.grid-progress {
    grid-area: progress-bar
}
</style>