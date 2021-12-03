<script>
import { store } from '../../../../store/store.js'
import { makeid } from '../../../../utils/utils.js'

export let room
export let event
export let key
export let isTopic
export let OP;


$: topicEvents = $store.allReplies[event.room_id]

$: events = (isTopic && !OP) ? topicEvents :  $store.events?.[room?.room_id]?.events

$: roomReactions =events?.filter(x => x.type == 'm.reaction')

$: eventReactions = roomReactions?.filter(x => x.content?.['m.relates_to']?.event_id == event.event_id)

$: thisReaction = eventReactions?.filter(x => x.content?.['m.relates_to']?.key == key)

$: reacted = thisReaction?.filter(x => x.sender == $store.active_account)[0]


$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

function react() {
    if(reacted) {
        store.removeEvent(room?.room_id, reacted)
        return
    }

    if(!reacted) {
        let content = {
            'm.relates_to': {
                event_id: event.event_id,
                key: key,
                rel_type: 'm.annotation',
            }
        }

        let tempid = makeid(32)

        let eventType = `m.reaction`

        matrix.sendEvent(room.room_id, eventType, content, tempid, (err, res) => {
            console.log(res)

        });
    }
}

function printInfo(e) {
    e.preventDefault()
    console.log(thisReaction)
}

</script>

<div class="reaction-item ph1 mr2 mv1" 
    on:contextmenu|self={printInfo}
    on:click={react}
    class:reacted={reacted}>
    <span>
        {@html key}
    </span>
    <span class="ml1">{thisReaction?.length}</span>
</div>

<style>
.reaction-item {
    border-radius: 9px;
    background-color: var(--background-2);
    border: 1px solid transparent;
    padding: 0.25rem 0.5rem;
    cursor: pointer;
    display: inline-block;
}
.reaction-item:hover {
    border: 1px solid var(--text-muted);
}
.reacted {
    border: 1px solid var(--green);
    background-color: var(--background-3);
}
.reacted:hover {
    border: 1px solid var(--green);
}
</style>
