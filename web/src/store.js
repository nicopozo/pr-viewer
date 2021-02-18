// store.js
import Vue from 'vue';
import Vuex from 'vuex';

// register vuex as a plugin with vue in Vue 2
Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        token: '',
        username: '',
    },
    mutations: {
        SET_TOKEN(state, newValue) {
            if (this.debug) console.log('setToken triggered with', newValue)
            state.token = newValue
        },
        SET_USERNAME(state, newValue) {
            if (this.debug) console.log('setUsername triggered with', newValue)
            state.username = newValue
        }
    },
    getters: {
        token: state => state.token,
        username: state => state.username,
    }
});