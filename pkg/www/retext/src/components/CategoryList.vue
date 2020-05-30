<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-4 category"
                 v-for="cat in categories"
                 :key="cat.name">
                <category-drop-zone :category="cat" @category-drop="associate($event)"/>
            </div>
        </div>
        <div class="row">
            <div class="col-md-12 category new-category">
                <category-drop-zone :category="newCat" @category-drop="newCategoryAssociate()"/>
            </div>
        </div>
    </div>
</template>

<script>
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
    // todo: goland says this is unused. should i not export?
    export default {
        name: 'CategoryList',
        components: {CategoryDropZone},
        props: ["channel"],
        data: () => {
            let newCat = {name: "New"}
            newCat.id = newCat.name;
            return {
                categories: [],
                newCat: newCat
            }
        },
        mounted() {
            this.axios.get("/category/list").then((res) => {
                this.categories = res.data.categories
            });
        },
        methods: {
            _actualAssociate: function(categoryID, words) {
                this.axios.post("/category/associate", {
                    key: words.documentID,
                    categoryID: categoryID,
                    text: words.text
                }).then((res) => {
                    // success toast or something
                    console.log("cl _aA success", res);
                }, (res) => {
                    // failed toast or something (or maybe on a finally()?
                    // failed to send
                    console.log("cl _aA success", res);
                });
            },
            associate: function(categoryID) {
                let x= this.channel.receive();
                console.log(x);
                x.then(words => {
                    this._actualAssociate(categoryID, words);
                }, (o) => {
                    console.log('failed', o);
                }).catch(() => {});
            },
            newCategoryAssociate: function() {
                let x= this.channel.receive();
                console.log(x);
                x.then(words => {
                    if(!words.documentID) {
                        return false;
                    }
                    let newCatName = prompt("Name of new category?");
                    if (newCatName === false) {
                        // don't grey out text if this is false... oh boy...
                        return false;
                    }

                    this.axios.post("/category/create", {category: newCatName})
                    .then((res) => {
                        let newCat = res.data;
                        this.categories.push(newCat);
                        this._actualAssociate(newCat.id, words);
                    });
                }, (o) => {
                    console.log('failed2', o);
                }, () => {});
            }
        }
    }
</script>

<style>
    .category:nth-child(n+1) {
        margin-top: 1%;
    }
</style>