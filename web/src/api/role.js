import request from '@/utils/request'

export const index = (params) => {
  console.log(params)
  return request({
    method: 'get',
    url: '/api/admin/roles',
    params
  })
}

export const store = (data) => {
  return request({
    method: 'post',
    url: `/api/admin/roles`,
    data
  })
}

export const show = (id) => {
  return request({
    method: 'get',
    url: `/api/admin/roles/${id}`
  })
}

export const update = (id, data) => {
  return request({
    method: 'put',
    url: `/api/admin/roles/${id}`,
    data
  })
}

export const destroy = (id) => {
  return request({
    method: 'delete',
    url: `/api/admin/roles/${id}`
  })
}

export const allRoles = (id) => {
  return request({
    method: 'get',
    url: `/api/admin/all_roles`
  })
}
