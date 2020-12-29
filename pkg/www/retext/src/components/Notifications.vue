<template>
    <transition-group name="notification-item" tag="p" class="notification-list">
        <div v-for="(n, index) in visibleNotifications" :key="n.id"
             class="notification-item"
        >
            <div class="item-visible">
                <div v-b-modal="'notification-' + index" class="item-first">{{ n.text }}</div>
                <div class="item-last"><i class="fa fa-times"
                        @click.stop="ignore(n.id)"></i></div>
            </div>
            <b-modal
                :title="n.title ? n.title : null"
                :id="'notification-' + index"
                @ok="resolve(n.id)"
                @cancel="ignore(n.id)"
                cancel-title="Ignore"
                :ok-title="n.resolveCallback?'Resolve':'OK'"
                :ok-only="!n.ignoreCallback"
            >
                <p>{{ n.text }}</p>
            </b-modal>
        </div>
    </transition-group>
</template>

<script>
import {mutations} from "@/store";
import {mapGetters} from "vuex";
export default {
    name: "Notifications",
    data: function () {
        return {
        };
    },
    computed: {
        ...mapGetters(["notifications"]),
        visibleNotifications() {
            return this.notifications.filter(n => !n.dismissed);
        },
    },
    methods: {
        dismiss(id) {
            this.$store.commit(mutations.notification.dismiss, id);
        },
        ignore(id) {
            this.$store.commit(mutations.notification.ignore, id);
        },
        resolve(id) {
            this.$store.commit(mutations.notification.resolve, id);
        },
    },
}
</script>

<style scoped>
.notification-list {
    width: 200px;
    height: 95%;
    bottom: 0;
    overflow: hidden;
    display: flex;
    flex-flow: column;
    justify-content: flex-end;
    pointer-events: none;
}

.notification-item {
    pointer-events: all;
    width: 100%;
    border: 1px solid var(--button-prim);
    border-radius: .25em;
    background: var(--bg-prim);
    display: inline-block;
}

/*this div/class added almost entirely to vertically align the i.fa-times*/
.item-visible {
    display: flex;
    align-items: center;
    width: 100%;
}

.notification-item:nth-child(n+1) {
    margin-top: 5px;
}

.item-first {
    flex-basis: 90%;
    padding-left: 5px;
    /*prevent dotted, white outline from appearing after selecting a div(happens in FF)*/
    outline: none;
}

.item-last {
    flex-basis: 10%;
    padding-right: 5px;
}
.item-last i {
    /*make the X(dismiss button) easier to click*/
    margin: 5px;
    cursor: pointer;
}


/*animations - see here https://vuejs.org/v2/guide/transitions.html*/
.notification-item {
    transition: all 1s;
    z-index: 2;
}
/*following the *-enter, *-leave-to, and *-leave-active are applied automatically through use of VueJS' transitions*/
/*noinspection CssUnusedSymbol*/
.notification-item-enter,
.notification-item-leave-to
{
    transform: translateX(200px);
    opacity: 0;
}
/*noinspection CssUnusedSymbol*/
.notification-item-leave-active {
    position: absolute;
    z-index: 1;
}
</style>
