<script>
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

import {debounce} from '../../utils/utils.js'

export let email;

let usernameInput;
let username = '';
let passwordInput;
let password;

function goBack() {
    dispatch('go-back', true)
}

onMount(() => {
    usernameInput.focus()
})

let checking = false;
let available = false;
let usernameAvailable = true;

function updateUsername(e) {
  const letters = /^[0-9a-zA-Z-]+$/;
  if(!e.key.match(letters)){
    e.preventDefault()
  }
  usernameAvailable = true
  available = false
  if(username.length === 0) {
      checking = false
      return
  }
}

function resetUsername() {
  usernameInput.value = ''
  username = ''
  available = false;
  usernameAvailable = true;
  usernameInput.focus()
}


async function checkUsername() {
    let endpoint = `/username/available`
    let data = {
        username: username,
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

function reset() {
  available = false
    checking = true
    if(usernameInput.value.length == 0) {
        resetUsername()
    }
  debounce(() =>{
    if(usernameInput.value.length < 3) {
      checking = false
      usernameAvailable = false
      available = false
      return
    }
    checkUsername().then((res) => {
      console.log(res)
        if(res?.available) {
          usernameAvailable = true
          available = true
        } else if(!res?.available) {
          usernameAvailable = false
          available = false
        }
        checking = false
    }).then(() => {
    })

  }, 500, this)
}

async function createUser() {
    let endpoint = `/signup/complete`
    let data = {
        username: username,
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

function create() {
    if(!usernameAvailable || !available) {
        alert("That username is not available.")
        usernameInput.focus()
        return
    }
    if(passwordInput?.value?.length < 8 ) {
        alert("That password is too short.")
        passwordInput.focus()
        return
    }
    /*
    createUser().then((res) => {
      console.log(res)
        if(res?.created == false) {
            error = true
        }
    }).then(() => {
    })
    */
        creating = true
        form.submit()
}

function onKeyPress(e) {
    if (e.charCode === 13) {
        create()
    }
}
let form;

let creating;

</script>

<div class="authbox pa3 gr-center flex flex-column w-100">
    <form method="POST" action="/signup/complete" bind:this={form} >
    <div class="relative">
        <input type="text" name="username"
        on:keypress={updateUsername}
        on:input={reset}
        bind:value={username}
        bind:this={usernameInput}
        autocomplete="off"
        minlength="3"
        placeholder="username" 
        class:oops={!usernameAvailable} 
        required>
        {#if checking}
          <div class="checking mh2 gr-default">
            <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
          </div>
        {/if}
        {#if !usernameAvailable}
          <div class="checking mh2 gr-default pointer" on:click={resetUsername}>
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
          </div>
        {/if}
        {#if available}
          <div class="checking mh2 gr-default">
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M21.03 5.72a.75.75 0 010 1.06l-11.5 11.5a.75.75 0 01-1.072-.012l-5.5-5.75a.75.75 0 111.084-1.036l4.97 5.195L19.97 5.72a.75.75 0 011.06 0z"></path></svg>
          </div>
        {/if}
    </div>
    <div class="pt3">
        <input type="password" name="password"
        bind:value={password}
        bind:this={passwordInput}
        on:keypress={onKeyPress}
        placeholder="password" required>
        <input type="email" name="email"
        bind:value={email} hidden>
    </div>
    </form>
    {#if error}
        <div class="pt3 alert lh-copy">
            There was an error. Your account could not be created. Try again later.
        </div>
    {/if}
    <div class="pt3 flex">
        <div class="gr-center">
            <div class="gr-center ico" on:click={goBack}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M10.78 19.03a.75.75 0 01-1.06 0l-6.25-6.25a.75.75 0 010-1.06l6.25-6.25a.75.75 0 111.06 1.06L5.81 11.5h14.44a.75.75 0 010 1.5H5.81l4.97 4.97a.75.75 0 010 1.06z"></path></svg>
            </div>
        </div>
        <div class="flex-one"></div>
        <div class="">
            <button class="" on:click={create} disabled={creating}>
                create account
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

