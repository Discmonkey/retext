import axios from 'axios';
import {Project} from "@/model/project";
import {Id} from "@/model/id";

export const mutations = {
    ADD_PROJECT: "addProject",
    SELECT_PROJECT: "selectProject",
}

export const actions = {
    LOAD_PROJECTS: "loadProjects",
    ADD_PROJECT: "addProject",
    SELECT_PROJECT: "projectId",
}

interface ProjectState {
    projects: Array<Project>;
    currentProject: Project;
}

type Commit = (s: string, a: any) => void;

export const ProjectModule = {
    state: {
        projects: [],
        currentProject: null,
    } as unknown as ProjectState,

    mutations: {

        [mutations.ADD_PROJECT](state: ProjectState, project: Project) {
            for (const projectOld of state.projects) {
                if (project.id === projectOld.id) {
                    return;
                }
            }

            state.projects.push(project);
        },

        [mutations.SELECT_PROJECT](state: ProjectState, projectId: Id) {
            const index = state.projects.findIndex(project => project.id === projectId);

            if (index >= 0) {
                state.currentProject = state.projects[index];
            }
        }

    },

    actions: {
        async [mutations.ADD_PROJECT]({commit}: {commit: Commit}, payload: Project) {
            const res = await axios.post("/project/create", payload);
            commit(mutations.ADD_PROJECT, res.data);
        },

        async [actions.LOAD_PROJECTS]({commit}: {commit: Commit}) {
            const res = await axios.get("/project/list")

            res.data.forEach((project: Project) => commit(mutations.ADD_PROJECT, project));
        },

        async [actions.SELECT_PROJECT]({commit, state}: {commit: Commit; state: ProjectState}, payload: Id) {
            if (state.currentProject != null && state.currentProject.id === payload) return;

            const index = state.projects.findIndex(project => project.id === payload);

            if (index === -1) {
                const res = await axios.get(`/project?projectId=${payload}`);
                commit(mutations.ADD_PROJECT, res.data);
            }

            commit(mutations.SELECT_PROJECT, payload);

        }
    },

    getters: {
        currentProject(state: ProjectState): Project {
            return state.currentProject;
        },

        projects(state: ProjectState): Array<Project> {
            return state.projects;
        }
    }
}