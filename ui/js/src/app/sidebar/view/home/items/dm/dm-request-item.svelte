<script>
import { store } from '../../../../../store/store.js'
import { user as userIcon } from '../../../../../utils/icons.js'

export let request;

function strip(id) {
    let x= id.split(":")[0]
    return x.substring(1)
}

$: displayNameExists = request?.sender?.display_name?.length > 0
$: avatarExists = request?.sender?.avatar_url?.length > 0

$: name = displayNameExists ? request?.sender.display_name :
    request?.sender?.username

$: avatar = `${homeServer}/_matrix/media/r0/download/${request?.sender?.avatar_url?.substring(6)}`


function reject() {
    store.rejectDMRequest(request.room_id)
}

function accept() {
    store.acceptDMRequest(request.room_id)
}

</script>

<div class="dm-request-item mh1 pb1 flex">
    <div class="pr2">
        {#if avatarExists}
            <div class="profile-avatar ncol bg-img"
                style="background-image: url({avatar});">
            </div>
        {:else}
            <div class="profile-avatar gr-default">
                <div class="log gr-center">
                    {@html userIcon}
                </div>
            </div>
        {/if}
    </div>
    <div class="flex-one flex flex-column gr-default" >
        <div class="gr-start-center clmp-1">
            {request?.sender?.username}
        </div>
    </div>
</div>

<div class="pad pl1 flex mb1">
    <div class="">
        <button on:click={accept}>Accept</button>
    </div>
    <div class="ml2">
        <button class="reject" on:click={reject}>Reject</button>
    </div>
</div>

<style>
.dm-request-item {
    padding: 0.25rem 0.5rem;
}
.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
}
.ncol {
    background-color: transparent;
}
.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.pad {
    margin-left: calc(32px + 1rem);
}
button {
    padding: 0.125rem 0.25rem;
}
.reject {
    background-color: var(--red);
}
</style>
