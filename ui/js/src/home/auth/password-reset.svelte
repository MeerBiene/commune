<script>
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

import {debounce} from '../../utils/utils.js'

export let email;

let passwordInput;
let password;
let repeatPasswordInput;
let repeatPassword;

function goBack() {
    dispatch('go-back', true)
}

function success() {
    dispatch('success', true)
}

onMount(() => {
    passwordInput.focus()
})



async function resetPassword() {
    let endpoint = `/password/reset`
    let data = {
        password: password,
        email: email,
    };
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let error;

function reset() {
    if(passwordInput?.value?.length < 8 ) {
        alert("That password is too short.")
        passwordInput.focus()
        return
    }
    if(repeatPasswordInput?.value?.length < 8 ) {
        alert("That password is too short.")
        passwordInput.focus()
        return
    }
    resetPassword().then((res) => {
      console.log(res)
        if(res?.reset) {
            success()
        } else {
        }
    }).then(() => {
    })
}

let form;

</script>

<div class="authbox pa3 gr-center flex flex-column w-100">
    <form method="POST" action="/signup/complete" bind:this={form} >
    <div class="">
        Set a new password for you Commune account.
    </div>
    <div class="pt3">
        <input type="password" name="password"
        bind:value={password}
        bind:this={passwordInput}
        placeholder="password" required>
        <input type="email" name="email"
        bind:value={email} hidden>
    </div>
    <div class="pt3">
        <input type="password" name="repeatpassword"
        bind:value={repeatPassword}
        bind:this={repeatPasswordInput}
        placeholder="repeat password" required>
        <input type="email" name="email"
        bind:value={email} hidden>
    </div>
    </form>
    {#if error}
        <div class="pt3 alert lh-copy">
            There was an error. Your password could not be reset. Try again later.
        </div>
    {/if}
    <div class="pt3 flex">
        <div class="flex-one"></div>
        <div class="">
            <button class="" on:click={reset}>
                reset password
            </button>
        </div>
    </div>
</div>


<style>
.oops {
  border: 1px solid red;
  box-shadow: inset 0 0 0 1px red;
}

.alert {
    color: red;
    font-weight: bold;
}

.checking {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
}

</style>

