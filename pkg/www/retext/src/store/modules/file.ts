import {Source} from "@/model/source";
import {Demo} from "@/model/demo";
import {WalkableSource} from "@/core/WalkableSource";
import {Id} from "@/model/id";
import vue from 'vue';
import {ModelFile, ModelFile as File} from "@/model/modelFile";
import {Commit, Dispatch} from "vuex";
import {API} from "@/core/API";

export const actions = {
    getFiles: "getFiles",
    postSource: "postSource",
    postDemo: "postDemo",
    getSource: "getSource",
    getDemo: "getDemo",
    selectSource: "selectSource",
    selectDemo: "selectDemo",
}

export const mutations = {
    addFile: "addFile",
    addSource: "addSource",
    addDemo: "addDemo",
    clearFiles: "clearFiles",
    selectSource: "selectSource",
    selectDemo: "selectDemo",
}

interface State {
    files: Array<File>,
    demos: {
        map: {[key: number]: Demo};
        current: Demo | null;
    },

    sources: {
        map: {[key: number]: Source};
        current: WalkableSource | null;
    },
}

export const Module = {

    state: {
        files: [],
        sources: {
            map: {},
            current: null,
        },
        demos: {
            map: {},
            current: null,
        },
    } as State,

    mutations: {
        [mutations.addDemo](state: State, payload: {demo: Demo, id: Id}) {
            if (payload.id in state.demos) {
                return;
            }

            vue.set(state.demos.map, payload.id, payload.demo);
        },

        [mutations.addFile](state: State, file: File) {
            const index = state.files.findIndex((f: File) => f.id == file.id);

            if (index == -1) {
                state.files.push(file);
            }
        },

        [mutations.addSource](state: State, payload: {source: Source, id: Id}) {
            if (payload.id in state.sources) {
                return;
            }

            vue.set(state.sources.map, payload.id, payload.source);
        },

        [mutations.clearFiles](state: State) {
            while (state.files.length > 0) {
                state.files.pop();
            }
        },

        [mutations.selectSource](state: State, payload: {source: Id}) {
            vue.set(state.sources, "current", new WalkableSource(state.sources.map[payload.source]));
        },

        [mutations.selectDemo](state: State, payload: {demo: Id}) {
            vue.set(state.demos, "current", state.demos.map[payload.demo]);
        }
    },

    actions: {
        async [actions.getFiles]({commit}: {commit: Commit}, payload: Id) {
            const files = await API.files.get(payload);

            files.forEach(file => {
                commit(mutations.addFile, file);
            })
        },

        async [actions.postSource]({commit}: {commit: Commit}, payload: {
            project: Id, formData: FormData
        }) {
            const uploaded = await API.source.post(payload.project, payload.formData);

            uploaded.forEach(f => commit(mutations.addFile, f));

        },

        async [actions.postDemo]({commit}: {commit: Commit}, payload: {
            project: Id, formData: FormData
        }) {
            const uploaded = await API.demo.post(payload.project, payload.formData);

            uploaded.forEach(f => commit(mutations.addFile, f));
        },

        async [actions.getDemo]({commit}: {commit: Commit}, payload: {
            fileId: Id
        }) {
            const parsed = await API.demo.get(payload.fileId);

            commit(mutations.addDemo, {
                demo: parsed, id: payload.fileId
            });

        },

        async [actions.getSource]({commit}: {commit: Commit}, payload: {
            fileId: Id
        }) {


            const parsed = await API.source.get(payload.fileId);

            commit(mutations.addSource, {
                source: parsed, id: payload.fileId
            });
        },

        async [actions.selectSource]({state, commit, dispatch}: {state: State; commit: Commit; dispatch: Dispatch}, payload: {
            fileId: Id
        }) {
            if (!(payload.fileId in state.sources.map)) {
                await dispatch(actions.getSource, payload);
            }

            commit(mutations.selectSource, {source: payload.fileId});
        },

        async [actions.selectDemo]({state, commit, dispatch}: {state: State; commit: Commit; dispatch: Dispatch}, payload: {
            fileId: Id
        }) {
            if (!(payload.fileId in state.demos.map)) {
                await dispatch(actions.getDemo, payload);
            }

            commit(mutations.selectDemo, {source: payload.fileId});
        }
    },

    getters: {
        source(state: State) {
            return state.sources.current;
        },

        demo(state: State) {
            return state.demos.current;
        },

        sources(state: State): Array<ModelFile> {
            return state.files.filter(file => file.type === 'KSOURCE');
        },

        demos(state: State): Array<ModelFile> {
            return state.files.filter(file => file.type === 'KDEMO');
        }
    }
}