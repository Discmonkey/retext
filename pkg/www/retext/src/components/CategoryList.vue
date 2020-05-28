<template>
    <div class="container">
        <div class="row">
            <div class="category" v-for="cat in categories" :key="cat.name">
                <category-drop-zone :category="cat" @category-drop="associate($event)"/>
            </div>
            <div class="category new-category">
                <category-drop-zone :category="{name:'New'}" @category-drop="newCategoryAssociate()"/>
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
    export default {
        name: 'CategoryList',
        components: {CategoryDropZone},
        props: ["channel"],
        data: () => {
            return {
                categories: [],
            }
        },
        mounted() {
            this.axios.get("/category/list").then((res) => {
                this.categories = res.data.categories
            });
            this.channel.receiveFunc = (r) => {
                let dropOn = document.querySelectorAll(".category-drop-zone:hover");
                if(!dropOn.length) {
                    // cancel
                }

                if(!r)
                    return;

                if(dropOn[0].classList.contains('new-category')) {
                    this.newCategoryAssociate(r);
                } else {
                    this.associate("", r)
                }
            }
            // document.addEventListener("mouseup", this._listenForMouseUp);
        },
        methods: {
            associate: function(categoryID) {
                this.channel.one({categoryID: categoryID}, (one, two, twoCb) => {
                    let words = two;

                    this.axios.post("/category/associate", {
                        key: words.documentID,
                        categoryID: categoryID,
                        text: words.text
                    }).then((res) => {
                        console.log(res);
                        if(twoCb)
                            twoCb(words);
                    }, (res) => {
                        // failed to send
                        console.log(res);
                    });
                });
            },
            /*associate: function(categoryID, r) {
                let r = this.channel.receive()

                if(!r)
                    return;

                this.axios.post("/category/associate", {
                    key: r.obj.documentID,
                    category: categoryID,
                    text: r.obj.text
                }).then((res) => {
                    console.log(res);
                    if(r.callback)
                        r.callback(true, r.obj);
                },(res) => {
                    console.log(res);
                    if(r.callback)
                        r.callback(false, r.obj);
                })
            },*/
            newCategoryAssociate: function(r) {
                console.log("new cat");
                let newCatName = prompt("Name of new category?");
                this.axios.post("/category/create", {category: newCatName})
                    .then((res) => {
                        let newCat = res.data;
                        this.categories.push(newCat);
                        this.associate(newCat.id, r);
                    });
            }
        },
        _listenForMouseUp: function() {
            let dropOn = document.querySelectorAll(".category-drop-zone:hover");
            if(dropOn.length) {
                console.log(dropOn)
            }
        },
        beforeDestroy() {
            // document.removeEventListener("mouseup", this._listenForMouseUp)
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    .category {
        border: 1px solid black;
    }

    .bg-navy{background-color:#001F3F}.bg-blue{background-color:#0074D9}.bg-aqua{background-color:#7FDBFF}.bg-teal{background-color:#39CCCC}.bg-olive{background-color:#3D9970}.bg-green{background-color:#2ECC40}.bg-lime{background-color:#01FF70}.bg-yellow{background-color:#FFDC00}.bg-orange{background-color:#FF851B}.bg-red{background-color:#FF4136}.bg-fuchsia{background-color:#F012BE}.bg-purple{background-color:#B10DC9}.bg-maroon{background-color:#85144B}.bg-white{background-color:#fff}.bg-gray{background-color:#aaa}.bg-silver{background-color:#ddd}.bg-black{background-color:#111}.navy{color:#001F3F}.blue{color:#0074D9}.aqua{color:#7FDBFF}.teal{color:#39CCCC}.olive{color:#3D9970}.green{color:#2ECC40}.lime{color:#01FF70}.yellow{color:#FFDC00}.orange{color:#FF851B}.red{color:#FF4136}.fuchsia{color:#F012BE}.purple{color:#B10DC9}.maroon{color:#85144B}.white{color:#fff}.silver{color:#ddd}.gray{color:#aaa}.black{color:#111}.border--navy{border-color:#001F3F}.border--blue{border-color:#0074D9}.border--aqua{border-color:#7FDBFF}.border--teal{border-color:#39CCCC}.border--olive{border-color:#3D9970}.border--green{border-color:#2ECC40}.border--lime{border-color:#01FF70}.border--yellow{border-color:#FFDC00}.border--orange{border-color:#FF851B}.border--red{border-color:#FF4136}.border--fuchsia{border-color:#F012BE}.border--purple{border-color:#B10DC9}.border--maroon{border-color:#85144B}.border--white{border-color:#fff}.border--gray{border-color:#aaa}.border--silver{border-color:#ddd}.border--black{border-color:#111}.fill-navy{fill:#001F3F}.fill-blue{fill:#0074D9}.fill-aqua{fill:#7FDBFF}.fill-teal{fill:#39CCCC}.fill-olive{fill:#3D9970}.fill-green{fill:#2ECC40}.fill-lime{fill:#01FF70}.fill-yellow{fill:#FFDC00}.fill-orange{fill:#FF851B}.fill-red{fill:#FF4136}.fill-fuchsia{fill:#F012BE}.fill-purple{fill:#B10DC9}.fill-maroon{fill:#85144B}.fill-white{fill:#fff}.fill-gray{fill:#aaa}.fill-silver{fill:#ddd}.fill-black{fill:#111}.stroke-navy{stroke:#001F3F}.stroke-blue{stroke:#0074D9}.stroke-aqua{stroke:#7FDBFF}.stroke-teal{stroke:#39CCCC}.stroke-olive{stroke:#3D9970}.stroke-green{stroke:#2ECC40}.stroke-lime{stroke:#01FF70}.stroke-yellow{stroke:#FFDC00}.stroke-orange{stroke:#FF851B}.stroke-red{stroke:#FF4136}.stroke-fuchsia{stroke:#F012BE}.stroke-purple{stroke:#B10DC9}.stroke-maroon{stroke:#85144B}.stroke-white{stroke:#fff}.stroke-gray{stroke:#aaa}.stroke-silver{stroke:#ddd}.stroke-black{stroke:#111}

</style>
