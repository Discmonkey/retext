import axios from 'axios';
import {ModelFile} from "@/model/modelFile";
import {Id} from "@/model/id";
import {Demo} from "@/model/demo";
import {Source} from "@/model/source";


const url = (endpoint: string, params: {[key: string]: any} = {}) => {

    if (Object.keys(params).length === 0) {
        return endpoint;
    }

    return endpoint + "?" + Object.keys(params).map(key => `${key}=${params[key]}`).join('&');
}


export const API = {

    files: {
        async get(project_id: Id): Promise<Array<ModelFile>> {
            const res = await axios.get(url("/files", {project_id}));

            return res.data;
        }
    },

    demo: {
        async get(file_id: Id): Promise<Demo> {
            const res = await axios.get(url("/demo", {file_id}));

            return res.data;
        },

        async post(project_id: Id, data: any): Promise<Array<ModelFile>> {
            const res = await axios.post(url("/demo", {project_id}),
                data,
                {
                    headers: {
                        "Content-Type": 'multipart/form-data'
                    }
                }
            )

            return res.data;
        }
    },

    source: {
        async get(file_id: Id): Promise<Source> {
            const res = await axios.get(url("/source", {file_id}));

            return res.data;
        },

        async post(project_id: Id, data: any): Promise<Array<ModelFile>> {
            const res = await axios.post(url("/source", {project_id}),
                data,
                {
                    headers: {
                        "Content-Type": 'multipart/form-data'
                    }
                }
            )

            return res.data;
        }
    }
}