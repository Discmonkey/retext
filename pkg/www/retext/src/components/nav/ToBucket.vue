<template>
    <div class="bucket-container">

    </div>
</template>

<script>
export default {
name: "ToBucket"
}
</script>

<style scoped>
    .bucket-container {
        display: grid;
        grid-template-columns: 5fr 2fr 2fr;
        grid-template-rows:  1fr  1fr;

        grid-template-areas:
            "header        .  view"
            "progress bar  .  view";

    }

    .grid-header {

    }
</style>