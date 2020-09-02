import Vue from "vue"
import Vuex from "vuex";

Vue.use(Vuex)

export const getters = {
    CONTAINERS: "containers",
    ID_TO_CONTAINER: "idToContainer",
}

export const mutations = {
    ADD_CONTAINER: "addContainer",
    ADD_CODE: "addCode",
}
export const actions = {
    CREATE_CONTAINER: "createContainer",
    CREATE_CODE: "createCode",
    INIT_CONTAINERS: "initContainers",
}

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
        [getters.CONTAINERS]: function(state) {
            return state.containers
        },
        [getters.ID_TO_CONTAINER]: function(state) {
            return state.idToContainer;
        }
    },
    mutations: {
        [mutations.ADD_CONTAINER]: function(state, container) {
            const codeContainer = prepareContainer(container);

            state.containers.push(codeContainer);
            state.idToContainer[codeContainer.containerId] = codeContainer;
        },
        [mutations.ADD_CODE]: function(state, {containerId, code}) {
            state.idToContainer[containerId].subcodes.push(code);
        }
    },
    actions: {
        [actions.CREATE_CONTAINER]: async function(context, {name}) {
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
        [actions.CREATE_CODE]: async function(context, {containerId, name}) {
            const code = await createCode(containerId, name);

            context.commit("addCode", {containerId, code});

            return code;
        },
        [actions.INIT_CONTAINERS]: async function(context) {
            Vue.axios.get("/code/list").then((res) => {
                const containers = res.data;

                for(const c of containers) {
                    context.commit("addContainer", c);
                }
            })
        },
    }
})