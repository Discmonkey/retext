import axios from 'axios';
import {ModelFile} from "@/model/modelFile";
import {Id} from "@/model/id";
import {Demo} from "@/model/demo";
import {Source} from "@/model/source";
import {Code} from "@/model/code";
import {CodeContainer} from "@/model/codeContainer";
import {DocumentText} from "@/model/documentText";
import {Insight} from "@/model/insight";


const url = (endpoint: string, params: {[key: string]: any} = {}) => {

    if (Object.keys(params).length === 0) {
        return endpoint;
    }

    return endpoint + "?" + Object.keys(params).map(key => `${key}=${params[key]}`).join('&');
}


export const API = {

    code: {
        async get(code_id: Id): Promise<Code> {
            const res = await axios.get(url("/code", {code_id}));

            return res.data;
        },

        async post(code: Code): Promise<Code> {
            const res = await axios.post(url("/code"), code);

            return res.data;
        }
    },

    code_container: {
        async post(project_id: Id): Promise<CodeContainer> {
            const res = await axios.post(url("/code_container", {project_id}));

            return res.data;
        }
    },

    code_containers: {
        async get(project_id: Id): Promise<Array<CodeContainer>> {
            const res = await axios.get(url("/code_containers", {project_id}));

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

    document_text: {
        async post(code_id: Id, text: DocumentText): Promise<DocumentText> {
            const res = await axios.post(url("/document_text", {code_id}), text);

            return res.data;
        },

        async delete(text_id: Id) {
            await axios.delete(url("/document_text", {text_id}));
        }
    },

    files: {
        async get(project_id: Id): Promise<Array<ModelFile>> {
            const res = await axios.get(url("/files", {project_id}));

            return res.data;
        }
    },

    insight: {
        async post(project_id: Id, value: string, text_ids: Array<Id>) {
            const res = await axios.post(url("/insight", {project_id}), {value, text_ids});

            return res.data;
        },
    },

    insights: {
        async get(project_id: Id): Promise<Array<Insight>> {
            const res = await axios.get(url("/insights", {project_id}));

            return res.data;
        },
    },

    insight_text: {
        async delete(insight_text_id: Id) {
            await axios.delete(url("/insight_text", {insight_text_id}));
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