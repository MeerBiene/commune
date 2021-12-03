<script>
import { onMount, tick, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()


let emailInput;
let email;

function goBack() {
    dispatch('go-back', true)
}

function showSignup() {
    dispatch('show-signup', true)
}

function anon() {
    dispatch('anon', true)
}


onMount(() => {
    if(codeInput) {
        codeInput.focus()
    }
})



let code;
let codeInput;

async function verifyVerifictionCode() {
    let endpoint = `/invite/valid`
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
        if(res?.valid && res?.valid) {
            showSignup()
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

<div class="authbox pa3 gr-center flex flex-column w-100">
    <div class="lh-copy">
        You need an invite code to join Commune.
    </div>
    <div class="pt3">
        <input type="text" name="code"
        bind:value={code}
        bind:this={codeInput}
        on:keypress={onCodeKeyPress}
        class:oops={invalid}
        on:input={reset}
        placeholder="invite code" required>
    </div>
        {#if invalid}
            <div class="pt3 invalid">
                Invalid code
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
            <button class="" 
                disabled={verifying}
                on:click={verifyCode}>
                {#if verifying}
                    checking...
                {:else}
                    check code
                {/if}
            </button>
        </div>
    </div>
</div>


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

