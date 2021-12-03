<script>
import { navigate } from 'svelte-navigator'
import { store } from '../../../../store/store.js'
import { user } from '../../../../utils/icons.js'
import { formatThreadTime, threadTimeAgo } from '../../../../utils/time.js'

export let thread;
export let last;

function go() {
    let pathname = `/${thread.server_alias}/${thread.alias}`
    navigate(pathname)
    store.updateActiveRoom(`/${thread.server_alias}`, pathname)
    store.updateActiveRooms(`/${thread.server_alias}`, pathname)
}

$: member = $store.allMembers[thread.sender]

$: hasAvatar = member?.avatar_url?.length > 0

$: avatar = `${homeServer}/_matrix/media/r0/download/${member?.avatar_url?.substring(6)}`

$: threadTime = formatThreadTime(thread?.thread_event?.origin_server_ts)
$: threadWhen = threadTimeAgo(thread?.thread_event?.origin_server_ts)

</script>

<div class="thread-item pa2" 
    class:mb3={!last}
    on:click={go}>
    <div class="fl-co">
        <div class="bold">
            {thread.name}
        </div>
        <div class="when mt2">
            <span class="">{threadWhen} {threadTime}</span>
        </div>
    </div>
    <div class="gr-center">
        {#if hasAvatar}
            <div class="profile-avatar bg-img gr-center"
                style="background-image: url({avatar});">
            </div>
        {:else}
            <div class="profile-avatar gr-default">
                <div class="log gr-center">
                    {@html user}
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
.thread-item {
    background-color: var(--background-3);
    border-radius: 7px;
    border: 1px solid transparent;
    cursor: pointer;
    transition: 0.1s;
    display: grid;
    grid-template-columns: 1fr auto;
}
.thread-item:hover {
    border: 1px solid var(--background-4);
}
.profile-avatar {
    width: 24px;
    height: 24px;
    background-color: var(--avatar);
    border-radius: 50%;
}
.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.when {
    font-size: 0.8rem;
}
</style>
