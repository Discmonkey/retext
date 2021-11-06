<template>
    <div class="container">
        <div class="row">
            <div class="col-md-12">
                <router-link :to="{name: 'Insights'}" style="float:right">Go to Insights list</router-link>
            </div>
        </div>
        <div class="row">
            <div class="col-md-9 insights">
                <textarea class="form-control" v-model="insightText" placeholder="insights here, no saving for now"></textarea>
            </div>
            <div class="col-md-3">
                <button class="btn btn-primary bold" id="add-button"
                        @click="createInsight()">
                    Save Insight
                </button>
            </div>
        </div>
        <div class="grid" v-if="container !== null" ref="container">

            <div class="item" v-for="(item, index) in texts" v-bind:key="index"
                 :class="{'item-selected': item.selected}"
                 @click="toggleSelected(item)"
            >
                <div class="item-header">

                </div>
                <div class="item-content">
                   {{item.text}}
                </div>

                <div class="item-footer">
                    <router-link :to="{name: 'Code', params: {documentId: item.document_id}}"> {{names[item.document_id]}}</router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script>

import Muuri from 'muuri';
import {actions} from "@/store";
import {API} from "@/core/API.ts";

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

        selectedTexts() {
            return this.texts.filter(i => i.selected);
        },

        projectId() {
            return parseInt(this.$route.params.projectId);
        }
    },
    name: "Bucket",
    data() {
        return {
            grid: null,
            insightText: "",
        }
    },

    methods: {
        async init() {
            if (this.container === null) {
                // files are needed to get the file names
                await this.$store.dispatch(actions.file.getFiles, this.projectId);
                await this.$store.dispatch(actions.code.INIT_CONTAINERS, this.projectId);
            }

            setTimeout(() => this.grid = new Muuri('.grid', {
                dragHandle: 'item-header', dragEnabled: true
            }), 50);
        },

        async createInsight() {
            const textIds = this.selectedTexts.map(t => t.id);

            await API.insight.post(this.projectId, this.insightText, textIds);
            // clear text after saving it
            this.insightText = "";
            for(let text of this.texts) {
                text.selected = false;
            }
        },

        toggleSelected(item) {
            this.$set(item, "selected", !item.selected);
        },
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
    margin: 20px;
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

.insights {
    width: 100%;
}

textarea {
    width: 100%;
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

.item-selected {
    box-shadow: 5px 5px 5px var(--blue);
}

</style>