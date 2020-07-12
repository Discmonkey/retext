import Vue from 'vue'
import VueRouter from 'vue-router'
import Code from '../views/Code'
import Upload from "../views/Upload";
import Demo from "../views/Demo";

Vue.use(VueRouter);

const routes = [
    {
        path: '/code/:documentID',
        name: 'Code',
        component: Code
    },
    {
        path: '/demo/:documentID',
        name: 'Demo',
        component: Demo
    },
    {
        path: '/upload',
        name: 'Upload',
        component: Upload,
        alias: "/"
    },
    {
        path: '/about',
        name: 'About',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
    }
];

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

export default router
