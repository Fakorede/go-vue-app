<template>
  <div>
    <Header />
    <div>
      <router-view />
    </div>
    <Footer />
  </div>
</template>

<script>
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import { store } from '@/store'

const getCookie = (name) => {
  return document.cookie.split("; ").reduce((r, v) => {
    return v.split("=")[0] === name 
      ? decodeURIComponent(v.split("=")[1]) 
      : r
  }, "")
}

export default {
  name: 'App',
  components: {
    Header,
    Footer,
  },
  data() {
    return {
      store
    }
  },
  beforeMount() {
    // check for cookie and update store
    let data = getCookie("_site_data")
    if (data !== "") {
      let cookieData = JSON.parse(data)
      store.token = cookieData.token.token
      store.user = {
        id: cookieData.user.id,
        first_name: cookieData.user.first_name,
        last_name: cookieData.user.last_name,
        email: cookieData.user.email,
      }
    }
  },
}
</script>

<style>

</style>
