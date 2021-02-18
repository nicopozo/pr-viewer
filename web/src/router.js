import Vue from 'vue'
import Router from 'vue-router'
import ListPullRequests from "./components/ListPullRequests.vue";
import Welcome from "./components/Welcome.vue"

Vue.use(Router)

export default new Router({
    //mode:"history",
    routes: [
        {
            path: '/',
            name: "Welcome",
            component: Welcome,
        },
        {
            path: '/list',
            name: "ListPullRequests",
            component: ListPullRequests,
        },
    ]
});