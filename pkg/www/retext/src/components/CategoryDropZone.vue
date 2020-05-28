<template>
    <span class="category-drop-zone"
        @mouseup="$emit('category-drop', category.id)">{{category.name}}</span>
</template>

<script>
    export default {
        name: "CategoryDropZone",
        props: ["category"]
    }
</script>

<style scoped>

</style>