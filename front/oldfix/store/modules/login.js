// initial state
const state = {
	is_login: false,
	token: '',
}
// getters
const getters = {}

// actions
const actions = {}

// mutations
const mutations = {
	login(state, auth) {
		state.is_login = true
		state.token = auth.token
	},
	loginOut(state) {
		state.is_login = false
		state.token = ''
	},
}

export default {
	namespaced: true,
	state,
	getters,
	actions,
	mutations
}
