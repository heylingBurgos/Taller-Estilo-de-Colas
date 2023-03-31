import { createStore } from 'vuex'

export default createStore({
  state: {
    user: {}
  },
  getters: {
  },
  mutations: {
    setUser(state, payload) {
        state.user = payload
    }
  },
  actions: {
    enviarVista({commit}, usuario) {
        commit('setUser', usuario)
    }
  },
  modules: {
  }
})