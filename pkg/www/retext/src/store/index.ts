import Vue from "vue"
import Vuex, {Store} from "vuex";
import {getColor, invertColor} from "@/core/Colors";
import {Id} from "@/model/id";
import {CodeContainer} from "@/model/codeContainer";
import {Code} from "@/model/code";
import * as project from "@/store/modules/project";
import * as document from "@/store/modules/file";
import * as code from "@/store/modules/code";

Vue.use(Vuex)

export const mutations = {
    code: code.mutations,
    file: document.mutations,
    project: project.mutations,
}
export const actions = {
    code: code.actions,
    file: document.actions,
    project: project.actions,
}

export const store = new Vuex.Store({
    state: {},

    getters: {},
    mutations: {},
    actions: {},

    modules: {
        ProjectModule: project.Module,
        DocumentModule: document.Module,
        CodeModule: code.Module,
    }
})