<template>
    <div :style="style">

    </div>
</template>

<script>
export default {
    name: "ColorBullet",
    props: ["color"],
    computed: {
        style() {
            return `border-radius: 50%; background-color: ${this.color}; width: 1.5em; height: 1.5em; margin-bottom: 2px`;
        }
    }
}
</script>

<style scoped>

</style>