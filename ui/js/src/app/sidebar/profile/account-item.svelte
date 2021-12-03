<script>
import { store } from '../../store/store.js'
export let account;

$: accountName = format(account.user_id)

$: selected = $store.active_account == account.user_id

function format() {
    return `${account.username}@${account.home_server}`
}

function switchToAccount() {
    if(account?.user_id != $store.active_account) {
        store.switchToAccount({
            user_id: account?.user_id,
        })
    }
}

</script>

<div class="account-item flex mb2" on:click={switchToAccount}>
    <div class="gr-center mr3 pointer">
        <div class="rad gr-default">
            {#if selected}
                <div class="radi gr-center"></div>
            {/if}
        </div>
    </div>
    <div class="gr-center flex-one pointer">
        {accountName}
    </div>
    <div class="gr-center">
        <button>Log Out</button>
    </div>
</div>

<style>
.account-item {

}

button {
    background-color: var(--background-3);
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
}

.rad {
    border-radius: 50%;
    border: 2px solid var(--white);
    width: 20px;
    height: 20px;
}

.radi {
    border-radius: 50%;
    border: 3px solid var(--background-2);
    background-color: var(--green);
    width: 100%;
    height: 100%;
}

.account-item:hover .rad{
    background-color: var(--green);
}

</style>
