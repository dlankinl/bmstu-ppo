import { createWebHistory, createRouter } from "vue-router";
import HomePage from './components/HomePage.vue';
import Login from "./components/Login.vue";
import Register from "./components/Register.vue";
import EntrepreneursPage from './components/EntrepreneursPage.vue';
import EntrepreneurPage from './components/EntrepreneurPage.vue';
import UpdateEntrepreneurPage from './components/UpdateEntrepreneurPage.vue';
import UpdateCompanyPage from './components/UpdateCompanyPage.vue';
import CompanyPage from './components/CompanyPage.vue';
import CompaniesPage from './components/CompaniesPage.vue';
import ProfilePage from "./components/ProfilePage.vue";
import CompanyCreatePage from './components/CompanyCreatePage.vue';
import UserSkillsPage from './components/UserSkillsPage.vue';
import SkillsPage from './components/SkillsPage.vue';
import UpdateSkillPage from './components/UpdateSkillPage.vue';
import CreateSkillPage from "./components/CreateSkillPage.vue";
import NotFound from "./components/NotFound.vue";
import ActivityFieldsPage from "./components/ActivityFieldsPage.vue";
import CreateActivityFieldPage from './components/CreateActivityFieldPage.vue';
import UpdateActivityFieldPage from './components/UpdateActivityFieldPage.vue';
import CreateFinReportPage from './components/CreateFinReportPage.vue';
import UpdateFinReportPage from './components/UpdateFinReportPage.vue';
import CompanyFinancialsPage from './components/CompanyFinancialsPage.vue';
import ContactsPage from "./components/ContactsPage.vue";
import CreateContactPage from './components/CreateContactPage.vue';
import UpdateContactPage from './components/UpdateContactPage.vue';
import ReviewPage from "./components/ReviewPage.vue";
import CreateReviewPage from './components/CreateReviewPage.vue';
import MyReviewsPage from './components/MyReviewsPage.vue';

const routes = [
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
    path: '/home',
    name: 'HomePage',
    component: HomePage,
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
    path: '/entrepreneurs/:id/contacts',
    name: 'ContactsPage',
    component: ContactsPage,
    props: true
  },
  {
    path: '/companies/create',
    name: 'CompanyCreatePage',
    component: CompanyCreatePage,
    props: true
  },
  {
    path: '/entrepreneurs/:id/skills',
    name: 'UserSkillsPage',
    component: UserSkillsPage,
    props: true
  },
  {
    path: '/skills',
    name: 'SkillsPage',
    component: SkillsPage,
    props: true
  },
  {
    path: '/skills/:id/update',
    name: 'UpdateSkillPage',
    component: UpdateSkillPage,
    props: true
  },
  {
    path: '/skills/create',
    name: 'CreateSkillPage',
    component: CreateSkillPage,
    props: true
  },
  {
    path: '/:notFound',
    component: NotFound
  },
  {
    path: '/activity-fields',
    name: 'ActivityFieldsPage',
    component: ActivityFieldsPage,
    props: true
  },
  {
    path: '/activity-fields/create',
    name: 'CreateActivityFieldPage',
    component: CreateActivityFieldPage,
    props: true
  },
  {
    path: '/activity-fields/:id/update',
    name: 'UpdateActivityFieldPage',
    component: UpdateActivityFieldPage,
    props: true
  },
  {
    path: '/companies/:id/financials/create',
    name: 'CreateFinReportPage',
    component: CreateFinReportPage,
    props: true
  },
  {
    path: '/companies/:id/financials',
    name: 'CompanyFinancialsPage',
    component: CompanyFinancialsPage,
    props: true
  },
  {
    path: '/financials/:id',
    name: 'UpdateFinReportPage',
    component: UpdateFinReportPage,
    props: true
  },
  {
    path: '/contacts/create',
    name: 'CreateContactPage',
    component: CreateContactPage,
    props: true
  },
  {
    path: '/contacts/:id/update',
    name: 'UpdateContactPage',
    component: UpdateContactPage,
    props: true
  },
  {
    path: '/entrepreneurs/:id/reviews',
    name: 'ReviewPage',
    component: ReviewPage,
    props: true
  },
  {
    path: '/entrepreneurs/:id/reviews/create',
    name: 'CreateReviewPage',
    component: CreateReviewPage,
    props: true
  },
  {
    path: '/profile/reviews',
    name: 'MyReviewsPage',
    component: MyReviewsPage,
    props: true
  }
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