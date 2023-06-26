const state = {
    msg:null,
}

const mutations = {
    SET_MSG: (state, msg) => {
        state.msg = msg
    },
}

const actions = {
  setMsg({ commit }, msg) {
    commit('SET_MSG', msg)
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
