<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-4 category"
                 v-for="cat in categories" :key="cat.id">
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
            newCat.id = 0;
            return {
                categories: [],
                newCat: newCat,
                currentCat: {},
            }
        },
        mounted() {
            this.axios.get("/category/list").then((res) => {
                this.categories = res.data
            });
        },
        methods: {
            _actualAssociate: function (categoryID, words, callback) {
                this.axios.post("/category/associate", {
                    key: words.documentID,
                    categoryID: categoryID,
                    text: words.text
                }).then((res) => {
                    callback();
                    // "success" toast or something
                    console.log("cl _aA success", res);
                }, (res) => {
                    // "an error occurred" toast or something
                    console.log("cl _aA failed", res);
                });
            },
            associate: function (categoryID) {
                if(!this.channel.isSending) {
                    return;
                }

                let packet = this.channel.receive();

                this._actualAssociate(categoryID, packet.data, packet.callback);
            },
            newCategoryAssociate: function () {
                if(!this.channel.isSending) {
                    return;
                }

                let packet = this.channel.receive();

                let newCatName = prompt("Name of new category?");
                if (newCatName === null) {
                    // prompt was cancelled
                    return false;
                }

                this.axios.post("/category/create", {category: newCatName})
                    .then((res) => {
                        let newCat = res.data;
                        this.categories.push(newCat);
                        this._actualAssociate(newCat.id, packet.data, packet.callback);
                    });

            }
        }
    }
</script>

<style>
    .category:nth-child(n+1) {
        margin-top: 1%;
    }
</style>