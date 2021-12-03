<script>
import { onMount, tick, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

import {debounce} from '../../utils/utils.js'


let showVerification = false;

let emailInput;
let email;

function goBack() {
    dispatch('go-back', true)
}

function showPasswordReset(email) {
    dispatch('password-reset', {
        email: email,
    })
}

onMount(() => {
    if(emailInput) {
        emailInput.focus()
    }
})


async function sendVerificationEmail() {
    let endpoint = `/password/verification`
    let data = {
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

let sending = false;

function onKeyPress(e) {
    if (e.charCode === 13) {
        sendEmail()
    }
}

let notfound = false;

function sendEmail() {
    if(emailInput?.value?.length == 0 ) {
        alert("You need to enter a valid email.")
        emailInput.focus()
        return
    }
    sending = true
    sendVerificationEmail().then((res) => {
      console.log(res)
        if(res?.emailed) {
            showVerification = true
        } else {
            notfound = true
        }
    }).then(() => {
        sending = false
    })
}

$: if(showVerification) {
    focusCodeInput()
}

async function focusCodeInput() {
    await tick()
    codeInput.focus()
}


let code;
let codeInput;

async function verifyVerifictionCode() {
    let endpoint = `/password/verification/validate`
    let data = {
        code: code,
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

function onCodeKeyPress(e) {
    if (e.charCode === 13) {
        verifyCode()
    }
}

let invalid = false;
let verifying;
function verifyCode() {
    if(codeInput?.value?.length == 0 ) {
        alert("That verification code is not the right length.")
        codeInput.focus()
        return
    }
    verifying = true
    verifyVerifictionCode().then((res) => {
      console.log(res)
        verifying = false
        if(res?.valid && res?.email) {
            showPasswordReset(res.email)
        } else {
            invalid = true
        }
    }).then(() => {
    })
}

function reset() {
    invalid = false
    verifying = false
}

</script>

{#if !showVerification}
<div class="authbox pa3 gr-center flex flex-column w-100">
    <div class="">
        Which account are you trying to recover?
    </div>
    <div class="pt3">
        <input type="email" name="email"
        bind:value={email}
        bind:this={emailInput}
        on:keypress={onKeyPress}
        placeholder="bob@bob.com" required>
    </div>
    {#if notfound}
    <div class="pt3 lh-copy">
        We couldn't find an account with that email.
    </div>
    {/if}
    <div class="pt3 flex">
        <div class="gr-center">
            <div class="gr-center ico" on:click={goBack}>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M10.78 19.03a.75.75 0 01-1.06 0l-6.25-6.25a.75.75 0 010-1.06l6.25-6.25a.75.75 0 111.06 1.06L5.81 11.5h14.44a.75.75 0 010 1.5H5.81l4.97 4.97a.75.75 0 010 1.06z"></path></svg>
            </div>
        </div>
        <div class="flex-one"></div>
        <div class="" class:sending={sending}>
            <button class="" 
                disabled={sending}
                on:click={sendEmail}>
                {#if sending}
                    emailing...
                {:else}
                    email verification code
                {/if}
            </button>
        </div>
    </div>
</div>
{/if}


{#if showVerification}
<div class="authbox pa3 gr-center flex flex-column w-100">
    <div class="lh-copy">
        We emailed you a verification code. Enter the code below to set up your account.
    </div>
    <div class="pt3">
        <input type="text" name="code"
        bind:value={code}
        bind:this={codeInput}
        on:keypress={onCodeKeyPress}
        class:oops={invalid}
        on:input={reset}
        placeholder="verification code" required>
    </div>
    <div class="pt3 flex">
        {#if invalid}
            <div class="invalid gr-center">
                Invalid code
            </div>
        {/if}
        <div class="flex-one"></div>
        <div class="">
            <button class="" 
                disabled={verifying}
                on:click={verifyCode}>
                {#if verifying}
                    verifying
                {:else}
                    verify code
                {/if}
            </button>
        </div>
    </div>
</div>
{/if}


<style>
.oops {
  border: 1px solid red;
  box-shadow: inset 0 0 0 1px red;
}

.invalid {
    color:red;
    font-weight: bold;
}

.checking {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
}

.sending {
    opacity: 0.7;
}

</style>

