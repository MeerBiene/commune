<script>
import {store } from '../../store/store.js'
import { onMount, tick, createEventDispatcher } from 'svelte'
import { useLocation } from 'svelte-navigator'
import {makeid} from '../../utils/utils.js'

import {EditorState, Selection} from "prosemirror-state"
import {EditorView} from "prosemirror-view"
import {Schema, DOMParser} from "prosemirror-model"
import {schema} from "./schema.js"
import {addMentionNodes, addTagNodes } from './plugins/mentions/index.js'
import {setup} from "./setup/index.js"
import marked from 'marked'

marked.setOptions({
  renderer: new marked.Renderer(),
  pedantic: false,
  gfm: true,
  breaks: false,
  smartLists: true,
  smartypants: false,
  xhtml: true
});

const dispatch =createEventDispatcher()

export let focus;
export let page;
export let room;
export let initial;
export let inline;
export let editing;
export let threadView;
export let topics;
export let show;
export let recording;

export let disableMentions;

$: roomAlias = (topics || show) ? room?.alias : page?.pathname?.split("/")[2]

$: id = makeid(8)

$: containerClass = `editor-${roomAlias}-${id}`
$: id = `in-${id}`

$: roomID = strip(room?.room_id)

function strip(id) {
    let x = id?.split(":")[0]
    return x?.substring(1)
}


$: if(store.focusEditorInRoomID == room.room_id) {
    if(!topics) {
        focusEditor()
    }
}

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: insertUser = $store.messages?.mentionUser?.room_id == room?.room_id 

$: if(insertUser && view && state) {
    let username = $store.messages.mentionUser.username

    let attrs = {
        username: username,
        user_id: `@${username}:${account.home_server}`
    };

    let node = view.state.schema.nodes["mention"].create(attrs)
    let from = view.state.selection.from
    let to = view.state.selection.to
    let tr = view.state.tr.replaceWith(from, to, node);
    tr = view.state.tr.replaceWith(from, to, node);
    tr = tr.insertText(" ")
    state = view.state.apply(tr);
    view.updateState(state);

    focusEditor()

    store.resetMentionUser()
}

let container;

let view;

export let placeholder = 'Message'

$: if(indicate && view && !topics) {
    setTimeout(() => {
        focusEditor()
    }, 10)
    setTimeout(() => {
        focusEditor()
    }, 100)
}

$: if((inline || editing) && view) {
    setTimeout(() => {
        const selection = Selection.atEnd(view.docView.node)
        const tr = view.state.tr.setSelection(selection)
        const state = view.state.apply(tr)
        view.updateState(state)
        focusEditor()
    }, 1)

}

export async function focusEditor() {
    await tick();
    view.focus()
    dispatch('focused', true)
}

export async function insertEmoji(key) {
    /*
    let emoji = schema.nodes.emoji
    await tick();
    console.log(view)
    let {$from} = view.state.selection, index = $from.index()
    if (view.dispatch)
      view.dispatch(state.tr.insert(emoji.create('bronto')))
    view.focus()
    let rendered =  twemoji.parse(key, {
        base: '/static/img/emoji/',
        className: 'emo',
        ext: '.svg',
        size: '72x72'
    })

    let el = document.createElement('img');
    el.innerHTML = rendered
    let img = el.querySelector('img')

    let em = {
        alt: key,
        src: img.src,
    }


    let node = view.state.schema.nodes["emoji"].create(em)
    let from = view.state.selection.from
    let to = view.state.selection.to
    */
    var tr = view.state.tr
    tr = tr.insertText(key)
    //tr = tr.insertText(" ")
    state = view.state.apply(tr);
    view.updateState(state);
    focusEditor()
    dispatchContent()
}


const location =  useLocation()

$: indicate = (threadView || topics) ? true : $location.pathname == page?.pathname

$: isMessageEvent = $location.pathname.split("/")[3] == `message`

let state;
let freshState;
let syncState;

onMount(() => {

    let nodes = schema.spec.nodes

    if(!disableMentions) {
        nodes = addTagNodes(addMentionNodes(nodes))
    }

    const mySchema = new Schema({
      nodes: nodes,
      marks: schema.spec.marks
    })

    let doc = DOMParser.fromSchema(mySchema).parse(content)

    if(initial) {
        let init = document.createElement("div");
        init.innerHTML = initial

        doc =  DOMParser.fromSchema(mySchema).parse(init)
    }


    state = EditorState.create({
    doc: doc,
      plugins: setup({
          schema: mySchema,
          placeholder: placeholder,
          containerClass: containerClass,
          inline: inline,
          topics: topics,
          show: show,
          disableMentions: disableMentions,
          room: room,
          id: id,
      })
    })
    freshState = state

    let viewOpts = {
        state: state,
    }

    viewOpts.dispatchTransaction =  (transaction) => {

        syncState = view.state.apply(transaction)
        view.updateState(syncState)
        dispatchContent()

        view.state.doc.descendants((node, pos) => {
            if(node.isTextblock) {
                //console.log(node)
            }
        })
    }

    view = new EditorView(container, viewOpts)

    if(!topics) {
        setTimeout(() => {
            setupListener()
        }, 10)
    }

    if(topics || show) {
        setTimeout(() => {
            setupFocusListener()
        }, 10)
    }

    setupLinkPasteListener()
    focusEditor()
})


