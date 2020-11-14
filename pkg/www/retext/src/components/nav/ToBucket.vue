<template>
    <div class="bucket-container">
        <div class="grid-header">
            <h4>{{codeContainer.main.name}}</h4>

            <i class="fa fa-trash"></i>
        </div>

        <div class="grid-view">
            <button class="btn btn-primary float-right" @click="goto(codeContainer.containerId)">
                View
            </button>
        </div>

        <div class="grid-progress">
            <b-progress  :max="100"
                            show-progress style="width: 100%" height="1.5rem">
                <b-progress-bar :value="percentage">
                    <span class="bar-label">Contains <strong>{{ percentage }}%</strong> of source files</span>
                </b-progress-bar>
            </b-progress>
        </div>
    </div>
</template>

<script>

export default {

    computed: {
        percentage() {
            return Math.floor(this.codeContainer.percentage * 100);
        }
    },
    name: "ToBucket",
    props: ["codeContainer"],
    methods: {
        goto(id) {
            this.$router.push(`/project/${this.$route.params.projectId}/view/${id}`);
        }
    }
}
</script>

<style scoped>

.bar-label {
    padding-left: 10px;
}
h4 {
    color: var(--gray);
    font-weight: bolder;
    text-transform: capitalize;

    margin: 0;
}

.bucket-container {
    display: grid;
    grid-template-columns: 5fr 2fr 1fr;
    grid-template-rows:
            1fr
            1fr;

    grid-template-areas:
            "header        .  view"
            "progress-bar  .  view";

    border: 2px solid var(--blue);
    border-radius: 5px;
    height: 100px;
    box-sizing: border-box;
    padding: 10px 10px 10px 10px;
}

.grid-header {]
    grid-area: header;
    display: flex;

    justify-content: space-between;
    align-items: center;
}

.grid-view {
    grid-area: view;
    display: flex;

    align-items: center;
    justify-items: center;
}

.grid-progress {
    grid-area: progress-bar;
    display: flex;
    align-items: center;
}
</style>