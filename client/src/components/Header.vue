<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">Navbar</a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <router-link class="nav-link active" aria-current="page" to="/">Home</router-link>
          </li>
          <li class="nav-item">
            <a 
              v-if="store.token" 
              @click="logout"
              class="nav-link" 
              href="javascript:void(0);" 
            >Logout</a>
            <router-link v-else class="nav-link" to="/login">Login</router-link>
          </li>
        </ul>

        <span class="navbar-text">
          {{ store.user.first_name ?? ''}}
        </span>
      </div>
    </div>
  </nav>
</template>

<script>
import { store } from '@/store'
import router from '@/router'
export default {
  data() {
    return {
      store
    }
  },
  methods: {
    logout() {
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({
          token: store.token
        })
      }

      fetch("http://localhost:8092/api/logout", requestOptions)
        .then(response => response.json())
        .then(() => {
          store.token = null
          store.user = {}
          document.cookie = '_site_data=; Path=/; ' + 
            'SameSite=Strict; Secure; ' +
            'Expires=Thu, 01 Jan 1970 00:00:01 GMT;'

          router.push({name: 'Login'})
        })
    }
  }
}
</script>
