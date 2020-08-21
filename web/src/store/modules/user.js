// import { login, logout, getInfo } from '@/api/user'
import { login, userInfo, logout } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { resetRouter } from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    name: '',
    avatar: '',
    id: null,
    last_login_at: null,
    email: null
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_ID: (state, id) => {
    state.id = id
  },
  SET_EMAIL: (state, email) => {
    state.email = email
  },
  SET_LAST_LOGIN_AT: (state, last_login_at) => {
    state.last_login_at = last_login_at
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      console.log('here')
      login({ username: username.trim(), password: password }).then(response => {
        console.log(response)
        const { data: { token, lifetime }} = response
        console.log(token, lifetime)
        commit('SET_TOKEN', token)
        setToken(token, lifetime)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      userInfo().then(response => {
        const { data } = response

        console.log(data)
        if (!data) {
          return reject('Verification failed, please Login again.')
        }

        // created_at: "0001-01-01T00:00:00Z"
        // deleted_at: null
        // email: "11@qq.com"
        // id: 2
        // last_login_at: "2020-08-18T22:30:31+08:00"
        // updated_at: "2020-08-18T22:30:31+08:00"
        // user_name: "bbc"
        const { user_name: name, avatar, id, email, last_login_at } = data

        commit('SET_ID', id)
        commit('SET_EMAIL', email)
        commit('SET_LAST_LOGIN_AT', last_login_at)
        commit('SET_NAME', name)
        commit('SET_AVATAR', avatar)
        resolve(data)
      }).catch(error => {
        removeToken() // must remove  token  first
        resetRouter()
        commit('RESET_STATE')
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        removeToken() // must remove  token  first
        resetRouter()
        commit('RESET_STATE')
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

