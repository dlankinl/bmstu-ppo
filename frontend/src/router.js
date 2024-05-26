import { createWebHistory, createRouter } from "vue-router";
import Home from "./components/Home.vue";
import Login from "./components/Login.vue";
import Register from "./components/Register.vue";
import EntrepreneursPage from './components/EntrepreneursPage.vue';
import EntrepreneurPage from './components/EntrepreneurPage.vue';
import UpdateEntrepreneurPage from './components/UpdateEntrepreneurPage.vue';
import UpdateCompanyPage from './components/UpdateCompanyPage.vue';
import CompanyPage from './components/CompanyPage.vue';
import CompaniesPage from './components/CompaniesPage.vue';
import ProfilePage from "./components/ProfilePage.vue";
import CompanyCreatePage from './components/CompanyCreatePage.vue'

const BoardAdmin = () => import("./components/BoardAdmin.vue")
const BoardModerator = () => import("./components/BoardModerator.vue")
const BoardUser = () => import("./components/BoardUser.vue")

const routes = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/home",
    component: Home,
  },
  {
    path: "/login",
    component: Login,
  },
  {
    path: "/register",
    component: Register,
  },
  {
    path: "/profile",
    name: "ProfilePage",
    component: ProfilePage,
  },
  {
    path: '/entrepreneurs',
    name: 'EntrepreneursPage',
    component: EntrepreneursPage,
    props: route => ({ page: parseInt(route.query.page) || 1 })
  },
  {
    path: '/entrepreneurs/:id',
    name: 'EntrepreneurPage',
    component: EntrepreneurPage,
    props: true
  },
  {
    path: '/entrepreneurs/:id/update',
    name: 'UpdateEntrepreneurPage',
    component: UpdateEntrepreneurPage,
    props: true
  },
  {
    path: '/companies/:id',
    name: 'CompanyPage',
    component: CompanyPage,
    props: true
  },
  {
    path: '/companies/:id/update',
    name: 'UpdateCompanyPage',
    component: UpdateCompanyPage,
    props: true
  },
  {
    path: '/entrepreneurs/:id/companies',
    name: 'CompaniesPage',
    component: CompaniesPage,
    props: true
  },
  {
    path: '/companies/create',
    name: 'CompanyCreatePage',
    component: CompanyCreatePage,
    props: true
  }
  // {
  //   path: '/entrepreneurs/:id/update',
  //   name: 'EntrepreneurUpdatePage',
  //   component: EntrepreneurUpdatePage,
  //   meta: { requiresAdmin: true }
  // }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const publicPages = ['/login', '/register', '/home', '/entrepreneurs'];
  const authRequired = !publicPages.includes(to.path) && !to.path.startsWith('/entrepreneurs/');
  const loggedIn = localStorage.getItem('user');

  if (authRequired && !loggedIn) {
    next('/login');
  } else {
    next();
  }
});

export default router;