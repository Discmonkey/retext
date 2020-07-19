<template>
    <div class="container-fluid pad-20">
        <div class="row">
            <div class="col-12 text-right">
                <button class="btn btn-primary bold add-height" id="add-button"
                        @text-drop="textDrop(defaultNewCat, $event.detail, $event)"
                        @click="createCategory(categories, null)">
                    Add a New Code
                </button>
            </div>

        </div>
<!--        @text-drop="textDrop(category, $event.detail, $event)" class="category-wrapper">-->
        <draggable v-model="categories" group="people" @start="drag=true" @end="drag=false">
            <div v-for="code in categories" :key="code.main.id">
                <div class="spacer">
                    <div :style="style(code.main.id)">
                        <div class="top-container margined">
                            <h5 class="code-title">{{code.main.name}}</h5>

                            <div class="btn btn-primary float-right no-events just-number self-right">
                                {{code.main.texts.length}}
                            </div>

                            <button class="btn btn-primary float-right">
                                <i class="fa fa-plus"></i>
                            </button>


                        </div>

                        <div v-for="subcategory in code.subcategories" v-bind:key="subcategory.id" class="subcategory margin-top">
                            <div class="item">
                                <div class="col-8 rborder">
                                    {{subcategory.name}}
                                </div>

                                <div class="col-4 center-text">
                                    {{subcategory.texts.length}}
                                </div>
                            </div>
                        </div>

                    </div>
                </div>
            </div>
        </draggable>
    </div>
</template>

<script>
    import Draggable from 'vuedraggable';
    import {getColor} from '@/core/Colors';
    // eslint-disable-next-line no-unused-vars
    let categoriesTypes = {
        categories: [{
            id: 0,
            name: "",
            texts: [{
                documentID: "",
                text: ""
            }]
        }]
    }
    function prepareCategory(mainCat) {
        for(let i = 0; i <= mainCat.subcategories.length; ++i) {
            if(mainCat.subcategories[i].id === mainCat.main) {
                mainCat.main = mainCat.subcategories.splice(i, 1)[0];
                break;
            }
        }
        return mainCat;
    }
    export default {
        name: 'CategoryList',
        components: {Draggable},
        data: () => {
            let newCat = {main: {name: "New", id: 0, texts: []}, subcategories: {}};
            return {
                categories: [],
                defaultNewCat: newCat,
            }
        },
        mounted() {
            this.axios.get("/category/list").then((res) => {
                let categories = res.data

                this.categories = [];
                for(let c of categories) {
                    this.categories.push(prepareCategory(c));
                }
            });
        },
        methods: {
            textDrop: function (parentCategory, packet, e) {
                e.stopPropagation(); // stop the even

                let category = packet.data.category;

                // unless an error happens, this function will get called
                let associate = (cat) => {
                    if(typeof cat === "boolean" || !cat)
                        return;
                    let c = cat;
                    if(parentCategory.main.id === 0) {
                        c = cat.main;
                    }
                    this._actuallyAssociate(c, packet.data.words, packet.callback);
                };

                if (parentCategory.main.id === 0) {
                    this.createCategory(this.categories, null).then(
                        associate
                    )
                } else if (!category) {
                    // dropped on the category-wrapper but not in a designated drop-zone.
                    // TODO: make the whole category-wrapper a drop zone? (see template above)
                    return false;
                } else if (category.id === 0) {
                    this.createCategory(parentCategory.subcategories, parentCategory.main.id).then(
                        associate
                    )
                } else {
                    associate(category);
                }
            },
            createCategory: function (categories, parentCategoryID) {
                let newCatName = prompt("Name of new category?");
                if (newCatName === null) {
                    // prompt was cancelled
                    return Promise.reject(true).then(() => {}, () => {});
                }

                return this.axios.post("/category/create", {
                    category: newCatName,
                    parentCategoryID: parentCategoryID
                }).then(function (res) {
                    let newCat = res.data
                    if(!parentCategoryID) {
                        newCat = prepareCategory(newCat);
                    }
                    categories.push(newCat);
                    return newCat;
                }, function () {
                    // todo: alert the user of failure
                    return false;
                });
            },
            _actuallyAssociate: function (category, words, callback) {
                this.axios.post("/category/associate", {
                    key: words.documentID,
                    categoryID: category.id,
                    text: words.text
                }).then(() => {
                    category.texts.push(words);
                    callback();
                    // todo: "success" toast or something
                }, () => {
                    // todo: "an error occurred" toast or something
                });
            },

            style(id) {
                return `border-radius: 10px; border: 3px solid ${getColor(id)}; padding: 20px;`;
            }
        }
    }
</script>

<style scoped>
    .category-wrapper {
        border: 2px solid blue;
        border-radius: .25em;
        margin: 5px;
        /*padding-bottom: 5px;*/
    }

    .bold {
        font-weight: bold;
    }

    #add-button {
        height: 50px;
    }

    .add-height {
        width: 200px
    }

    .category-wrapper, .subcategory {
        margin-bottom: 1%;
    }

    .parent-category-option {
        padding: 0;
        margin-top: 50%;
        line-height: 1em;
        width: 1em;
        text-align: center;
    }

    .overflow {
        overflow: auto;
    }

    .pad-20 {
        padding: 20px;
    }

    .rounded-series * {
        border: 1px solid grey;
        height: 100%;
        text-align: center;
    }

    .rounded-series > :first-child {
        text-align: left;
        border-bottom-left-radius: .25em;
        border-top-left-radius: .25em;
    }

    .rounded-series > :last-child {
        border-bottom-right-radius: .25em;
        border-top-right-radius: .25em;
    }

    .spacer {
        margin-top: 5px;
        margin-bottom: 5px;
    }

    .text-right {
        text-align: right;
    }

    .wide {
        width: 100%;
    }

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