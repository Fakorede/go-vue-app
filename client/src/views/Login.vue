<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">Login</h1>
                <hr>
                <form-tag @myevent="submitHandler" name="myform" event="myevent">

                    <text-input
                        v-model="email"
                        label="Email"
                        type="email"
                        name="email"
                        autocomplete="off"
                        required="true">
                    </text-input>

                    <text-input
                        v-model="password"
                        label="Password"
                        type="password"
                        name="password"
                        autocomplete="off"
                        required="true">
                    </text-input>

                    <hr>
                    <input type="submit" class="btn btn-primary" value="Login">
                </form-tag>
            </div>
        </div>
    </div>
</template>

<script>
import FormTag from '@/components/forms/FormTag.vue'
import TextInput from '@/components/forms/TextInput.vue'
import { store } from '@/store'
import router from '@/router'
import notie from 'notie'

export default {
    name: 'login',
    components: {
        FormTag,
        TextInput,
    },
    data() {
        return {
            email: "",
            password: "",
            store,
        }
    },
    methods: {
        submitHandler() {
            const payload = {
                "username": this.email,
                "password": this.password
            }

            const requestOptions = {
                method: "POST",
                body: JSON.stringify(payload)
            }

            fetch("http://localhost:8092/api/login", requestOptions)
                .then(response => response.json())
                .then((response) => {
                    if (response.error) {
                        console.log("Error: ",response.message)
                        notie.alert({
                            type: 'error',
                            text: response.message
                        })
                    } else {
                        store.token = response.data.token.token
                        store.user = {
                            id: response.data.user.id,
                            first_name: response.data.user.first_name,
                            last_name: response.data.user.last_name,
                            email: response.data.user.email,
                        }

                        // save to cookie
                        let date = new Date()
                        let expiryDays = 1
                        date.setTime(date.getTime() + (expiryDays * 24 * 60 * 60 * 1000))

                        const expires = "expires=" + date.toUTCString()
                        document.cookie = "_site_data="
                            + JSON.stringify(response.data)
                            + "; "
                            + expires
                            + "; path=/; SameSite=strict; Secure;"

                        router.push({name: 'Home'})
                    }
                })
        }
    },
}
</script>