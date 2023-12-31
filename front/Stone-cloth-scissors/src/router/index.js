import Vue from 'vue'
import VueRouter from 'vue-router'
import home from '../views/home.vue'
import game from '../views/game.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: home
  },
  {
    path: '/game',
    name: 'game',
    component: game
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
