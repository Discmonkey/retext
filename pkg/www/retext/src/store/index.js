import Vue from "vue"
import Vuex from "vuex";

Vue.use(Vuex)

function prepareContainer(backendCodeContainer) {
    const main = backendCodeContainer.subcodes.shift();
    return {
        containerId: backendCodeContainer.containerId,
        main,
        subcodes: backendCodeContainer.subcodes
    }
}
async function createCode(containerId, name) {
    return (await Vue.axios.post("/code/create", {
                code: name,
                containerId
            })).data
}

export const store = new Vuex.Store({
    state: {
        containers: [],
        idToContainer: {}
    },
    getters: {
        containers(state) {
            return state.containers
        },
        idToContainer(state) {
            return state.idToContainer;
        }
    },
    mutations: {
        addContainer(state, container) {
            const codeContainer = prepareContainer(container);

            state.containers.push(codeContainer);
            state.idToContainer[codeContainer.containerId] = codeContainer;
        },
        addCode(state, {containerId, code}) {
            state.idToContainer[containerId].subcodes.push(code);
        }
    },
    actions: {
        async createContainer(context, {name}) {
            const res = await Vue.axios.post("/code/container/create");
            const containerId = res.data.ContainerId;

            const code = await createCode(containerId, name);

            const newContainer = {
                containerId: containerId,
                main: code,
                subcodes: [],
            }

            context.commit("addContainer", newContainer);

            return code;
        },
        async createCode(context, {containerId, name}) {
            const code = await createCode(containerId, name);

            context.commit("addCode", {containerId, code});

            return code;
        },
        initContainers: (context) => {
            Vue.axios.get("/code/list").then((res) => {
                const containers = res.data;

                for(const c of containers) {
                    context.commit("addContainer", c);
                }
            })
        }
    }
})