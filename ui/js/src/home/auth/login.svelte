<script>
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()


export let passwordReset;
export let desktop;

let usernameInput;
let username;
let passwordInput;
let password;

let loggingIn = false;

function showSignup() {
    dispatch('show-signup', true)
}

function showPassword() {
    dispatch('show-password', true)
}

onMount(() => {
    usernameInput.focus()
})

function onKeyPress(e) {
    if (e.charCode === 13) {
        login()
    }
}

let form;

async function validateLogin() {
    let endpoint = `/login/validate`
    let data = {
        username: username,
        password: password,
    };
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Content-Type': 'applicationjson'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let loginerror;

function login() {
    if(usernameInput.value.length == 0) {
        alert("That username doesn't look right.")
        usernameInput.focus()
        return
    }
    if(passwordInput.value.length < 8) {
        alert("That password is too short.")
        passwordInput.focus()
        return
    }

    if(desktop) {
        console.log(username, password)
        apiLogin()
        return
    }
    //form.submit()
    loggingIn = true
    validateLogin().then((res) => {
      console.log(res)
        if(res?.valid) {
            loginerror = false
            form.submit()
        } else {
            loginerror = true
            loggingIn = false
        }
    }).then(() => {
    })

}

async function fetchIdentity() {
    let endpoint = `http://localhost.com:8989/api/v0/login`
    let data = {
        username: username,
        password: password,
    };
    let resp = await window.__TAURI__.http.fetch(endpoint, {
        method: 'POST',
        body: {
            type: 'Json',
            payload: data,
        },
        headers:{
            'Content-Type': 'text/plain'
        }
    })
    const ret = await resp
    return Promise.resolve(ret)
}

function apiLogin() {
    fetchIdentity().then((res) => {
      console.log(res)
    }).then(() => {
    })
}

</script>


<div class="authbox pa3 gr-center flex flex-column w-100">
    {#if passwordReset}
    <div class="pb3">
        Your password was successfully reset.
    </div>
    {/if}
    <form bind:this={form} method="POST" action="/login">
    <div class="">
        <input type="text" name="username"
        bind:value={username}
        bind:this={usernameInput}
        placeholder="username"/>
    </div>
    <div class="pt3">
        <input type="password" name="password"
        bind:value={password}
        bind:this={passwordInput}
        on:keypress={onKeyPress}
        placeholder="password"/>
    </div>
    </form>
    {#if loginerror}
    <div class="red pt3 lh-copy">
        Username or password did not match.
    </div>
    {/if}
    <div class="pt3 flex">
        {#if loggingIn}
            <div class="ml2 gr-center spinner-s mr2">
            </div>
        {:else}
            <div class="gr-center">
                <button class="" on:click={login}>log in</button>
            </div>
        {/if}
        <div class="flex-one"></div>
        <div class="gr-center" on:click={showSignup}>
            <span class="link">sign up</span>
        </div>
    </div>
</div>
<div class="pt3 flex gr-center">
    <span class="link smaller" on:click={showPassword}>
        Forgot password?
    </span>
</div>
