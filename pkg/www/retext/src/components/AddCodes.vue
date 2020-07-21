<template>

</template>

<script>
export default {
    name: "Code",
    props: ["color", "code"],
    computed: {
        style() {
            return `border-radius: 10px; border: 3px solid ${this.color}; padding: 20px;`;
        }
    },

    mounted() {
        console.log(this.code);
    },

    methods: {
        getTextsLength(category) {
            let length = category.main.texts.length;

            if (category.subcategories) {
                for (let subCat of category.subcategories) {
                    length += subCat.texts.length;
                }
            }

            return length;
        }
    }
}
</script>

<style scoped>

    .center-text {
        text-align: center;
    }
    .code-title {
        text-transform: capitalize;
        font-weight: bolder;
        margin: 0;
    }

    .just-number {
        background-color: white;
        color: black;
        border: none;
        font-weight: bolder;
    }

    .text-right {
        text-align: right;
    }

    .margin-top {
        margin-top: 20px;
    }

    .no-events {
        pointer-events:none;
    }

    .rborder {
        border-right: 1px solid gray;
    }

    .subcategory {
        border-radius: 3px;
        color: gray;
        padding: 3px;
        margin-bottom: 5px;
        border: 1px solid gray;
    }

    .top-container {
        display: flex;
        align-content: center;
        align-items: center;
    }

    .self-right {
        margin-left: auto;
    }

</style>