<template>
    <div class="container-fluid">
        <div class="row">
            <button class="btn btn-primary offset-md-5 col-md-7"
                    @text-drop="textDrop(defaultNewCat, $event.detail, $event)"
                    @click="createCategory(categories, null)"
            >Add a New Code
            </button>

        </div>
        <div v-for="category in categories" :key="category.main.id"
             @text-drop="textDrop(category, $event.detail, $event)"
             class="col-md-12 category-wrapper">
            <div class="row">
                <div class="col-sm-7"><b>
                    <category-drop-zone :category="category.main">{{category.main.name}}</category-drop-zone>
                </b></div>
                <div class="col-sm-5 row">
                    <div class="col-md-4">
                        <category-drop-zone
                                :category="defaultNewCat.main"
                        >
                            <div
                                @click="createCategory(category.subcategories, category.main.id)"
                                class="btn-primary parent-category-option">+</div>
                        </category-drop-zone>
                    </div>
                    <div class="col-md-4">
                        <div class="btn-primary parent-category-option">D</div>
                    </div>
                    <div class="col-md-4">
                        <div class="btn-primary parent-category-option">{{getTextsLength(category)}}</div>
                    </div>
                </div>
            </div>
            <div style="margin-left: 7%" class="row">
                <category-drop-zone v-for="subCat in category.subcategories" :key="subCat.id"
                                    class="row rounded-series subcategory" style="width:100%"
                                    :category="subCat"
                >
                    <span class="col-md-9">{{subCat.name}}</span>
                    <span class="col-md-3">{{subCat.texts.length}}</span>
                </category-drop-zone>
            </div>
        </div>
    </div>
</template>

<script>
    // todo: add back colors
    import CategoryDropZone from "./CategoryDropZone";

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
        components: {CategoryDropZone},
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
                }, function (res) {
                    // todo: alert the user of failure
                    console.log("an error occurred", res);
                    return false;
                });
            },
            _actuallyAssociate: function (category, words, callback) {
                this.axios.post("/category/associate", {
                    key: words.documentID,
                    categoryID: category.id,
                    text: words.text
                }).then((res) => {
                    category.texts.push(words);
                    callback();
                    // "success" toast or something
                    console.log("cl _aA success", res);
                }, (res) => {
                    // "an error occurred" toast or something
                    console.log("cl _aA failed", res);
                });
            },
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
    .category-wrapper {
        border: 2px solid blue;
        border-radius: .25em;
        margin: 5px;
        /*padding-bottom: 5px;*/
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
</style>