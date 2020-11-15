<template>
    <div class="container">
        <div class="grid" v-if="container !== null">
            <div class="item" v-for="(item, index) in texts" v-bind:key="index">
                <div class="item-header">

                </div>
                <div class="item-content">
                   {{item.text}}
                </div>

                <div class="item-footer">
                    <router-link :to="`/project/${projectId}/code/${item.document_id}`"> {{names[item.document_id]}}</router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script>

import Muuri from 'muuri';
import {actions} from "@/store";

export default {

    computed: {
        names() {
            return this.$store.getters.fileNames;
        },
        container() {
            const id = parseInt(this.$route.params.codeId);
            if (id in this.$store.getters.idToContainer) {
                return this.$store.getters.idToContainer[id];
            } else {
                return null;
            }
        },

        texts() {
            if (this.container === null) {
                return [];
            }

            return [this.container.main, ...this.container.subcodes].reduce((a, b) => a.concat(b.texts), []);
        },

        projectId() {
            return parseInt(this.$route.params.projectId);
        }
    },
    name: "Bucket",
    data() {
        return {
            grid: null,
        }
    },

    methods: {
        async init() {
            if (this.container === null) {
                // files are needed to get the file names
                await this.$store.dispatch(actions.file.getFiles, this.projectId);
                await this.$store.dispatch(actions.INIT_CONTAINERS);
            }

            setTimeout(() => this.grid = new Muuri('.grid', {
                dragHandle: 'item-header', dragEnabled: true
            }), 50);
        }
    },

    mounted() {
        this.init();
    },
}
</script>

<style scoped>

h4 {
    margin: 0;
}

.grid {
    position: relative;
}
.item {
    display: block;
    position: absolute;
    width: 500px;
    height: 250px;
    margin: 5px;
    z-index: 1;
    border-radius: 10px;
    border: 2px solid var(--blue );
    padding: 10px;
    box-sizing: border-box;
}

.item.muuri-item-dragging {
    z-index: 3;
}

.item.muuri-item-releasing {
    z-index: 2;
}

.item.muuri-item-hidden {
    z-index: 0;
}

.item-content {
    position: relative;
    width: 100%;
    height: calc(100% - 2em);
    display: flex;
    flex-direction: column;
    overflow-y: scroll;
    box-sizing: border-box;
    padding-bottom: 5px;
    padding-top: 5px;
}

.item-content::-webkit-scrollbar {
    display: none;  /* Safari and Chrome */
}

.item::-webkit-scrollbar {
    display: none;  /* Safari and Chrome */
}
.item-header {
    width: 100%;
    height: 1em;
    cursor: pointer;
}

.item-footer {
    margin-top: 5px;
    height: 1em;
    text-align: right;
}

</style>