<script>
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

import {debounce} from '../../utils/utils.js'

export let email;

let usernameInput;
$: username = generateUsername(13)
let passwordInput;
let password;

function goBack() {
    dispatch('go-back', true)
}

onMount(() => {
    passwordInput.focus()
})

function generateUsername(length) {
    var result           = '';
    var characters       = '0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
    let username  = `user${result}`
   return username;
}

async function createUser() {
    let endpoint = `/signup/complete`
    let data = {
        username: username,
        password: password,
        anonymous: true,
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
        form.submit()
}

let form;

</script>

<div class="authbox pa3 gr-center flex flex-column w-100">
    <form method="POST" action="/signup/complete" bind:this={form} >
    <div class="">
        That's a nice username, right?
    </div>
    <div class="pt3">
        <input type="text" name="username"
        bind:value={username}
        disabled>
    </div>
    <div class="pt3">
        <input type="password" name="password"
        bind:value={password}
        bind:this={passwordInput}
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
            <button class="" on:click={create}>
                create anynomous account
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