function setupLinkPasteListener() {
    let el = document.querySelector(`.${containerClass}`)
    if(!el) {
        return
    }
    el.addEventListener('paste', (e) => {
        e.preventDefault()
        window.syncEventsPosition(room?.room_id)
            let expression = /(http|ftp|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])?/g
            let regex = new RegExp(expression);
     
            let paste = (event.clipboardData || window.clipboardData).getData('text');
            let matches = paste.match(regex);
            if(matches && matches.length > 0) {
                for(let i=0;i<matches.length; i++) {
                    if(i == 9) {
                        break
                    }
                    dispatch('new-link', {
                        href: matches[i],
                    })
                    window.syncEventsPosition(room?.room_id)
                }
            }
    }); 

}


function setupFocusListener() {
    let el = document.querySelector(`.${containerClass}`)
    el.addEventListener('focus', (e) => {
        dispatch('focused', true)
    })
    el.addEventListener('blur', (e) => {
        dispatch('focused', false)
    })
}

let keys = {enter: false, shift:false};

function dispatchContent() {
    let content = container.innerText

    let cc = document.querySelector(`.${containerClass}`).cloneNode( true );
    if(cc) {
        /*
        let imgs = cc.querySelectorAll('img')
        if(imgs?.length > 0) {
            imgs.forEach(img => {
                let key = img.alt
                var textnode = document.createTextNode(key);
                img.parentNode.replaceChild(textnode, img)
            })
        }
        */
        let spans = cc.querySelectorAll('span')
        if(spans?.length > 0) {
            spans.forEach(span => {
                span.removeAttribute('contenteditable')
            })
        }
    }

    let len = view.state.doc.textContent.length
    if(len != 0 ) {
        let xx = content.replaceAll(' ', '')
        if(xx.length == 0) {
            return
        }
    }
    //let markdown = marked(cc.innerHTML)
    dispatch("sync", {
        plain_text: content,
        html: cc.innerHTML,
        length: view.state.doc.textContent.length
    })
    cc = null
}

function setupListener() {
    container.addEventListener('keydown', (e) => {
        if(inline) {
            if(e.key == 'Escape') {
                dispatch('cancel-edit', true)
            }
        }
        if(e.key == 'Shift') {
            keys.shift = true
        }
        if(e.key == 'Enter') {
            keys.enter = true
            if(keys.shift) {
                window.syncEventsPosition(room?.room_id)
                return
            } else {
                if(window.insertedEmoji || window.insertedMention || window.insertedTag) {
                    return
                }
                let length = view.state.doc.textContent.length


                if(editing && length == 0) {
                    return
                }

                dispatch('enter', true)
                if(length < 2000) {
                    view.updateState(freshState);
                }

            }
        }
    })
    container.addEventListener('keyup', (e) => {
        if(e.key == 'Shift') {
            keys.shift = false
        }
        if(e.key == 'Enter') {
            keys.enter = false
        }
    })
}

export function reset() {
    view.updateState(freshState);
}

$: {
    if(!keys.shift && keys.enter) {
        //view.updateState(newState);
    }
}

let content;

let _content = {
    html: null,
    plain_text: null,
    length: 0,
};

export function getContent() {
    updateContent(state)
    return _content
}

let updateContent = (state) => {

    _content = {
        plain_text: document.querySelector(`.${containerClass}`).innerText,
        length: state.doc.textContent.length,
    }

    let html = document.querySelector(`.${containerClass}`).innerHTML
    let plc = document.querySelector(`.${containerClass}`).querySelector('.placeholder')
    if(plc) {
        return
    }

    //let html = getHTML(state, state.doc.type.schema)
    let div = document.createElement("div");
    div.innerHTML = html
    _content.html = div.innerHTML

}


$: noDisplay = !indicate && !inline && !threadView && !topics && !show
$: display = indicate && !inline && !topics

</script>

<div class="input-root"
    id={id}
    class:no-dis={(noDisplay || recording) && !isMessageEvent}
    class:dis={display && !recording}>
<div class="input-box" 
    class:pr3={!topics}
    class:topics={topics || show}
    class:active-editor={indicate}
    bind:this={container}></div>


</div>

<div class="no-dis" bind:this={content}></div>


<style>

.input-root {
    display: grid;
    grid-template-rows: repeat();
}


.topics {
    height: 100%;
    display: grid;
    overflow: hidden;
    max-height: 40vh;
}

.dis {
    display: grid;
}
.no-dis {
    display: none;
}
</style>
