<template>
    <div class="container-fluid">
        <div class="row">
            <button class="btn btn-primary offset-md-5 col-md-7"
                    @text-drop="textDrop(newCat, $event.detail, $event)"
                    @click="createCategory(categories, null)"
            >Add a New Code
            </button>

        </div>
        <div @text-drop="textDrop(category, $event.detail, $event)"
             class="col-md-12 category-wrapper"
             v-for="category in categories" :key="category.id">
            <div class="row">
                <div class="col-sm-7"><b>
                    <category-drop-zone :category="category">{{category.name}}</category-drop-zone>
                </b></div>
                <div class="col-sm-5 row">
                    <div class="col-md-4">
                        <category-drop-zone
                                :category="newCat"
                        >
                            <div
                                @click="createCategory(category.subcategories, category.id)"
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
                <category-drop-zone class="row rounded-series subcategory" style="width:100%"
                                    v-for="subCat in category.subcategories" :key="subCat.id"
                                    :category="subCat"
                >
                    <span class="col-md-9">{{subCat.name}}</span>
                    <span class="col-md-3">{{getTextsLength(subCat)}}</span>
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
            id: "",
            name: "",
            texts: [{
                documentID: "",
                text: ""
            }]
        }]
    }
    export default {
        name: 'CategoryList',
        components: {CategoryDropZone},
        props: ["channel"],
        data: () => {
            let newCat = {name: "New", id: 0, texts: []};
            return {
                categories: [],
                newCat: newCat,
            }
        },
        mounted() {
            this.axios.get("/category/list").then((res) => {
                this.categories = res.data
                // todo: change back-end to return this response...
                this.categories = [{
                    "id": 1,
                    "name": "test #1",
                    "texts": [{
                        "documentID": "test_doc.txt",
                        "text": "nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of"
                    }],
                    "subcategories": [{
                        "id": 2,
                        "name": "cat 2",
                        "texts": [{
                            "documentID": "test_doc.txt",
                            "text": "men, living and dead, who struggled here, have consecrated"
                        }, {
                            "documentID": "test_doc.txt",
                            "text": "men, living and dead, who struggled here, have consecrated"
                        }]
                    }, {
                        "id": 3,
                        "name": "three",
                        "texts": [{
                            "documentID": "test_doc.txt",
                            "text": "they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us -- that from these honored dead we take increased devotion"
                        }]
                    }]
                }];
            });
        },
        methods: {
            textDrop: function (parentCategory, packet, e) {
                e.stopPropagation(); // stop the even

                let category = packet.data.category;

                console.log(parentCategory, packet);

                // unless an error happens, this function will get called
                let associate = (cat) => {
                    this._actualAssociate(cat, packet.data.words, packet.callback);
                };

                if (parentCategory.id === 0) {
                    this.createCategory(this.categories, null).then(
                        associate,
                        this._cancelCreateCategory
                    )
                } else if (!category) {
                    // dropped on the category-wrapper but not in a designated drop-zone.
                    // todo: make the whole category-wrapper a drop zone?
                    return false;
                } else if (category.id === 0) {
                    this.createCategory(parentCategory.subcategories, parentCategory.id).then(
                        associate,
                        this._cancelCreateCategory
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
                    // todo: remove this if subcat...
                    if (!newCat.subcategories) {
                        newCat.subcategories = [];
                    }
                    categories.push(newCat);
                    return newCat;
                }, function (res) {
                    // todo: alert the user of failure
                    console.log("an error occurred", res);
                    return false;
                });
            },
            _actualAssociate: function (category, words, callback) {
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
            _cancelCreateCategory: function(userCancelled) {
                if(userCancelled) {
                    // no worries?
                } else {
                    // the use should have been notified already in createCategory(); do nothing
                }
                // why this function: prevent console errors for uncaught promises
            },
            getTextsLength(category) {
                let length = category.texts.length;

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