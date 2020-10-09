import Vue from 'vue'
import VueRouter from 'vue-router'
import Code from '@/views/Code'
import Upload from "@/views/Upload";
import Demo from "@/views/Demo";
import Home from "@/Home";
import Project from "@/views/Project";

Vue.use(VueRouter);

// Set the landing page by changing homePath below.
// homePath gets "injected" into `routes` below programmatically.
// const homePath = '/projects';

const routes = [
    {
        path: '/projects',
        name: 'Projects',
        component: Project,
        alias: '/'
    },
    {
        path: '/project/:projectId',
        name: 'project',
        component: Home,
        children: [
            {
                path: 'code/:documentId',
                name: 'Code',
                component: Code
            },
            {
                path: 'demo/:documentId',
                name: 'Demo',
                component: Demo
            },
            {
                path: 'upload',
                name: 'Upload',
                component: Upload,
            },
        ],
    },
];


const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

router.afterEach((to, from) => {
    console.log(to, from);
})

export default router
