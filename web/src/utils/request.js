import axios from 'axios'
import { Message } from 'element-ui'
import store from '@/store'
import router from '../router'
import { getToken, removeToken } from '@/utils/auth'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent

    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = getToken()
    }
    return config
  },
  error => {
    // do something with request error
    console.log('request err :' + error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const res = response.data

    return Promise.resolve(res)
  },
  error => {
    console.log('error.data', error.response)
    if (error.response.status === 422) {
      const e = error.response.data.error
      for (const key in e) {
        if (Object.prototype.hasOwnProperty.call(e, key)) {
          console.log(e[key])

          Message({
            message: e[key],
            type: 'error',
            duration: 5 * 1000
          })
        }
      }
    } else if (error.response.status === 401) {
      const token = getToken()
      if (token) {
        Message.error('登陆过期，请重新登陆')
        removeToken()
      } else {
        Message.error('用户名或者密码错误')
      }
      router.push({ name: 'login' })
    } else {
      Message.error(error.response.data.msg)
    }

    return Promise.reject(error)
  }
)

export default service
