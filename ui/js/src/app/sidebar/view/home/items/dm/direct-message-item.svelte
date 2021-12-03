<script>
import {onMount} from 'svelte'
import { store } from '../../../../../store/store.js'
import { user as userIcon } from '../../../../../utils/icons.js'
import { navigate, useLocation } from 'svelte-navigator'
import DirectMessageThread from './direct-message-thread.svelte'

export let item;
export let sender;

const location = useLocation()

$: pathname = $location.pathname

$: notMe = Object.fromEntries(
    Object.entries(item?.members).filter(([key, _]) => key != $store.active_account) )

$: others = Object.entries(notMe)
$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]


$: member = $store.allMembers[item?.sender]

$: username =  strip(item?.sender)

function strip(id) {
    let x= id.split(":")[0]
    return x.substring(1)
}

$: displayNameExists = member?.display_name?.length > 0
$: avatarExists = member?.avatar_url?.length > 0

$: name = displayNameExists ? member.display_name : username

$: avatar = `${homeServer}/_matrix/media/r0/download/${member?.avatar_url?.substring(6)}`

function names() {
    let n = []
    for (const [id, user] of Object.entries(notMe)) {
        if(user.display_name?.length > 0) {
            n.push(user.display_name)
        } else {
            n.push(strip(id))
        }
    }
    return n.join(", ")
}

$: multiple = others?.length > 1

function go() {
    navigate(`/messages/${item.alias}`, {replace:true})
    store.updateActiveHomePage(`/messages/${item.alias}`)
    store.updateActiveDirectMessages(`/messages/${item.alias}`)
    setTimeout(() => {
        store.resetDMnotification(item.room_id)
    }, 1000)
    /*
    if($store.isMobile && $store.mobileViewToggled) {
        $store.mobileViewToggled = false
    }
    */
}

$: indicate = sender == `${item.alias}`

$: if(indicate) {
    document.title = item.name
    localStorage.setItem('last-active-page', `/messages/${item.alias}`)
}


onMount(() => {
    if(indicate) {
        store.updateActiveHomePage(`/messages/${item.alias}`)
        store.updateActiveDirectMessages(`/messages/${item.alias}`)
        setTimeout(() => {
            store.resetDMnotification(item.room_id)
        }, 1000)
    }
})

function print(e) {
    e.preventDefault()
    console.log(item)
}

$: children = item.rooms
//$: children = $store.allRooms.filter(x => x.thread_in_room_id == item.room_id)

$: hasChildren = children?.length > 0
$: single = children?.length == 1

$: dmRoom = account?.dm_notifications[item.room_id]

</script>

<div class="direct-message-item-container pb1 mh1" 
    on:contextmenu={print}
    on:click={go}>
    <div class="direct-message-item flex" 
        class:active={indicate}>
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
        <div class="flex-one flex flex-column gr-default" 
            class:usa={indicate}>
            <div class="gr-start-center clmp-1">
                {names()}
            </div>
            {#if multiple}
                <div class="mem gr-start-center mute">
                    {others.length} members
                </div>
            {/if}
        </div>
        {#if dmRoom}
            <div class="gr-default">
                <div class="dm-count gr-center">
                    {dmRoom.count}
                </div>
            </div>
        {/if}
    </div>

</div>

    {#if hasChildren}
        <div class="flex room-item-threads mb2">

            <div class="flex-one flex-column flex">
                {#each children as thread, i (thread.room_id)}
                    <DirectMessageThread 
                    last={i == children?.length - 1}
                    single={single}
                    thread={thread} />
                {/each}

            </div>


        </div>
    {/if}

<style>
.direct-message-item {
    padding: 0.25rem 0.5rem;
    cursor: pointer;
    border-radius: 4px;
}
.direct-message-item:hover {
    background: var(--background-8);
}
.active {
    background: var(--background-8);
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
.usa {
    color: var(--white);
}
.spine {
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-left: 2px solid var(--text-muted);
    position: relative;
}

.spine:after {
    background: var(--text-muted);
    content: "";
    height: 2px;
    position: absolute;
    top: 16px;
    width: 10px;
}

.spine-last {
    margin-top: -15px;
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-radius: 0 0 0 4px;
}

.spine-last:after {
    top: 0;
    margin-top: 0;
    margin-left: 0;
    width: 10px;
    height: 18px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
    content: "";
    position: absolute;
}

.spine-single {
    margin-left: 15px;
    width: 0.5rem;
    height: 20px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
}
.thread-name {
    border-radius: 5px;
    cursor: pointer;
    color: var(--text-light);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
}

.thread-name:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.thread-name-ac {
    border-radius: 5px;
    cursor: pointer;
    color: var(--white);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
    background-color: var(--background-6);
}

.thread-name-ac:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.room-item-threads {
    padding: 0 0 0 0.5rem;
}

.mem {
    margin-top: 0.125rem;
    font-size: 0.8rem;
}

.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.dm-count {
    min-width: 16px;
    background: crimson;
    border-radius: 4px;
    padding: 0.125rem 0.25rem;
    font-size: 0.7rem;
    font-weight: bold;
    color: var(--white);
    line-height: 0.7rem;
    text-align: center;
}

</style>
