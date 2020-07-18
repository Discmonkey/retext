<template>
    <div>
    <h3>hello </h3>
    <table class="table table-striped">
        <thead>
            <tr>
                <th scope="col" v-for="(val, idx) in columns" v-bind:key="idx">
                    {{val}}
                </th>
            </tr>
        </thead>

        <tbody>
            <tr v-for="(row, idx) in rows" v-bind:key="idx">
                <td v-for="(val, idx) in row" v-bind:key="idx">{{val}}</td>
            </tr>
        </tbody>
    </table>
    </div>
</template>

<script>
export default {
    name: "DemoDisplay",
    props: ["filename"],
    data() {
        return {
            columns: [],
            rows: [],
        }
    },
    methods: {
        loadDocument: function (documentName) {
            this.axios.get(`/file/load?key=${documentName}`).then(res => {
                if (!res.data) {
                    return;
                }

                console.log(res.data);

                let columns = res.data.Attributes.Columns;
                let values = res.data.Attributes.Values;

                for (let column of columns) {
                    this.columns.push(column);

                    let counter = 0;
                    for (let value of values[column]) {
                        if (this.rows.length <= counter) {
                            this.rows.push([]);
                        }

                        this.rows[counter].push(value);
                        counter++;
                    }

                    counter += 1;
                }


            }).catch(e => (console.error(e)));
        }
    },

    mounted() {
        this.loadDocument(this.filename);
    }
}
</script>

<style scoped>

</style>