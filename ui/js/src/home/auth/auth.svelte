<script>
import { onMount, tick } from 'svelte'
import Signup from './signup.svelte'
import Login from './login.svelte'
import ForgotPassword from './forgot-password.svelte'
import PasswordReset from './password-reset.svelte'
import VerifyEmail from './verify-email.svelte'
import InviteCode from './invite-code.svelte'

export let desktop;

let page = "login"

$: login = page == "login"
$: signup = page == "signup"
$: password = page == "password"
$: passwordReset = page == "password-reset"
$: verifyEmail = page == "verify-email"
$: inviteCode = page == "invite-code"


let usernameInput;
let emailInput;

onMount(() => {
})


function showLogin() {
    page = "login"
}

let email = ''
let passwordResetEmail = ''

function showSignup(e) {
    email = e.detail.email
    page = "signup"
}

function showPassword() {
    page = "password"
}

function showPasswordReset(e) {
    passwordResetEmail = e.detail.email
    page = "password-reset"
}

function showEmailSent() {
    page = "email-sent"
}

function showVerifyEmail() {
    page = "verify-email"
}

let resetPassword = false;

function showForgotPassword() {
    resetPassword = true
    page = "verify-email"
}

let passChanged = false;

function resetSuccess() {
    passChanged = true
    showLogin()
}

function showInviteCode() {
    page = 'invite-code'
}


async function focusUsername() {
    await tick()
    usernameInput.focus()
}

async function focusEmail() {
    await tick()
    emailInput.focus()
}

</script>

{#if login}
    <Login 
    desktop={desktop}
    passwordReset={passChanged}
    on:show-password={showPassword}
    on:forgot-password={showForgotPassword}
    on:show-signup={showSignup} />
{/if}


{#if signup}
    <Signup 
    email={email}
    on:email-sent={showEmailSent} 
    on:go-back={showLogin} />
{/if}

{#if password}
    <ForgotPassword 
    on:password-reset={showPasswordReset}
    on:go-back={showLogin} />
{/if}

{#if passwordReset}
    <PasswordReset 
    on:success={resetSuccess}
    email={passwordResetEmail}/>
{/if}

{#if verifyEmail}
    <VerifyEmail 
    on:go-back={showLogin} 
    on:verified={showSignup} />
{/if}


{#if inviteCode}
    <InviteCode 
    on:go-back={showLogin} 
    on:show-signup={showVerifyEmail} />
{/if}

