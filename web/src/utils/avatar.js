import _ from 'lodash'

const randomSetAvatar = () => {
  const avatarRange = [1, 10]
  const path = 'avatars/' + _.random(avatarRange[0], avatarRange[1]) + '.png'
  localStorage.setItem('avatar', path)

  return path
}

const getAvatar = () => {
  localStorage.getItem('avatar')
}

const removeAvatar = () => {
  localStorage.removeItem('avatar')
}

export {
  randomSetAvatar,
  getAvatar,
  removeAvatar
}
