import { createStore } from 'vuex'

export default createStore({
  state: {
    turnos: {}
  },
  getters: {
  },
  mutations: {
    setTurnos(state, payload) {
      state.turnos = payload
      console.log(state.turnos)
    }
  },
  actions: {
    async fetchData({commit}) {
      try {
        const response = await fetch('http://localhost:4000/')
        const data = await response.json;
        commit('setTurnos', data)
      } catch (error) {
        console.log(error)
      }
    }
  },
  modules: {
  }
})
